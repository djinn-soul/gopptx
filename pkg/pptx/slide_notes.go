package pptx

import "strings"

// WithNotes sets speaker notes for this slide.
func (s SlideContent) WithNotes(notes string) SlideContent {
	s.Notes = strings.TrimSpace(notes)
	return s
}
