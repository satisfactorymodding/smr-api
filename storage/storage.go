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

func DeleteVersion(ctx context.Context, modID string, name string, versionID string) bool {
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

func DeleteMod(ctx context.Context, modID string, name string, versionID string) bool {
	if storage == nil {
		return false
	}

	cleanName := cleanModName(name)

	query := postgres.GetModVersion(ctx, modID, versionID)

	if len(query.Arch) != 0 {
		for _, link := range query.Arch {
			if success := DeleteModArch(ctx, modID, cleanName, versionID, link.Platform); !success {
				return false
			}
		}
	} else {
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)

		if err := storage.Delete(key); err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to delete version")
			return false
		}
	}

	return true
}

func DeleteModArch(ctx context.Context, modID string, name string, versionID string, platform string) bool {
	if storage == nil {
		return false
	}

	cleanName := cleanModName(name)
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+platform+"-"+versionID)

	if err := storage.Delete(key); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to delete version link")
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

func SeparateMod(ctx context.Context, body []byte, modID, name string, versionID string, modVersion string) (bool, string) {
	//read combined file
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))

	if err != nil {
		return false, ""
	}

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

	cleanName := cleanModName(name)

	// Add some files to the archive.
	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, ".pdb") || strings.HasPrefix(file.Name, ".debug") || strings.Contains(file.Name, modVersion) {
			continue
		}

		if strings.Contains(file.Name, "LinuxServer") {
			err = WriteZipFile(ctx, file, "LinuxServer", zipWriterLinuxServer)

			if err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to write zip to LinuxServer smod")
				return false, ""
			}

			LinuxServer = true
		}

		if strings.Contains(file.Name, "WindowsServer") {
			err = WriteZipFile(ctx, file, "WindowsServer", zipWriterWin64Server)

			if err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to write zip to Win64Server smod")
				return false, ""
			}

			Win64Server = true
		}

		if strings.Contains(file.Name, "WindowsNoEditor") {
			err = WriteZipFile(ctx, file, "WindowsNoEditor", zipWriterWin64Client)

			if err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to write zip to WindowsNoEditor smod")
				return false, ""
			}

			Win64Client = true
		}
	}

	//Close files
	zipWriterLinuxServer.Close()
	zipWriterWin64Server.Close()
	zipWriterWin64Client.Close()

	//Write to mod_link and upload new smaller smod file
	if LinuxServer {
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-LinuxServer-"+modVersion)
		platform := "LinuxServer"

		err = WriteModArch(ctx, key, versionID, platform, bufLinuxServer)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("Failed to save LinuxServer smod")
			return false, ""
		}
	}

	if Win64Server {
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-Win64Server-"+modVersion)
		platform := "Win64Server"

		err = WriteModArch(ctx, key, versionID, platform, bufWin64Server)

		if err != nil {
			log.Ctx(ctx).Err(err).Msg("Failed to save Win64Server smod")
			return false, ""
		}
	}
	if Win64Client {
		key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-WindowsNoEditor-"+modVersion)
		platform := "WindowsNoEditor"

		err = WriteModArch(ctx, key, versionID, platform, bufWin64Client)

		if err != nil {
			log.Ctx(ctx).Err(err).Msg("Failed to save WindowsNoEditor smod")
			return false, ""
		}
	}

	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-Combined-"+modVersion)
	platform := "Combined"

	combined := bytes.NewBuffer(body)

	err = WriteModArch(ctx, key, versionID, platform, combined)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to save Combined smod versionID: " + versionID)
		return false, ""
	}

	key = fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-WindowsNoEditor-"+modVersion)

	return true, key
}

func WriteZipFile(ctx context.Context, file *zip.File, platform string, zipWriter *zip.Writer) error {
	fileName := strings.ReplaceAll(file.Name, platform+"/", "")
	zipFile, err := zipWriter.Create(fileName)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create smod file for " + platform)
		return errors.Wrap(err, "Failed to open smod file for "+platform)
	}

	rawFile, err := file.Open()

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to open smod file for " + platform)
		return errors.Wrap(err, "Failed to open smod file for "+platform)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rawFile)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to read from buffer for " + platform)
		return errors.Wrap(err, "Failed to read from buffer for "+platform)
	}

	_, err = zipFile.Write(buf.Bytes())

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to write to smod file: " + platform)
		return errors.Wrap(err, "Failed to write smod file for "+platform)
	}

	return nil
}

func WriteModArch(ctx context.Context, key string, versionID string, platform string, buffer *bytes.Buffer) error {
	_, err := storage.Put(ctx, key, bytes.NewReader(buffer.Bytes()))

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to write smod: " + key)
		return errors.Wrap(err, "Failed to load smod:"+key)
	}

	hash := sha256.New()
	hash.Write(buffer.Bytes())

	dbModArch := &postgres.ModArch{
		ModVersionArchID: versionID,
		Platform:         platform,
		Key:              key,
		Hash:             hex.EncodeToString(hash.Sum(nil)),
		Size:             int64(len(buffer.Bytes())),
	}

	_, err = postgres.CreateModArch(ctx, dbModArch)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create ModArch: " + versionID + "-" + platform)
		return errors.Wrap(err, "Failed to create ModArch: "+versionID+"-"+platform)
	}

	return nil
}
