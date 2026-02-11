package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/elements"

type (
	// PlaceholderType defines the type of placeholder.
	PlaceholderType = elements.PlaceholderType
	// Placeholder represents a placeholder shape on a slide layout or master.
	Placeholder = elements.Placeholder
)

const (
	PlaceholderTypeTitle      = elements.PlaceholderTypeTitle
	PlaceholderTypeBody       = elements.PlaceholderTypeBody
	PlaceholderTypeCentrTitle = elements.PlaceholderTypeCentrTitle
	PlaceholderTypeSubTitle   = elements.PlaceholderTypeSubTitle
	PlaceholderTypeDt         = elements.PlaceholderTypeDt
	PlaceholderTypeSldNum     = elements.PlaceholderTypeSldNum
	PlaceholderTypeFtr        = elements.PlaceholderTypeFtr
	PlaceholderTypeHdr        = elements.PlaceholderTypeHdr
	PlaceholderTypeObj        = elements.PlaceholderTypeObj
	PlaceholderTypeChart      = elements.PlaceholderTypeChart
	PlaceholderTypeTbl        = elements.PlaceholderTypeTbl
	PlaceholderTypeClipArt    = elements.PlaceholderTypeClipArt
	PlaceholderTypeDgm        = elements.PlaceholderTypeDgm
	PlaceholderTypeMedia      = elements.PlaceholderTypeMedia
	PlaceholderTypeSldImg     = elements.PlaceholderTypeSldImg
	PlaceholderTypePic        = elements.PlaceholderTypePic
)
