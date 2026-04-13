package presentation

import (
	"archive/zip"
	"strconv"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/notes"
)

func WritePresentationPackage(
	zw *zip.Writer,
	meta Metadata,
	slides []elements.SlideContent,
	slideCount int,
) error {
	pw := pptxxml.NewPackageWriter()

	mediaCatalog, mediaErr := media.BuildMediaCatalog(slides, meta.NotesMaster)
	if mediaErr != nil {
		return mediaErr
	}

	chartParts := BuildChartParts(slides)
	smartArtParts := BuildSmartArtParts(slides)
	notesParts := notes.BuildRenderedNotesParts(slides)
	effectiveMasters := getEffectiveMasters(meta)
	masterCount := len(effectiveMasters)
	notesThemeIndex := getNotesThemeIndex(len(notesParts) > 0)
	authors, commentsBySlide, commentSlideIndices := prepareComments(meta, slides)
	hasVBA := meta.VBA.IsMacroEnabled()

	if err := addBasicPropertyFiles(
		pw, meta, slides, slideCount, len(notesParts), ChartPartCount(chartParts), SmartArtPartCount(smartArtParts),
		notesParts, masterCount, notesThemeIndex, mediaCatalog.ImageExtensions(),
		authors, commentSlideIndices, hasVBA,
	); err != nil {
		return err
	}
	addLayoutFiles(pw, masterCount)
	addMasterFiles(pw, effectiveMasters, mediaCatalog)
	addThemeFiles(pw, meta.Theme, masterCount)
	addNotesMasterFiles(pw, meta, masterCount, notesThemeIndex, mediaCatalog)
	addHandoutMasterFiles(pw, meta, masterCount)

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
	if err := renderSlides(
		pw,
		meta,
		slides,
		mediaCatalog,
		chartBySlide,
		smartArtBySlide,
		notesTargets,
		masterCount,
		commentsBySlide,
	); err != nil {
		return err
	}

	if len(authors) > 0 {
		pw.AddPart("ppt/commentAuthors.xml", pptxxml.CommentAuthorsXML(authors))
	}
	if err := writeCustomXMLParts(pw, meta.CustomXML); err != nil {
		return err
	}
	if err := writeVBAParts(pw, meta.VBA); err != nil {
		return err
	}
	for i, f := range meta.EmbeddedFonts {
		path := "ppt/fonts/font" + strconv.Itoa(i+1) + ".fntdata"
		pw.AddBinaryPart(path, f.Data)
	}

	return pw.WriteTo(zw)
}

func WritePackageFiles(
	zw *zip.Writer,
	meta Metadata,
	slides []elements.SlideContent,
	slideCount int,
) error {
	return WritePresentationPackage(zw, meta, slides, slideCount)
}
