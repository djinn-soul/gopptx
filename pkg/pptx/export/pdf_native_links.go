package export

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// PowerPoint's own PDF export keeps a deck's links clickable. The reader has
// always recovered them — a shape's click action, a run's hyperlink — and the
// native renderer used to drop every one of them on the floor, so an exported
// deck's links were dead. These turn them into PDF link annotations.

// setNativePDFDocumentInfo fills the document information dictionary, which is
// what a viewer shows as the file's title and what a search index reads. It was
// left empty, so an exported deck opened as an untitled document.
func setNativePDFDocumentInfo(pdf *gopdf.GoPdf, title string) {
	pdf.SetInfo(gopdf.PdfInfo{
		Title:        strings.TrimSpace(title),
		Creator:      pdfProducerName,
		Producer:     pdfProducerName,
		CreationDate: time.Now(),
	})
}

// pdfProducerName is how the exporter names itself in the document information.
const pdfProducerName = "gopptx"

// slideOutlineTitle names a slide's bookmark. A slide with no title of its own
// is still worth a bookmark, so it is numbered instead.
func slideOutlineTitle(slide elements.SlideContent, index int) string {
	if title := strings.TrimSpace(slide.Title); title != "" {
		return title
	}
	return "Slide " + strconv.Itoa(index)
}

// slidePDFAnchor names the page a slide was drawn on, so a link to slide N can
// find it. gopdf resolves an internal link by anchor name.
func slidePDFAnchor(slideNumber int) string {
	return "slide-" + strconv.Itoa(slideNumber)
}

// addPDFHyperlink puts a link annotation over the given box, in the renderer's
// own top-left coordinates. A link the PDF format cannot express — running a
// program, say — is left out rather than pointed somewhere wrong.
func addPDFHyperlink(pdf *gopdf.GoPdf, link *action.Hyperlink, x, y, w, h float64) {
	if link == nil || w <= 0 || h <= 0 {
		return
	}
	if link.Action.Type == action.HyperlinkActionSlide {
		if link.Action.SlideNumber == 0 {
			return
		}
		pdf.AddInternalLink(slidePDFAnchor(int(link.Action.SlideNumber)), x, y, w, h)
		return
	}
	target := pdfHyperlinkTarget(link.Action)
	if target == "" {
		return
	}
	pdf.AddExternalLink(target, x, y, w, h)
}

// pdfHyperlinkTarget is the URI a PDF viewer should open, or empty for an action
// with no URI form.
func pdfHyperlinkTarget(target action.HyperlinkAction) string {
	switch target.Type {
	case action.HyperlinkActionURL:
		return strings.TrimSpace(target.URL)
	case action.HyperlinkActionEmail:
		return mailtoURI(target.EmailAddress, target.EmailSubject)
	case action.HyperlinkActionFile:
		return fileURI(target.FilePath)
	default:
		// Slide jumps are internal links, and the remaining actions — running a
		// program, the show-control jumps — have no PDF equivalent.
		return ""
	}
}

func mailtoURI(address, subject string) string {
	address = strings.TrimSpace(address)
	if address == "" {
		return ""
	}
	uri := "mailto:" + address
	if subject = strings.TrimSpace(subject); subject != "" {
		uri += "?subject=" + url.QueryEscape(subject)
	}
	return uri
}

func fileURI(filePath string) string {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return ""
	}
	if strings.Contains(filePath, "://") {
		return filePath
	}
	return "file:///" + strings.TrimPrefix(strings.ReplaceAll(filePath, `\`, "/"), "/")
}

// shapeClickAction is the link a shape is clicked through. ClickAction is the
// current field; Hyperlink is the older one kept for source compatibility, so it
// only applies when the shape carries no click action of its own.
func shapeClickAction(clickAction, legacy *action.Hyperlink) *action.Hyperlink {
	if clickAction != nil {
		return clickAction
	}
	return legacy
}
