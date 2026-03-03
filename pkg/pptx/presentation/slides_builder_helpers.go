package presentation

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func (b *slidePartBuilder) allocatePlaceholderChartRels(overrides []shapes.PlaceholderContent) []pptxxml.ChartRel {
	rels := make([]pptxxml.ChartRel, 0)
	for _, override := range overrides {
		if override.Chart == nil {
			continue
		}
		rels = append(rels, pptxxml.ChartRel{RID: b.nextRID()})
	}
	return rels
}

func (b *slidePartBuilder) buildTitleSpec(slide elements.SlideContent) pptxxml.TitleSpec {
	return pptxxml.TitleSpec{
		Text:      slide.Title,
		SizePt:    slide.TitleSize,
		Color:     slide.TitleColor,
		Bold:      slide.TitleBold,
		Italic:    slide.TitleItalic,
		Underline: slide.TitleUnderline,
		Align:     slide.TitleAlign,
		Font:      slide.TitleFont,
	}
}

func (b *slidePartBuilder) buildContentStyleSpec(slide elements.SlideContent) pptxxml.ContentStyleSpec {
	return pptxxml.ContentStyleSpec{
		SizePt:    slide.ContentSize,
		Color:     slide.ContentColor,
		Bold:      slide.ContentBold,
		Italic:    slide.ContentItalic,
		Underline: slide.ContentUnderline,
		VAlign:    slide.ContentVAlign,
	}
}

func (b *slidePartBuilder) mapBackground(bg *elements.SlideBackground) string {
	if bg == nil {
		return ""
	}
	if bg.Type == elements.SlideBackgroundPicture && bg.PictureFill != nil {
		mediaName, ok := b.catalog.MediaNameForImage(*bg.PictureFill)
		if ok {
			rid := b.nextRID()
			b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
			return rid
		}
	}
	return ""
}

func (b *slidePartBuilder) mapTransition(slide elements.SlideContent) string {
	if slide.Transition == nil {
		return ""
	}

	opt, ok := slide.Transition.(transitions.TransitionOptions)
	if !ok || opt.Sound == nil {
		return elements.SlideTransitionXML(slide)
	}

	path := opt.Sound.RelID
	if strings.HasPrefix(path, "file:") {
		path = strings.TrimPrefix(path, "file:")
	}

	mediaName, ok := b.catalog.MediaNameForImage(shapes.Image{Path: path})
	if !ok {
		return elements.SlideTransitionXML(slide)
	}

	rid := b.nextRID()
	b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))

	newOpt := opt
	newSound := *opt.Sound
	newSound.RelID = rid
	newOpt.Sound = &newSound

	return newOpt.XML()
}
