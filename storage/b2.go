package storage

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/satisfactorymodding/smr-api/redis"
	"github.com/satisfactorymodding/smr-api/redis/jobs"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type B2 struct {
	//Client    *b2.Client
	//Bucket    *b2.Bucket
	BaseURL   string
	S3Client  *s3.S3
	S3Session *session.Session
	Config    Config
}

func initializeB2(ctx context.Context, config Config) *B2 {
	//ctx := context.Background()
	//client, err := b2.NewClient(ctx, config.OldKey, config.OldSecret)
	//
	//if err != nil {
	//	log.Ctx(ctx).Error().Msgf("Failed to create client %s", err.Error())
	//	return nil
	//}
	//
	//bucket, err := client.Bucket(ctx, config.OldBucket)
	//
	//if err != nil {
	//	log.Ctx(ctx).Error().Msgf("Failed to create session %s", err.Error())
	//	return nil
	//}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Key, config.Secret, ""),
		Endpoint:         aws.String("https://s3." + config.Region + ".backblazeb2.com"),
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create S3 session")
		return nil
	}

	s3Client := s3.New(newSession)

	return &B2{
		//Client:    client,
		//Bucket:    bucket,
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

	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(b2o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})

	if err != nil {
		return cleanedKey, errors.Wrap(err, "failed to upload file")
	}

	jobs.SubmitJobCopyObjectToOldBucketTask(ctx, key)

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

func (b2o *B2) CopyObjectFromOldBucket(key string) error {
	//cleanedKey := strings.TrimPrefix(key, "/")
	//obj := b2o.Bucket.Object(cleanedKey)
	//reader := obj.NewReader(context.Background())
	//
	//log.Ctx(ctx).Info("Copying file from old bucket: ", key)
	//
	//uploader := s3manager.NewUploader(b2o.S3Session)
	//
	//_, err := uploader.Upload(&s3manager.UploadInput{
	//	Body:   reader,
	//	Bucket: aws.String(b2o.Config.Bucket),
	//	Key:    aws.String(cleanedKey),
	//})
	//
	//return err
	return errors.New("oh no")
}

func (b2o *B2) CopyObjectToOldBucket(key string) error {
	//cleanedKey := strings.TrimPrefix(key, "/")
	//
	//log.Ctx(ctx).Info("Copying file to old bucket: ", key)
	//
	//object, err := b2o.S3Client.GetObject(&s3.GetObjectInput{
	//	Bucket: aws.String(b2o.Config.Bucket),
	//	Key:    aws.String(cleanedKey),
	//})
	//
	//if err != nil {
	//	return err
	//}
	//
	//obj := b2o.Bucket.Object(cleanedKey)
	//
	//writer := obj.NewWriter(context.Background())
	//_, err = io.Copy(writer, object.Body)
	//
	//if err != nil {
	//	_ = writer.Close()
	//	return err
	//}
	//
	//err = writer.Close()
	//
	//if err != nil {
	//	_ = writer.Close()
	//	return err
	//}
	//
	//return err
	return nil
}

func (b2o *B2) ScheduleCopyAllObjectsFromOldBucket(scheduler func(string)) {
	//iterator := b2o.Bucket.List(context.Background())
	//for iterator.Next() {
	//	scheduler(iterator.Object().Name())
	//}
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
