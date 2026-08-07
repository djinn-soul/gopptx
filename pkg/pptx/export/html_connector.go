//nolint:mnd // Arrowhead geometry uses fixed proportions of the line width.
package export

import (
	"fmt"
	"html"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// Connectors were dropped from the HTML export entirely, so a diagram built out
// of boxes and arrows came out as boxes. They are drawn into the slide's own SVG
// alongside the shapes.

// renderConnectorsSVG draws a slide's connectors, or returns empty when it has
// none.
func renderConnectorsSVG(connectors []shapes.Connector) string {
	if len(connectors) == 0 {
		return ""
	}
	var body strings.Builder
	for _, connector := range connectors {
		body.WriteString(renderConnector(connector))
	}
	if body.Len() == 0 {
		return ""
	}
	return `<svg class="slide-svg" viewBox="0 0 1000 562.5" preserveAspectRatio="xMidYMid meet"` +
		` xmlns="http://www.w3.org/2000/svg">` + "\n" + body.String() + "</svg>\n"
}

func renderConnector(connector shapes.Connector) string {
	if connector.IsDecorative {
		return ""
	}
	x1 := emuToPx(int64(connector.StartX))
	y1 := emuToPx(int64(connector.StartY))
	x2 := emuToPx(int64(connector.EndX))
	y2 := emuToPx(int64(connector.EndY))

	color := ensureHashOrDefault(connector.Line.Color, "#000000")
	width := emuToPx(int64(connector.Line.Width))
	if width <= 0 {
		width = defaultConnectorWidthPx
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<path d="%s" fill="none" stroke="%s" stroke-width="%.2f"%s/>`+"\n",
		connectorPath(connector, x1, y1, x2, y2), color, width, connectorDashAttr(connector))

	// The heads are drawn as filled triangles rather than markers, so the SVG
	// needs no defs and stays self-contained.
	if connector.EndArrow != "" && connector.EndArrow != shapes.ArrowTypeNone {
		sb.WriteString(arrowHead(x2, y2, x1, y1, width, color))
	}
	if connector.StartArrow != "" && connector.StartArrow != shapes.ArrowTypeNone {
		sb.WriteString(arrowHead(x1, y1, x2, y2, width, color))
	}
	if label := strings.TrimSpace(connector.Label); label != "" {
		fmt.Fprintf(&sb,
			`<text x="%.2f" y="%.2f" font-size="11px" fill="%s" text-anchor="middle">%s</text>`+"\n",
			(x1+x2)/2, (y1+y2)/2-4, color, html.EscapeString(label))
	}
	return sb.String()
}

// connectorPath is the line itself: straight for a straight connector, and a
// right-angled dog-leg for an elbow one, which is how PowerPoint routes it.
func connectorPath(connector shapes.Connector, x1, y1, x2, y2 float64) string {
	if strings.Contains(strings.ToLower(connector.Type), "bent") ||
		strings.Contains(strings.ToLower(connector.Type), "elbow") {
		midX := (x1 + x2) / 2
		return fmt.Sprintf("M %.2f %.2f L %.2f %.2f L %.2f %.2f L %.2f %.2f",
			x1, y1, midX, y1, midX, y2, x2, y2)
	}
	if strings.Contains(strings.ToLower(connector.Type), "curve") {
		return fmt.Sprintf("M %.2f %.2f C %.2f %.2f %.2f %.2f %.2f %.2f",
			x1, y1, (x1+x2)/2, y1, (x1+x2)/2, y2, x2, y2)
	}
	return fmt.Sprintf("M %.2f %.2f L %.2f %.2f", x1, y1, x2, y2)
}

func connectorDashAttr(connector shapes.Connector) string {
	if connector.Line.Dash == "" || connector.Line.Dash == shapes.LineDashSolid {
		return ""
	}
	return ` stroke-dasharray="4 2"`
}

// arrowHead is a triangle at (tipX, tipY) pointing away from (fromX, fromY).
func arrowHead(tipX, tipY, fromX, fromY, width float64, color string) string {
	angle := math.Atan2(tipY-fromY, tipX-fromX)
	length := math.Max(arrowHeadMinPx, width*arrowHeadLengthFactor)
	spread := length * arrowHeadSpreadFactor

	baseX := tipX - length*math.Cos(angle)
	baseY := tipY - length*math.Sin(angle)
	perpX := -math.Sin(angle) * spread
	perpY := math.Cos(angle) * spread

	return fmt.Sprintf(`<polygon points="%.2f,%.2f %.2f,%.2f %.2f,%.2f" fill="%s"/>`+"\n",
		tipX, tipY,
		baseX+perpX, baseY+perpY,
		baseX-perpX, baseY-perpY,
		color)
}

const (
	// defaultConnectorWidthPx is the hairline a connector with no stated width
	// is drawn with, in the SVG's own coordinate space.
	defaultConnectorWidthPx = 1.0
	// arrowHeadMinPx keeps a head visible on a hairline connector.
	arrowHeadMinPx = 6.0
	// arrowHeadLengthFactor and arrowHeadSpreadFactor size the head against the
	// line width, as PowerPoint's medium arrowhead does.
	arrowHeadLengthFactor = 4.0
	arrowHeadSpreadFactor = 0.5
)
