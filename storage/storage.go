package storage

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Storage interface {
	Get(key string) (io.ReadCloser, error)
	Put(ctx context.Context, key string, body io.ReadSeeker) (string, error)
	SignGet(key string) (string, error)
	SignPut(key string) (string, error)
	StartMultipartUpload(key string) error
	UploadPart(key string, part int64, data io.ReadSeeker) error
	CompleteMultipartUpload(key string) error
	Rename(from string, to string) error
	Delete(key string) error
	Meta(key string) (*ObjectMeta, error)
	List(key string) ([]Object, error)
}

type ObjectMeta struct {
	ContentLength *int64
	ContentType   *string
}

type Object struct {
	Key          *string
	LastModified *time.Time
}

type Config struct {
	Type     string `json:"type"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	BaseURL  string `json:"base_url"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region"`
}

var storage Storage

func InitializeStorage(ctx context.Context) {
	baseConfig := Config{
		Type:     viper.GetString("storage.type"),
		Bucket:   viper.GetString("storage.bucket"),
		Key:      viper.GetString("storage.key"),
		Secret:   viper.GetString("storage.secret"),
		BaseURL:  viper.GetString("storage.base_url"),
		Endpoint: viper.GetString("storage.endpoint"),
		Region:   viper.GetString("storage.region"),
	}

	storage = configToStorage(ctx, baseConfig)

	if storage == nil {
		panic("Failed to initialize storage!")
	}

	log.Info().Msgf("Storage initialized: %s", baseConfig.Type)
}

func configToStorage(ctx context.Context, config Config) Storage {
	switch config.Type {
	case "wasabi":
		return initializeWasabi(ctx, config)
	case "b2":
		return initializeB2(ctx, config)
	case "s3":
		return initializeS3(ctx, config)
	}

	panic("Unknown storage type: " + viper.GetString("storage.type"))
}

func StartUploadMultipartMod(ctx context.Context, modID string, name string, versionID string) (bool, string) {
	if storage == nil {
		return false, ""
	}

	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	if err := StartMultipartUpload(key); err != nil {
		log.Err(err).Msg("failed to upload mod")
		return false, ""
	}

	return true, fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename))
}

func UploadMultipartMod(ctx context.Context, modID string, name string, versionID string, part int64, data io.ReadSeeker) (bool, string) {
	if storage == nil {
		return false, ""
	}

	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	if err := UploadPart(key, part, data); err != nil {
		log.Err(err).Msg("failed to upload mod")
		return false, ""
	}

	return true, fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename))
}

func CompleteUploadMultipartMod(ctx context.Context, modID string, name string, versionID string) (bool, string) {
	if storage == nil {
		return false, ""
	}

	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	if err := CompleteMultipartUpload(key); err != nil {
		log.Err(err).Msg("failed to upload mod")
		return false, ""
	}

	return true, fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename))
}

func UploadModLogo(ctx context.Context, modID string, data io.ReadSeeker) (bool, string) {
	if storage == nil {
		return false, ""
	}

	key := fmt.Sprintf("/images/mods/%s/logo.webp", modID)

	key, err := storage.Put(ctx, key, data)
	if err != nil {
		log.Err(err).Msg("failed to upload mod logo")
		return false, ""
	}

	return true, key
}

func UploadUserAvatar(ctx context.Context, userID string, data io.ReadSeeker) (bool, string) {
	if storage == nil {
		return false, ""
	}

	key := fmt.Sprintf("/images/users/%s/avatar.webp", userID)

	err := retry.Do(
		func() error {
			var err error
			key, err = storage.Put(ctx, key, data)
			return errors.Wrap(err, "failed to upload user avatar")
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			log.Err(err).Msgf("failed to upload user avatar, retrying [%d]", n)
		}),
	)
	if err != nil {
		log.Err(err).Msg("failed to upload user avatar")
		return false, ""
	}

	return true, key
}

func GenerateDownloadLink(key string) string {
	if storage == nil {
		return ""
	}

	url, err := storage.SignGet(key)
	if err != nil {
		return ""
	}

	return url
}

func StartMultipartUpload(key string) error {
	if storage == nil {
		return errors.New("storage not initialized")
	}

	return errors.Wrap(storage.StartMultipartUpload(key), "failed to start multipart upload")
}

func UploadPart(key string, part int64, data io.ReadSeeker) error {
	if storage == nil {
		return errors.New("storage not initialized")
	}

	return errors.Wrap(storage.UploadPart(key, part, data), "failed to upload part")
}

func CompleteMultipartUpload(key string) error {
	if storage == nil {
		return errors.New("storage not initialized")
	}

	return errors.Wrap(storage.CompleteMultipartUpload(key), "failed to complete multipart upload")
}

func CopyObjectFromOldBucket(key string) error {
	// Ignored
	return nil
}

func CopyObjectToOldBucket(key string) error {
	// Ignored
	return nil
}

func ScheduleCopyAllObjectsFromOldBucket(scheduler func(string)) {
	// Ignored
}

func Get(key string) (io.ReadCloser, error) {
	if storage == nil {
		return nil, errors.New("storage not initialized")
	}

	get, err := storage.Get(key)
	return get, errors.Wrap(err, "failed to get object")
}

func GetMod(modID string, name string, versionID string) (io.ReadCloser, error) {
	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	return Get(key)
}

func RenameVersion(ctx context.Context, modID string, name string, versionID string, version string) (bool, string) {
	if storage == nil {
		return false, ""
	}

	cleanName := cleanModName(name)

	from := fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(cleanName)+"-"+versionID)
	to := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+version)

	log.Info().Msgf("Renaming file from %s to %s", from, to)

	if err := storage.Rename(from, to); err != nil {
		log.Err(err).Msg("failed to rename version")
		return false, ""
	}

	fromUnescaped := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)
	if err := storage.Delete(fromUnescaped); err != nil {
		log.Err(err).Msg("failed to delete version")
		return false, ""
	}

	return true, fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(cleanName)+"-"+version)
}

func DeleteMod(ctx context.Context, modID string, name string, versionID string) bool {
	if storage == nil {
		return false
	}

	cleanName := cleanModName(name)

	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)

	log.Info().Str("key", key).Msg("deleting version")
	if err := storage.Delete(key); err != nil {
		log.Err(err).Msg("failed to delete version")
		return false
	}

	return true
}

func DeleteModTarget(ctx context.Context, modID string, name string, versionID string, target string) bool {
	if storage == nil {
		return false
	}

	cleanName := cleanModName(name)
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+target+"-"+versionID)

	log.Info().Str("key", key).Msg("deleting mod target")
	if err := storage.Delete(key); err != nil {
		log.Err(err).Msg("failed to delete version target")
		return false
	}

	return true
}

func ModVersionMeta(ctx context.Context, modID string, name string, versionID string) *ObjectMeta {
	if storage == nil {
		return nil
	}

	cleanName := cleanModName(name)

	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)

	meta, err := storage.Meta(key)
	if err != nil {
		log.Err(err).Msg("failed to delete version")
		return nil
	}

	return meta
}

func cleanModName(name string) string {
	cleanName := strings.ReplaceAll(name, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "\\", "_")
	cleanName = strings.ReplaceAll(cleanName, ":", "_")
	cleanName = strings.ReplaceAll(cleanName, "*", "_")
	cleanName = strings.ReplaceAll(cleanName, "?", "_")
	cleanName = strings.ReplaceAll(cleanName, "\"", "_")
	cleanName = strings.ReplaceAll(cleanName, "<", "_")
	cleanName = strings.ReplaceAll(cleanName, ">", "_")
	cleanName = strings.ReplaceAll(cleanName, "|", "_")
	cleanName = strings.ReplaceAll(cleanName, ";", "_")
	return strings.ReplaceAll(cleanName, "/", "_")
}

var encodeMapping = map[string]string{
	"\"": "%22",
	"#":  "%23",
	"&":  "%26",
	"+":  "%2B",
	",":  "%2C",
	"<":  "%3C",
	">":  "%3E",
	"?":  "%3F",
	"[":  "%5B",
	"\\": "%5C",
	"]":  "%5D",
	"^":  "%5E",
	"`":  "%60",
	"{":  "%7B",
	"|":  "%7C",
	"}":  "%7D",
}

func EncodeName(name string) string {
	// Must be first
	result := strings.ReplaceAll(name, "%", "%25")
	for k, v := range encodeMapping {
		result = strings.ReplaceAll(result, k, v)
	}
	return result
}

func SeparateModTarget(ctx context.Context, body []byte, modID, name, modVersion, target string) (bool, string, string, int64) {
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return false, "", "", 0
	}

	cleanName := cleanModName(name)

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, file := range zipReader.File {
		if !strings.HasPrefix(file.Name, target+"/") && file.Name != target+"/" {
			continue
		}

		err = copyModFileToArchZip(file, zipWriter, strings.TrimPrefix(file.Name, target+"/"))

		if err != nil {
			log.Err(err).Msg("failed to add file to " + target + " archive")
			return false, "", "", 0
		}
	}

	zipWriter.Close()

	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+target+"-"+modVersion)

	_, err = storage.Put(ctx, key, bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Err(err).Msg("failed to save " + target + " archive")
		return false, "", "", 0
	}

	hash := sha256.New()
	hash.Write(buf.Bytes())

	return true, key, hex.EncodeToString(hash.Sum(nil)), int64(buf.Len())
}

func copyModFileToArchZip(file *zip.File, zipWriter *zip.Writer, newName string) error {
	fileHeader := file.FileHeader
	fileHeader.Name = newName

	zipFile, err := zipWriter.CreateHeader(&fileHeader)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	rawFile, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer rawFile.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rawFile)

	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	_, err = zipFile.Write(buf.Bytes())

	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}

func DeleteOldModAssets(modReference string, before time.Time) {
	list, err := storage.List(fmt.Sprintf("assets/mods/%s", modReference))
	if err != nil {
		log.Err(err).Msg("failed to list assets")
		return
	}

	for _, object := range list {
		if object.Key == nil {
			continue
		}

		if object.LastModified == nil || object.LastModified.Before(before) {
			if err := storage.Delete(*object.Key); err != nil {
				log.Err(err).Str("key", *object.Key).Msg("failed deleting old asset")
				return
			}
		}
	}
}

func UploadModAsset(ctx context.Context, modReference string, path string, data []byte) {
	if storage == nil {
		return
	}

	key := fmt.Sprintf("/assets/mods/%s/%s", modReference, strings.TrimPrefix(path, "/"))

	_, err := storage.Put(ctx, key, bytes.NewReader(data))
	if err != nil {
		log.Err(err).Str("path", path).Msg("failed to upload mod asset")
	}
}

func ListModAssets(modReference string) ([]string, error) {
	if storage == nil {
		return nil, errors.New("no storage defined")
	}

	list, err := storage.List(fmt.Sprintf("assets/mods/%s", modReference))
	if err != nil {
		return nil, errors.New("failed to list assets")
	}

	out := make([]string, len(list))
	for i, object := range list {
		if object.Key == nil {
			continue
		}

		out[i] = *object.Key
	}

	sort.Strings(out)

	return out, nil
}
