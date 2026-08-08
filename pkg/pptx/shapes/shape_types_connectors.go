package shapes

//go:generate go run ../../../cmd/gen_shape_presets . ../../../python/gopptx/shape_types.py

// Connector presets. These are ST_ShapeType members like any other preset, and
// the Python side has carried the numbered variants for a while; Go named only
// the three ConnectorType aliases, so generating the Python surface from Go
// would have dropped them.
const (
	// ShapeTypeStraightConnector1 is a straight line between two points.
	ShapeTypeStraightConnector1 = "straightConnector1"

	// ShapeTypeBentConnector2 starts the elbow connectors, which route in
	// right angles. The number is how many segments the route has.
	ShapeTypeBentConnector2 = "bentConnector2"
	ShapeTypeBentConnector3 = "bentConnector3"
	ShapeTypeBentConnector4 = "bentConnector4"
	ShapeTypeBentConnector5 = "bentConnector5"

	// ShapeTypeCurvedConnector2 starts the curved connectors, same numbering.
	ShapeTypeCurvedConnector2 = "curvedConnector2"
	ShapeTypeCurvedConnector3 = "curvedConnector3"
	ShapeTypeCurvedConnector4 = "curvedConnector4"
	ShapeTypeCurvedConnector5 = "curvedConnector5"
)
