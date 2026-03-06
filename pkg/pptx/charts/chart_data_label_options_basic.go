package charts

import "strings"

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c BarChart) WithDataLabelPosition(position string) BarChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes data-label content fields.
func (c BarChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) BarChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c LineChart) WithDataLabelPosition(position string) LineChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes data-label content fields.
func (c LineChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) LineChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c PieChart) WithDataLabelPosition(position string) PieChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes pie/doughnut data-label content fields.
func (c PieChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) PieChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c DoughnutChart) WithDataLabelPosition(position string) DoughnutChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes pie/doughnut data-label content fields.
func (c DoughnutChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) DoughnutChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}
