package editor

// Command operation names shared between bridge clients and Go dispatcher.
const (
	OpSlideCount              = "slide_count"
	OpAddSlide                = "add_slide"
	OpRemoveSlide             = "remove_slide"
	OpMoveSlide               = "move_slide"
	OpDuplicateSlide          = "duplicate_slide"
	OpGetMetadata             = "get_metadata"
	OpUpdateChartData         = "update_chart_data"
	OpListSlideCharts         = "list_slide_charts"
	OpListSlideLayouts        = "list_slide_layouts"
	OpRebindSlideLayout       = "rebind_slide_layout"
	OpCloneLayoutMasterFamily = "clone_layout_master_family"
	OpAddSection              = "add_section"
	OpRemoveSection           = "remove_section"
	OpRenameSection           = "rename_section"
	OpGetCoreProperties       = "get_core_properties"
	OpSetCoreProperties       = "set_core_properties"
	OpApplyTheme              = "apply_theme"
	OpSetSlideSize            = "set_slide_size"
	OpListShapes              = "list_shapes"
	OpAddShape                = "add_shape"
	OpRemoveShape             = "remove_shape"
	OpUpdateShape             = "update_shape"
	OpGetNotes                = "get_notes"
	OpSetNotes                = "set_notes"
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
	OpAddSection,
	OpRemoveSection,
	OpRenameSection,
	OpGetCoreProperties,
	OpSetCoreProperties,
	OpApplyTheme,
	OpSetSlideSize,
	OpListShapes,
	OpAddShape,
	OpRemoveShape,
	OpUpdateShape,
	OpGetNotes,
	OpSetNotes,
}
