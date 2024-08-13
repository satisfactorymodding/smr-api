package storage

import (
	"bytes"
	"context"
	"io"
	"log/slog"
)

type Wrapper struct {
	Reader Storage
	Writer Storage
}

func initializeWrapper(reader Storage, writer Storage) *Wrapper {
	return &Wrapper{
		Reader: reader,
		Writer: writer,
	}
}

func (w *Wrapper) Get(key string) (io.ReadCloser, error) {
	get, err := w.Writer.Get(key)

	if err != nil {
		return w.Reader.Get(key)
	}

	return get, nil
}

func (w *Wrapper) Put(ctx context.Context, key string, body io.ReadSeeker) (string, error) {
	return w.Writer.Put(ctx, key, body)
}

func (w *Wrapper) SignGet(key string) (string, error) {
	return w.Reader.SignGet(key)
}

func (w *Wrapper) SignPut(key string) (string, error) {
	return w.Writer.SignPut(key)
}

func (w *Wrapper) StartMultipartUpload(key string) error {
	return w.Writer.StartMultipartUpload(key)
}

func (w *Wrapper) UploadPart(key string, part int64, data io.ReadSeeker) error {
	return w.Writer.UploadPart(key, part, data)
}

func (w *Wrapper) CompleteMultipartUpload(key string) error {
	return w.Writer.CompleteMultipartUpload(key)
}

func (w *Wrapper) Rename(from string, to string) error {
	// check if file exists in Write
	// if does not exist, copy it from Read
	obj, err := w.Writer.Get(from)
	if obj != nil {
		defer obj.Close()
	}

	if err != nil || obj == nil {
		slog.Warn("file did not exist, copying", slog.String("from", from))
		get, err := w.Reader.Get(from)
		if err != nil {
			return err
		}

		defer get.Close()
		all, err := io.ReadAll(get)
		if err != nil {
			return err
		}

		if _, err := w.Writer.Put(context.Background(), from, bytes.NewReader(all)); err != nil {
			return err
		}
	}

	return w.Writer.Rename(from, to)
}

func (w *Wrapper) Delete(key string) error {
	err := w.Writer.Delete(key)
	if err != nil {
		slog.Error("file deletion failed", slog.String("key", key), slog.Any("err", err))
	}
	return err
}

func (w *Wrapper) Meta(key string) (*ObjectMeta, error) {
	get, err := w.Writer.Meta(key)

	if err != nil {
		return w.Reader.Meta(key)
	}

	return get, nil
}

func (w *Wrapper) List(prefix string) ([]Object, error) {
	get, err := w.Writer.List(prefix)

	if err != nil {
		return w.Reader.List(prefix)
	}

	return get, nil
}
