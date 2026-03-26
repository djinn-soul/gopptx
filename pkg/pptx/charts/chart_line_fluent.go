package charts

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

// WithAltText sets the alternative text for accessibility.
func (c LineChart) WithAltText(text string) LineChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c LineChart) WithDecorative(enabled bool) LineChart {
	c.IsDecorative = enabled
	return c
}

// Position sets chart position in EMU.
func (c LineChart) Position(x styling.Length, y styling.Length) LineChart {
	return c.withBounds(x, y, c.CX, c.CY)
}

// Size sets chart size in EMU.
func (c LineChart) Size(cx styling.Length, cy styling.Length) LineChart {
	return c.withBounds(c.X, c.Y, cx, cy)
}

func (c LineChart) withBounds(
	x styling.Length,
	y styling.Length,
	cx styling.Length,
	cy styling.Length,
) LineChart {
	c.X = x
	c.Y = y
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c LineChart) WithTitle(title string) LineChart {
	c.Title = title
	return c
}

// WithLineColor sets the line color using RGB hex.
func (c LineChart) WithLineColor(color string) LineChart {
	c.LineColor = NormalizeHexColor(color)
	return c
}
