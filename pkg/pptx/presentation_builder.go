package pptx

import (
	"fmt"
	"os"
)

// PresentationBuilder provides a fluent API for creating presentations.
type PresentationBuilder struct {
	title    string
	metadata PresentationMetadata
	slides   []SlideContent
}

// NewPresentationBuilder creates a new presentation builder with a title.
func NewPresentationBuilder(title string) *PresentationBuilder {
	return &PresentationBuilder{
		title: title,
	}
}

// WithMetadata sets the presentation metadata.
func (b *PresentationBuilder) WithMetadata(meta PresentationMetadata) *PresentationBuilder {
	b.metadata = meta
	return b
}

// AddSlide adds a slide to the presentation.
func (b *PresentationBuilder) AddSlide(slide SlideContent) *PresentationBuilder {
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

// Build compiles the presentation into a PPTX byte slice.
func (b *PresentationBuilder) Build() ([]byte, error) {
	// Metadata title overrides builder title if both are present.
	if b.metadata.Title == "" {
		b.metadata.Title = b.title
	}

	return CreateWithMetadata(b.metadata, b.slides)
}

// WriteToFile builds the presentation and writes it to the specified file path.
func (b *PresentationBuilder) WriteToFile(path string) error {
	data, err := b.Build()
	if err != nil {
		return fmt.Errorf("failed to build presentation: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
