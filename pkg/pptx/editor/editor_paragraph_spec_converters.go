package editor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// editorParagraphToSpec converts editor paragraph pointers to XML rendering spec.
func editorParagraphToSpec(p *common.Paragraph) (pptxxml.BulletParagraphSpec, error) {
	if p == nil {
		return pptxxml.BulletParagraphSpec{}, nil
	}
	spec := pptxxml.BulletParagraphSpec{}
	if err := applyParagraphIndents(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	if err := applyParagraphTabs(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	if err := applyParagraphAlignmentAndLevel(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	if err := applyParagraphLineSpacing(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	if err := applyParagraphSpacing(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	if err := applyParagraphBullet(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	if err := applyParagraphBulletFormatting(&spec, p); err != nil {
		return pptxxml.BulletParagraphSpec{}, err
	}
	return spec, nil
}

func applyParagraphIndents(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if p.Indent != nil {
		spec.LeftIndent = int64(*p.Indent)
	}
	if p.Hanging == nil {
		return nil
	}
	if *p.Hanging < 0 {
		return errors.New("paragraph.hanging must be >= 0")
	}
	spec.HangingIndent = int64(-*p.Hanging)
	return nil
}

func applyParagraphTabs(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if len(p.TabStops) == 0 {
		return nil
	}
	spec.TabStops = make([]int64, 0, len(p.TabStops))
	for _, pos := range p.TabStops {
		if pos < 0 {
			return errors.New("paragraph.tab_stops values must be >= 0")
		}
		spec.TabStops = append(spec.TabStops, int64(pos))
	}
	return nil
}

func applyParagraphAlignmentAndLevel(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if p.Alignment != nil {
		normalized, err := normalizeParagraphAlignment(*p.Alignment)
		if err != nil {
			return err
		}
		spec.Align = normalized
	}
	if p.Level == nil {
		return nil
	}
	if *p.Level < 0 || *p.Level > 8 {
		return errors.New("paragraph.level must be between 0 and 8")
	}
	spec.Level = *p.Level
	return nil
}

func applyParagraphLineSpacing(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if p.LineSpacingPct != nil && p.LineSpacingPts != nil {
		return errors.New("paragraph line spacing cannot set both percent and points")
	}
	if p.LineSpacingPct != nil {
		if *p.LineSpacingPct < 0 {
			return errors.New("paragraph.line_spacing_pct must be >= 0")
		}
		if *p.LineSpacingPct%paragraphPct != 0 {
			return errors.New("paragraph.line_spacing_pct must be divisible by 1000")
		}
		spec.LineSpacingPct = *p.LineSpacingPct / paragraphPct
	}
	if p.LineSpacingPts == nil {
		return nil
	}
	if *p.LineSpacingPts < 0 {
		return errors.New("paragraph.line_spacing_pts must be >= 0")
	}
	if *p.LineSpacingPts%paragraphPtUnit != 0 {
		return errors.New("paragraph.line_spacing_pts must be divisible by 100")
	}
	spec.LineSpacingPts = *p.LineSpacingPts / paragraphPtUnit
	return nil
}

func applyParagraphSpacing(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if value, ok, err := paragraphPointsValue(p.SpaceBeforePts, "paragraph.space_before_pts"); err != nil {
		return err
	} else if ok {
		spec.SpaceBeforePt = value
	}
	if value, ok, err := paragraphPointsValue(p.SpaceAfterPts, "paragraph.space_after_pts"); err != nil {
		return err
	} else if ok {
		spec.SpaceAfterPt = value
	}
	return nil
}

func paragraphPointsValue(raw *int, field string) (int, bool, error) {
	if raw == nil {
		return 0, false, nil
	}
	if *raw < 0 {
		return 0, false, fmt.Errorf("%s must be >= 0", field)
	}
	if *raw%paragraphPtUnit != 0 {
		return 0, false, fmt.Errorf("%s must be divisible by 100", field)
	}
	return *raw / paragraphPtUnit, true, nil
}

func applyParagraphBullet(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if p.BulletStyle != nil {
		style, err := normalizeParagraphBulletStyle(*p.BulletStyle)
		if err != nil {
			return err
		}
		spec.BulletStyle = style
	}
	if p.BulletChar != nil {
		spec.BulletChar = *p.BulletChar
	}
	if spec.BulletStyle == bulletStyleCustom && strings.TrimSpace(spec.BulletChar) == "" {
		return errors.New("paragraph.bullet_char is required when paragraph.bullet_style=custom")
	}
	return nil
}

func normalizeParagraphBulletStyle(raw string) (string, error) {
	switch style := strings.ToLower(strings.TrimSpace(raw)); style {
	case bulletStyleNone, bulletStyleBullet, bulletStyleCustom, bulletStyleNumber,
		bulletStyleLetterLower, bulletStyleLetterUpper, bulletStyleRomanLower, bulletStyleRomanUpper:
		return style, nil
	default:
		return "", fmt.Errorf("unsupported paragraph.bullet_style %q", raw)
	}
}

func applyParagraphBulletFormatting(spec *pptxxml.BulletParagraphSpec, p *common.Paragraph) error {
	if p.BulletColor != nil {
		color, err := editorshape.NormalizeHexColor(*p.BulletColor)
		if err != nil {
			return fmt.Errorf("paragraph.bullet_color: %w", err)
		}
		spec.BulletColor = color
	}
	if p.BulletSizePct == nil {
		return nil
	}
	if *p.BulletSizePct < 0 {
		return errors.New("paragraph.bullet_size_pct must be >= 0")
	}
	spec.BulletSize = *p.BulletSizePct
	return nil
}

func normalizeParagraphAlignment(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "l", "left":
		return "l", nil
	case alignCtr, alignCenter, "middle":
		return alignCtr, nil
	case "r", "right":
		return "r", nil
	case alignJust, "justify":
		return alignJust, nil
	case alignDist, "distribute":
		return alignDist, nil
	case "thaidist":
		return "thaiDist", nil
	case "justlow":
		return "justLow", nil
	default:
		return "", fmt.Errorf("unsupported paragraph.alignment %q", raw)
	}
}
