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
	OpUpdateChartFormatting   = "update_chart_formatting"
	OpGetChartState           = "get_chart_state"
	OpListSlideCharts         = "list_slide_charts"
	OpGetSlideLayoutRef       = "get_slide_layout_ref"
	OpListSlideLayouts        = "list_slide_layouts"
	OpListSlideMasters        = "list_slide_masters"
	OpListMasterLayouts       = "list_master_layouts"
	OpRebindSlideLayout       = "rebind_slide_layout"
	OpCloneLayoutMasterFamily = "clone_layout_master_family"
	OpAddSlideMaster          = "add_slide_master"
	OpRemoveSlideMaster       = "remove_slide_master"
	OpAddSlideLayout          = "add_slide_layout"
	OpRemoveSlideLayout       = "remove_slide_layout"
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
	OpGetShapeTextState       = "get_shape_text_state"
	OpGetShapeRuns            = "get_shape_runs"
	OpSetShapeRuns            = "set_shape_runs"
	OpUpdateShapeRunText      = "update_shape_run_text"
	OpAppendShapeRun          = "append_shape_run"
	OpAddShape                = "add_shape"
	OpAddTextbox              = "add_textbox"
	OpAddConnector            = "add_connector"
	OpAddGroupShape           = "add_group_shape"
	OpGroupShapes             = "group_shapes"
	OpUngroupShapes           = "ungroup_shapes"
	OpBuildFreeform           = "build_freeform"
	OpAddImage                = "add_image"
	OpRemoveShape             = "remove_shape"
	OpUpdateShape             = "update_shape"
	OpMoveShapeToFront        = "move_shape_to_front"
	OpMoveShapeToBack         = "move_shape_to_back"
	OpGetNotes                = "get_notes"
	OpNotesSlideExists        = "notes_slide_exists"
	OpSetNotes                = "set_notes"
	OpSetModifyPassword       = "set_modify_password"
	OpSetMarkAsFinal          = "set_mark_as_final"
	OpAddTable                = "add_table"
	OpGetTable                = "get_table"
	OpMergeTableCells         = "merge_table_cells"
	OpSplitTableCell          = "split_table_cell"
	OpUpdateTableFlags        = "update_table_flags"
	OpUpdateTableCell         = "update_table_cell"
	OpSetTableStyle           = "set_table_style"
	OpSetTableRowHeight       = "set_table_row_height"
	OpSetTableColumnWidth     = "set_table_column_width"
	OpAddCustomXML            = "add_custom_xml"
	OpListCustomXML           = "list_custom_xml"
	OpRemoveCustomXML         = "remove_custom_xml"
	OpAddVba                  = "add_vba"
	OpListPlaceholders        = "list_placeholders"
	OpSetPlaceholderContent   = "set_placeholder_content"
	OpGetImageMetadata        = "get_image_metadata"
	OpAddVideo                = "add_video"
	OpAddAudio                = "add_audio"
	OpAddOLEObject            = "add_ole_object"
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
		OpUpdateChartFormatting,
		OpGetChartState,
		OpListSlideCharts,
		OpGetSlideLayoutRef,
		OpListSlideLayouts,
		OpListSlideMasters,
		OpListMasterLayouts,
		OpRebindSlideLayout,
		OpCloneLayoutMasterFamily,
		OpAddSlideMaster,
		OpRemoveSlideMaster,
		OpAddSlideLayout,
		OpRemoveSlideLayout,
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
		OpGetShapeTextState,
		OpGetShapeRuns,
		OpSetShapeRuns,
		OpUpdateShapeRunText,
		OpAppendShapeRun,
		OpAddShape,
		OpAddTextbox,
		OpAddConnector,
		OpAddGroupShape,
		OpGroupShapes,
		OpUngroupShapes,
		OpBuildFreeform,
		OpAddImage,
		OpRemoveShape,
		OpUpdateShape,
		OpMoveShapeToFront,
		OpMoveShapeToBack,
		OpGetNotes,
		OpNotesSlideExists,
		OpSetNotes,
		OpSetModifyPassword,
		OpSetMarkAsFinal,
		OpAddTable,
		OpGetTable,
		OpMergeTableCells,
		OpSplitTableCell,
		OpUpdateTableFlags,
		OpUpdateTableCell,
		OpSetTableStyle,
		OpSetTableRowHeight,
		OpSetTableColumnWidth,
		OpAddCustomXML,
		OpListCustomXML,
		OpRemoveCustomXML,
		OpAddVba,
		OpListPlaceholders,
		OpSetPlaceholderContent,
		OpGetImageMetadata,
		OpAddVideo,
		OpAddAudio,
		OpAddOLEObject,
	}
}
