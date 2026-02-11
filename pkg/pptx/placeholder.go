package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/shapes"

type (
	// PlaceholderType defines the type of placeholder.
	PlaceholderType = shapes.PlaceholderType
	// Placeholder represents a placeholder shape on a slide layout or master.
	Placeholder = shapes.Placeholder
)

const (
	PlaceholderTypeTitle      = shapes.PlaceholderTypeTitle
	PlaceholderTypeBody       = shapes.PlaceholderTypeBody
	PlaceholderTypeCentrTitle = shapes.PlaceholderTypeCentrTitle
	PlaceholderTypeSubTitle   = shapes.PlaceholderTypeSubTitle
	PlaceholderTypeDt         = shapes.PlaceholderTypeDt
	PlaceholderTypeSldNum     = shapes.PlaceholderTypeSldNum
	PlaceholderTypeFtr        = shapes.PlaceholderTypeFtr
	PlaceholderTypeHdr        = shapes.PlaceholderTypeHdr
	PlaceholderTypeObj        = shapes.PlaceholderTypeObj
	PlaceholderTypeChart      = shapes.PlaceholderTypeChart
	PlaceholderTypeTbl        = shapes.PlaceholderTypeTbl
	PlaceholderTypeClipArt    = shapes.PlaceholderTypeClipArt
	PlaceholderTypeDgm        = shapes.PlaceholderTypeDgm
	PlaceholderTypeMedia      = shapes.PlaceholderTypeMedia
	PlaceholderTypeSldImg     = shapes.PlaceholderTypeSldImg
	PlaceholderTypePic        = shapes.PlaceholderTypePic
)
