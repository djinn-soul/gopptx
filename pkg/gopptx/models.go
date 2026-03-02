package gopptx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

const (
	NSP             = "http://schemas.openxmlformats.org/presentationml/2006/main"
	defaultFilePerm = 0o600
)

// Presentation represents the main presentation XML component.
type Presentation struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main presentation"`
	Title   string   `xml:"-"`
	Slides  []*Slide `xml:"sldIdLst>sldId"`
}

// AddSlide adds a new blank slide to the presentation.
func (p *Presentation) AddSlide() *Slide {
	slide := &Slide{}
	p.Slides = append(p.Slides, slide)
	return slide
}

// Save writes the presentation to a .pptx file.
//
// The method serializes all slide data including bullets, images, shapes, connectors,
// transitions, notes, placeholder overrides, and animations into a valid PPTX file.
func (p *Presentation) Save(path string) error {
	slideCount := len(p.Slides)
	if slideCount < 1 {
		return errors.New("at least one slide is required")
	}

	slideTitle := p.Title
	if slideTitle == "" {
		slideTitle = "Presentation"
	}
	slides := make([]elements.SlideContent, 0, slideCount)
	for i, slide := range p.Slides {
		if slide == nil {
			slide = &Slide{}
		}
		slides = append(slides, slide.toSlideContent(i))
	}
	data, err := pptx.CreateWithSlides(slideTitle, slides)
	if err != nil {
		return err
	}
	return writeFileAtomically(path, data, defaultFilePerm)
}

// Slide represents a single slide XML component.
type Slide struct {
	XMLName              xml.Name                    `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
	Title                string                      `xml:"-"`
	Bullets              []string                    `xml:"-"`
	Images               []shapes.Image              `xml:"-"`
	Shapes               []shapes.Shape              `xml:"-"`
	Connectors           []shapes.Connector          `xml:"-"`
	Transition           transitions.SlideTransition `xml:"-"`
	Notes                string                      `xml:"-"`
	NotesBody            []elements.Paragraph        `xml:"-"`
	PlaceholderOverrides []shapes.PlaceholderContent `xml:"-"`
	Animations           []animations.Animation      `xml:"-"`
	Table                *tables.Table               `xml:"-"`
	BarChart             *charts.BarChart            `xml:"-"`
	SmartArt             []smartart.SmartArt         `xml:"-"`
}

func (s *Slide) AddBullet(text string) {
	s.Bullets = append(s.Bullets, text)
}

func (s *Slide) AddImage(img shapes.Image) {
	s.Images = append(s.Images, img)
}

func (s *Slide) AddShape(shape shapes.Shape) {
	s.Shapes = append(s.Shapes, shape)
}

func (s *Slide) AddConnector(conn shapes.Connector) {
	s.Connectors = append(s.Connectors, conn)
}

func (s *Slide) SetTransition(t transitions.SlideTransition) {
	s.Transition = t
}

func (s *Slide) SetTable(table tables.Table) {
	tmp := table
	s.Table = &tmp
}

func (s *Slide) SetBarChart(chart charts.BarChart) {
	tmp := chart
	s.BarChart = &tmp
}

func (s *Slide) AddSmartArt(sa smartart.SmartArt) {
	s.SmartArt = append(s.SmartArt, sa)
}

func (s *Slide) SetNotes(notes string) {
	s.Notes = notes
	p := elements.NewParagraph()
	p.Runs = append(p.Runs, elements.NewRun(notes))
	s.NotesBody = []elements.Paragraph{p}
}

func (s *Slide) SetRichNotes(body []elements.Paragraph) {
	s.NotesBody = append([]elements.Paragraph(nil), body...)
	var sb strings.Builder
	for i, p := range body {
		for _, run := range p.Runs {
			sb.WriteString(run.Text)
		}
		if i < len(body)-1 {
			sb.WriteString("\n")
		}
	}
	s.Notes = sb.String()
}

func (s *Slide) AddNoteParagraph(p elements.Paragraph) {
	s.NotesBody = append(s.NotesBody, p)
	if s.Notes != "" {
		s.Notes += "\n"
	}
	for _, run := range p.Runs {
		s.Notes += run.Text
	}
}

func (s *Slide) AddPlaceholderOverride(content shapes.PlaceholderContent) {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, content)
}

func (s *Slide) AddPlaceholderText(index int, textValue string) {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderTextType(index),
		Text:  textValue,
	})
}

func (s *Slide) AddPlaceholderImage(index int, img shapes.Image) {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderImageType(index),
		Image: &img,
	})
}

func (s *Slide) AddAnimation(anim animations.Animation) {
	s.Animations = append(s.Animations, anim)
}

func (s *Slide) AddAnimationDefinition(def animations.AnimationDefinition) {
	if def == nil {
		return
	}
	s.Animations = append(s.Animations, def.ToAnimation())
}

func (s *Slide) AddAnimationSequence(spacingMS uint32, defs ...animations.AnimationDefinition) {
	if len(defs) == 0 {
		return
	}
	for i, def := range defs {
		if def == nil {
			continue
		}
		anim := def.ToAnimation()
		if i == 0 {
			anim.Trigger = animations.AnimationOnClick
		} else if anim.Trigger == animations.AnimationOnClick {
			anim.Trigger = animations.AnimationAfterPrevious
		}
		if spacingMS > 0 {
			anim.DelayMS = spacingMS * uint32(i)
		}
		s.Animations = append(s.Animations, anim)
	}
}

func (s *Slide) toSlideContent(index int) elements.SlideContent {
	title := s.Title
	if title == "" {
		title = fmt.Sprintf("Slide %d", index+1)
	}

	slide := elements.NewSlide(title)
	for _, bullet := range s.Bullets {
		slide = slide.AddBullet(bullet)
	}
	if len(s.Images) > 0 {
		slide.Images = append(slide.Images, s.Images...)
	}
	if len(s.Shapes) > 0 {
		slide.Shapes = append(slide.Shapes, s.Shapes...)
	}
	if len(s.Connectors) > 0 {
		slide.Connectors = append(slide.Connectors, s.Connectors...)
	}
	if s.Transition != nil {
		slide.Transition = s.Transition
	}
	if s.Notes != "" {
		slide.Notes = s.Notes
	}
	if len(s.NotesBody) > 0 {
		slide.NotesBody = append(slide.NotesBody, s.NotesBody...)
	}
	if len(s.PlaceholderOverrides) > 0 {
		slide.PlaceholderOverrides = append(slide.PlaceholderOverrides, s.PlaceholderOverrides...)
	}
	if len(s.Animations) > 0 {
		slide.Animations = append(slide.Animations, s.Animations...)
	}
	if s.Table != nil {
		tmp := *s.Table
		slide.Table = &tmp
	}
	if s.BarChart != nil {
		tmp := *s.BarChart
		slide.Chart = &tmp
	}
	if len(s.SmartArt) > 0 {
		slide.SmartArtDiagrams = append(slide.SmartArtDiagrams, s.SmartArt...)
	}
	return slide
}

func placeholderTextType(index int) string {
	if index == 0 {
		return "title"
	}
	return "body"
}

func placeholderImageType(index int) string {
	if index == 0 {
		return "title"
	}
	return "pic"
}

func writeFileAtomically(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".gopptx-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	committed := false
	defer func() {
		_ = tmpFile.Close()
		if !committed {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err = tmpFile.Write(content); err != nil {
		return err
	}
	if err = tmpFile.Chmod(perm); err != nil {
		return err
	}
	if err = tmpFile.Sync(); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	if _, err = os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return err
	}
	committed = true
	return nil
}
