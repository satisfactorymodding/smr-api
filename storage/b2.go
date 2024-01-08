package storage

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/satisfactorymodding/smr-api/redis"
)

type B2 struct {
	BaseURL   string
	S3Client  *s3.S3
	S3Session *session.Session
	Config    Config
}

func initializeB2(ctx context.Context, config Config) *B2 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Key, config.Secret, ""),
		Endpoint:         aws.String("https://s3." + config.Region + ".backblazeb2.com"),
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Err(err).Msg("failed to create S3 session")
		return nil
	}

	s3Client := s3.New(newSession)

	return &B2{
		BaseURL:   config.BaseURL,
		S3Client:  s3Client,
		S3Session: newSession,
		Config:    config,
	}
}

func (b2o *B2) Get(key string) (io.ReadCloser, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	object, err := b2o.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(b2o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get object")
	}

	return object.Body, nil
}

func (b2o *B2) Put(ctx context.Context, key string, body io.ReadSeeker) (string, error) {
	cleanedKey := strings.TrimPrefix(key, "/")
	uploader := s3manager.NewUploader(b2o.S3Session)

	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(b2o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return cleanedKey, errors.Wrap(err, "failed to upload file")
	}

	return key, nil
}

func (b2o *B2) SignGet(key string) (string, error) {
	// Public Bucket
	cleanedKey := strings.TrimPrefix(key, "/")
	return fmt.Sprintf("%s/file/%s/%s", b2o.BaseURL, b2o.Config.Bucket, cleanedKey), nil
}

func (b2o *B2) SignPut(key string) (string, error) {
	// Unsupported at the moment
	return "", errors.New("Unsupported")
}

func (b2o *B2) StartMultipartUpload(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	upload, err := b2o.S3Client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(b2o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create multipart upload")
	}

	redis.StoreMultipartUploadID(cleanedKey, *upload.UploadId)

	return nil
}

func (b2o *B2) UploadPart(key string, part int64, data io.ReadSeeker) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	id := redis.GetMultipartUploadID(cleanedKey)

	response, err := b2o.S3Client.UploadPart(&s3.UploadPartInput{
		Body:       data,
		Bucket:     aws.String(b2o.Config.Bucket),
		Key:        aws.String(cleanedKey),
		PartNumber: aws.Int64(part),
		UploadId:   aws.String(id),
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload part")
	}

	redis.StoreMultipartCompletedPart(cleanedKey, *response.ETag, int(part))

	return nil
}

func (b2o *B2) CompleteMultipartUpload(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	id := redis.GetMultipartUploadID(cleanedKey)
	parts := redis.GetAndClearMultipartCompletedParts(cleanedKey)
	completedParts := make([]*s3.CompletedPart, len(parts))

	for part, etag := range parts {
		partInt, _ := strconv.ParseInt(part, 10, 64)
		completedParts[partInt-1] = &s3.CompletedPart{ETag: aws.String(etag), PartNumber: &partInt}
	}

	_, err := b2o.S3Client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(b2o.Config.Bucket),
		Key:             aws.String(cleanedKey),
		MultipartUpload: &s3.CompletedMultipartUpload{Parts: completedParts},
		UploadId:        aws.String(id),
	})

	return errors.Wrap(err, "failed to complete multipart upload")
}

func (b2o *B2) Rename(from string, to string) error {
	cleanedKey := strings.TrimPrefix(to, "/")

	_, err := b2o.S3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(b2o.Config.Bucket),
		CopySource: aws.String(b2o.Config.Bucket + from),
		Key:        aws.String(cleanedKey),
	})

	return errors.Wrap(err, "failed to copy object")
}

func (b2o *B2) Delete(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")

	for i := 0; i < 10; i++ {
		versions, err := b2o.S3Client.ListObjectVersions(&s3.ListObjectVersionsInput{
			Bucket:    aws.String(b2o.Config.Bucket),
			KeyMarker: aws.String(cleanedKey),
			Prefix:    aws.String(cleanedKey),
		})
		if err != nil {
			return errors.Wrap(err, "failed to list object versions")
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
			return nil
		}

		_, err = b2o.S3Client.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(b2o.Config.Bucket),
			Delete: &s3.Delete{
				Objects: objects,
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to delete objects")
		}
	}

	return nil
}

func (b2o *B2) Meta(key string) (*ObjectMeta, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	data, err := b2o.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(b2o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get object meta")
	}

	return &ObjectMeta{
		ContentLength: data.ContentLength,
		ContentType:   data.ContentType,
	}, nil
}

func (b2o *B2) List(prefix string) ([]Object, error) {
	out := make([]Object, 0)

	err := b2o.S3Client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(b2o.Config.Bucket),
		Prefix: aws.String(prefix),
	}, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, obj := range output.Contents {
			out = append(out, Object{
				Key:          obj.Key,
				LastModified: obj.LastModified,
			})
		}
		return true
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list objects")
	}

	return out, nil
}
