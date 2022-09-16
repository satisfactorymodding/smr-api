package util

import (
	"context"
	// GIF Support
	_ "image/gif"
	// JPEG Support
	_ "image/jpeg"
	// PNG Support
	_ "image/png"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/smr-api/util/converter"
)

func LinkToWebp(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid url")
	}

	imageAsBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("invalid url")
	}

	return converter.ConvertAnyImageToWebp(ctx, imageAsBytes)
}
