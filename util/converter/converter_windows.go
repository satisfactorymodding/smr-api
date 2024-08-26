package converter

import (
	"bytes"
	"context"
	"fmt"
	"image"
	// GIF Support
	_ "image/gif"
	// JPEG Support
	_ "image/jpeg"
	// PNG Support
	_ "image/png"
	"log/slog"

	"github.com/Vilsol/slox"
	"github.com/chai2010/webp"
	"github.com/galdor/go-thumbhash"
)

func ConvertAnyImageToWebp(ctx context.Context, imageAsBytes []byte) ([]byte, string, error) {
	imageData, imageType, err := image.Decode(bytes.NewReader(imageAsBytes))
	if err != nil {
		message := "error converting image to webp"
		slox.Error(ctx, message, slog.Any("err", err))
		return nil, "", fmt.Errorf("%s: %w", message, err)
	}

	result := bytes.NewBuffer(make([]byte, 0))

	if imageType == "gif" {
		message := "converting gif to webp not supported on windows"
		slox.Error(ctx, message, slog.Any("err", err))
		return nil, "", fmt.Errorf("%s: %w", message, err)
	}

	if err := webp.Encode(result, imageData, nil); err != nil {
		return nil, "", fmt.Errorf("error converting image to webp: %w", err)
	}

	hash := thumbhash.EncodeImage(imageData)
	thumbHash := base64.StdEncoding.EncodeToString(hash)

	return result.Bytes(), thumbHash, nil
}

func DecodeAny(data []byte) (image.Image, error) {
	imageData, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		decode, err := webp.Decode(bytes.NewReader(data))
		if err == nil {
			return decode, nil
		}
		return nil, fmt.Errorf("error decoding image webp: %w", err)
	}
	return imageData, nil
}
