package editorcommon

// TableStyleDefinition describes a custom table style entry.
type TableStyleDefinition struct {
	StyleID string `json:"style_id,omitempty"`
	Name    string `json:"name"`
}

// TableStyleInfo is a lightweight table style listing record.
type TableStyleInfo struct {
	StyleID string `json:"style_id"`
	Name    string `json:"name"`
}
