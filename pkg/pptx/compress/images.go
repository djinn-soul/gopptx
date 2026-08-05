package compress

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"strings"
)

// optimizeImage re-encodes one media part. It returns the original bytes
// unchanged when the image cannot be decoded or when re-encoding would make
// the part larger, so compression never inflates a package.
func optimizeImage(data []byte, name string, level Level, quality int) ([]byte, bool, bool) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data, false, false
	}

	resized := false
	if maxDim := level.MaxImageDimension(); maxDim > 0 {
		if scaled, ok := scaleDown(src, maxDim); ok {
			src = scaled
			resized = true
		}
	}

	encoded, err := encodeImage(src, name, quality)
	if err != nil || len(encoded) == 0 {
		return data, false, false
	}
	if len(encoded) >= len(data) && !resized {
		return data, false, false
	}
	return encoded, true, resized
}

func encodeImage(img image.Image, name string, quality int) ([]byte, error) {
	var buf bytes.Buffer
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, ".png") {
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		if err := encoder.Encode(&buf, img); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// scaleDown shrinks img so its longest edge is at most maxDim. It reports
// false when the image already fits.
func scaleDown(img image.Image, maxDim int) (image.Image, bool) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return img, false
	}
	longest := max(w, h)
	if longest <= maxDim {
		return img, false
	}

	ratio := float64(maxDim) / float64(longest)
	newW := int(float64(w) * ratio)
	newH := int(float64(h) * ratio)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}
	return boxResize(img, newW, newH), true
}

// boxResize averages each source block into one destination pixel. It is
// slower than nearest-neighbour but avoids the aliasing that makes downscaled
// screenshots unreadable.
func boxResize(src image.Image, newW, newH int) image.Image {
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	for y := range newH {
		y0 := bounds.Min.Y + y*srcH/newH
		y1 := bounds.Min.Y + (y+1)*srcH/newH
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for x := range newW {
			x0 := bounds.Min.X + x*srcW/newW
			x1 := bounds.Min.X + (x+1)*srcW/newW
			if x1 <= x0 {
				x1 = x0 + 1
			}
			dst.Set(x, y, averageBlock(src, x0, y0, x1, y1))
		}
	}
	return dst
}

func averageBlock(src image.Image, x0, y0, x1, y1 int) color.RGBA {
	var sr, sg, sb, sa uint64
	count := uint64(0)
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			sr += uint64(r)
			sg += uint64(g)
			sb += uint64(b)
			sa += uint64(a)
			count++
		}
	}
	if count == 0 {
		return color.RGBA{}
	}
	return color.RGBA{
		R: to8BitChannel(sr, count),
		G: to8BitChannel(sg, count),
		B: to8BitChannel(sb, count),
		A: to8BitChannel(sa, count),
	}
}

// to8BitChannel averages a run of 16-bit colour samples down to one 8-bit
// channel. RGBA() promises samples of at most 0xffff, so the average can never
// exceed 0xff once scaled; the clamp keeps that true for any caller that hands
// back an out-of-range image.
func to8BitChannel(sum, count uint64) uint8 {
	const to8Bit = 257 // 16-bit channel -> 8-bit
	if count == 0 {
		return 0
	}
	scaled := sum / count / to8Bit
	if scaled > math.MaxUint8 {
		return math.MaxUint8
	}
	return uint8(scaled)
}
