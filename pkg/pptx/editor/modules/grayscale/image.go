package grayscale

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
)

const jpegQuality = 100

// ImageBytes converts image bytes to grayscale and returns the encoded output.
func ImageBytes(data []byte, format string) ([]byte, string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}

	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			gray, ok := color.GrayModel.Convert(pixel).(color.Gray)
			if !ok {
				continue
			}
			alpha, ok := color.AlphaModel.Convert(pixel).(color.Alpha)
			if !ok {
				continue
			}
			dst.Set(x, y, color.NRGBA{R: gray.Y, G: gray.Y, B: gray.Y, A: alpha.A})
		}
	}

	normalized := strings.ToLower(strings.TrimSpace(format))
	var out bytes.Buffer
	switch normalized {
	case "jpg", "jpeg":
		if err := jpeg.Encode(&out, dst, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return nil, "", fmt.Errorf("encode jpeg: %w", err)
		}
		return out.Bytes(), "jpeg", nil
	default:
		if err := png.Encode(&out, dst); err != nil {
			return nil, "", fmt.Errorf("encode png: %w", err)
		}
		return out.Bytes(), "png", nil
	}
}
