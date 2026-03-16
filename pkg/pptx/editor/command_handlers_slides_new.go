package editor

import (
	"encoding/json"
)

// handleDuplicateSlideAfter clones a slide and appends it immediately after.
//
// Payload: {"slide_index": N}.
// Response: {"new_index": M}.
func handleDuplicateSlideAfter(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	newIdx, err := e.DuplicateSlideAfter(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]int{"new_index": newIdx}, nil
}

// handleValidate runs structural validation and returns issues.
//
// Payload: {} (empty).
// Response: {"issues": [...], "issue_count": N}.
func handleValidate(e *PresentationEditor, _ json.RawMessage) (any, error) {
	issues := e.Validate()
	return map[string]any{
		"issues":      issues,
		"issue_count": len(issues),
	}, nil
}

// handleRepair attempts to automatically repair structural issues.
//
// Payload: {} (empty).
// Response: {"repaired_count": N, "unrepaired_count": M, "repaired": [...], "unrepaired": [...]}.
func handleRepair(e *PresentationEditor, _ json.RawMessage) (any, error) {
	result, err := e.Repair()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"repaired_count":   len(result.IssuesRepaired),
		"unrepaired_count": len(result.IssuesUnrepaired),
		"repaired":         result.IssuesRepaired,
		"unrepaired":       result.IssuesUnrepaired,
	}, nil
}
