package pptxxml

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func notesParagraphsXML(paragraphs []text.Paragraph) string {
	if len(paragraphs) == 0 {
		return `
<a:p><a:endParaRPr lang="en-US"/></a:p>`
	}
	var b strings.Builder
	defaultStyle := ContentStyleSpec{SizePt: defaultNotesFontSizePt}

	for _, paragraph := range paragraphs {
		styleSpec := convertNotesStyle(paragraph.Style)
		runSpecs := make([]TextRunSpec, len(paragraph.Runs))
		for i, run := range paragraph.Runs {
			runSpecs[i] = convertNotesRun(run)
		}
		b.WriteString(bulletParagraphRuns(runSpecs, styleSpec, defaultStyle))
	}
	return b.String()
}
