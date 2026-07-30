package media

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"sync"
)

const (
	posterWidth  = 640
	posterHeight = 360
	// Play triangle bounds in pixels, centred on the poster.
	posterTriangleLeft   = 269
	posterTriangleRight  = 371
	posterTriangleTop    = 119
	posterTriangleBottom = 241
)

//nolint:gochecknoglobals // the poster is deterministic; encode it once
var (
	defaultPosterOnce  sync.Once
	defaultPosterBytes []byte

	// posterPanelColor is a dark slate backdrop; posterGlyphColor is the play
	// triangle drawn on it.
	//nolint:gochecknoglobals,mnd // color channel values, immutable
	posterPanelColor = color.RGBA{R: 0x1F, G: 0x25, B: 0x33, A: 0xFF}
	//nolint:gochecknoglobals,mnd // color channel values, immutable
	posterGlyphColor = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

// DefaultVideoPosterPNG returns the poster frame used when a video insert call
// supplies none.
//
// It is a visible placeholder — a dark panel with a play triangle — because a
// video shape renders as its poster until played. The previous 1x1 transparent
// pixel made every video added without a poster invisible on the slide.
func DefaultVideoPosterPNG() []byte {
	defaultPosterOnce.Do(func() {
		defaultPosterBytes = renderDefaultPoster()
	})
	out := make([]byte, len(defaultPosterBytes))
	copy(out, defaultPosterBytes)
	return out
}

func renderDefaultPoster() []byte {
	panel := posterPanelColor
	glyph := posterGlyphColor

	img := image.NewRGBA(image.Rect(0, 0, posterWidth, posterHeight))
	for y := range posterHeight {
		for x := range posterWidth {
			img.Set(x, y, panel)
		}
	}
	drawPlayTriangle(img, glyph)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		// Encoding an in-memory RGBA image cannot fail; fall back to an empty
		// poster rather than panicking in a library call.
		return nil
	}
	return buf.Bytes()
}

// drawPlayTriangle fills a right-pointing triangle centred on the poster.
func drawPlayTriangle(img *image.RGBA, fill color.Color) {
	left, right := posterTriangleLeft, posterTriangleRight
	top, bottom := posterTriangleTop, posterTriangleBottom
	midY := (top + bottom) / 2

	for x := left; x <= right; x++ {
		// The triangle narrows linearly from the left edge to the tip.
		progress := float64(x-left) / float64(right-left)
		halfHeight := int(float64(midY-top) * (1 - progress))
		for y := midY - halfHeight; y <= midY+halfHeight; y++ {
			img.Set(x, y, fill)
		}
	}
}
