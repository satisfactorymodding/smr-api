package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Vilsol/slox"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/redis"
)

type S3 struct {
	BaseURL   string
	S3Client  *s3.S3
	S3Session *session.Session
	Config    Config
}

func initializeS3(ctx context.Context, config Config) *S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Key, config.Secret, ""),
		Endpoint:         aws.String(config.Endpoint),
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		slox.Error(ctx, "failed to create S3 session", slog.Any("err", err))
		return nil
	}

	s3Client := s3.New(newSession)

	return &S3{
		BaseURL:   config.BaseURL,
		S3Client:  s3Client,
		S3Session: newSession,
		Config:    config,
	}
}

func (s3o *S3) Get(key string) (io.ReadCloser, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	object, err := s3o.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object.Body, nil
}

func (s3o *S3) Put(ctx context.Context, key string, body io.ReadSeeker) (string, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	uploader := s3manager.NewUploader(s3o.S3Session)

	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(viper.GetString("storage.bucket")),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return cleanedKey, fmt.Errorf("failed to upload file: %w", err)
	}

	return key, nil
}

func (s3o *S3) SignGet(key string) (string, error) {
	// Public Bucket
	cleanedKey := strings.TrimPrefix(key, "/")
	return fmt.Sprintf(viper.GetString("storage.keypath"), s3o.BaseURL, viper.GetString("storage.bucket"), cleanedKey), nil
}

func (s3o *S3) SignPut(_ string) (string, error) {
	// Unsupported at the moment
	return "", errors.New("Unsupported")
}

func (s3o *S3) StartMultipartUpload(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	upload, err := s3o.S3Client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(viper.GetString("storage.bucket")),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart upload: %w", err)
	}

	redis.StoreMultipartUploadID(cleanedKey, *upload.UploadId)

	return nil
}

func (s3o *S3) UploadPart(key string, part int64, data io.ReadSeeker) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	id := redis.GetMultipartUploadID(cleanedKey)

	response, err := s3o.S3Client.UploadPart(&s3.UploadPartInput{
		Body:       data,
		Bucket:     aws.String(viper.GetString("storage.bucket")),
		Key:        aws.String(cleanedKey),
		PartNumber: aws.Int64(part),
		UploadId:   aws.String(id),
	})
	if err != nil {
		return fmt.Errorf("failed to upload part: %w", err)
	}

	redis.StoreMultipartCompletedPart(cleanedKey, *response.ETag, int(part))

	return nil
}

func (s3o *S3) CompleteMultipartUpload(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	id := redis.GetMultipartUploadID(cleanedKey)
	parts := redis.GetAndClearMultipartCompletedParts(cleanedKey)
	completedParts := make([]*s3.CompletedPart, len(parts))

	for part, etag := range parts {
		partInt, _ := strconv.ParseInt(part, 10, 64)
		completedParts[partInt-1] = &s3.CompletedPart{ETag: aws.String(etag), PartNumber: &partInt}
	}

	_, err := s3o.S3Client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(viper.GetString("storage.bucket")),
		Key:             aws.String(cleanedKey),
		MultipartUpload: &s3.CompletedMultipartUpload{Parts: completedParts},
		UploadId:        aws.String(id),
	})
	if err != nil {
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}

	return nil
}

func (s3o *S3) Rename(from string, to string) error {
	cleanedKey := strings.TrimPrefix(to, "/")

	_, err := s3o.S3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(viper.GetString("storage.bucket")),
		CopySource: aws.String(viper.GetString("storage.bucket") + from),
		Key:        aws.String(cleanedKey),
	})
	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

func (s3o *S3) Delete(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")

	for i := 0; i < 10; i++ {
		versions, err := s3o.S3Client.ListObjectVersions(&s3.ListObjectVersionsInput{
			Bucket:    aws.String(viper.GetString("storage.bucket")),
			KeyMarker: aws.String(cleanedKey),
			Prefix:    aws.String(cleanedKey),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") {
				_, err = s3o.S3Client.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String(viper.GetString("storage.bucket")),
					Key:    aws.String(cleanedKey),
				})

				if err != nil {
					return fmt.Errorf("failed to delete objects: %w", err)
				}

				return nil
			}

			return fmt.Errorf("failed to list object versions: %w", err)
		}

		objects := make([]*s3.ObjectIdentifier, len(versions.Versions)+len(versions.DeleteMarkers))

		for i, version := range versions.Versions {
			objects[i] = &s3.ObjectIdentifier{
				Key:       version.Key,
				VersionId: version.VersionId,
			}
		}

		for i, marker := range versions.DeleteMarkers {
			objects[i+len(versions.Versions)] = &s3.ObjectIdentifier{
				Key:       marker.Key,
				VersionId: marker.VersionId,
			}
		}

		if len(objects) == 0 {
			_, err = s3o.S3Client.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(viper.GetString("storage.bucket")),
				Key:    aws.String(cleanedKey),
			})

			if err != nil {
				return fmt.Errorf("failed to delete objects: %w", err)
			}

			return nil
		}

		_, err = s3o.S3Client.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(viper.GetString("storage.bucket")),
			Delete: &s3.Delete{
				Objects: objects,
			},
		})

		if err != nil {
			return fmt.Errorf("failed to delete objects: %w", err)
		}
	}

	return nil
}

func (s3o *S3) Meta(key string) (*ObjectMeta, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	data, err := s3o.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(viper.GetString("storage.bucket")),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object meta: %w", err)
	}

	return &ObjectMeta{
		ContentLength: data.ContentLength,
		ContentType:   data.ContentType,
	}, nil
}

func (s3o *S3) List(prefix string) ([]Object, error) {
	objects, err := s3o.S3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(viper.GetString("storage.bucket")),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	out := make([]Object, len(objects.Contents))
	for i, obj := range objects.Contents {
		out[i] = Object{
			Key:          obj.Key,
			LastModified: obj.LastModified,
		}
	}

	return out, nil
}
