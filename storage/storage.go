package storage

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/Vilsol/slox"
	"github.com/avast/retry-go/v3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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
	Keypath  string `json:"keypath"`
}

type storageKey struct{}

func InitializeStorage(ctx context.Context) context.Context {
	baseConfig := Config{
		Type:     viper.GetString("storage.type"),
		Bucket:   viper.GetString("storage.bucket"),
		Key:      viper.GetString("storage.key"),
		Secret:   viper.GetString("storage.secret"),
		BaseURL:  viper.GetString("storage.base_url"),
		Endpoint: viper.GetString("storage.endpoint"),
		Region:   viper.GetString("storage.region"),
		Keypath:  viper.GetString("storage.keypath"),
	}

	storage := configToStorage(ctx, baseConfig)

	if storage == nil {
		panic("Failed to initialize storage!")
	}

	if viper.IsSet("storage.reader.type") {
		reader := configToStorage(ctx, Config{
			Type:     viper.GetString("storage.reader.type"),
			Bucket:   viper.GetString("storage.reader.bucket"),
			Key:      viper.GetString("storage.reader.key"),
			Secret:   viper.GetString("storage.reader.secret"),
			BaseURL:  viper.GetString("storage.reader.base_url"),
			Endpoint: viper.GetString("storage.reader.endpoint"),
			Region:   viper.GetString("storage.reader.region"),
			Keypath:  viper.GetString("storage.reader.keypath"),
		})

		storage = initializeWrapper(reader, storage)
	}

	slox.Info(ctx, "storage initialized", slog.String("type", baseConfig.Type))

	return context.WithValue(ctx, storageKey{}, storage)
}

func Client(ctx context.Context) Storage {
	c := ctx.Value(storageKey{})
	if c == nil {
		return nil
	}
	return c.(Storage)
}

func TransferContext(source context.Context, target context.Context) context.Context {
	c := source.Value(storageKey{})
	if c == nil {
		return target
	}
	return context.WithValue(target, storageKey{}, c)
}

func configToStorage(ctx context.Context, config Config) Storage {
	if config.Type == "s3" {
		return initializeS3(ctx, config)
	}

	panic("Unknown storage type: " + viper.GetString("storage.type"))
}

func StartUploadMultipartMod(ctx context.Context, modID string, name string, versionID string) (string, error) {
	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	if err := StartMultipartUpload(ctx, key); err != nil {
		slox.Error(ctx, "failed to upload mod", slog.Any("err", err))
		return "", fmt.Errorf("failed to upload mod: %w", err)
	}

	return fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename)), nil
}

func UploadMultipartMod(ctx context.Context, modID string, name string, versionID string, part int64, data io.ReadSeeker) (string, error) {
	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	if err := UploadPart(ctx, key, part, data); err != nil {
		slox.Error(ctx, "failed to upload mod", slog.Any("err", err))
		return "", fmt.Errorf("failed to upload mod: %w", err)
	}

	return fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename)), nil
}

func CompleteUploadMultipartMod(ctx context.Context, modID string, name string, versionID string) (string, error) {
	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	if err := CompleteMultipartUpload(ctx, key); err != nil {
		slox.Error(ctx, "failed to upload mod", slog.Any("err", err))
		return "", fmt.Errorf("failed to upload mod: %w", err)
	}

	return fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename)), nil
}

func UploadModLogo(ctx context.Context, modID string, data io.ReadSeeker) (string, error) {
	key := fmt.Sprintf("/images/mods/%s/logo.webp", modID)

	key, err := Client(ctx).Put(ctx, key, data)
	if err != nil {
		slox.Error(ctx, "failed to upload mod logo", slog.Any("err", err))
		return "", fmt.Errorf("failed to upload mod logo: %w", err)
	}

	return key, nil
}

func UploadUserAvatar(ctx context.Context, userID string, data io.ReadSeeker) (string, error) {
	key := fmt.Sprintf("/images/users/%s/avatar.webp", userID)

	err := retry.Do(
		func() error {
			var err error
			key, err = Client(ctx).Put(ctx, key, data)
			if err != nil {
				return fmt.Errorf("failed to upload user avatar: %w", err)
			}
			return nil
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			slox.Error(ctx, "failed to upload user avatar, retrying", slog.Any("err", err), slog.Any("n", n))
		}),
	)
	if err != nil {
		slox.Error(ctx, "failed to upload user avatar", slog.Any("err", err))
		return "", fmt.Errorf("failed to upload user avatar: %w", err)
	}

	return key, nil
}

func GenerateDownloadLink(ctx context.Context, key string) string {
	url, err := Client(ctx).SignGet(key)
	if err != nil {
		return ""
	}

	return url
}

func StartMultipartUpload(ctx context.Context, key string) error {
	if err := Client(ctx).StartMultipartUpload(key); err != nil {
		return fmt.Errorf("failed to start multipart upload: %w", err)
	}

	return nil
}

func UploadPart(ctx context.Context, key string, part int64, data io.ReadSeeker) error {
	if err := Client(ctx).UploadPart(key, part, data); err != nil {
		return fmt.Errorf("failed to upload part: %w", err)
	}

	return nil
}

func CompleteMultipartUpload(ctx context.Context, key string) error {
	if err := Client(ctx).CompleteMultipartUpload(key); err != nil {
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}

	return nil
}

func Get(ctx context.Context, key string) (io.ReadCloser, error) {
	get, err := Client(ctx).Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return get, nil
}

func GetMod(ctx context.Context, modID string, name string, versionID string) (io.ReadCloser, error) {
	cleanName := cleanModName(name)

	filename := cleanName + "-" + versionID
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	return Get(ctx, key)
}

func RenameVersion(ctx context.Context, modID string, name string, versionID string, version string) (string, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "RenameVersion")
	defer span.End()

	cleanName := cleanModName(name)

	from := fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(cleanName)+"-"+versionID)
	to := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+version)

	slox.Info(ctx, "renaming file", slog.String("from", from), slog.String("to", to))

	if err := Client(ctx).Rename(from, to); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slox.Error(ctx, "failed to rename version", slog.Any("err", err))
		return "", fmt.Errorf("failed to rename version: %w", err)
	}

	fromUnescaped := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)
	if err := Client(ctx).Delete(fromUnescaped); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slox.Error(ctx, "failed to delete version", slog.Any("err", err))
		return "", fmt.Errorf("failed to delete version: %w", err)
	}

	return fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(cleanName+"-"+version)), nil
}

func DeleteMod(ctx context.Context, modID string, name string, versionID string) error {
	cleanName := cleanModName(name)

	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+versionID)

	slox.Info(ctx, "deleting version", slog.String("key", key))
	if err := Client(ctx).Delete(key); err != nil {
		slox.Error(ctx, "failed to delete version", slog.Any("err", err))
		return fmt.Errorf("failed to delete version: %w", err)
	}

	return nil
}

func DeleteModTarget(ctx context.Context, modID string, name string, versionID string, target string) error {
	cleanName := cleanModName(name)
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, cleanName+"-"+target+"-"+versionID)

	slox.Info(ctx, "deleting mod target", slog.String("key", key))
	if err := Client(ctx).Delete(key); err != nil {
		slox.Error(ctx, "failed to delete version target", slog.Any("err", err))
		return fmt.Errorf("failed to delete version target: %w", err)
	}

	return nil
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

func SeparateModTarget(ctx context.Context, body []byte, modID, name, modVersion, target string) (string, string, int64, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "SeparateModTarget")
	defer span.End()

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return "", "", 0, fmt.Errorf("failed to create zip reader: %w", err)
	}

	cleanName := cleanModName(name)

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, file := range zipReader.File {
		if !strings.HasPrefix(file.Name, target+"/") {
			continue
		}
		trimmedName := strings.TrimPrefix(file.Name, target+"/")
		if len(trimmedName) == 0 {
			continue
		}

		err = copyModFileToArchZip(file, zipWriter, trimmedName)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			slox.Error(ctx, "failed to add file to archive", slog.Any("err", err), slog.String("target", target))
			return "", "", 0, fmt.Errorf("failed to add file to archive: %w", err)
		}
	}

	zipWriter.Close()

	filename := cleanName + "-" + target + "-" + modVersion
	key := fmt.Sprintf("/mods/%s/%s.smod", modID, filename)

	_, err = Client(ctx).Put(ctx, key, bytes.NewReader(buf.Bytes()))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slox.Error(ctx, "failed to save archive", slog.Any("err", err), slog.String("target", target))
		return "", "", 0, fmt.Errorf("failed to save archive: %w", err)
	}

	hash := sha256.New()
	hash.Write(buf.Bytes())

	encodedKey := fmt.Sprintf("/mods/%s/%s.smod", modID, EncodeName(filename))
	return encodedKey, hex.EncodeToString(hash.Sum(nil)), int64(buf.Len()), nil
}

func copyModFileToArchZip(file *zip.File, zipWriter *zip.Writer, newName string) error {
	fileHeader := file.FileHeader
	fileHeader.Name = newName

	zipFile, err := zipWriter.CreateHeader(&fileHeader)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	rawFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer rawFile.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rawFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	_, err = zipFile.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func DeleteOldModAssets(ctx context.Context, modReference string, before time.Time) {
	list, err := Client(ctx).List(fmt.Sprintf("assets/mods/%s", modReference))
	if err != nil {
		slox.Error(ctx, "failed to list assets", slog.Any("err", err))
		return
	}

	for _, object := range list {
		if object.Key == nil {
			continue
		}

		if object.LastModified == nil || object.LastModified.Before(before) {
			if err := Client(ctx).Delete(*object.Key); err != nil {
				slox.Error(ctx, "failed deleting old asset", slog.Any("err", err), slog.String("key", *object.Key))
				return
			}
		}
	}
}

func UploadModAsset(ctx context.Context, modReference string, path string, data []byte) {
	slox.Info(ctx, "uploading asset", slog.String("mod_reference", modReference), slog.String("asset", path), slog.Int("size", len(data)))

	key := fmt.Sprintf("/assets/mods/%s/%s", modReference, strings.TrimPrefix(path, "/"))

	_, err := Client(ctx).Put(ctx, key, bytes.NewReader(data))
	if err != nil {
		slox.Error(ctx, "failed to upload mod asset", slog.Any("err", err), slog.String("path", path))
	}
}

func ListModAssets(ctx context.Context, modReference string) ([]string, error) {
	slox.Info(ctx, "listing assets", slog.String("mod_reference", modReference))

	list, err := Client(ctx).List(fmt.Sprintf("assets/mods/%s", modReference))
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
