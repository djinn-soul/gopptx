package editorcommon

// ChartSelector identifies a slide chart by index and/or relationship ID.
type ChartSelector struct {
	Index *int   `json:"index,omitempty"`
	RelID string `json:"rel_id,omitempty"`
}

// ChartSeriesData carries one chart series worth of input data.
type ChartSeriesData struct {
	// Hidden keeps the series in the embedded workbook but out of the plot, so
	// the data behind the chart can be a superset of what it draws.
	Hidden     bool      `json:"hidden,omitempty"`
	Name       *string   `json:"name,omitempty"`
	Categories []string  `json:"categories,omitempty"`
	Values     []float64 `json:"values,omitempty"`
	XValues    []float64 `json:"x_values,omitempty"`
	YValues    []float64 `json:"y_values,omitempty"`
	Sizes      []float64 `json:"sizes,omitempty"`
}

// ChartDataUpdate is the complete chart update payload.
type ChartDataUpdate struct {
	Categories           []string          `json:"categories,omitempty"`
	MultiLevelCategories [][]string        `json:"multi_level_categories,omitempty"`
	Series               []ChartSeriesData `json:"series,omitempty"`
}

// ChartDataBatchItem describes one chart update in a batch payload.
type ChartDataBatchItem struct {
	ChartSelector ChartSelector   `json:"chart_selector"`
	Data          ChartDataUpdate `json:"data"`
}

// DataLabelOffset positions one data label relative to its default spot.
// X and Y are fractions of the chart area, as PowerPoint's own "drag the label"
// gesture records them.
type DataLabelOffset struct {
	// SeriesIndex is the zero-based c:ser index; defaults to the first series.
	SeriesIndex int `json:"series_index,omitempty"`
	// PointIndex is the zero-based data point the label belongs to.
	PointIndex int `json:"point_index"`
	// X and Y are offsets from the label's default position.
	X *float64 `json:"x,omitempty"`
	Y *float64 `json:"y,omitempty"`
}

// ChartFormatUpdate is a partial formatting patch for an existing chart part.
type ChartFormatUpdate struct {
	ShowTitle    *bool    `json:"show_title,omitempty"`
	Title        *string  `json:"title,omitempty"`
	TitleOverlay *bool    `json:"title_overlay,omitempty"`
	TitleX       *float64 `json:"title_x,omitempty"`
	TitleY       *float64 `json:"title_y,omitempty"`
	// DataLabelOffsets nudges individual data labels. Doughnut and pie charts
	// reject most c:dLblPos values, so a manual layout is the only way to move
	// one of their labels.
	DataLabelOffsets []DataLabelOffset `json:"data_label_offsets,omitempty"`
	// DataLabelPoints formats individual data labels: number format, font and
	// display flags on one c:dLbl.
	DataLabelPoints []ChartDataLabelPoint `json:"data_label_points,omitempty"`
	// Trendlines replaces every c:trendline on each series it addresses.
	Trendlines []ChartTrendline `json:"trendlines,omitempty"`
	// AppendTrendlines adds c:trendline elements without rebuilding existing
	// trendlines, preserving unsupported formatting and extension children.
	AppendTrendlines []ChartTrendline `json:"append_trendlines,omitempty"`
	// ClearTrendlineSeries removes every c:trendline from the listed series,
	// which an empty Trendlines list cannot express.
	ClearTrendlineSeries []int `json:"clear_trendline_series,omitempty"`
	// ErrorBars replaces every c:errBars on each series it addresses.
	ErrorBars []ChartErrorBars `json:"error_bars,omitempty"`
	// ClearErrorBarSeries removes every c:errBars from the listed series.
	ClearErrorBarSeries []int `json:"clear_error_bar_series,omitempty"`
	// DataPoints merges per-point formatting into the series it addresses.
	DataPoints []ChartDataPoint `json:"data_points,omitempty"`
	// ClearDataPointSeries drops every c:dPt from the listed series.
	ClearDataPointSeries []int `json:"clear_data_point_series,omitempty"`
	// SeriesInverts toggles c:invertIfNegative per series.
	SeriesInverts []ChartSeriesInvert `json:"series_invert_if_negative,omitempty"`
	// DataTable shows or hides the c:dTable grid under the plot area.
	DataTable *ChartDataTable `json:"data_table,omitempty"`
	// PlotAreaLine formats the outline in the plot area's own c:spPr.
	PlotAreaLine            *ChartLineFormat `json:"plot_area_line,omitempty"`
	PlotVisibleOnly         *bool            `json:"plot_visible_only,omitempty"`
	ShowLegend              *bool            `json:"show_legend,omitempty"`
	LegendPosition          *string          `json:"legend_position,omitempty"`
	LegendOverlay           *bool            `json:"legend_overlay,omitempty"`
	ShowDataLabels          *bool            `json:"show_data_labels,omitempty"`
	DataLabelPosition       *string          `json:"data_label_position,omitempty"`
	DataLabelShowLegendKey  *bool            `json:"data_label_show_legend_key,omitempty"`
	DataLabelShowValue      *bool            `json:"data_label_show_value,omitempty"`
	DataLabelShowCategory   *bool            `json:"data_label_show_category,omitempty"`
	DataLabelShowSeriesName *bool            `json:"data_label_show_series_name,omitempty"`
	DataLabelShowPercent    *bool            `json:"data_label_show_percent,omitempty"`
	DataLabelShowBubbleSize *bool            `json:"data_label_show_bubble_size,omitempty"`
	DataLabelNumberFormat   *string          `json:"data_label_number_format,omitempty"`
	DataLabelFormatLinked   *bool            `json:"data_label_format_linked,omitempty"`
	DataLabelWordWrap       *bool            `json:"data_label_word_wrap,omitempty"`
	// DataLabelFillColor and DataLabelBorder format the label box itself, the
	// c:spPr of the plot-wide c:dLbls (upstream #662, #716).
	DataLabelFillColor *string          `json:"data_label_fill_color,omitempty"`
	DataLabelNoFill    *bool            `json:"data_label_no_fill,omitempty"`
	DataLabelBorder    *ChartLineFormat `json:"data_label_border,omitempty"`
	// SeriesFormats sets the fill, line and marker of whole series, which is
	// what recolours the markers of a line or scatter series (upstream #872).
	SeriesFormats []ChartSeriesFormat `json:"series_formats,omitempty"`
	// SeriesLines draws the c:serLines connectors of a stacked bar or
	// bar-of-pie chart (upstream #846).
	SeriesLines              *ChartSeriesLines `json:"series_lines,omitempty"`
	ChartGrouping            *string           `json:"chart_grouping,omitempty"`
	GapWidth                 *int              `json:"gap_width,omitempty"`
	Overlap                  *int              `json:"overlap,omitempty"`
	CategoryAxisTickLabelPos *string           `json:"category_axis_tick_label_pos,omitempty"`
	ValueAxisTickLabelPos    *string           `json:"value_axis_tick_label_pos,omitempty"`
	CategoryAxisMajorGrid    *bool             `json:"category_axis_major_gridlines,omitempty"`
	ValueAxisMajorGrid       *bool             `json:"value_axis_major_gridlines,omitempty"`
	CategoryAxisMinorGrid    *bool             `json:"category_axis_minor_gridlines,omitempty"`
	ValueAxisMinorGrid       *bool             `json:"value_axis_minor_gridlines,omitempty"`
	// Gridline formats write the c:spPr of a c:majorGridlines or
	// c:minorGridlines. Setting one turns that gridline on, since an unstyled
	// absent gridline has nothing to format (upstream #984).
	CategoryAxisMajorGridFormat *ChartLineFormat `json:"category_axis_major_gridline_format,omitempty"`
	CategoryAxisMinorGridFormat *ChartLineFormat `json:"category_axis_minor_gridline_format,omitempty"`
	ValueAxisMajorGridFormat    *ChartLineFormat `json:"value_axis_major_gridline_format,omitempty"`
	ValueAxisMinorGridFormat    *ChartLineFormat `json:"value_axis_minor_gridline_format,omitempty"`
	CategoryAxisCrosses         *string          `json:"category_axis_crosses,omitempty"`
	ValueAxisCrosses            *string          `json:"value_axis_crosses,omitempty"`
	CategoryAxisHasTitle        *bool            `json:"category_axis_has_title,omitempty"`
	ValueAxisHasTitle           *bool            `json:"value_axis_has_title,omitempty"`
	CategoryAxisTitle           *string          `json:"category_axis_title,omitempty"`
	ValueAxisTitle              *string          `json:"value_axis_title,omitempty"`
	CategoryAxisMinimumScale    *float64         `json:"category_axis_minimum_scale,omitempty"`
	CategoryAxisMaximumScale    *float64         `json:"category_axis_maximum_scale,omitempty"`
	ValueAxisMinimumScale       *float64         `json:"value_axis_minimum_scale,omitempty"`
	ValueAxisMaximumScale       *float64         `json:"value_axis_maximum_scale,omitempty"`
	CategoryAxisMajorUnit       *float64         `json:"category_axis_major_unit,omitempty"`
	CategoryAxisMinorUnit       *float64         `json:"category_axis_minor_unit,omitempty"`
	ValueAxisMajorUnit          *float64         `json:"value_axis_major_unit,omitempty"`
	ValueAxisMinorUnit          *float64         `json:"value_axis_minor_unit,omitempty"`
	CategoryAxisNumberFormat    *string          `json:"category_axis_number_format,omitempty"`
	ValueAxisNumberFormat       *string          `json:"value_axis_number_format,omitempty"`
	// CategoryAxisTickMarkSkip and CategoryAxisLabelAlignment exist only on
	// CT_CatAx; CT_DateAx and CT_ValAx have no such children.
	CategoryAxisTickMarkSkip   *int    `json:"category_axis_tick_mark_skip,omitempty"`
	CategoryAxisLabelAlignment *string `json:"category_axis_label_alignment,omitempty"`
	// CategoryAxisVisible and ValueAxisVisible write <c:delete>: the axis
	// element stays so the series keep referring to it, and the flag decides
	// whether PowerPoint draws it (upstream #473, #852).
	CategoryAxisVisible *bool `json:"category_axis_visible,omitempty"`
	ValueAxisVisible    *bool `json:"value_axis_visible,omitempty"`
	// Tick label rotation in degrees, the "Custom angle" of an axis label
	// (upstream #329).
	CategoryAxisTickLabelRotation *float64 `json:"category_axis_tick_label_rotation,omitempty"`
	ValueAxisTickLabelRotation    *float64 `json:"value_axis_tick_label_rotation,omitempty"`
	// ValueAxisCrossBetween is "between" or "midCat": whether the category
	// axis crosses between tick marks or on them (upstream #349). It lives on
	// CT_ValAx even though it describes the category axis.
	ValueAxisCrossBetween    *string `json:"value_axis_cross_between,omitempty"`
	CategoryAxisFormatLinked *bool   `json:"category_axis_format_linked,omitempty"`
	ValueAxisFormatLinked    *bool   `json:"value_axis_format_linked,omitempty"`
	CameraPreset             *string `json:"camera_preset,omitempty"`
	CameraFieldOfView        *int    `json:"camera_field_of_view,omitempty"`
	LightRig                 *string `json:"light_rig,omitempty"`
	LightDirection           *string `json:"light_direction,omitempty"`
	LightRigRevolution       *bool   `json:"light_rig_revolution,omitempty"`
}

// ChartAxisState is a read snapshot for one chart axis.
type ChartAxisState struct {
	Present       bool     `json:"present"`
	TickLabelPos  string   `json:"tick_label_pos,omitempty"`
	MajorGridline bool     `json:"major_gridline,omitempty"`
	MinorGridline bool     `json:"minor_gridline,omitempty"`
	Crosses       string   `json:"crosses,omitempty"`
	HasTitle      bool     `json:"has_title,omitempty"`
	Title         string   `json:"title,omitempty"`
	MinimumScale  *float64 `json:"minimum_scale,omitempty"`
	MaximumScale  *float64 `json:"maximum_scale,omitempty"`
	MajorUnit     *float64 `json:"major_unit,omitempty"`
	MinorUnit     *float64 `json:"minor_unit,omitempty"`
	NumberFormat  string   `json:"number_format,omitempty"`
	FormatLinked  *bool    `json:"format_linked,omitempty"`
	// TickMarkSkip and LabelAlignment are category-axis-only children.
	TickMarkSkip   *int   `json:"tick_mark_skip,omitempty"`
	LabelAlignment string `json:"label_alignment,omitempty"`
	// Visible is false when the axis carries <c:delete val="1"/>.
	Visible *bool `json:"visible,omitempty"`
	// TickLabelRotation is the axis label angle in degrees.
	TickLabelRotation *float64 `json:"tick_label_rotation,omitempty"`
	// CrossBetween is "between" or "midCat", and is set only on a value axis.
	CrossBetween string `json:"cross_between,omitempty"`
	// Gridline formats mirror the c:spPr of the gridline elements, and are nil
	// when the gridline carries no explicit style.
	MajorGridlineFormat *ChartLineFormat `json:"major_gridline_format,omitempty"`
	MinorGridlineFormat *ChartLineFormat `json:"minor_gridline_format,omitempty"`
}

// SlideChartRef describes a chart relationship discovered on a slide.
type SlideChartRef struct {
	Index     int
	RelID     string
	ChartPart string
	ShapeID   int
}

// SlideLayoutInfo describes one available slide layout part.
type SlideLayoutInfo struct {
	Part       string
	Name       string
	MasterPart string
	// LayoutID is the p:sldLayoutId/@id the master lists this layout under, which
	// is what python-pptx's SlideMaster.get_layout(slide_layout_id) keys on. Zero
	// when the layout is not referenced from a master's sldLayoutIdLst.
	LayoutID     int
	Shapes       []string
	Placeholders []PlaceholderInfo
}

// PlaceholderInfo describes a placeholder in a layout or master.
type PlaceholderInfo struct {
	Type  string
	Index int
	Name  string
	X     float64
	Y     float64
	CX    float64
	CY    float64
}

// NotesShapeInfo describes one shape discovered on a notes slide.
type NotesShapeInfo struct {
	ID                int
	Name              string
	Type              string
	Text              string
	X                 float64
	Y                 float64
	CX                float64
	CY                float64
	PlaceholderIndex  int
	PlaceholderType   string
	SupportsTextFrame bool `json:"supports_text_frame"`
}

// SlideMasterInfo describes one available slide master part.
type SlideMasterInfo struct {
	Part         string
	Shapes       []string
	Placeholders []PlaceholderInfo
}

// SlideMasterCloneResult summarizes an in-package layout/master clone operation.
type SlideMasterCloneResult struct {
	MasterPart string
	ThemePart  string
	LayoutMap  map[string]string
}

// ChartDataSource says where a chart's numbers come from: an embedded workbook,
// an external (linked) one, or nothing at all. A caller needs this to tell a
// linked chart from an embedded one before choosing how to update it
// (upstream #115).
type ChartDataSource struct {
	ChartPart string `json:"chart_part"`
	// Kind is "embedded", "external" or "none".
	Kind  string `json:"kind"`
	RelID string `json:"rel_id,omitempty"`
	// Target is the raw relationship target; PartPath is that target resolved
	// to a package part, and is empty for an external link.
	Target   string `json:"target,omitempty"`
	PartPath string `json:"part_path,omitempty"`
	// AutoUpdate mirrors <c:autoUpdate>, which asks PowerPoint to refresh the
	// link when the deck opens.
	AutoUpdate *bool `json:"auto_update,omitempty"`
}
