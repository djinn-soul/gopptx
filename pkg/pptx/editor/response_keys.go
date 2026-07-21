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
)

// shapeTypeRect is the default preset geometry reported for shapes with no
// explicit geometry.
const shapeTypeRect = "rect"

// contentTypesPath is the OPC content-types part.
const contentTypesPath = "[Content_Types].xml"
