package shapes

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"math"
	"os"

	// Registered so DecodeConfig recognises the formats a deck can embed.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Every image constructor demanded a width and a height, so placing a picture
// at its natural size, or at a given width without distorting it, meant
// decoding the file yourself first. The decoder was already a dependency —
// the editor and the URL fetcher both use it — just not from here.

// pixelsPerInch is the resolution OOXML assumes for an image with no stated
// density: 96 dpi, i.e. 9525 EMU per pixel.
const pixelsPerInch = 96

// ErrUnknownImageSize is returned when the intrinsic size cannot be read,
// either because the bytes are not an image or because the placement carries a
// URL that has not been fetched yet.
var ErrUnknownImageSize = errors.New("image size unavailable")

// PixelsToLength converts a pixel count at 96 dpi into EMU.
func PixelsToLength(pixels int) styling.Length {
	return styling.Inches(float64(pixels) / pixelsPerInch)
}

// ImageSizePixels reads the intrinsic pixel size of encoded image bytes,
// returning width then height.
func ImageSizePixels(data []byte) (int, int, error) {
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrUnknownImageSize, err)
	}
	return config.Width, config.Height, nil
}

// NaturalSize returns the placement's intrinsic size in EMU, reading the file
// from Path when the bytes are not already loaded. A URL-sourced image has no
// size until it is fetched, and reports ErrUnknownImageSize.
func (img Image) NaturalSize() (styling.Length, styling.Length, error) {
	data := img.Data
	if len(data) == 0 && img.Path != "" {
		fileData, err := os.ReadFile(img.Path)
		if err != nil {
			return 0, 0, fmt.Errorf("read image %s: %w", img.Path, err)
		}
		data = fileData
	}
	if len(data) == 0 {
		return 0, 0, ErrUnknownImageSize
	}

	width, height, err := ImageSizePixels(data)
	if err != nil {
		return 0, 0, err
	}
	return PixelsToLength(width), PixelsToLength(height), nil
}

// AspectRatio is width divided by height, from the placement's own extent when
// it has one and from the image itself otherwise.
func (img Image) AspectRatio() (float64, error) {
	if img.CX > 0 && img.CY > 0 {
		return float64(img.CX) / float64(img.CY), nil
	}
	cx, cy, err := img.NaturalSize()
	if err != nil {
		return 0, err
	}
	if cy == 0 {
		return 0, ErrUnknownImageSize
	}
	return float64(cx) / float64(cy), nil
}

// AtNaturalSize sets the extent to the image's intrinsic size.
func (img Image) AtNaturalSize() (Image, error) {
	cx, cy, err := img.NaturalSize()
	if err != nil {
		return img, err
	}
	img.CX, img.CY = cx, cy
	return img, nil
}

// ScaleToWidth sets the width and derives the height from the aspect ratio.
func (img Image) ScaleToWidth(width styling.Length) (Image, error) {
	ratio, err := img.AspectRatio()
	if err != nil {
		return img, err
	}
	img.CX = width
	img.CY = styling.Length(math.Round(float64(width) / ratio))
	return img, nil
}

// ScaleToHeight sets the height and derives the width from the aspect ratio.
func (img Image) ScaleToHeight(height styling.Length) (Image, error) {
	ratio, err := img.AspectRatio()
	if err != nil {
		return img, err
	}
	img.CY = height
	img.CX = styling.Length(math.Round(float64(height) * ratio))
	return img, nil
}

// FitWithin scales the image to the largest size that fits the box without
// distorting it — the "scale to fit" a caller usually wants for a placeholder.
func (img Image) FitWithin(maxWidth, maxHeight styling.Length) (Image, error) {
	ratio, err := img.AspectRatio()
	if err != nil {
		return img, err
	}
	width := maxWidth
	height := styling.Length(math.Round(float64(width) / ratio))
	if height > maxHeight {
		height = maxHeight
		width = styling.Length(math.Round(float64(height) * ratio))
	}
	img.CX, img.CY = width, height
	return img, nil
}

// NewImageAtNaturalSize places a file at (x, y) in its own size.
func NewImageAtNaturalSize(path string, x, y styling.Length) (Image, error) {
	return Image{Path: path, X: x, Y: y}.AtNaturalSize()
}

// NewImageFromBytesAtNaturalSize places raw bytes at (x, y) in their own size.
func NewImageFromBytesAtNaturalSize(data []byte, format string, x, y styling.Length) (Image, error) {
	return Image{Data: data, Format: format, X: x, Y: y}.AtNaturalSize()
}

// NewImageScaledToWidth places a file at (x, y) at the given width, keeping its
// aspect ratio.
func NewImageScaledToWidth(path string, x, y, width styling.Length) (Image, error) {
	return Image{Path: path, X: x, Y: y}.ScaleToWidth(width)
}
