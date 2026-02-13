package editor

// Command operation names shared between bridge clients and Go dispatcher.
const (
	OpSlideCount     = "slide_count"
	OpAddSlide       = "add_slide"
	OpRemoveSlide    = "remove_slide"
	OpMoveSlide      = "move_slide"
	OpDuplicateSlide = "duplicate_slide"
	OpGetMetadata    = "get_metadata"
	OpUpdateChartData = "update_chart_data"
	OpListSlideCharts = "list_slide_charts"
	OpListSlideLayouts = "list_slide_layouts"
	OpRebindSlideLayout = "rebind_slide_layout"
	OpCloneLayoutMasterFamily = "clone_layout_master_family"
)

// SupportedOps is the canonical list of operations accepted by ExecuteCommand.
var SupportedOps = []string{
	OpSlideCount,
	OpAddSlide,
	OpRemoveSlide,
	OpMoveSlide,
	OpDuplicateSlide,
	OpGetMetadata,
	OpUpdateChartData,
	OpListSlideCharts,
	OpListSlideLayouts,
	OpRebindSlideLayout,
	OpCloneLayoutMasterFamily,
}
