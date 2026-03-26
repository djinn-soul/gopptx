package charts

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

// WithAltText sets the alternative text for accessibility.
func (c BarChart) WithAltText(text string) BarChart {
	c.AltText = text
	return c
}

// WithDecorative marks the chart as decorative (ignored by screen readers).
func (c BarChart) WithDecorative(enabled bool) BarChart {
	c.IsDecorative = enabled
	return c
}

// Position sets chart position in EMU.
func (c BarChart) Position(x styling.Length, y styling.Length) BarChart {
	c.X = x
	c.Y = y
	return c
}

// Size sets chart size in EMU.
func (c BarChart) Size(cx styling.Length, cy styling.Length) BarChart {
	c.CX = cx
	c.CY = cy
	return c
}

// WithTitle sets the chart title.
func (c BarChart) WithTitle(title string) BarChart {
	c.Title = title
	return c
}

// WithBarColor sets the bar fill color using RGB hex.
func (c BarChart) WithBarColor(color string) BarChart {
	c.BarColor = NormalizeHexColor(color)
	return c
}
