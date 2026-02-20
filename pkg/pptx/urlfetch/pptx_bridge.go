package urlfetch

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// presentationCreate generates PPTX bytes from a title and slide list.
func presentationCreate(title string, slides []elements.SlideContent) ([]byte, error) {
	return pptx.CreateWithSlides(title, slides)
}
