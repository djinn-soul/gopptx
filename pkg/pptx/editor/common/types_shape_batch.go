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
