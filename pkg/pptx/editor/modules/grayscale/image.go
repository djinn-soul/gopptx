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
			r, g, b, a := img.At(x, y).RGBA()
			luma := uint8((299*int(r>>8) + 587*int(g>>8) + 114*int(b>>8) + 500) / 1000)
			dst.Set(x, y, color.NRGBA{R: luma, G: luma, B: luma, A: uint8(a >> 8)})
		}
	}

	normalized := strings.ToLower(strings.TrimSpace(format))
	var out bytes.Buffer
	switch normalized {
	case "jpg", "jpeg":
		if err := jpeg.Encode(&out, dst, &jpeg.Options{Quality: 100}); err != nil {
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
