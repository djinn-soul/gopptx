package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

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

func addNotesMasterFiles(
	pw *pptxxml.PackageWriter,
	meta Metadata,
	masterCount, notesThemeIndex int,
	mc *media.Catalog,
) {
	if notesThemeIndex == 0 {
		return
	}

	var backgroundRID string
	var mediaName []string
	if meta.NotesMaster != nil && meta.NotesMaster.Background != nil &&
		meta.NotesMaster.Background.Type == elements.SlideBackgroundPicture &&
		meta.NotesMaster.Background.PictureFill != nil {
		if name, ok := mc.MediaNameForImage(*meta.NotesMaster.Background.PictureFill); ok {
			backgroundRID = "rId2"
			mediaName = []string{"../media/" + name}
		}
	}
	spec := elements.MapNotesMasterToSpec(meta.NotesMaster, backgroundRID)
	pw.AddPart("ppt/notesMasters/notesMaster1.xml", pptxxml.NotesMaster(spec))
	pw.AddPart(
		"ppt/notesMasters/_rels/notesMaster1.xml.rels",
		pptxxml.NotesMasterRelationships(notesThemeIndex, mediaName),
	)

	if notesThemeIndex > masterCount {
		pw.AddPart(fmt.Sprintf("ppt/theme/theme%d.xml", notesThemeIndex), pptxxml.Theme(mapThemeToSpec(meta.Theme)))
	}
}

func addHandoutMasterFiles(pw *pptxxml.PackageWriter, meta Metadata, masterCount int) {
	if meta.HandoutMaster == nil {
		return
	}
	pw.AddPart("ppt/handoutMasters/handoutMaster1.xml", meta.HandoutMaster.GenerateXML())
	pw.AddPart(
		"ppt/handoutMasters/_rels/handoutMaster1.xml.rels",
		handout.RelationshipsXML(masterCount),
	)
}
