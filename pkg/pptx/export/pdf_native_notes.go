package export

import (
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// No export path could emit speaker notes, so a deck's notes were lost on the
// way to PDF even though the model carries them. They go on their own page
// after the slide they belong to, which keeps the slide pages a faithful
// rendering and the notes readable.

const (
	notesPageMarginPt   = 54 // 0.75 inch
	notesHeadingSizePt  = 16
	notesBodySizePt     = 12
	notesLineSpacing    = 1.35
	frontTitleSizePt    = 32
	frontSubtitleSizePt = 14
)

// renderPDFNotesPage writes the notes for one slide. It returns false when the
// slide has none, so the caller can skip the page entirely.
func renderPDFNotesPage(pdf *gopdf.GoPdf, slide elements.SlideContent, slideNumber int, page pageSize) bool {
	paragraphs := notesParagraphs(slide)
	if len(paragraphs) == 0 {
		return false
	}

	pdf.AddPage()
	textWidth := page.WidthPt - 2*notesPageMarginPt
	y := float64(notesPageMarginPt)

	setPDFTextFontWithHint(pdf, notesHeadingSizePt, true, false, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetX(notesPageMarginPt)
	pdf.SetY(y)
	heading := "Notes — slide " + strconv.Itoa(slideNumber)
	if title := strings.TrimSpace(slide.Title); title != "" {
		heading += ": " + title
	}
	_ = pdf.Cell(nil, heading)
	y += pdfLineHeight(notesHeadingSizePt) * notesLineSpacing

	setPDFTextFontWithHint(pdf, notesBodySizePt, false, false, "")
	lineHeight := pdfLineHeight(notesBodySizePt) * notesLineSpacing
	for _, paragraph := range paragraphs {
		for _, line := range wrapPDFTextWithMetrics(pdf, paragraph, textWidth) {
			if y+lineHeight > page.HeightPt-notesPageMarginPt {
				return true
			}
			pdf.SetX(notesPageMarginPt)
			pdf.SetY(y)
			_ = pdf.Cell(nil, line)
			y += lineHeight
		}
		y += lineHeight / 2
	}
	return true
}

// renderPDFFrontmatter writes a title page for the deck.
func renderPDFFrontmatter(pdf *gopdf.GoPdf, title string, slideCount int, page pageSize) {
	pdf.AddPage()

	setPDFTextFontWithHint(pdf, frontTitleSizePt, true, false, "")
	pdf.SetTextColor(0, 0, 0)
	deckTitle := strings.TrimSpace(title)
	if deckTitle == "" {
		deckTitle = "Presentation"
	}
	textWidth := page.WidthPt - 2*notesPageMarginPt
	y := page.HeightPt/2 - pdfLineHeight(frontTitleSizePt)

	for _, line := range wrapPDFTextWithMetrics(pdf, deckTitle, textWidth) {
		pdf.SetX(alignedTextX(pdf, line, notesPageMarginPt, textWidth, "ctr"))
		pdf.SetY(y)
		_ = pdf.Cell(nil, line)
		y += pdfLineHeight(frontTitleSizePt) * notesLineSpacing
	}

	setPDFTextFontWithHint(pdf, frontSubtitleSizePt, false, false, "")
	subtitle := strconv.Itoa(slideCount) + " slides"
	if slideCount == 1 {
		subtitle = "1 slide"
	}
	pdf.SetX(alignedTextX(pdf, subtitle, notesPageMarginPt, textWidth, "ctr"))
	pdf.SetY(y)
	_ = pdf.Cell(nil, subtitle)

	setPDFTextFontWithHint(pdf, defaultFontSize, false, false, "")
}
