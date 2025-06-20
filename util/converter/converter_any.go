package converter

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"log/slog"

	webp "github.com/HugoSmits86/nativewebp"
	"github.com/Vilsol/slox"
	"github.com/galdor/go-thumbhash"

	// JPEG Support
	_ "image/jpeg"
	// PNG Support
	_ "image/png"
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
		allFrames, err := gif.DecodeAll(bytes.NewReader(imageAsBytes))
		if err != nil {
			message := "error converting image to webp"
			slox.Error(ctx, message, slog.Any("err", err))
			return nil, "", fmt.Errorf("%s: %w", message, err)
		}

		palettesToImages := make([]image.Image, len(allFrames.Image))
		durations := make([]uint, len(allFrames.Image))
		disposals := make([]uint, len(allFrames.Image))

		for i, paletted := range allFrames.Image {
			palettesToImages[i] = paletted
			durations[i] = max(10, uint(allFrames.Delay[i])/10)
			disposals[i] = 0
		}

		err = webp.EncodeAll(result, &webp.Animation{
			Images:          palettesToImages,
			Durations:       durations,
			Disposals:       disposals,
			LoopCount:       uint16(allFrames.LoopCount),
			BackgroundColor: 0xffffffff,
		}, nil)
		if err != nil {
			return nil, "", fmt.Errorf("error converting image to webp: %w", err)
		}

		hash := thumbhash.EncodeImage(imageData)
		thumbHash := base64.StdEncoding.EncodeToString(hash)

		return result.Bytes(), thumbHash, nil
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

	decode, err := gif.Decode(bytes.NewReader(data))
	if err == nil {
		return decode, nil
	}

	return nil, fmt.Errorf("error decoding image: %w", err)
}
