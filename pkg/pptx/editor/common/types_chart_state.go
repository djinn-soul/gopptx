package editorcommon

// ChartState is a read snapshot for chart-level object model traversal.
type ChartState struct {
	ChartStyle *int                `json:"chart_style,omitempty"`
	CategoryAx ChartAxisState      `json:"category_axis"`
	ValueAx    ChartAxisState      `json:"value_axis"`
	Series     []ChartSeriesData   `json:"series,omitempty"`
	Scene3D    ChartScene3DState   `json:"scene3d"`
	DataLabels ChartDataLabelState `json:"data_labels"`
	// Trendlines and ErrorBars mirror the update payload shape, one entry per
	// c:trendline and c:errBars respectively.
	Trendlines []ChartTrendline `json:"trendlines,omitempty"`
	ErrorBars  []ChartErrorBars `json:"error_bars,omitempty"`
	DataPoints []ChartDataPoint `json:"data_points,omitempty"`
	DataTable  *ChartDataTable  `json:"data_table,omitempty"`
	// PlotAreaLine is the outline in the plot area's direct c:spPr.
	PlotAreaLine *ChartLineFormat `json:"plot_area_line,omitempty"`
	// DataLabelPoints is every c:dLbl that carries its own number format, font
	// or display flags.
	DataLabelPoints []ChartDataLabelPoint `json:"data_label_points,omitempty"`
	// SeriesFormats is the fill, line and marker of each series that carries
	// any of them; series with no explicit formatting are omitted.
	SeriesFormats []ChartSeriesFormat `json:"series_formats,omitempty"`
	// SeriesLines is set when the first plot draws c:serLines.
	SeriesLines *ChartSeriesLines `json:"series_lines,omitempty"`
}

// ChartDataLabelState is the persisted data-label state for the first chart plot.
type ChartDataLabelState struct {
	Present        bool   `json:"present"`
	Position       string `json:"position,omitempty"`
	ShowValue      bool   `json:"show_value,omitempty"`
	ShowCategory   bool   `json:"show_category,omitempty"`
	ShowSeriesName bool   `json:"show_series_name,omitempty"`
	NumberFormat   string `json:"number_format,omitempty"`
	FormatLinked   *bool  `json:"format_linked,omitempty"`
	WordWrap       *bool  `json:"word_wrap,omitempty"`
}

// ChartScene3DState is a read snapshot for chart-level 3D scene settings.
type ChartScene3DState struct {
	CameraPreset       string `json:"camera_preset,omitempty"`
	CameraFieldOfView  int    `json:"camera_field_of_view,omitempty"`
	LightRig           string `json:"light_rig,omitempty"`
	LightDirection     string `json:"light_direction,omitempty"`
	LightRigRevolution bool   `json:"light_rig_revolution,omitempty"`
}
