package storage

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
)

type Wasabi struct {
	S3Client *s3.S3
	Bucket   *string
}

func initializeWasabi(ctx context.Context, config Config) *Wasabi {
	bucket := aws.String(config.Bucket)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Key, config.Secret, ""),
		Endpoint:         aws.String(config.Endpoint),
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create session")
		return nil
	}

	s3Client := s3.New(newSession)

	return &Wasabi{
		Bucket:   bucket,
		S3Client: s3Client,
	}
}

func (wasabi *Wasabi) Get(key string) (io.ReadCloser, error) {
	obj, err := wasabi.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: wasabi.Bucket,
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get object")
	}

	return obj.Body, nil
}

func (wasabi *Wasabi) Put(ctx context.Context, key string, body io.ReadSeeker) (string, error) {
	_, err := wasabi.S3Client.PutObject(&s3.PutObjectInput{
		Body:   body,
		Bucket: wasabi.Bucket,
		Key:    aws.String(key),
	})

	if err != nil {
		return key, errors.Wrap(err, "failed to put object")
	}

	return key, nil
}

func (wasabi *Wasabi) SignGet(key string) (string, error) {
	req, _ := wasabi.S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: wasabi.Bucket,
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", errors.Wrap(err, "failed to sign url")
	}

	return urlStr, nil
}

func (wasabi *Wasabi) SignPut(key string) (string, error) {
	req, _ := wasabi.S3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: wasabi.Bucket,
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", errors.Wrap(err, "failed to sign url")
	}

	return urlStr, nil
}

func (wasabi *Wasabi) StartMultipartUpload(key string) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) UploadPart(key string, part int64, data io.ReadSeeker) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) CompleteMultipartUpload(key string) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) CopyObjectFromOldBucket(key string) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) CopyObjectToOldBucket(key string) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) ScheduleCopyAllObjectsFromOldBucket(scheduler func(string)) {
}

func (wasabi *Wasabi) Rename(from string, to string) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) Delete(key string) error {
	return errors.New("Unsupported")
}

func (wasabi *Wasabi) Meta(key string) (*ObjectMeta, error) {
	return nil, errors.New("Unsupported")
}
