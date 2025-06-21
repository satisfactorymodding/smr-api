package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/Vilsol/slox"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/satisfactorymodding/smr-api/redis"
)

type S3 struct {
	BaseURL  string
	S3Client *s3.Client
	Config   Config
}

func initializeS3(ctx context.Context, config Config) *S3 {
	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.Key, config.Secret, "")),
		awsconfig.WithBaseEndpoint(config.Endpoint),
		awsconfig.WithRegion(config.Region),
	)
	if err != nil {
		slox.Error(ctx, "failed to create S3 session", slog.Any("err", err))
		return nil
	}

	s3Client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	return &S3{
		BaseURL:  config.BaseURL,
		S3Client: s3Client,
		Config:   config,
	}
}

func (s3o *S3) Get(key string) (io.ReadCloser, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	object, err := s3o.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object.Body, nil
}

func (s3o *S3) Put(ctx context.Context, key string, body io.ReadSeeker) (string, error) {
	ctx, span := otel.Tracer("ficsit-app").Start(ctx, "Put")
	defer span.End()

	span.SetAttributes(attribute.String("key", key))

	cleanedKey := strings.TrimPrefix(key, "/")

	uploader := manager.NewUploader(s3o.S3Client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Body:   body,
		Bucket: aws.String(s3o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		span.RecordError(err)
		return cleanedKey, fmt.Errorf("failed to upload file: %w", err)
	}

	return key, nil
}

func (s3o *S3) SignGet(key string) (string, error) {
	// Public Bucket
	cleanedKey := strings.TrimPrefix(key, "/")
	return fmt.Sprintf(s3o.Config.Keypath, s3o.BaseURL, s3o.Config.Bucket, cleanedKey), nil
}

func (s3o *S3) SignPut(_ string) (string, error) {
	// Unsupported at the moment
	return "", errors.New("Unsupported")
}

func (s3o *S3) StartMultipartUpload(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	upload, err := s3o.S3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s3o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})
	if err != nil {
		return fmt.Errorf("failed to create multipart upload: %w", err)
	}

	redis.StoreMultipartUploadID(cleanedKey, *upload.UploadId)

	return nil
}

func (s3o *S3) UploadPart(key string, part int32, data io.ReadSeeker) error {
	cleanedKey := strings.TrimPrefix(key, "/")
	id := redis.GetMultipartUploadID(cleanedKey)

	response, err := s3o.S3Client.UploadPart(context.TODO(), &s3.UploadPartInput{
		Body:       data,
		Bucket:     aws.String(s3o.Config.Bucket),
		Key:        aws.String(cleanedKey),
		PartNumber: aws.Int32(part),
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
	parts := redis.GetMultipartCompletedParts(cleanedKey)
	completedParts := make([]types.CompletedPart, len(parts))

	for part, etag := range parts {
		partInt, _ := strconv.ParseInt(part, 10, 32)
		completedParts[partInt-1] = types.CompletedPart{ETag: aws.String(etag), PartNumber: aws.Int32(int32(partInt))}
	}

	_, err := s3o.S3Client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(s3o.Config.Bucket),
		Key:             aws.String(cleanedKey),
		MultipartUpload: &types.CompletedMultipartUpload{Parts: completedParts},
		UploadId:        aws.String(id),
	})
	if err != nil {
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}

	redis.ClearMultipartCompletedParts(cleanedKey)

	return nil
}

func (s3o *S3) Rename(from string, to string) error {
	cleanedKey := strings.TrimPrefix(to, "/")

	_, err := s3o.S3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(s3o.Config.Bucket),
		CopySource: aws.String(s3o.Config.Bucket + from),
		Key:        aws.String(cleanedKey),
	})
	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

func (s3o *S3) Delete(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")

	// Check up to 10 object pages
	for range 10 {
		versions, err := s3o.S3Client.ListObjectVersions(context.TODO(), &s3.ListObjectVersionsInput{
			Bucket:    aws.String(s3o.Config.Bucket),
			KeyMarker: aws.String(cleanedKey),
			Prefix:    aws.String(cleanedKey),
		})
		if err != nil {
			if strings.Contains(err.Error(), "NotImplemented") {
				_, err = s3o.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
					Bucket: aws.String(s3o.Config.Bucket),
					Key:    aws.String(cleanedKey),
				})
				if err != nil {
					return fmt.Errorf("failed to delete objects: %w", err)
				}

				return nil
			}

			return fmt.Errorf("failed to list object versions: %w", err)
		}

		objects := make([]types.ObjectIdentifier, len(versions.Versions)+len(versions.DeleteMarkers))

		for i, version := range versions.Versions {
			objects[i] = types.ObjectIdentifier{
				Key:       version.Key,
				VersionId: version.VersionId,
			}
		}

		for i, marker := range versions.DeleteMarkers {
			objects[i+len(versions.Versions)] = types.ObjectIdentifier{
				Key:       marker.Key,
				VersionId: marker.VersionId,
			}
		}

		if len(objects) == 0 {
			_, err = s3o.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: aws.String(s3o.Config.Bucket),
				Key:    aws.String(cleanedKey),
			})
			if err != nil {
				return fmt.Errorf("failed to delete objects: %w", err)
			}

			return nil
		}

		_, err = s3o.S3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(s3o.Config.Bucket),
			Delete: &types.Delete{
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

	data, err := s3o.S3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s3o.Config.Bucket),
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
	out := make([]Object, 0)

	var marker *string
	for {
		objects, err := s3o.S3Client.ListObjects(context.TODO(), &s3.ListObjectsInput{
			Bucket:  aws.String(s3o.Config.Bucket),
			Marker:  marker,
			MaxKeys: aws.Int32(math.MaxInt32),
			Prefix:  aws.String(prefix),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		marker = objects.NextMarker

		for _, obj := range objects.Contents {
			out = append(out, Object{
				Key:          obj.Key,
				LastModified: obj.LastModified,
			})
		}

		if objects.IsTruncated == nil || !*objects.IsTruncated {
			break
		}
	}

	return out, nil
}
