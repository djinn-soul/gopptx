package shapes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

// PlaceholderType defines the type of placeholder.
type PlaceholderType string

const (
	PlaceholderTypeTitle      PlaceholderType = "title"
	PlaceholderTypeBody       PlaceholderType = "body"
	PlaceholderTypeCentrTitle PlaceholderType = "ctrTitle"
	PlaceholderTypeSubTitle   PlaceholderType = "subTitle"
	PlaceholderTypeDt         PlaceholderType = "dt"
	PlaceholderTypeSldNum     PlaceholderType = "sldNum"
	PlaceholderTypeFtr        PlaceholderType = "ftr"
	PlaceholderTypeHdr        PlaceholderType = "hdr"
	PlaceholderTypeObj        PlaceholderType = "obj"
	PlaceholderTypeChart      PlaceholderType = "chart"
	PlaceholderTypeTbl        PlaceholderType = "tbl"
	PlaceholderTypeClipArt    PlaceholderType = "clipArt"
	PlaceholderTypeDgm        PlaceholderType = "dgm"
	PlaceholderTypeMedia      PlaceholderType = "media"
	PlaceholderTypeSldImg     PlaceholderType = "sldImg"
	PlaceholderTypePic        PlaceholderType = "pic"
)

// Placeholder represents a placeholder shape on a slide layout or master.
type Placeholder struct {
	Type  PlaceholderType
	Index int
	Name  string

	// underlying shape reference or geometry for layout inheritance
	X, Y, CX, CY styling.Length
}

// InsertPicture inserts an image into the placeholder.
func (p *Placeholder) InsertPicture(imagePath string) Image {
	return Image{
		Path:        imagePath,
		X:           p.X,
		Y:           p.Y,
		CX:          p.CX,
		CY:          p.CY,
		Placeholder: p,
	}
}

// InsertPictureFromBytes inserts an image from bytes into the placeholder.
func (p *Placeholder) InsertPictureFromBytes(data []byte, format string) Image {
	return Image{
		Data:        data,
		Format:      format,
		X:           p.X,
		Y:           p.Y,
		CX:          p.CX,
		CY:          p.CY,
		Placeholder: p,
	}
}

// PlaceholderTarget identifies a placeholder to override.
type PlaceholderTarget struct {
	Type  string
	Index int
	Name  string
}

// PlaceholderTextStyle describes text formatting overrides for a placeholder.
type PlaceholderTextStyle struct {
	SizePt    *int
	Color     *string
	Bold      *bool
	Italic    *bool
	Underline *string
	Align     *string
	Font      *string
}

// PlaceholderOverrideOptions defines geometry and style overrides for a placeholder.
type PlaceholderOverrideOptions struct {
	X, Y, CX, CY *styling.Length
	TextStyle    *PlaceholderTextStyle
	ForceRect    *bool
}

// Validate ensures override options are coherent and safe to render.
func (o *PlaceholderOverrideOptions) Validate() error {
	if o == nil {
		return nil
	}
	coords := []*styling.Length{o.X, o.Y, o.CX, o.CY}
	setCount := 0
	for _, c := range coords {
		if c != nil {
			setCount++
		}
	}
	if setCount != 0 && setCount != len(coords) {
		return errors.New("placeholder override geometry requires X, Y, CX, and CY together")
	}
	if setCount == len(coords) {
		if o.CX.Emu() <= 0 || o.CY.Emu() <= 0 {
			return errors.New("placeholder override geometry requires positive CX and CY")
		}
	}
	if o.TextStyle != nil {
		return o.TextStyle.Validate()
	}
	return nil
}

// PlaceholderContent describes overridden content for a slide layout placeholder.
type PlaceholderContent struct {
	Index int
	Type  string
	Text  string
	Image *Image
	Table *tables.Table
	Chart charts.ChartDefinition

	// Extension: Layout/Style Overrides
	Target   *PlaceholderTarget
	Override *PlaceholderOverrideOptions
}

// ValidateOverride validates style/geometry override payload only.
func (p PlaceholderContent) ValidateOverride() error {
	if p.Override == nil {
		return nil
	}
	if err := p.Override.Validate(); err != nil {
		return fmt.Errorf("placeholder override %d: %w", p.Index, err)
	}
	return nil
}

// Validate ensures placeholder text style values are valid.
func (s *PlaceholderTextStyle) Validate() error {
	if s == nil {
		return nil
	}
	if s.SizePt != nil && (*s.SizePt < 1 || *s.SizePt > 400) {
		return errors.New("placeholder text size must be between 1 and 400 pt")
	}
	if s.Color != nil && !common.IsHexColor(*s.Color) {
		return errors.New("placeholder text color must be 6-digit RGB hex")
	}
	if s.Align != nil {
		align := strings.TrimSpace(*s.Align)
		if align != "l" && align != "ctr" && align != "r" && align != "just" {
			return fmt.Errorf("invalid placeholder text alignment %q", align)
		}
	}
	if s.Font != nil && strings.TrimSpace(*s.Font) == "" {
		return errors.New("placeholder text font cannot be empty")
	}
	return nil
}

// InsertText inserts text into the placeholder.
func (p *Placeholder) InsertText(text string) PlaceholderContent {
	return PlaceholderContent{
		Index: p.Index,
		Type:  string(p.Type),
		Text:  text,
	}
}

// InsertTable inserts a table into the placeholder.
func (p *Placeholder) InsertTable(table tables.Table) PlaceholderContent {
	return PlaceholderContent{
		Index: p.Index,
		Type:  string(p.Type),
		Table: &table,
	}
}

// InsertPictureToSlide returns a PlaceholderContent for the image.
func (p *Placeholder) InsertPictureToSlide(image Image) PlaceholderContent {
	return PlaceholderContent{
		Index: p.Index,
		Type:  string(p.Type),
		Image: &image,
	}
}
