package markdown

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

const (
	imageLayoutStartXInches = 0.8
	imageLayoutGapXInches   = 4.2
	imageLayoutGapYInches   = 2.15
	imageLayoutWidthInches  = 3.8
	imageLayoutHeightInches = 1.95
	imageLayoutMinTopInches = 2.2
	imageLayoutMaxTopInches = 4.6
	imageLayoutMaxTopY      = 7.0
)

func (p *markdownASTParser) nextEmbeddedImageFrame() (styling.Length, styling.Length, styling.Length, styling.Length) {
	baseY := p.computeEmbeddedImageBaseYInches()
	slot := p.imagePlacementCount
	col := slot % 2
	row := slot / 2

	x := imageLayoutStartXInches + float64(col)*imageLayoutGapXInches
	y := baseY + float64(row)*imageLayoutGapYInches
	maxY := imageLayoutMaxTopY - imageLayoutHeightInches
	if y > maxY {
		y = maxY
	}

	return styling.Inches(x),
		styling.Inches(y),
		styling.Inches(imageLayoutWidthInches),
		styling.Inches(imageLayoutHeightInches)
}

func (p *markdownASTParser) computeEmbeddedImageBaseYInches() float64 {
	if p.current == nil {
		return imageLayoutMinTopInches
	}
	bulletOffset := float64(len(p.current.Bullets)) * 0.23
	if bulletOffset > 1.6 {
		bulletOffset = 1.6
	}
	tableOffset := 0.0
	if p.current.Table != nil {
		tableOffset = 0.5
	}
	base := imageLayoutMinTopInches + bulletOffset + tableOffset
	if base > imageLayoutMaxTopInches {
		base = imageLayoutMaxTopInches
	}
	return base
}
