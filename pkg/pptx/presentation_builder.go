package pptx

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// PresentationBuilder provides a fluent API for creating presentations.
type PresentationBuilder struct {
	title           string
	metadata        Metadata
	slides          []SlideContent
	showSlideNumber bool
	footerText      string
	showDateTime    bool
}

// NewPresentationBuilder creates a new presentation builder with a title.
func NewPresentationBuilder(title string) *PresentationBuilder {
	return &PresentationBuilder{
		title: title,
	}
}

// WithMetadata sets the presentation metadata.
func (b *PresentationBuilder) WithMetadata(meta Metadata) *PresentationBuilder {
	b.metadata = meta
	return b
}

// WithSlideSize sets the slide dimensions for the presentation.
func (b *PresentationBuilder) WithSlideSize(size SlideSize) *PresentationBuilder {
	b.metadata.SlideSize = size
	return b
}

// WithTheme sets the theme (colors and fonts) for the presentation.
func (b *PresentationBuilder) WithTheme(theme styling.Theme) *PresentationBuilder {
	b.metadata.Theme = &theme
	return b
}

// WithMaster sets the slide master (background, footer, shapes) for the presentation.
func (b *PresentationBuilder) WithMaster(master *elements.SlideMaster) *PresentationBuilder {
	b.metadata.Master = master
	return b
}

// AddSlide adds a slide to the presentation.
func (b *PresentationBuilder) AddSlide(slide SlideContent) *PresentationBuilder {
	if b.showSlideNumber {
		slide = slide.WithSlideNumber(true)
	}
	b.slides = append(b.slides, slide)
	return b
}

// AddTitleSlide adds a title-only slide to the presentation.
func (b *PresentationBuilder) AddTitleSlide(title string) *PresentationBuilder {
	return b.AddSlide(NewSlide(title))
}

// AddBulletSlide adds a slide with title and bullet points.
func (b *PresentationBuilder) AddBulletSlide(title string, bullets []string) *PresentationBuilder {
	slide := NewSlide(title)
	for _, bullet := range bullets {
		slide = slide.AddBullet(bullet)
	}
	return b.AddSlide(slide)
}

// AddShapesSlide adds a slide with title and multiple shapes.
func (b *PresentationBuilder) AddShapesSlide(title string, shapes ...Shape) *PresentationBuilder {
	slide := NewSlide(title)
	for _, s := range shapes {
		slide = slide.AddShape(s)
	}
	return b.AddSlide(slide)
}

// AddCustomXML adds a custom XML part to the presentation.
func (b *PresentationBuilder) AddCustomXML(content string) *PresentationBuilder {
	b.metadata.CustomXML = append(b.metadata.CustomXML, CustomXMLPart{Content: content})
	return b
}

// WithSlideNumbers enables or disables slide number display for all slides in the presentation.
func (b *PresentationBuilder) WithSlideNumbers(show bool) *PresentationBuilder {
	b.showSlideNumber = show
	for i := range b.slides {
		b.slides[i] = b.slides[i].WithSlideNumber(show)
	}
	return b
}

// WithFooter sets the footer text for all slides in the presentation.
func (b *PresentationBuilder) WithFooter(text string) *PresentationBuilder {
	b.footerText = text
	return b
}

// WithDateTime enables or disables date/time display for all slides in the presentation.
func (b *PresentationBuilder) WithDateTime(show bool) *PresentationBuilder {
	b.showDateTime = show
	return b
}

// WithModifyPassword sets a "Password to Modify" on the presentation.
func (b *PresentationBuilder) WithModifyPassword(password string) *PresentationBuilder {
	b.metadata.Protection.ModifyPassword = password
	return b
}

// WithMarkAsFinal sets the "Mark as Final" property on the presentation.
func (b *PresentationBuilder) WithMarkAsFinal(final bool) *PresentationBuilder {
	b.metadata.Protection.MarkAsFinal = final
	return b
}

// WithSignaturesEnabled enables the digital signature origin part placeholder.
func (b *PresentationBuilder) WithSignaturesEnabled(enabled bool) *PresentationBuilder {
	b.metadata.Protection.SignaturesEnabled = enabled
	return b
}

// Build compiles the presentation into a PPTX byte slice.
func (b *PresentationBuilder) Build() ([]byte, error) {
	// Metadata title overrides builder title if both are present.
	if b.metadata.Title == "" {
		b.metadata.Title = b.title
	}

	b.metadata.FooterText = b.footerText
	b.metadata.ShowDateTime = b.showDateTime

	return CreateWithMetadata(b.metadata, b.slides)
}

// WriteToFile builds the presentation and writes it to the specified file path.
func (b *PresentationBuilder) WriteToFile(path string) error {
	data, err := b.Build()
	if err != nil {
		return fmt.Errorf("failed to build presentation: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}
