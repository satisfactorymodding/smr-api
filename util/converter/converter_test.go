package converter

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	webp "github.com/HugoSmits86/nativewebp"
	"github.com/MarvinJWendt/testza"
)

func TestConvertAnyImageToWebp_PNG(t *testing.T) {
	// Create a simple PNG image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			img.Set(x, y, color.RGBA{uint8(x * 255 / 100), uint8(y * 255 / 100), 128, 255})
		}
	}

	// Convert to PNG bytes
	var pngData []byte
	{
		tmpFile, err := os.CreateTemp("", "test*.png")
		testza.AssertNoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = png.Encode(tmpFile, img)
		testza.AssertNoError(t, err)

		pngData, err = os.ReadFile(tmpFile.Name())
		testza.AssertNoError(t, err)
	}

	ctx := context.Background()
	webpData, thumbHash, err := ConvertAnyImageToWebp(ctx, pngData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, webpData)
	testza.AssertNotEqual(t, "", thumbHash)
	testza.AssertTrue(t, len(webpData) > 0)

	// Verify the thumbhash is valid base64
	_, err = base64.StdEncoding.DecodeString(thumbHash)
	testza.AssertNoError(t, err)

	// Verify the result is valid WebP
	_, err = webp.Decode(bytes.NewReader(webpData))
	testza.AssertNoError(t, err)
}

func TestConvertAnyImageToWebp_JPEG(t *testing.T) {
	// Create a simple JPEG image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			img.Set(x, y, color.RGBA{uint8(x * 255 / 100), uint8(y * 255 / 100), 128, 255})
		}
	}

	// Convert to JPEG bytes
	var jpegData []byte
	{
		tmpFile, err := os.CreateTemp("", "test*.jpg")
		testza.AssertNoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 80})
		testza.AssertNoError(t, err)

		jpegData, err = os.ReadFile(tmpFile.Name())
		testza.AssertNoError(t, err)
	}

	ctx := context.Background()
	webpData, thumbHash, err := ConvertAnyImageToWebp(ctx, jpegData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, webpData)
	testza.AssertNotEqual(t, "", thumbHash)
	testza.AssertTrue(t, len(webpData) > 0)

	// Verify the thumbhash is valid base64
	_, err = base64.StdEncoding.DecodeString(thumbHash)
	testza.AssertNoError(t, err)

	// Verify the result is valid WebP
	_, err = webp.Decode(bytes.NewReader(webpData))
	testza.AssertNoError(t, err)
}

func TestConvertAnyImageToWebp_GIF(t *testing.T) {
	gifPath := filepath.Join("..", "..", "tests", "testdata", "img", "eglite.gif")
	gifData, err := os.ReadFile(gifPath)
	testza.AssertNoError(t, err)
	testza.AssertTrue(t, len(gifData) > 0)

	ctx := context.Background()
	webpData, thumbHash, err := ConvertAnyImageToWebp(ctx, gifData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, webpData)
	testza.AssertNotEqual(t, "", thumbHash)
	testza.AssertTrue(t, len(webpData) > 0)

	_, err = base64.StdEncoding.DecodeString(thumbHash)
	testza.AssertNoError(t, err)

	// Use snapshot testing for the WebP output to detect any changes in conversion
	err = testza.SnapshotCreateOrValidate(t, "eglite_gif_to_webp", webpData)
	testza.AssertNoError(t, err)
	err = testza.SnapshotCreateOrValidate(t, "eglite_gif_thumbhash", thumbHash)
	testza.AssertNoError(t, err)

	decodedImg, err := DecodeAny(webpData)
	if err != nil {
		testza.AssertTrue(t, len(webpData) >= 12)
		testza.AssertEqual(t, "RIFF", string(webpData[0:4]))
		testza.AssertEqual(t, "WEBP", string(webpData[8:12]))
	} else {
		testza.AssertNotNil(t, decodedImg)
	}
}

func TestConvertAnyImageToWebp_AnimatedGIF(t *testing.T) {
	// Create a simple animated GIF
	var gifData []byte
	{
		tmpFile, err := os.CreateTemp("", "test*.gif")
		testza.AssertNoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Create a simple 2-frame animated GIF
		palette := color.Palette{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, // Black
			color.RGBA{0xff, 0x00, 0x00, 0xff}, // Red
			color.RGBA{0x00, 0xff, 0x00, 0xff}, // Green
		}

		frame1 := image.NewPaletted(image.Rect(0, 0, 10, 10), palette)
		frame2 := image.NewPaletted(image.Rect(0, 0, 10, 10), palette)

		// Fill frames with different colors
		for y := range 10 {
			for x := range 10 {
				frame1.SetColorIndex(x, y, 1) // Red
				frame2.SetColorIndex(x, y, 2) // Green
			}
		}

		anim := &gif.GIF{
			Image: []*image.Paletted{frame1, frame2},
			Delay: []int{50, 50}, // 50 centiseconds each
		}

		err = gif.EncodeAll(tmpFile, anim)
		testza.AssertNoError(t, err)

		gifData, err = os.ReadFile(tmpFile.Name())
		testza.AssertNoError(t, err)
	}

	ctx := context.Background()
	webpData, thumbHash, err := ConvertAnyImageToWebp(ctx, gifData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, webpData)
	testza.AssertNotEqual(t, "", thumbHash)
	testza.AssertTrue(t, len(webpData) > 0)

	_, err = base64.StdEncoding.DecodeString(thumbHash)
	testza.AssertNoError(t, err)

	// Use snapshot testing for the animated WebP output
	err = testza.SnapshotCreateOrValidate(t, "animated_gif_to_webp", webpData)
	testza.AssertNoError(t, err)
	err = testza.SnapshotCreateOrValidate(t, "animated_gif_thumbhash", thumbHash)
	testza.AssertNoError(t, err)

	decodedImg, err := DecodeAny(webpData)
	if err != nil {
		testza.AssertTrue(t, len(webpData) >= 12)
		testza.AssertEqual(t, "RIFF", string(webpData[0:4]))
		testza.AssertEqual(t, "WEBP", string(webpData[8:12]))
	} else {
		testza.AssertNotNil(t, decodedImg)
	}
}

func TestConvertAnyImageToWebp_InvalidData(t *testing.T) {
	ctx := context.Background()
	invalidData := []byte("not an image")

	webpData, thumbHash, err := ConvertAnyImageToWebp(ctx, invalidData)

	testza.AssertTrue(t, err != nil)
	testza.AssertNil(t, webpData)
	testza.AssertEqual(t, "", thumbHash)
}

func TestDecodeAny_PNG(t *testing.T) {
	// Create a simple PNG image
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := range 50 {
		for x := range 50 {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
		}
	}

	// Convert to PNG bytes
	var pngData []byte
	{
		tmpFile, err := os.CreateTemp("", "test*.png")
		testza.AssertNoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = png.Encode(tmpFile, img)
		testza.AssertNoError(t, err)

		pngData, err = os.ReadFile(tmpFile.Name())
		testza.AssertNoError(t, err)
	}

	decodedImg, err := DecodeAny(pngData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, decodedImg)
	testza.AssertEqual(t, image.Rect(0, 0, 50, 50), decodedImg.Bounds())
}

func TestDecodeAny_JPEG(t *testing.T) {
	// Create a simple JPEG image
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := range 50 {
		for x := range 50 {
			img.Set(x, y, color.RGBA{0, 255, 0, 255}) // Green
		}
	}

	// Convert to JPEG bytes
	var jpegData []byte
	{
		tmpFile, err := os.CreateTemp("", "test*.jpg")
		testza.AssertNoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 80})
		testza.AssertNoError(t, err)

		jpegData, err = os.ReadFile(tmpFile.Name())
		testza.AssertNoError(t, err)
	}

	decodedImg, err := DecodeAny(jpegData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, decodedImg)
	testza.AssertEqual(t, image.Rect(0, 0, 50, 50), decodedImg.Bounds())
}

func TestDecodeAny_GIF(t *testing.T) {
	// Read the test GIF file
	gifPath := filepath.Join("..", "..", "tests", "testdata", "img", "eglite.gif")
	gifData, err := os.ReadFile(gifPath)
	testza.AssertNoError(t, err)

	decodedImg, err := DecodeAny(gifData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, decodedImg)
	testza.AssertTrue(t, decodedImg.Bounds().Dx() > 0)
	testza.AssertTrue(t, decodedImg.Bounds().Dy() > 0)
}

func TestDecodeAny_WebP(t *testing.T) {
	// Create a simple image and convert to WebP
	img := image.NewRGBA(image.Rect(0, 0, 30, 30))
	for y := range 30 {
		for x := range 30 {
			img.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue
		}
	}

	// Convert to WebP bytes
	var webpData []byte
	{
		tmpFile, err := os.CreateTemp("", "test*.webp")
		testza.AssertNoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = webp.Encode(tmpFile, img, nil)
		testza.AssertNoError(t, err)

		webpData, err = os.ReadFile(tmpFile.Name())
		testza.AssertNoError(t, err)
	}

	decodedImg, err := DecodeAny(webpData)

	testza.AssertNoError(t, err)
	testza.AssertNotNil(t, decodedImg)
	testza.AssertEqual(t, image.Rect(0, 0, 30, 30), decodedImg.Bounds())
}

func TestDecodeAny_InvalidData(t *testing.T) {
	invalidData := []byte("not an image")

	decodedImg, err := DecodeAny(invalidData)

	testza.AssertTrue(t, err != nil)
	testza.AssertNil(t, decodedImg)
}

func TestDecodeAny_EmptyData(t *testing.T) {
	emptyData := []byte{}

	decodedImg, err := DecodeAny(emptyData)

	testza.AssertTrue(t, err != nil)
	testza.AssertNil(t, decodedImg)
}

func BenchmarkConvertAnyImageToWebp_PNG(b *testing.B) {
	// Create a test PNG image
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := range 200 {
		for x := range 200 {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	var pngData []byte
	{
		tmpFile, err := os.CreateTemp("", "bench*.png")
		if err != nil {
			b.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = png.Encode(tmpFile, img)
		if err != nil {
			b.Fatal(err)
		}

		pngData, err = os.ReadFile(tmpFile.Name())
		if err != nil {
			b.Fatal(err)
		}
	}

	ctx := context.Background()
	b.ResetTimer()

	for range b.N {
		_, _, err := ConvertAnyImageToWebp(ctx, pngData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeAny_PNG(b *testing.B) {
	// Create a test PNG image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	var pngData []byte
	{
		tmpFile, err := os.CreateTemp("", "bench*.png")
		if err != nil {
			b.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err = png.Encode(tmpFile, img)
		if err != nil {
			b.Fatal(err)
		}

		pngData, err = os.ReadFile(tmpFile.Name())
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for range b.N {
		_, err := DecodeAny(pngData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
