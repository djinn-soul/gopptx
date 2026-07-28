package editorcommon

// Per-series chart annotations: trendlines, error bars, per-point formatting,
// and negative-value inversion. Each is addressed by zero-based series index.

// ChartDataTable is the c:dTable grid drawn under the plot area.
//
// Show false removes it; the flags default to on when the chart has no data
// table yet, and otherwise keep whatever the chart already had.
type ChartDataTable struct {
	Show                 bool  `json:"show"`
	ShowHorizontalBorder *bool `json:"show_horizontal_border,omitempty"`
	ShowVerticalBorder   *bool `json:"show_vertical_border,omitempty"`
	ShowOutline          *bool `json:"show_outline,omitempty"`
	ShowKeys             *bool `json:"show_keys,omitempty"`
	FontSizePt           *int  `json:"font_size_pt,omitempty"`
}

// ChartTrendline describes one c:trendline fitted to a chart series.
//
// Order applies only to the poly type and Period only to movingAvg; PowerPoint
// repairs a file that carries either on any other trendline type.
type ChartTrendline struct {
	// SeriesIndex is the zero-based c:ser index; defaults to the first series.
	SeriesIndex int `json:"series_index,omitempty"`
	// Type is one of linear, poly, exp, log, movingAvg, power.
	Type            string   `json:"type"`
	Order           *int     `json:"order,omitempty"`
	Period          *int     `json:"period,omitempty"`
	Name            *string  `json:"name,omitempty"`
	Forward         *float64 `json:"forward,omitempty"`
	Backward        *float64 `json:"backward,omitempty"`
	Intercept       *float64 `json:"intercept,omitempty"`
	DisplayRSquared *bool    `json:"display_r_squared,omitempty"`
	DisplayEquation *bool    `json:"display_equation,omitempty"`
	LineColor       *string  `json:"line_color,omitempty"`
	LineWidthEMU    *int     `json:"line_width_emu,omitempty"`
	LineDash        *string  `json:"line_dash,omitempty"`
}

// ChartErrorBars describes one c:errBars attached to a chart series.
//
// Custom error bars carry PlusReference/MinusReference formulas; every other
// ValueType carries the scalar Value instead. Scatter and bubble series take
// one entry per Direction, so X and Y bars are separate c:errBars elements.
type ChartErrorBars struct {
	// SeriesIndex is the zero-based c:ser index; defaults to the first series.
	SeriesIndex int `json:"series_index,omitempty"`
	// BarType is one of both, minus, plus.
	BarType string `json:"bar_type"`
	// ValueType is one of cust, fixedVal, percentage, stdDev, stdErr.
	ValueType string `json:"value_type"`
	// Direction is x or y; omit it on charts with a single error-bar axis.
	Direction      *string  `json:"direction,omitempty"`
	Value          *float64 `json:"value,omitempty"`
	PlusReference  *string  `json:"plus_reference,omitempty"`
	MinusReference *string  `json:"minus_reference,omitempty"`
	NoEndCap       *bool    `json:"no_end_cap,omitempty"`
	LineColor      *string  `json:"line_color,omitempty"`
}

// ChartDataPoint is per-point formatting for one c:dPt in a series.
//
// Explosion applies to pie and doughnut series only; Bubble3D to bubble series.
type ChartDataPoint struct {
	// SeriesIndex is the zero-based c:ser index; defaults to the first series.
	SeriesIndex int `json:"series_index,omitempty"`
	// PointIndex is the zero-based data point this formatting applies to.
	PointIndex       int     `json:"point_index"`
	FillColor        *string `json:"fill_color,omitempty"`
	LineColor        *string `json:"line_color,omitempty"`
	LineWidthEMU     *int    `json:"line_width_emu,omitempty"`
	InvertIfNegative *bool   `json:"invert_if_negative,omitempty"`
	Bubble3D         *bool   `json:"bubble_3d,omitempty"`
	Explosion        *int    `json:"explosion,omitempty"`
	// Marker fields recolour the point's marker on a scatter, line or radar
	// series. The c:spPr above formats the segment leading to the point on
	// those charts, not the marker itself, so a scatter point needs these.
	MarkerFillColor *string `json:"marker_fill_color,omitempty"`
	MarkerLineColor *string `json:"marker_line_color,omitempty"`
	// MarkerSymbol is one of circle, dash, diamond, dot, none, plus, square,
	// star, triangle, x, auto.
	MarkerSymbol *string `json:"marker_symbol,omitempty"`
	// MarkerSize is the marker width in points, 2 to 72.
	MarkerSize *int `json:"marker_size,omitempty"`
}

// ChartDataLabelPoint formats the label of a single data point, the c:dLbl.
//
// PowerPoint keeps a point's number format and label font on the label rather
// than on the point, so this is separate from ChartDataPoint. A c:dLbl inherits
// nothing from the surrounding c:dLbls: whatever display flag it omits is off,
// so the flags of the label being replaced — else the series, else the plot —
// are carried over, and the Show fields override them.
type ChartDataLabelPoint struct {
	// SeriesIndex is the zero-based c:ser index; defaults to the first series.
	SeriesIndex int `json:"series_index,omitempty"`
	// PointIndex is the zero-based data point whose label this is.
	PointIndex int `json:"point_index"`
	// NumberFormat is a c:numFmt format code, such as "0.0%".
	NumberFormat *string `json:"number_format,omitempty"`
	// FormatLinked mirrors c:numFmt/@sourceLinked; it defaults to false when
	// NumberFormat is set, because a linked label ignores its own code.
	FormatLinked   *bool   `json:"format_linked,omitempty"`
	FontColor      *string `json:"font_color,omitempty"`
	FontSizePt     *int    `json:"font_size_pt,omitempty"`
	FontBold       *bool   `json:"font_bold,omitempty"`
	ShowValue      *bool   `json:"show_value,omitempty"`
	ShowCategory   *bool   `json:"show_category,omitempty"`
	ShowSeriesName *bool   `json:"show_series_name,omitempty"`
	ShowPercent    *bool   `json:"show_percent,omitempty"`
	ShowLegendKey  *bool   `json:"show_legend_key,omitempty"`
	// Delete removes the label for this point, which is how PowerPoint hides
	// one label of a series that otherwise shows them.
	Delete *bool `json:"delete,omitempty"`
}

// ChartSeriesInvert toggles c:invertIfNegative on one series.
//
// The flag alone leaves PowerPoint drawing negative points in the inverse of
// the series fill, which is usually white; NegativeFillColor writes an explicit
// c:dPt fill for each negative point so they stay readable.
type ChartSeriesInvert struct {
	// SeriesIndex is the zero-based c:ser index; defaults to the first series.
	SeriesIndex       int     `json:"series_index,omitempty"`
	InvertIfNegative  bool    `json:"invert_if_negative"`
	NegativeFillColor *string `json:"negative_fill_color,omitempty"`
}
