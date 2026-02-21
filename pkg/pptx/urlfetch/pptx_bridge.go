package urlfetch

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// presentationCreate generates PPTX bytes from a title and slide list.
func presentationCreateWithMetadata(
	title, creator string,
	slides []elements.SlideContent,
) ([]byte, error) {
	meta := pptx.Metadata{
		Metadata: pptx.MetadataFields{
			Title:   title,
			Creator: creator,
		},
	}
	return pptx.CreateWithMetadata(meta, slides)
}
