package elements

import (
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

// WithDefaultBulletStyle sets the base style for new bullets.
func (s SlideContent) WithDefaultBulletStyle(style ParagraphStyle) SlideContent {
	s.DefaultBulletStyle = style
	return s
}

// WithBulletStyle sets the bullet style for all bullets on this slide.
func (s SlideContent) WithBulletStyle(style ParagraphStyle) SlideContent {
	s.DefaultBulletStyle = style
	for i := range s.BulletStyles {
		s.BulletStyles[i] = style
	}
	return s
}

// WithTransition sets the transition for the slide.
func (s SlideContent) WithTransition(t transitions.SlideTransition) SlideContent {
	s.Transition = t
	return s
}

// WithTransitionOptions sets built-in transition options.
func (s SlideContent) WithTransitionOptions(opt transitions.TransitionOptions) SlideContent {
	s.Transition = opt
	return s
}

// WithMorphTransition sets a Morph transition using PowerPoint's default mode.
// We intentionally omit the option attribute to maximize compatibility.
func (s SlideContent) WithMorphTransition() SlideContent {
	s.Transition = transitions.TransitionOptions{
		Type: transitions.TransitionMorph,
	}
	return s
}

// WithMorphTransitionOptions sets a Morph transition with explicit options.
func (s SlideContent) WithMorphTransitionOptions(option transitions.MorphOption) SlideContent {
	s.Transition = transitions.TransitionOptions{
		Type:        transitions.TransitionMorph,
		MorphOption: option,
	}
	return s
}

// WithTransitionSound sets a sound file for the slide transition.
func (s SlideContent) WithTransitionSound(path string) SlideContent {
	// If transition is nil or not options, default to cut.
	opt, ok := s.Transition.(transitions.TransitionOptions)
	if !ok {
		opt = transitions.TransitionOptions{Type: transitions.TransitionCut}
	}
	if opt.Sound == nil {
		opt.Sound = &transitions.TransitionSound{}
	}
	// Store the path in RelID temporarily; it will be resolved to a relation ID
	// during package writing.
	opt.Sound.RelID = "file:" + path
	opt.Sound.Name = filepath.Base(path)
	s.Transition = opt
	return s
}

// WithBulletStyleName sets primary bullet style by name (e.g. BulletStyleNumber).
func (s SlideContent) WithBulletStyleName(styleName string) SlideContent {
	style := s.DefaultBulletStyle
	style.BulletStyle = NormalizeBulletStyle(styleName)
	return s.WithBulletStyle(style)
}

// WithLayout sets the slide layout (supports canonical and compatibility aliases).
func (s SlideContent) WithLayout(layout string) SlideContent {
	s.Layout = NormalizeSlideLayout(layout)
	return s
}

// WithBackgroundColor sets the slide background as RGB hex.
func (s SlideContent) WithBackgroundColor(color string) SlideContent {
	normalized := common.NormalizeHexColor(color)
	bg := NewSolidBackground(normalized)
	s.Background = &bg
	return s
}

// WithBackground sets a complex background for the slide.
func (s SlideContent) WithBackground(bg SlideBackground) SlideContent {
	s.Background = &bg
	return s
}

// WithGradientBackground sets a gradient background for the slide.
func (s SlideContent) WithGradientBackground(gradient shapes.ShapeGradientFill) SlideContent {
	bg := NewGradientBackground(gradient)
	return s.WithBackground(bg)
}

// WithPictureBackground sets a picture background for the slide using image data.
func (s SlideContent) WithPictureBackground(img shapes.Image) SlideContent {
	bg := NewPictureBackground(img)
	return s.WithBackground(bg)
}

// WithTitleSize sets the title font size in points.
func (s SlideContent) WithTitleSize(size int) SlideContent {
	s.TitleSize = size
	return s
}

// WithTitleColor sets the title color as RGB hex.
func (s SlideContent) WithTitleColor(color string) SlideContent {
	s.TitleColor = common.NormalizeHexColor(color)
	return s
}

// WithTitleBold sets whether the title is bold.
func (s SlideContent) WithTitleBold(bold bool) SlideContent {
	s.TitleBold = bold
	return s
}

// WithTitleItalic sets whether the title is italic.
func (s SlideContent) WithTitleItalic(italic bool) SlideContent {
	s.TitleItalic = italic
	return s
}

// WithTitleUnderline sets whether the title is underlined.
func (s SlideContent) WithTitleUnderline(underline bool) SlideContent {
	s.TitleUnderline = underline
	return s
}

// WithTitleAlign sets the horizontal alignment of the title (l|ctr|r|just).
func (s SlideContent) WithTitleAlign(align string) SlideContent {
	s.TitleAlign = strings.ToLower(strings.TrimSpace(align))
	return s
}

// WithTitleFont sets the typeface for the slide title (e.g., "Consolas").
func (s SlideContent) WithTitleFont(font string) SlideContent {
	s.TitleFont = strings.TrimSpace(font)
	return s
}

// WithContentSize sets the content font size in points.
func (s SlideContent) WithContentSize(size int) SlideContent {
	s.ContentSize = size
	return s
}

// WithContentColor sets the content color as RGB hex.
func (s SlideContent) WithContentColor(color string) SlideContent {
	s.ContentColor = common.NormalizeHexColor(color)
	return s
}

// WithContentBold sets whether the content is bold.
func (s SlideContent) WithContentBold(bold bool) SlideContent {
	s.ContentBold = bold
	return s
}

// WithContentItalic sets whether the content is italic.
func (s SlideContent) WithContentItalic(italic bool) SlideContent {
	s.ContentItalic = italic
	return s
}

// WithContentUnderline sets whether the content is underlined.
func (s SlideContent) WithContentUnderline(underline bool) SlideContent {
	s.ContentUnderline = underline
	return s
}

// WithContentVAlign sets the vertical alignment of the main content (t|ctr|b).
func (s SlideContent) WithContentVAlign(align string) SlideContent {
	s.ContentVAlign = strings.ToLower(strings.TrimSpace(align))
	return s
}

// WithTitleOnlyLayout sets the layout to title_only.
func (s SlideContent) WithTitleOnlyLayout() SlideContent {
	s.Layout = SlideLayoutTitleOnly
	return s
}

// WithBlankLayout sets the layout to blank.
func (s SlideContent) WithBlankLayout() SlideContent {
	s.Layout = SlideLayoutBlank
	return s
}

// WithCenteredTitleLayout sets the layout to centered_title.
func (s SlideContent) WithCenteredTitleLayout() SlideContent {
	s.Layout = SlideLayoutCenteredTitle
	return s
}

// WithTitleAndBigContentLayout sets the layout to title_and_big_content.
func (s SlideContent) WithTitleAndBigContentLayout() SlideContent {
	s.Layout = SlideLayoutTitleAndBigContent
	return s
}

// WithTwoColumnLayout sets the layout to two_column.
func (s SlideContent) WithTwoColumnLayout() SlideContent {
	return s.WithLayout(SlideLayoutTwoColumn)
}

// WithTitleAndContentLayout sets the layout to title_and_content.
func (s SlideContent) WithTitleAndContentLayout() SlideContent {
	return s.WithLayout(SlideLayoutTitleAndContent)
}
