package pptx

import (
	"errors"
	"fmt"
)

// Edit builds the presentation and reopens it for editing, bridging the two
// halves of the API.
//
// gopptx has a create path and an edit path with different shapes:
// PresentationBuilder assembles immutable SlideContent values, while
// Presentation wraps a mutable, index-addressed editor. Getting from one to
// the other used to mean serializing to disk and opening the file again:
//
//	if err := b.WriteToFile(path); err != nil { ... }
//	prs, err := pptx.Open(path)
//
// Edit does that in memory instead, so composing a deck and then editing it is
// a single call and needs no temporary file:
//
//	prs, err := pptx.NewPresentationBuilder("Report").
//	        AddSlide(pptx.NewSlide("Overview")).
//	        Edit()
//	if err != nil { ... }
//	defer prs.Close()
//	prs.SetSlideTitle(0, "Revised")
//
// The returned Presentation owns its own copy of the deck; further calls on the
// builder do not affect it. Because it was never read from a file, it has no
// path: use SaveAs, SaveToBytes or SaveToWriter rather than Save.
func (b *PresentationBuilder) Edit() (*Presentation, error) {
	if b == nil {
		return nil, errors.New("presentation builder is nil")
	}

	data, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("build presentation: %w", err)
	}

	prs, err := OpenFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("open built presentation: %w", err)
	}
	return prs, nil
}
