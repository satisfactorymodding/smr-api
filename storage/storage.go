package storage

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/avast/retry-go/v3"
	"github.com/satisfactorymodding/smr-api/db/postgres"

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
}

type ObjectMeta struct {
	ContentLength *int64
	ContentType   *string
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

	log.Ctx(ctx).Info().Msgf("Storage initialized: %s", baseConfig.Type)
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
		log.Ctx(ctx).Err(err).Msg("failed to upload mod")
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
		log.Ctx(ctx).Err(err).Msg("failed to upload mod")
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
		log.Ctx(ctx).Err(err).Msg("failed to upload mod")
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
		log.Ctx(ctx).Err(err).Msg("failed to upload mod logo")
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
			log.Ctx(ctx).Err(err).Msgf("failed to upload user avatar, retrying [%d]", n)
		}),
	)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to upload user avatar")
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

	log.Ctx(ctx).Info().Msgf("Renaming file from %s to %s", from, to)

	if err := storage.Rename(from, to); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to rename version")
		return false, ""
	}

	fromUnescaped := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)
	if err := storage.Delete(fromUnescaped); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to delete version")
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

	if err := storage.Delete(key); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to delete version")
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
		log.Ctx(ctx).Err(err).Msg("failed to delete version")
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

func SeparateMod(ctx context.Context, body []byte, modID, name string, versionID string) error {

	//read combined file

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))

	if err != nil {
		return errors.New("invalid zip archive")
	}

	// Create a buffer to write our archive to.
	fmt.Println("Beginning Splitting Operations")
	bufLinuxServer := new(bytes.Buffer)
	bufWin64Server := new(bytes.Buffer)
	bufWin64Client := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriterLinuxServer := zip.NewWriter(bufLinuxServer)
	zipWriterWin64Server := zip.NewWriter(bufWin64Server)
	zipWriterWin64Client := zip.NewWriter(bufWin64Client)

	var LinuxServer = false
	var Win64Server = false
	var Win64Client = false

	// Add some files to the archive.
	for _, file := range zipReader.File {

		if !strings.Contains(file.Name, "pdb") && !strings.Contains(file.Name, "debug") {

			if strings.Contains(file.FileHeader.Name, "LinuxServer") {
				file.FileHeader.Name = strings.ReplaceAll(file.FileHeader.Name, "LinuxServer/", "")
				zipWriterLinuxServer.Copy(file)
				LinuxServer = true
			}

			if strings.Contains(file.FileHeader.Name, "WindowsServer") {
				file.FileHeader.Name = strings.ReplaceAll(file.FileHeader.Name, "WindowsServer/", "")
				zipWriterWin64Server.Copy(file)
				Win64Server = true
			}

			if strings.Contains(file.FileHeader.Name, "WindowsNoClient") {
				file.FileHeader.Name = strings.ReplaceAll(file.FileHeader.Name, "WindowsNoClient/", "")
				zipWriterWin64Client.Copy(file)
				Win64Client = true
			}
		}
	}

	// Make sure to check the error on Close.
	err = zipWriterLinuxServer.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = zipWriterWin64Server.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = zipWriterWin64Client.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	//Write to mod_link and upload new smaller smod file
	if LinuxServer {
		cleanName := cleanModName(name)
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-LinuxServer-"+versionID)
		_, err = storage.Put(ctx, key, bytes.NewReader(bufLinuxServer.Bytes()))

		if err != nil {
			fmt.Println(err)
			return err
		}

		hash := sha256.New()
		_, err = hash.Write(bufLinuxServer.Bytes())

		dbModLink := &postgres.ModLink{
			ModVersionLinkID: versionID,
			Link:             key,
			Hash:             hex.EncodeToString(hash.Sum(nil)),
			Size:             int64(len(bufLinuxServer.Bytes())),
		}

		postgres.CreateModLink(ctx, dbModLink)
	}

	if Win64Server {
		cleanName := cleanModName(name)
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-Win64Server-"+versionID)
		_, err = storage.Put(ctx, key, bytes.NewReader(bufWin64Server.Bytes()))

		hash := sha256.New()
		_, err = hash.Write(bufWin64Server.Bytes())

		if err != nil {
			fmt.Println(err)
			return err
		}

		dbModLink := &postgres.ModLink{
			ModVersionLinkID: versionID,
			Link:             key,
			Hash:             hex.EncodeToString(hash.Sum(nil)),
			Size:             int64(len(bufWin64Server.Bytes())),
		}

		postgres.CreateModLink(ctx, dbModLink)
	}
	if Win64Client {
		cleanName := cleanModName(name)
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-WindowNoEditor-"+versionID)
		_, err = storage.Put(ctx, key, bytes.NewReader(bufWin64Client.Bytes()))

		hash := sha256.New()
		_, err = hash.Write(bufWin64Client.Bytes())

		if err != nil {
			fmt.Println(err)
			return err
		}

		dbModLink := &postgres.ModLink{
			ModVersionLinkID: versionID,
			Link:             key,
			Hash:             hex.EncodeToString(hash.Sum(nil)),
			Size:             int64(len(bufWin64Client.Bytes())),
		}

		postgres.CreateModLink(ctx, dbModLink)
	}

	return nil

}
