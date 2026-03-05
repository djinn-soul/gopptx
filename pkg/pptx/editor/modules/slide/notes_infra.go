package slide

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type HasPartFn func(string) bool
type GetPartFn func(string) ([]byte, bool)
type SetPartFn func(string, []byte)

func EnsureNotesMasterThemePart(hasPart HasPartFn, getPart GetPartFn, setPart SetPartFn) {
	if hasPart("ppt/theme/theme2.xml") {
		return
	}
	theme1, ok := getPart("ppt/theme/theme1.xml")
	if !ok {
		return
	}
	setPart("ppt/theme/theme2.xml", CloneBytes(theme1))
}

func EnsureNotesInfrastructure(
	hasPart HasPartFn,
	setPart SetPartFn,
	nonSlideRels []common.EditorRelationship,
	nextRelIDNum int,
	notesMasterThemeIndex int,
) ([]common.EditorRelationship, int) {
	if hasPart("ppt/notesMasters/notesMaster1.xml") {
		return nonSlideRels, nextRelIDNum
	}

	setPart("ppt/notesMasters/notesMaster1.xml", []byte(pptxxml.NotesMaster(nil)))
	setPart(
		"ppt/notesMasters/_rels/notesMaster1.xml.rels",
		[]byte(pptxxml.NotesMasterRelationships(notesMasterThemeIndex, nil)),
	)

	for _, rel := range nonSlideRels {
		if rel.Type == common.RelTypeNotesMaster {
			return nonSlideRels, nextRelIDNum
		}
	}

	nonSlideRels = append(nonSlideRels, common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", nextRelIDNum),
		Type:   common.RelTypeNotesMaster,
		Target: "notesMasters/notesMaster1.xml",
	})
	return nonSlideRels, nextRelIDNum + 1
}
