package editor

// Response and payload map keys shared across the command handlers. These are
// part of the JSON command protocol, so they are defined once here rather than
// repeated as literals at each construction site.
const (
	keySlides      = "slides"
	keyShapes      = "shapes"
	keyPlaceholder = "placeholders"
	keyIndex       = "index"
	keyType        = "type"
	keyName        = "name"
	keySlideIndex  = "slide_index"
	keyShapeIDs    = "shape_ids"
	keyShapeID     = "shape_id"
	keyLeft        = "left"
	keyTop         = "top"
	keyWidth       = "width"
	keyHeight      = "height"
)

// shapeTypeRect is the default preset geometry reported for shapes with no
// explicit geometry.
const shapeTypeRect = "rect"

// targetModeExternal marks a relationship whose target is a URL rather than a
// part inside the package.
const targetModeExternal = "External"

// contentTypesPath is the OPC content-types part.
const contentTypesPath = "[Content_Types].xml"
