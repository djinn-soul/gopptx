package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// TextRun describes a single piece of text with uniform styling.
	TextRun = elements.TextRun
)

func NewTextRun(text string) TextRun {
	return elements.NewTextRun(text)
}

func NormalizeTextRuns(runs []TextRun) []TextRun {
	return elements.NormalizeTextRuns(runs)
}
