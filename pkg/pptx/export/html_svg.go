//nolint:mnd // SVG geometry uses fixed ratios/constants for deterministic PPT-like rendering.
package export

import (
	"fmt"
	"html"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// PPTX slides default to 9144000x5143500 EMU (10x5.625 inches).
// We map 9144000 EMU to 1000 pixels for a 16:9 1000x562 internal coordinate space.
const (
	emuPerInch = 914400.0
	pxPerInch  = 100.0 // arbitrary scaling for internal SVG coordinates
)

func emuToPx(emu int64) float64 {
	return (float64(emu) / emuPerInch) * pxPerInch
}

func renderShapesSVG(shapeList []shapes.Shape) string {
	if len(shapeList) == 0 {
		return ""
	}

	var sb strings.Builder
	// ViewBox matches 10x5.625 inches at 100px/inch -> 1000x562.5
	// We'll use 1000x562 rounded, but precise floats in viewBox.
	sb.WriteString(
		`<svg class="slide-svg" viewBox="0 0 1000 562.5" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">` + "\n",
	)

	// Collect unique gradient defs to avoid output duplication
	gradientDefs := make(map[string]string)
	for i, shape := range shapeList {
		if shape.GradientFill != nil {
			defID := fmt.Sprintf("grad-%d", i)
			gradientDefs[defID] = buildGradientDef(defID, *shape.GradientFill)
		}
	}

	if len(gradientDefs) > 0 {
		sb.WriteString("<defs>\n")
		for _, defStr := range gradientDefs {
			sb.WriteString(defStr)
		}
		sb.WriteString("</defs>\n")
	}

	// Render shapes
	for i, shape := range shapeList {
		sb.WriteString(renderShape(shape, i))
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

func ensureHash(color string) string {
	if color == "" {
		return ""
	}
	if !strings.HasPrefix(color, "#") {
		return "#" + color
	}
	return color
}

func buildGradientDef(id string, grad shapes.ShapeGradientFill) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, `<linearGradient id="%s" x1="0%%" y1="0%%" x2="100%%" y2="0%%">`+"\n", id)
	for _, stop := range grad.Stops {
		opacityStr := "1.0"
		if stop.Transparency != nil {
			opacityStr = fmt.Sprintf("%.2f", 1.0-*stop.Transparency)
		}
		color := ensureHash(stop.Color)
		fmt.Fprintf(&sb, `  <stop offset="%d%%" stop-color="%s" stop-opacity="%s" />`+"\n",
			stop.PositionPct,
			color,
			opacityStr)
	}
	sb.WriteString(`</linearGradient>` + "\n")
	return sb.String()
}

//nolint:funlen // Shape SVG emission keeps geometry cases in one switch for readability of output mapping.
func renderShape(s shapes.Shape, index int) string {
	x := emuToPx(int64(s.X))
	y := emuToPx(int64(s.Y))
	cx := emuToPx(int64(s.CX))
	cy := emuToPx(int64(s.CY))

	// Fill extraction
	fill := `fill="rgba(0,0,0,0)"` // default transparent
	if s.Fill != nil {
		opacity := 1.0
		if s.Fill.Transparency != nil {
			opacity = 1.0 - *s.Fill.Transparency
		}
		color := ensureHash(s.Fill.Color)
		if opacity < 1.0 {
			fill = fmt.Sprintf(`fill="%s" fill-opacity="%.2f"`, color, opacity)
		} else {
			fill = fmt.Sprintf(`fill="%s"`, color)
		}
	} else if s.GradientFill != nil {
		fill = fmt.Sprintf(`fill="url(#grad-%d)"`, index)
	}

	// Line extraction
	stroke := ""
	if s.Line != nil && s.Line.Width > 0 {
		strokeWidth := emuToPx(int64(s.Line.Width))
		strokeColor := ensureHash(s.Line.Color)
		stroke = fmt.Sprintf(`stroke="%s" stroke-width="%.2f"`, strokeColor, strokeWidth)
		if s.Line.Dash != shapes.LineDashSolid {
			// simple dash array mapping
			stroke += ` stroke-dasharray="4 2"`
		}
	}

	// Transform extraction
	transform := ""
	if s.RotationDeg != nil && *s.RotationDeg != 0 {
		centerXPx := x + cx/2
		centerYPx := y + cy/2
		transform = fmt.Sprintf(` transform="rotate(%d, %.2f, %.2f)"`, *s.RotationDeg, centerXPx, centerYPx)
	}

	var sb strings.Builder

	// Primitives
	switch s.Type {
	case shapes.ShapeTypeRectangle:
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" %s %s%s/>`,
			x,
			y,
			cx,
			cy,
			fill,
			stroke,
			transform)
	case shapes.ShapeTypeRoundedRectangle:
		rx := math.Min(cx, cy) * 0.1 // rough approx
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" rx="%.2f" %s %s%s/>`,
			x,
			y,
			cx,
			cy,
			rx,
			fill,
			stroke,
			transform)
	case shapes.ShapeTypeEllipse:
		rx := cx / 2.0
		ry := cy / 2.0
		cxCenter := x + rx
		cyCenter := y + ry
		fmt.Fprintf(&sb, `<ellipse cx="%.2f" cy="%.2f" rx="%.2f" ry="%.2f" %s %s%s/>`,
			cxCenter,
			cyCenter,
			rx,
			ry,
			fill,
			stroke,
			transform)
	case shapes.ShapeTypeTriangle:
		points := fmt.Sprintf("%.2f,%.2f %.2f,%.2f %.2f,%.2f", x+(cx/2), y, x, y+cy, x+cx, y+cy)
		fmt.Fprintf(&sb, `<polygon points="%s" %s %s%s/>`, points, fill, stroke, transform)
	case shapes.ShapeTypeRightArrow:
		aw := cx * 0.5 // arrow head width
		bw := cx - aw  // body width
		hh := cy * 0.5 // half height
		bh := cy * 0.5 // body height
		points := fmt.Sprintf("%.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f",
			x, y+(cy-bh)/2,
			x+bw, y+(cy-bh)/2,
			x+bw, y,
			x+cx, y+hh,
			x+bw, y+cy,
			x+bw, y+cy-(cy-bh)/2,
			x, y+cy-(cy-bh)/2)
		fmt.Fprintf(&sb, `<polygon points="%s" %s %s%s/>`, points, fill, stroke, transform)
	case shapes.ShapeTypeLeftArrow:
		aw := cx * 0.5
		hh := cy * 0.5
		bh := cy * 0.5
		points := fmt.Sprintf("%.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f",
			x+aw, y+(cy-bh)/2,
			x+cx, y+(cy-bh)/2,
			x+cx, y+cy-(cy-bh)/2,
			x+aw, y+cy-(cy-bh)/2,
			x+aw, y+cy,
			x, y+hh,
			x+aw, y)
		fmt.Fprintf(&sb, `<polygon points="%s" %s %s%s/>`, points, fill, stroke, transform)
	default:
		// Fallback for exotic shapes
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" %s %s%s/>`,
			x,
			y,
			cx,
			cy,
			fill,
			stroke,
			transform)
	}
	sb.WriteString("\n")

	// Text over shape
	if s.Text != "" {
		textColor := "#000000" // Default fallback; actual shape text extraction could be richer
		fontSize := 14.0       // Arbitrary px for shape text fallback
		textX := x + cx/2
		textY := y + cy/2
		escaped := html.EscapeString(s.Text)
		fmt.Fprintf(
			&sb,
			`<text x="%.2f" y="%.2f" font-size="%.1fpx" fill="%s" text-anchor="middle" dominant-baseline="middle"%s>%s</text>`+"\n",
			textX,
			textY,
			fontSize,
			textColor,
			transform,
			escaped,
		)
	}

	return sb.String()
}
