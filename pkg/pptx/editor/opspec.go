package editor

//go:generate go run ../../../cmd/gen_ops opspec.go ../../../python/gopptx/ops.py ../../../python/gopptx/ops.pyi

// Command operation names shared between bridge clients and Go dispatcher.
const (
	OpBatchExecute            = "batch_execute"
	OpSlideCount              = "slide_count"
	OpAddSlide                = "add_slide"
	OpRemoveSlide             = "remove_slide"
	OpMoveSlide               = "move_slide"
	OpDuplicateSlide          = "duplicate_slide"
	OpGetMetadata             = "get_metadata"
	OpUpdateChartData         = "update_chart_data"
	OpListSlideCharts         = "list_slide_charts"
	OpListSlideLayouts        = "list_slide_layouts"
	OpListSlideMasters        = "list_slide_masters"
	OpListMasterLayouts       = "list_master_layouts"
	OpRebindSlideLayout       = "rebind_slide_layout"
	OpCloneLayoutMasterFamily = "clone_layout_master_family"
	OpAddSection              = "add_section"
	OpRemoveSection           = "remove_section"
	OpRenameSection           = "rename_section"
	OpGetSections             = "get_sections"
	OpGetCoreProperties       = "get_core_properties"
	OpSetCoreProperties       = "set_core_properties"
	OpApplyTheme              = "apply_theme"
	OpSetSlideSize            = "set_slide_size"
	OpSetSlideTitle           = "set_slide_title"
	OpMergeFromFile           = "merge_from_file"
	OpUpdateSlide             = "update_slide"
	OpAddChart                = "add_chart"
	OpListSlides              = "list_slides"
	OpFindAndReplace          = "find_and_replace"
	OpSearchShapes            = "search_shapes"
	OpGetAuthors              = "get_authors"
	OpAddAuthor               = "add_author"
	OpGetComments             = "get_comments"
	OpAddComment              = "add_comment"
	OpRemoveComment           = "remove_comment"
	OpListShapes              = "list_shapes"
	OpAddShape                = "add_shape"
	OpAddImage                = "add_image"
	OpRemoveShape             = "remove_shape"
	OpUpdateShape             = "update_shape"
	OpGetNotes                = "get_notes"
	OpSetNotes                = "set_notes"
	OpSetModifyPassword       = "set_modify_password"
	OpSetMarkAsFinal          = "set_mark_as_final"
)

// SupportedOps returns the canonical list of operations accepted by ExecuteCommand.
func SupportedOps() []string {
	return []string{
		OpBatchExecute,
		OpSlideCount,
		OpAddSlide,
		OpRemoveSlide,
		OpMoveSlide,
		OpDuplicateSlide,
		OpGetMetadata,
		OpUpdateChartData,
		OpListSlideCharts,
		OpListSlideLayouts,
		OpListSlideMasters,
		OpListMasterLayouts,
		OpRebindSlideLayout,
		OpCloneLayoutMasterFamily,
		OpAddSection,
		OpRemoveSection,
		OpRenameSection,
		OpGetSections,
		OpGetCoreProperties,
		OpSetCoreProperties,
		OpApplyTheme,
		OpSetSlideSize,
		OpSetSlideTitle,
		OpMergeFromFile,
		OpUpdateSlide,
		OpAddChart,
		OpListSlides,
		OpFindAndReplace,
		OpSearchShapes,
		OpGetAuthors,
		OpAddAuthor,
		OpGetComments,
		OpAddComment,
		OpRemoveComment,
		OpListShapes,
		OpAddShape,
		OpAddImage,
		OpRemoveShape,
		OpUpdateShape,
		OpGetNotes,
		OpSetNotes,
		OpSetModifyPassword,
		OpSetMarkAsFinal,
	}
}
