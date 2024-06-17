package converter

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log/slog"

	"github.com/Vilsol/slox"
	"github.com/chai2010/webp"
	// GIF Support
	_ "image/gif"
	// JPEG Support
	_ "image/jpeg"
	// PNG Support
	_ "image/png"
)

func ConvertAnyImageToWebp(ctx context.Context, imageAsBytes []byte) ([]byte, error) {
	imageData, imageType, err := image.Decode(bytes.NewReader(imageAsBytes))
	if err != nil {
		message := "error converting image to webp"
		slox.Error(ctx, message, slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", message, err)
	}

	result := bytes.NewBuffer(make([]byte, 0))

	if imageType == "gif" {
		message := "converting gif to webp not supported on windows"
		slox.Error(ctx, message, slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", message, err)
	}

	if err := webp.Encode(result, imageData, nil); err != nil {
		return nil, fmt.Errorf("error converting image to webp: %w", err)
	}

	return result.Bytes(), nil
}
