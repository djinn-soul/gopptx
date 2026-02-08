package pptx

import (
	"archive/zip"
	"fmt"

	"github.com/djinn09/goppt/internal/pptxxml"
)

func writeNotesFiles(zw *zip.Writer, parts []renderedNotesPart) error {
	if len(parts) > 0 {
		if err := writeFile(zw, "ppt/theme/theme2.xml", pptxxml.Theme()); err != nil {
			return err
		}
	}

	for _, part := range parts {
		path := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", part.slideNumber)
		if err := writeFile(zw, path, part.slideXML); err != nil {
			return err
		}

		relsPath := fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", part.slideNumber)
		if err := writeFile(zw, relsPath, part.relsXML); err != nil {
			return err
		}
	}
	return nil
}
