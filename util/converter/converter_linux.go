package converter

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"

	"log/slog"

	"github.com/Vilsol/slox"
	"github.com/chai2010/webp"
	"github.com/galdor/go-thumbhash"
	giftowebp "github.com/sizeofint/gif-to-webp"
	webpd "github.com/tidbyt/go-libwebp/webp"

	// GIF Support
	_ "image/gif"
	// JPEG Support
	_ "image/jpeg"
	// PNG Support
	_ "image/png"
)

var converter = giftowebp.NewConverter()

func ConvertAnyImageToWebp(ctx context.Context, imageAsBytes []byte) ([]byte, string, error) {
	imageData, imageType, err := image.Decode(bytes.NewReader(imageAsBytes))
	if err != nil {
		message := "error converting image to webp"
		slox.Error(ctx, message, slog.Any("err", err))
		return nil, "", fmt.Errorf("%s: %w", message, err)
	}

	result := bytes.NewBuffer(make([]byte, 0))

	if imageType == "gif" {
		webpBin, err := converter.Convert(imageAsBytes)
		if err != nil {
			message := "error converting image to webp"
			slox.Error(ctx, message, slog.Any("err", err))
			return nil, "", fmt.Errorf("%s: %w", message, err)
		}

		return webpBin, "", nil
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
	if err == nil {
		return imageData, nil
	}

	imageData, err = webp.Decode(bytes.NewReader(data))
	if err == nil {
		return imageData, nil
	}

	decoder, err := webpd.NewAnimationDecoder(data)
	if err == nil {
		decode, err := decoder.Decode()
		if err == nil {
			return decode.Image[0], nil
		}
	}

	return nil, fmt.Errorf("error decoding image: %w", err)
}
