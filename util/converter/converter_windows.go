package converter

import (
	"bytes"
	"github.com/chai2010/webp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func ConvertAnyImageToWebp(imageAsBytes []byte) ([]byte, error) {
	imageData, _, err := image.Decode(bytes.NewReader(imageAsBytes))

	if err != nil {
		err := errors.Wrap(err, "error converting image to webp")
		log.Ctx(ctx).Error(err)
		return nil, err
	}

	result := bytes.NewBuffer(make([]byte, 0))

	if err := webp.Encode(result, imageData, nil); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
