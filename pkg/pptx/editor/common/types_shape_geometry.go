package editorcommon

// FreeformGeometry is the read model for a shape's <a:custGeom>. Custom
// geometry used to be write-only: AddFreeformShape emitted <a:custGeom> and
// nothing parsed it back, so a freeform authored in PowerPoint (or by gopptx)
// could not be inspected (upstream #1020).
type FreeformGeometry struct {
	Paths []FreeformPath `json:"paths,omitempty"`
}

// FreeformPath is one <a:path> inside the geometry's <a:pathLst>.
type FreeformPath struct {
	// W and H are the path's own coordinate space, which is independent of the
	// shape extent the path is scaled into.
	W int `json:"w,omitempty"`
	H int `json:"h,omitempty"`
	// Fill and Stroke mirror the optional <a:path> attributes; empty means the
	// OOXML default ("norm" and stroked).
	Fill     string            `json:"fill,omitempty"`
	Stroke   *bool             `json:"stroke,omitempty"`
	Segments []FreeformSegment `json:"segments,omitempty"`
}

// Segment kinds emitted in FreeformSegment.Type, matching the DrawingML tag
// names so a caller can map back to the source XML without a lookup table.
const (
	FreeformSegmentMoveTo     = "moveTo"
	FreeformSegmentLineTo     = "lnTo"
	FreeformSegmentCubicBezTo = "cubicBezTo"
	FreeformSegmentQuadBezTo  = "quadBezTo"
	FreeformSegmentArcTo      = "arcTo"
	FreeformSegmentClose      = "close"
)

// FreeformSegment is one drawing command in a path. Points carries the command's
// <a:pt> children in document order: one for moveTo and lnTo, two for quadBezTo,
// three for cubicBezTo, none for arcTo and close.
type FreeformSegment struct {
	Type   string          `json:"type"`
	Points []FreeformPoint `json:"points,omitempty"`
	// Arc parameters, set only when Type is arcTo.
	WidthRadius  *int     `json:"width_radius,omitempty"`
	HeightRadius *int     `json:"height_radius,omitempty"`
	StartAngle   *float64 `json:"start_angle_deg,omitempty"`
	SwingAngle   *float64 `json:"swing_angle_deg,omitempty"`
}

// FreeformPoint is one <a:pt>. A coordinate may be a guide formula rather than a
// number (for example x="wd2"); in that case X or Y stays 0 and the raw
// attribute is preserved in XExpr or YExpr so nothing is silently invented.
type FreeformPoint struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	XExpr string `json:"x_expr,omitempty"`
	YExpr string `json:"y_expr,omitempty"`
}

// PictureFill is the read model for an <a:blipFill> used as a shape fill.
// RelID is the raw r:embed relationship id; ImagePart is that relationship
// resolved to a package part path, and is empty when the fill links to an
// external image (r:link) or the relationship is missing.
type PictureFill struct {
	RelID     string `json:"rel_id,omitempty"`
	ImagePart string `json:"image_part,omitempty"`
	// External is the target of an r:link fill, which is not embedded in the
	// package.
	External string `json:"external,omitempty"`
	// Mode is "stretch" or "tile"; empty when the fill declares neither.
	Mode string `json:"mode,omitempty"`
	// Crop holds the <a:srcRect> insets in percent, when present.
	Crop *PictureFillCrop `json:"crop,omitempty"`
}

// PictureFillCrop holds <a:srcRect> insets as fractions of the source image.
type PictureFillCrop struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}
