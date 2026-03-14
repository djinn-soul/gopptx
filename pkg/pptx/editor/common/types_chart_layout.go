package editorcommon

// ChartSelector identifies a slide chart by index and/or relationship ID.
type ChartSelector struct {
	Index *int   `json:"index,omitempty"`
	RelID string `json:"rel_id,omitempty"`
}

// ChartSeriesData carries one chart series worth of input data.
type ChartSeriesData struct {
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

// ChartFormatUpdate is a partial formatting patch for an existing chart part.
type ChartFormatUpdate struct {
	ShowTitle                *bool   `json:"show_title,omitempty"`
	Title                    *string `json:"title,omitempty"`
	TitleOverlay             *bool   `json:"title_overlay,omitempty"`
	PlotVisibleOnly          *bool   `json:"plot_visible_only,omitempty"`
	ShowLegend               *bool   `json:"show_legend,omitempty"`
	LegendPosition           *string `json:"legend_position,omitempty"`
	LegendOverlay            *bool   `json:"legend_overlay,omitempty"`
	ShowDataLabels           *bool   `json:"show_data_labels,omitempty"`
	DataLabelPosition        *string `json:"data_label_position,omitempty"`
	DataLabelShowLegendKey   *bool   `json:"data_label_show_legend_key,omitempty"`
	DataLabelShowValue       *bool   `json:"data_label_show_value,omitempty"`
	DataLabelShowCategory    *bool   `json:"data_label_show_category,omitempty"`
	DataLabelShowSeriesName  *bool   `json:"data_label_show_series_name,omitempty"`
	DataLabelShowPercent     *bool   `json:"data_label_show_percent,omitempty"`
	DataLabelShowBubbleSize  *bool   `json:"data_label_show_bubble_size,omitempty"`
	CategoryAxisTickLabelPos *string `json:"category_axis_tick_label_pos,omitempty"`
	ValueAxisTickLabelPos    *string `json:"value_axis_tick_label_pos,omitempty"`
	CategoryAxisMajorGrid    *bool   `json:"category_axis_major_gridlines,omitempty"`
	ValueAxisMajorGrid       *bool   `json:"value_axis_major_gridlines,omitempty"`
	CategoryAxisCrosses      *string `json:"category_axis_crosses,omitempty"`
	ValueAxisCrosses         *string `json:"value_axis_crosses,omitempty"`
	CameraPreset             *string `json:"camera_preset,omitempty"`
	CameraFieldOfView        *int    `json:"camera_field_of_view,omitempty"`
	LightRig                 *string `json:"light_rig,omitempty"`
	LightDirection           *string `json:"light_direction,omitempty"`
	LightRigRevolution       *bool   `json:"light_rig_revolution,omitempty"`
}

// ChartAxisState is a read snapshot for one chart axis.
type ChartAxisState struct {
	Present       bool   `json:"present"`
	TickLabelPos  string `json:"tick_label_pos,omitempty"`
	MajorGridline bool   `json:"major_gridline,omitempty"`
	Crosses       string `json:"crosses,omitempty"`
}

// ChartState is a read snapshot for chart-level object model traversal.
type ChartState struct {
	ChartStyle *int              `json:"chart_style,omitempty"`
	CategoryAx ChartAxisState    `json:"category_axis"`
	ValueAx    ChartAxisState    `json:"value_axis"`
	Series     []ChartSeriesData `json:"series,omitempty"`
	Scene3D    ChartScene3DState `json:"scene3d"`
}

// ChartScene3DState is a read snapshot for chart-level 3D scene settings.
type ChartScene3DState struct {
	CameraPreset       string `json:"camera_preset,omitempty"`
	CameraFieldOfView  int    `json:"camera_field_of_view,omitempty"`
	LightRig           string `json:"light_rig,omitempty"`
	LightDirection     string `json:"light_direction,omitempty"`
	LightRigRevolution bool   `json:"light_rig_revolution,omitempty"`
}

// SlideChartRef describes a chart relationship discovered on a slide.
type SlideChartRef struct {
	Index     int
	RelID     string
	ChartPart string
}

// SlideLayoutInfo describes one available slide layout part.
type SlideLayoutInfo struct {
	Part         string
	Name         string
	MasterPart   string
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
	ID               int
	Name             string
	Type             string
	Text             string
	X                float64
	Y                float64
	CX               float64
	CY               float64
	PlaceholderIndex int
	PlaceholderType  string
	HasTextFrame     bool
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
