package editor

import (
	"fmt"
	"path"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func (e *PresentationEditor) renderSlideNotesTarget(
	slide elements.SlideContent,
	slideNumber int,
	existingNotesTarget string,
) string {
	notesTarget := strings.TrimSpace(existingNotesTarget)
	if strings.TrimSpace(slide.Notes) == "" {
		return notesTarget
	}

	e.ensureNotesInfrastructure()
	slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	notesPath := e.ensureSlideNotesPart(slidePath)
	e.parts.Set(notesPath, []byte(pptxxml.NotesSlide(editorslide.EditorNotesBody(slide))))

	notesRelsPath := common.SlideRelsPartName(notesPath)
	e.parts.Set(notesRelsPath, []byte(pptxxml.NotesSlideRelationships(slideNumber)))
	return "../notesSlides/" + path.Base(notesPath)
}

func (e *PresentationEditor) ensureSlideNotesPart(slidePath string) string {
	notesPath, ok := e.notesInventory[slidePath]
	if ok {
		return notesPath
	}
	notesPath = fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", e.nextNotesNum)
	e.nextNotesNum++
	e.notesInventory[slidePath] = notesPath
	return notesPath
}
