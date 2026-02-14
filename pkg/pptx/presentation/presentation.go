package presentation

import (
	"archive/zip"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/notes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// PresentationMetadata defines non-content properties of a PPTX.
type PresentationMetadata struct {
	common.PresentationMetadata
	Theme       *styling.Theme
	Master      *elements.SlideMaster
	Masters     []*elements.SlideMaster
	NotesMaster *elements.NotesMaster
}

// SlideSize defines presentation dimensions in EMUs.
type SlideSize = common.SlideSize

// Default slide sizes.
var (
	SlideSize4x3  = common.SlideSize4x3
	SlideSize16x9 = common.SlideSize16x9
)

func WritePackageFiles(zw *zip.Writer, meta PresentationMetadata, slides []elements.SlideContent, slideCount int) error {
	pw := pptxxml.NewPackageWriter()

	mediaCatalog, err := media.BuildMediaCatalog(slides)
	if err != nil {
		return err
	}

	effectiveMasters := meta.Masters
	if len(effectiveMasters) == 0 && meta.Master != nil {
		effectiveMasters = []*elements.SlideMaster{meta.Master}
	}
	if len(effectiveMasters) == 0 {
		effectiveMasters = []*elements.SlideMaster{elements.NewMaster()}
	}
	// Current slide layout selection API binds slides to the canonical layout set only.
	// Until per-slide master/layout binding is implemented, emitting additional masters
	// creates unreferenced master/layout families that PowerPoint can reject.
	if len(effectiveMasters) > 1 {
		effectiveMasters = effectiveMasters[:1]
	}
	masterCount := len(effectiveMasters)

	chartParts := BuildChartParts(slides)
	chartBySlide := chartPartBySlide(chartParts)
	notesParts := notes.BuildRenderedNotesParts(slides)
	notesTargets := notes.NotesTargetBySlide(notesParts)
	hasNotes := len(notesParts) > 0

	files := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", pptxxml.ContentTypes(slideCount, mediaCatalog.ImageExtensions(), len(chartParts), notes.NotesSlideNumbers(notesParts), hasNotes, len(meta.CustomXML), masterCount)},
		{"_rels/.rels", pptxxml.RootRelationships()},
		{"ppt/_rels/presentation.xml.rels", pptxxml.PresentationRelationships(slideCount, hasNotes, len(meta.CustomXML), masterCount)},
		{"ppt/presentation.xml", pptxxml.Presentation(meta.Title, slideCount, hasNotes, meta.SlideSize.Width, meta.SlideSize.Height, masterCount)},
		{"ppt/theme/theme1.xml", pptxxml.Theme(mapThemeToSpec(meta.Theme))},
		{"docProps/core.xml", pptxxml.CoreProperties(pptxxml.CorePropertiesInfo{
			Title:       meta.Title,
			Subject:     meta.Subject,
			Creator:     meta.Creator,
			Description: meta.Description,
		})},
		{"docProps/app.xml", pptxxml.AppProperties(slideCount, len(notesParts), meta.SlideSize.Width, meta.SlideSize.Height)},
	}

	layoutXML := []string{
		pptxxml.SlideLayoutTitleAndContent(),
		pptxxml.SlideLayoutTitleOnly(),
		pptxxml.SlideLayoutBlank(),
		pptxxml.SlideLayoutCenteredTitle(),
		pptxxml.SlideLayoutTitleAndBigContent(),
		pptxxml.SlideLayoutTwoColumn(),
	}
	for masterNum := 1; masterNum <= masterCount; masterNum++ {
		for layoutIdx := 1; layoutIdx <= len(layoutXML); layoutIdx++ {
			layoutName := fmt.Sprintf("slideLayout%d.xml", layoutIdx)
			if masterNum > 1 {
				layoutName = fmt.Sprintf("slideLayout%d_m%d.xml", layoutIdx, masterNum)
			}
			files = append(files, struct {
				name    string
				content string
			}{fmt.Sprintf("ppt/slideLayouts/%s", layoutName), layoutXML[layoutIdx-1]})
			files = append(files, struct {
				name    string
				content string
			}{fmt.Sprintf("ppt/slideLayouts/_rels/%s.rels", layoutName), pptxxml.SlideLayoutRelationships(masterNum)})
		}
	}

	for i, master := range effectiveMasters {
		masterNum := i + 1
		targets, refs := buildMasterImageInfo(master, mediaCatalog)
		spec := mapMasterToSpec(master, refs)
		if spec != nil {
			spec.MasterIndex = masterNum
		}
		files = append(files, struct {
			name    string
			content string
		}{fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", masterNum), pptxxml.SlideMaster(spec)})
		files = append(files, struct {
			name    string
			content string
		}{fmt.Sprintf("ppt/slideMasters/_rels/slideMaster%d.xml.rels", masterNum), pptxxml.SlideMasterRelationships(targets, masterNum)})
	}

	if hasNotes {
		notesMasterSpec := elements.MapNotesMasterToSpec(meta.NotesMaster)
		files = append(files,
			struct {
				name    string
				content string
			}{"ppt/notesMasters/notesMaster1.xml", pptxxml.NotesMaster(notesMasterSpec)},
			struct {
				name    string
				content string
			}{"ppt/notesMasters/_rels/notesMaster1.xml.rels", pptxxml.NotesMasterRelationships()},
			struct {
				name    string
				content string
			}{"ppt/theme/theme2.xml", pptxxml.Theme(mapThemeToSpec(meta.Theme))},
		)
	}

	for _, item := range files {
		pw.AddPart(item.name, item.content)
	}

	if err := writeMediaFiles(pw, mediaCatalog); err != nil {
		return err
	}
	if err := writeChartFiles(pw, chartParts); err != nil {
		return err
	}
	if err := notes.WriteNotesFiles(pw, notesParts); err != nil {
		return err
	}

	if err := renderSlides(pw, meta, slides, mediaCatalog, chartBySlide, notesTargets); err != nil {
		return err
	}

	if err := writeCustomXMLParts(pw, meta.CustomXML); err != nil {
		return err
	}

	return pw.WriteTo(zw)
}
