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

const (
	minMasterCountWithNativeNotesTheme = 2
	singleMasterNotesThemeIndex        = 2
)

// Metadata defines non-content properties of a PPTX.
type Metadata struct {
	common.Metadata

	Theme       *styling.Theme
	Master      *elements.SlideMaster
	Masters     []*elements.SlideMaster
	NotesMaster *elements.NotesMaster
}

// SlideSize defines presentation dimensions in EMUs.
type SlideSize = common.SlideSize

// GetSlideSize4x3 returns the standard 4:3 slide size.
func GetSlideSize4x3() SlideSize {
	return common.GetSlideSize4x3()
}

// GetSlideSize16x9 returns the standard 16:9 widescreen slide size.
func GetSlideSize16x9() SlideSize {
	return common.GetSlideSize16x9()
}

func WritePackageFiles(
	zw *zip.Writer,
	meta Metadata,
	slides []elements.SlideContent,
	slideCount int,
) error {
	pw := pptxxml.NewPackageWriter()

	mediaCatalog, mediaErr := media.BuildMediaCatalog(slides)
	if mediaErr != nil {
		return mediaErr
	}

	chartParts := BuildChartParts(slides)
	smartArtParts := BuildSmartArtParts(slides)
	notesParts := notes.BuildRenderedNotesParts(slides)
	effectiveMasters := getEffectiveMasters(meta)
	masterCount := len(effectiveMasters)
	notesThemeIndex := getNotesThemeIndex(len(notesParts) > 0, masterCount)

	addBasicPropertyFiles(
		pw, meta, slideCount, len(notesParts), len(chartParts), len(smartArtParts),
		notesParts, masterCount, notesThemeIndex, mediaCatalog.ImageExtensions(),
	)
	addLayoutFiles(pw, masterCount)
	addMasterFiles(pw, effectiveMasters, mediaCatalog)
	addThemeFiles(pw, meta.Theme, masterCount)
	addNotesMasterFiles(pw, meta, masterCount, notesThemeIndex)

	if err := writeMediaFiles(pw, mediaCatalog); err != nil {
		return err
	}
	if err := writeChartFiles(pw, chartParts); err != nil {
		return err
	}
	if err := writeSmartArtFiles(pw, smartArtParts); err != nil {
		return err
	}
	if err := notes.WriteNotesFiles(pw, notesParts); err != nil {
		return err
	}

	chartBySlide := chartPartBySlide(chartParts)
	smartArtBySlide := smartArtPartBySlide(smartArtParts)
	notesTargets := notes.TargetBySlide(notesParts)
	if err := renderSlides(pw, meta, slides, mediaCatalog, chartBySlide, smartArtBySlide, notesTargets, masterCount); err != nil {
		return err
	}

	if err := writeCustomXMLParts(pw, meta.CustomXML); err != nil {
		return err
	}

	return pw.WriteTo(zw)
}

func getEffectiveMasters(meta Metadata) []*elements.SlideMaster {
	if len(meta.Masters) > 0 {
		return meta.Masters
	}
	if meta.Master != nil {
		return []*elements.SlideMaster{meta.Master}
	}
	return []*elements.SlideMaster{elements.NewMaster()}
}

func getNotesThemeIndex(hasNotes bool, masterCount int) int {
	if !hasNotes {
		return 0
	}
	if masterCount >= minMasterCountWithNativeNotesTheme {
		return masterCount + 1
	}
	return singleMasterNotesThemeIndex
}

func addBasicPropertyFiles(
	pw *pptxxml.PackageWriter,
	meta Metadata,
	slideCount, notesPartCount, chartPartCount, smartArtPartCount int,
	notesParts []notes.RenderedNotesPart,
	masterCount, notesThemeIndex int,
	mediaExtensions []string,
) {
	hasNotes := notesPartCount > 0
	pw.AddPart("[Content_Types].xml", pptxxml.ContentTypes(
		slideCount, mediaExtensions, chartPartCount, smartArtPartCount,
		notes.SlideNumbers(notesParts), hasNotes,
		len(meta.CustomXML), masterCount, notesThemeIndex,
	))
	pw.AddPart("_rels/.rels", pptxxml.RootRelationships())
	pw.AddPart(
		"ppt/_rels/presentation.xml.rels",
		pptxxml.PresentationRelationships(slideCount, hasNotes, len(meta.CustomXML), masterCount),
	)
	pw.AddPart(
		"ppt/presentation.xml",
		pptxxml.Presentation(
			meta.Title, slideCount, hasNotes,
			meta.SlideSize.Width, meta.SlideSize.Height, masterCount,
		),
	)
	pw.AddPart("docProps/core.xml", pptxxml.CoreProperties(pptxxml.CorePropertiesInfo{
		Title: meta.Title, Subject: meta.Subject, Creator: meta.Creator, Description: meta.Description,
	}))
	pw.AddPart(
		"docProps/app.xml",
		pptxxml.AppProperties(slideCount, notesPartCount, meta.SlideSize.Width, meta.SlideSize.Height),
	)
}

func addLayoutFiles(pw *pptxxml.PackageWriter, masterCount int) {
	layoutXMLs := []string{
		pptxxml.SlideLayoutTitleAndContent(),
		pptxxml.SlideLayoutTitleOnly(),
		pptxxml.SlideLayoutBlank(),
		pptxxml.SlideLayoutCenteredTitle(),
		pptxxml.SlideLayoutTitleAndBigContent(),
		pptxxml.SlideLayoutTwoColumn(),
	}
	for masterNum := 1; masterNum <= masterCount; masterNum++ {
		for i, xml := range layoutXMLs {
			idx := (masterNum-1)*len(layoutXMLs) + (i + 1)
			name := fmt.Sprintf("slideLayout%d.xml", idx)
			pw.AddPart(fmt.Sprintf("ppt/slideLayouts/%s", name), xml)
			pw.AddPart(fmt.Sprintf("ppt/slideLayouts/_rels/%s.rels", name), pptxxml.SlideLayoutRelationships(masterNum))
		}
	}
}

func addMasterFiles(pw *pptxxml.PackageWriter, masters []*elements.SlideMaster, mc *media.Catalog) {
	for i, master := range masters {
		masterNum := i + 1
		targets, refs := buildMasterImageInfo(master, mc)
		spec := mapMasterToSpec(master, refs)
		if spec != nil {
			spec.MasterIndex = masterNum
		}
		pw.AddPart(fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", masterNum), pptxxml.SlideMaster(spec))
		pw.AddPart(
			fmt.Sprintf("ppt/slideMasters/_rels/slideMaster%d.xml.rels", masterNum),
			pptxxml.SlideMasterRelationships(targets, masterNum, masterNum),
		)
	}
}

func addThemeFiles(pw *pptxxml.PackageWriter, theme *styling.Theme, masterCount int) {
	themeXML := pptxxml.Theme(mapThemeToSpec(theme))
	for i := 1; i <= masterCount; i++ {
		pw.AddPart(fmt.Sprintf("ppt/theme/theme%d.xml", i), themeXML)
	}
}

func addNotesMasterFiles(pw *pptxxml.PackageWriter, meta Metadata, masterCount, notesThemeIndex int) {
	if notesThemeIndex == 0 {
		return
	}
	spec := elements.MapNotesMasterToSpec(meta.NotesMaster)
	pw.AddPart("ppt/notesMasters/notesMaster1.xml", pptxxml.NotesMaster(spec))
	pw.AddPart("ppt/notesMasters/_rels/notesMaster1.xml.rels", pptxxml.NotesMasterRelationships(notesThemeIndex))
	if notesThemeIndex > masterCount {
		pw.AddPart(fmt.Sprintf("ppt/theme/theme%d.xml", notesThemeIndex), pptxxml.Theme(mapThemeToSpec(meta.Theme)))
	}
}
