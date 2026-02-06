package pptx

import "fmt"

// Image describes one image placement on a slide.
type Image struct {
	Path string
	X    int64
	Y    int64
	CX   int64
	CY   int64
}

// NewImage creates an image placement with explicit EMU coordinates and size.
func NewImage(path string, x int64, y int64, cx int64, cy int64) Image {
	return Image{
		Path: path,
		X:    x,
		Y:    y,
		CX:   cx,
		CY:   cy,
	}
}

func validateImage(image Image, slideIndex int, imageIndex int) error {
	if image.Path == "" {
		return fmt.Errorf("slide %d image %d path cannot be empty", slideIndex, imageIndex)
	}
	if image.X < 0 || image.Y < 0 {
		return fmt.Errorf("slide %d image %d position cannot be negative", slideIndex, imageIndex)
	}
	if image.CX <= 0 || image.CY <= 0 {
		return fmt.Errorf("slide %d image %d size must be > 0", slideIndex, imageIndex)
	}
	return nil
}
