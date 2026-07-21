package export

// Image file extensions and normalized format names shared across the image
// readers and exporters in this package.
const (
	extPNG  = ".png"
	extJPEG = ".jpeg"
	extGIF  = ".gif"

	formatJPEG = "jpeg"
	formatEMF  = "emf"
)

// Preset geometry names produced by editorTypeToPreset.
const (
	presetRect       = "rect"
	presetRoundRect  = "roundRect"
	presetEllipse    = "ellipse"
	presetRightArrow = "rightArrow"
	presetLeftArrow  = "leftArrow"
)

// OOXML theme color slot names.
const (
	themeSlotDk1 = "dk1"
	themeSlotLt1 = "lt1"
	themeSlotDk2 = "dk2"
	themeSlotLt2 = "lt2"
)

// adjustValue2 is the second shape adjustment handle name.
const adjustValue2 = "adj2"

// anchorBottom is the bottom vertical anchor value.
const anchorBottom = "bottom"
