package converter

import (
	"bytes"
	"context"
	"image"

	"github.com/chai2010/webp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

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
		log.Err(err).Msg(message)
		return nil, errors.Wrap(err, message)
	}

	result := bytes.NewBuffer(make([]byte, 0))

	if imageType == "gif" {
		message := "converting gif to webp not supported on windows"
		log.Err(err).Msg(message)
		return nil, errors.Wrap(err, message)
	}

	if err := webp.Encode(result, imageData, nil); err != nil {
		return nil, errors.Wrap(err, "error converting image to webp")
	}

	return result.Bytes(), nil
}
