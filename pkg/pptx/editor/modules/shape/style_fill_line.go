package shape

import (
	"errors"
	"fmt"
	"math"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	gradientPercentScale       = 100.0
	gradientPositionScaleStyle = 1000.0
	maxGradientPercent         = 100.0
)

func RenderShapeStyleXML(
	fill *common.ShapeFill,
	line *common.ShapeLine,
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) (string, error) {
	var style strings.Builder
	fillXML, err := RenderFillXML(fill)
	if err != nil {
		return "", err
	}
	style.WriteString(fillXML)
	lineXML, err := RenderLineXML(line)
	if err != nil {
		return "", err
	}
	style.WriteString(lineXML)
	effectsXML, err := RenderEffectsXML(shadow, glow, blur, softEdge, reflection)
	if err != nil {
		return "", err
	}
	style.WriteString(effectsXML)
	return style.String(), nil
}
func RenderLineXML(line *common.ShapeLine) (string, error) {
	if line == nil {
		return "", nil
	}
	lnAttrs, err := renderLineAttrs(line)
	if err != nil {
		return "", err
	}
	lineColor, err := renderLineColor(line)
	if err != nil {
		return "", err
	}
	lineDash, err := renderLineDash(line)
	if err != nil {
		return "", err
	}
	lineArrows, err := renderLineArrows(line)
	if err != nil {
		return "", err
	}
	if lineColor == "" && lineDash == "" && lineArrows == "" {
		return `<a:ln` + lnAttrs + `/>`, nil
	}
	return renderLineElement(lnAttrs, lineDash, lineColor, lineArrows), nil
}
func renderLineAttrs(line *common.ShapeLine) (string, error) {
	if line.WidthEmu == nil {
		return "", nil
	}
	if *line.WidthEmu <= 0 {
		return "", errors.New("line.width_emu must be > 0")
	}
	return fmt.Sprintf(` w="%d"`, *line.WidthEmu), nil
}
func renderLineColor(line *common.ShapeLine) (string, error) {
	if line.Color == nil {
		return "", nil
	}
	color, err := NormalizeHexColor(*line.Color)
	if err != nil {
		return "", fmt.Errorf("line.color: %w", err)
	}
	return color, nil
}
func renderLineDash(line *common.ShapeLine) (string, error) {
	if line.DashStyle == nil {
		return "", nil
	}
	dash, err := NormalizeLineDashStyle(*line.DashStyle)
	if err != nil {
		return "", fmt.Errorf("line.dash_style: %w", err)
	}
	return dash, nil
}
func renderLineArrows(line *common.ShapeLine) (string, error) {
	var b strings.Builder
	head, err := renderLineArrowEnd(
		line.StartArrow,
		line.StartArrowWidth,
		line.StartArrowLength,
		"headEnd",
		"line.start_arrow",
		"line.start_arrow_width",
		"line.start_arrow_length",
	)
	if err != nil {
		return "", err
	}
	tail, err := renderLineArrowEnd(
		line.EndArrow,
		line.EndArrowWidth,
		line.EndArrowLength,
		"tailEnd",
		"line.end_arrow",
		"line.end_arrow_width",
		"line.end_arrow_length",
	)
	if err != nil {
		return "", err
	}
	b.WriteString(head)
	b.WriteString(tail)
	return b.String(), nil
}
func renderLineArrowEnd(
	arrowRaw, widthRaw, lengthRaw *string,
	tagName, arrowField, widthField, lengthField string,
) (string, error) {
	if arrowRaw == nil && widthRaw == nil && lengthRaw == nil {
		return "", nil
	}
	if arrowRaw == nil {
		return "", fmt.Errorf("%s is required when %s or %s is set", arrowField, widthField, lengthField)
	}
	arrow, err := NormalizeArrowType(*arrowRaw)
	if err != nil {
		return "", fmt.Errorf("%s: %w", arrowField, err)
	}
	var attrs strings.Builder
	attrs.WriteString(` type="`)
	attrs.WriteString(arrow)
	attrs.WriteString(`"`)
	if widthRaw != nil {
		width, widthErr := NormalizeArrowSize(*widthRaw)
		if widthErr != nil {
			return "", fmt.Errorf("%s: %w", widthField, widthErr)
		}
		attrs.WriteString(` w="`)
		attrs.WriteString(width)
		attrs.WriteString(`"`)
	}
	if lengthRaw != nil {
		length, lengthErr := NormalizeArrowSize(*lengthRaw)
		if lengthErr != nil {
			return "", fmt.Errorf("%s: %w", lengthField, lengthErr)
		}
		attrs.WriteString(` len="`)
		attrs.WriteString(length)
		attrs.WriteString(`"`)
	}
	return `<a:` + tagName + attrs.String() + `/>`, nil
}
func renderLineElement(lnAttrs, lineDash, lineColor, lineArrows string) string {
	var style strings.Builder
	style.WriteString(`<a:ln`)
	style.WriteString(lnAttrs)
	style.WriteString(`>`)
	if lineDash != "" {
		style.WriteString(`<a:prstDash val="`)
		style.WriteString(lineDash)
		style.WriteString(`"/>`)
	}
	if lineColor != "" {
		style.WriteString(`<a:solidFill><a:srgbClr val="`)
		style.WriteString(lineColor)
		style.WriteString(`"/></a:solidFill>`)
	}
	style.WriteString(lineArrows)
	style.WriteString(`</a:ln>`)
	return style.String()
}
func RenderFillXML(fill *common.ShapeFill) (string, error) {
	if fill == nil {
		return "", nil
	}
	modeCount := 0
	if fill.Solid != nil {
		modeCount++
	}
	if fill.Gradient != nil {
		modeCount++
	}
	if fill.Pattern != nil {
		modeCount++
	}
	if fill.Background != nil {
		modeCount++
	}
	if modeCount > 1 {
		return "", errors.New("fill.solid, fill.gradient, fill.pattern, and fill.background are mutually exclusive")
	}
	if fill.Transparency != nil && fill.Solid == nil {
		return "", errors.New("fill.transparency requires fill.solid")
	}
	if fill.Solid != nil {
		color, err := NormalizeHexColor(*fill.Solid)
		if err != nil {
			return "", fmt.Errorf("fill.solid: %w", err)
		}
		colorXML, err := renderFillSolidColorXML(color, fill.Transparency)
		if err != nil {
			return "", err
		}
		return `<a:solidFill>` + colorXML + `</a:solidFill>`, nil
	}
	if fill.Background != nil {
		if !*fill.Background {
			return "", errors.New("fill.background must be true when provided")
		}
		return `<a:noFill/>`, nil
	}
	if fill.Gradient != nil {
		return renderGradientFillXML(fill.Gradient)
	}
	if fill.Pattern != nil {
		return renderPatternFillXML(fill.Pattern)
	}
	return "", nil
}
func renderFillSolidColorXML(color string, transparency *float64) (string, error) {
	if transparency == nil {
		return `<a:srgbClr val="` + color + `"/>`, nil
	}
	if *transparency < 0 || *transparency > 1 {
		return "", errors.New("fill.transparency must be between 0.0 and 1.0")
	}
	alpha := int(math.Round((1.0 - *transparency) * float64(ooxmlPercentScale)))
	return fmt.Sprintf(`<a:srgbClr val="%s"><a:alpha val="%d"/></a:srgbClr>`, color, alpha), nil
}
func renderGradientFillXML(gradient *common.GradientFill) (string, error) {
	if gradient == nil {
		return "", nil
	}
	stops := gradient.Stops
	if len(stops) == 0 {
		return "", errors.New("fill.gradient.stops must contain at least 1 stop")
	}
	var b strings.Builder
	b.WriteString(`<a:gradFill><a:gsLst>`)
	for i := range stops {
		stop := stops[i]
		color, err := NormalizeHexColor(stop.Color)
		if err != nil {
			return "", fmt.Errorf("fill.gradient.stops[%d].color: %w", i, err)
		}
		pos := 0.0
		if stop.PositionPct != nil {
			pos = *stop.PositionPct
		} else if len(stops) > 1 {
			pos = float64(i) * (gradientPercentScale / float64(len(stops)-1))
		}
		if pos < 0.0 || pos > maxGradientPercent {
			return "", fmt.Errorf("fill.gradient.stops[%d].position_pct must be between 0 and 100", i)
		}
		b.WriteString(
			fmt.Sprintf(
				`<a:gs pos="%d"><a:srgbClr val="%s"/></a:gs>`,
				int(math.Round(pos*gradientPositionScaleStyle)),
				color,
			),
		)
	}
	b.WriteString(`</a:gsLst>`)
	if gradient.AngleDeg != nil {
		rotation, err := normalizeRotation(*gradient.AngleDeg)
		if err != nil {
			return "", fmt.Errorf("fill.gradient.angle_deg: %w", err)
		}
		b.WriteString(fmt.Sprintf(`<a:lin ang="%d" scaled="1"/>`, rotation))
	}
	b.WriteString(`</a:gradFill>`)
	return b.String(), nil
}
