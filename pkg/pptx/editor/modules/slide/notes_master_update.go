package slide

import (
	"os"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func UpdateNotesMasterParts(
	master *elements.NotesMaster,
	setPart SetPartFn,
	notesMasterThemeIndex int,
	registerImage RegisterImageFn,
) error {
	backgroundRID, mediaNames, err := ResolveNotesMasterBackgroundMedia(
		master,
		os.ReadFile,
		registerImage,
	)
	if err != nil {
		return err
	}

	spec := elements.MapNotesMasterToSpec(master, backgroundRID)
	setPart("ppt/notesMasters/notesMaster1.xml", []byte(pptxxml.NotesMaster(spec)))
	setPart(
		"ppt/notesMasters/_rels/notesMaster1.xml.rels",
		[]byte(pptxxml.NotesMasterRelationships(notesMasterThemeIndex, mediaNames)),
	)
	return nil
}
