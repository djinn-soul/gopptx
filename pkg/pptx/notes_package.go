package pptx

import (
	"archive/zip"
	"fmt"
)

func writeNotesFiles(zw *zip.Writer, parts []renderedNotesPart) error {
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
