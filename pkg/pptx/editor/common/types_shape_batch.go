package editorcommon

// TextboxInsert describes one textbox to add in a bulk slide insert operation.
type TextboxInsert struct {
	Left    float64 `json:"left"`
	Top     float64 `json:"top"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	ShapeID *int    `json:"shape_id,omitempty"`
	Text    string  `json:"text,omitempty"`
}

// ConnectorInsert describes one connector to add in a bulk slide insert operation.
type ConnectorInsert struct {
	ShapeUpdate

	ConnectorType string  `json:"connector_type"`
	BeginX        float64 `json:"begin_x"`
	BeginY        float64 `json:"begin_y"`
	EndX          float64 `json:"end_x"`
	EndY          float64 `json:"end_y"`
	ShapeID       *int    `json:"shape_id,omitempty"`
}
