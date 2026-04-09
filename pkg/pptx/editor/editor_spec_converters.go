package editor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	alignCenter     = "center"
	alignCtr        = "ctr"
	alignJust       = "just"
	alignDist       = "dist"
	bulletStyleNone = "none"
	paragraphPtUnit = 100
	paragraphPct    = 1000
	defaultBodyPrIn = 457200
)

// editorParagraphToSpec converts editor paragraph pointers to XML rendering spec.
func editorParagraphToSpec(p *common.Paragraph) (pptxxml.BulletParagraphSpec, error) {
	if p == nil {
		return pptxxml.BulletParagraphSpec{}, nil
	}
	spec := pptxxml.BulletParagraphSpec{}
	if p.Indent != nil {
		spec.LeftIndent = int64(*p.Indent)
	}
	if p.Hanging != nil {
		if *p.Hanging < 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.hanging must be >= 0")
		}
		spec.HangingIndent = int64(-*p.Hanging)
	}
	if len(p.TabStops) > 0 {
		spec.TabStops = make([]int64, 0, len(p.TabStops))
		for _, pos := range p.TabStops {
			if pos < 0 {
				return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.tab_stops values must be >= 0")
			}
			spec.TabStops = append(spec.TabStops, int64(pos))
		}
	}
	if p.Alignment != nil {
		normalized, err := normalizeParagraphAlignment(*p.Alignment)
		if err != nil {
			return pptxxml.BulletParagraphSpec{}, err
		}
		spec.Align = normalized
	}
	if p.Level != nil {
		if *p.Level < 0 || *p.Level > 8 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.level must be between 0 and 8")
		}
		spec.Level = *p.Level
	}
	if p.LineSpacingPct != nil {
		if *p.LineSpacingPct < 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.line_spacing_pct must be >= 0")
		}
		if *p.LineSpacingPct%paragraphPct != 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.line_spacing_pct must be divisible by 1000")
		}
		spec.LineSpacingPct = *p.LineSpacingPct / paragraphPct
	}
	if p.LineSpacingPts != nil {
		if *p.LineSpacingPts < 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.line_spacing_pts must be >= 0")
		}
		if *p.LineSpacingPts%paragraphPtUnit != 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.line_spacing_pts must be divisible by 100")
		}
		spec.LineSpacingPts = *p.LineSpacingPts / paragraphPtUnit
	}
	if p.LineSpacingPct != nil && p.LineSpacingPts != nil {
		return pptxxml.BulletParagraphSpec{}, errors.New("paragraph line spacing cannot set both percent and points")
	}
	if p.SpaceBeforePts != nil {
		if *p.SpaceBeforePts < 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.spcBef must be >= 0")
		}
		if *p.SpaceBeforePts%paragraphPtUnit != 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.space_before_pts must be divisible by 100")
		}
		spec.SpaceBeforePt = *p.SpaceBeforePts / paragraphPtUnit
	}
	if p.SpaceAfterPts != nil {
		if *p.SpaceAfterPts < 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.spcAft must be >= 0")
		}
		if *p.SpaceAfterPts%paragraphPtUnit != 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.space_after_pts must be divisible by 100")
		}
		spec.SpaceAfterPt = *p.SpaceAfterPts / paragraphPtUnit
	}
	if p.BulletStyle != nil {
		rawStyle := strings.ToLower(strings.TrimSpace(*p.BulletStyle))
		switch rawStyle {
		case bulletStyleNone, "bullet", "custom", "number", "letter_lower", "letter_upper", "roman_lower", "roman_upper":
			spec.BulletStyle = rawStyle
		default:
			return pptxxml.BulletParagraphSpec{}, fmt.Errorf("unsupported paragraph.bullet_style %q", *p.BulletStyle)
		}
	}
	if p.BulletChar != nil {
		spec.BulletChar = *p.BulletChar
	}
	if spec.BulletStyle == "custom" && strings.TrimSpace(spec.BulletChar) == "" {
		return pptxxml.BulletParagraphSpec{}, errors.New(
			"paragraph.bullet_char is required when paragraph.bullet_style=custom",
		)
	}
	if p.BulletColor != nil {
		color, err := editorshape.NormalizeHexColor(*p.BulletColor)
		if err != nil {
			return pptxxml.BulletParagraphSpec{}, fmt.Errorf("paragraph.bullet_color: %w", err)
		}
		spec.BulletColor = color
	}
	if p.BulletSizePct != nil {
		if *p.BulletSizePct < 0 {
			return pptxxml.BulletParagraphSpec{}, errors.New("paragraph.bullet_size_pct must be >= 0")
		}
		spec.BulletSize = *p.BulletSizePct
	}
	return spec, nil
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

func editorTextFrameToSpec(tf *common.TextFrame) (*pptxxml.TextFrameSpec, error) {
	if tf == nil {
		return nil, nil
	}
	spec := &pptxxml.TextFrameSpec{
		MarginLeft:   defaultBodyPrIn,
		MarginRight:  defaultBodyPrIn,
		MarginTop:    defaultBodyPrIn,
		MarginBottom: defaultBodyPrIn,
		Wrap:         "square",
		Anchor:       "ctr",
		AutoFit:      "spAutoFit",
	}
	if tf.MarginTop != nil {
		spec.MarginTop = int64(*tf.MarginTop)
	}
	if tf.MarginBottom != nil {
		spec.MarginBottom = int64(*tf.MarginBottom)
	}
	if tf.MarginLeft != nil {
		spec.MarginLeft = int64(*tf.MarginLeft)
	}
	if tf.MarginRight != nil {
		spec.MarginRight = int64(*tf.MarginRight)
	}
	if tf.WordWrap != nil {
		if *tf.WordWrap {
			spec.Wrap = "square"
		} else {
			spec.Wrap = "none"
		}
	}
	if tf.VerticalAlign != nil && *tf.VerticalAlign != "" {
		anchor, err := normalizeTextFrameVerticalAlign(*tf.VerticalAlign)
		if err != nil {
			return nil, err
		}
		spec.Anchor = anchor
	}
	if tf.Orientation != nil && *tf.Orientation != "" {
		orientation, err := normalizeTextFrameOrientation(*tf.Orientation)
		if err != nil {
			return nil, err
		}
		spec.Orientation = orientation
	}
	if tf.Columns != nil && *tf.Columns < minTextFrameColumns {
		return nil, fmt.Errorf("text_frame.columns must be >= %d", minTextFrameColumns)
	}
	if tf.Columns != nil {
		spec.NumCol = *tf.Columns
	}
	if tf.Rotation != nil {
		rotation, err := normalizeTextFrameRotation(*tf.Rotation)
		if err != nil {
			return nil, err
		}
		spec.Rotation = &rotation
	}
	if tf.AutoFitType != nil {
		switch strings.ToLower(strings.TrimSpace(*tf.AutoFitType)) {
		case "normal":
			spec.AutoFit = "normAutoFit"
		case "shape":
			spec.AutoFit = "spAutoFit"
		case bulletStyleNone:
			spec.AutoFit = "none"
		default:
			spec.AutoFit = "spAutoFit"
		}
	} else if tf.AutoFit != nil {
		if *tf.AutoFit {
			spec.AutoFit = "spAutoFit"
		} else {
			spec.AutoFit = "none"
		}
	}
	return spec, nil
}
