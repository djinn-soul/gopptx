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
	OpUpdateChartDataBatch    = "update_chart_data_batch"
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
	OpGetSlideTextStates      = "get_slide_text_states"
	OpGetShapeTextState       = "get_shape_text_state"
	OpGetShapeRuns            = "get_shape_runs"
	OpSetShapeRuns            = "set_shape_runs"
	OpSetSlideShapeRuns       = "set_slide_shape_runs"
	OpUpdateDeckRunTexts      = "update_deck_run_texts"
	OpUpdateSlideRunTexts     = "update_slide_run_texts"
	OpUpdateShapeRunText      = "update_shape_run_text"
	OpAppendShapeRun          = "append_shape_run"
	OpAddShape                = "add_shape"
	OpAddTextbox              = "add_textbox"
	OpAddTextboxes            = "add_textboxes"
	OpAddConnectors           = "add_connectors"
	OpReserveShapeIDs         = "reserve_shape_ids"
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
	OpSetNotesShapeText       = "set_notes_shape_text"
	OpSetNotesShapeProps      = "set_notes_shape_props"
	OpSetModifyPassword       = "set_modify_password"
	OpSetMarkAsFinal          = "set_mark_as_final"
	OpAddTable                = "add_table"
	OpGetTable                = "get_table"
	OpMergeTableCells         = "merge_table_cells"
	OpSplitTableCell          = "split_table_cell"
	OpUpdateTableFlags        = "update_table_flags"
	OpUpdateTableCell         = "update_table_cell"
	OpSetTableStyle           = "set_table_style"
	OpDefineTableStyle        = "define_table_style"
	OpListTableStyles         = "list_table_styles"
	OpSetTableRowHeight       = "set_table_row_height"
	OpSetTableColumnWidth     = "set_table_column_width"
	OpAddCustomXML            = "add_custom_xml"
	OpListCustomXML           = "list_custom_xml"
	OpRemoveCustomXML         = "remove_custom_xml"
	OpAddVba                  = "add_vba"
	OpExportPDF               = "export_pdf"
	OpExportHTML              = "export_html"
	OpListPlaceholders        = "list_placeholders"
	OpSetPlaceholderContent   = "set_placeholder_content"
	OpGetImageMetadata        = "get_image_metadata"
	OpAddVideo                = "add_video"
	OpAddAudio                = "add_audio"
	OpAddOLEObject            = "add_ole_object"
	OpMarkdownToSlides        = "markdown_to_slides"
	OpURLFetchToSlides        = "url_fetch_to_slides"
	OpAddMermaidShape         = "add_mermaid_shape"
	OpAddSmartArt             = "add_smartart"
	OpAddAnimation            = "add_animation"
	OpSetSlideTransition      = "set_slide_transition"
	OpBuildStatusTemplate     = "build_status_template"
	OpBuildSimpleTemplate     = "build_simple_template"
	OpBuildProposalTemplate   = "build_proposal_template"
	OpBuildTrainingTemplate   = "build_training_template"
	OpBuildTechnicalTemplate  = "build_technical_template"

	// OpUpdateSmartArt, OpSetSlideBackground, OpSetSlideHeaderFooter, OpGetHandoutMaster,
	// OpUpdateHandoutMaster, and OpHasDigitalSignature extend parity with the Go API.
	OpUpdateSmartArt       = "update_smartart"
	OpSetSlideBackground   = "set_slide_background"
	OpSetSlideHeaderFooter = "set_slide_header_footer"
	OpGetHandoutMaster     = "get_handout_master"
	OpUpdateHandoutMaster  = "update_handout_master"
	OpHasDigitalSignature  = "has_digital_signature"

	// OpDuplicateSlideAfter and subsequent constants bridge PresentationEditor methods that had no JSON op.
	OpDuplicateSlideAfter   = "duplicate_slide_after"
	OpMoveShapeToIndex      = "move_shape_to_index"
	OpValidate              = "validate"
	OpRepair                = "repair"
	OpListSlideImages       = "list_slide_images"
	OpSwapImageByIndex      = "swap_image_by_index"
	OpSwapImageByRelID      = "swap_image_by_rel_id"
	OpGetLayoutShapes       = "get_layout_shapes"
	OpGetMasterShapes       = "get_master_shapes"
	OpGetLayoutPlaceholders = "get_layout_placeholders"
	OpGetMasterPlaceholders = "get_master_placeholders"
	OpSetGlobalThemePreset  = "set_global_theme_preset"
	OpSetThemeColorScheme   = "set_theme_color_scheme"
	OpSetThemeFontScheme    = "set_theme_font_scheme"
	OpGetThemeInventory     = "get_theme_inventory"
	OpListNotesShapes       = "list_notes_shapes"
	OpListNotesPlaceholders = "list_notes_placeholders"
	OpUpdateNotesMaster     = "update_notes_master"
	OpMergeFromEditor       = "merge_from_editor"
)

// SupportedOps returns the canonical list of operations accepted by ExecuteCommand.
// supportedSlideAndMetaOps returns ops for slide management, metadata, charts, and layout.
func supportedSlideAndMetaOps() []string {
	return []string{
		OpBatchExecute,
		OpSlideCount,
		OpAddSlide,
		OpRemoveSlide,
		OpMoveSlide,
		OpDuplicateSlide,
		OpGetMetadata,
		OpListSlides,
		OpUpdateSlide,
		OpSetSlideTitle,
		OpMergeFromFile,
		OpGetCoreProperties,
		OpSetCoreProperties,
		OpApplyTheme,
		OpSetSlideSize,
		OpAddSection,
		OpRemoveSection,
		OpRenameSection,
		OpGetSections,
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
		OpUpdateChartData,
		OpUpdateChartDataBatch,
		OpUpdateChartFormatting,
		OpGetChartState,
		OpListSlideCharts,
		OpAddChart,
		OpDuplicateSlideAfter,
		OpValidate,
		OpRepair,
		OpGetLayoutShapes,
		OpGetMasterShapes,
		OpGetLayoutPlaceholders,
		OpGetMasterPlaceholders,
		OpSetGlobalThemePreset,
		OpSetThemeColorScheme,
		OpSetThemeFontScheme,
		OpGetThemeInventory,
		OpBuildStatusTemplate,
		OpBuildSimpleTemplate,
		OpBuildProposalTemplate,
		OpBuildTrainingTemplate,
		OpBuildTechnicalTemplate,
	}
}

// supportedContentOps returns ops for shapes, text, tables, notes, media, and export.
func supportedContentOps() []string {
	return []string{
		OpFindAndReplace,
		OpSearchShapes,
		OpGetAuthors,
		OpAddAuthor,
		OpGetComments,
		OpAddComment,
		OpRemoveComment,
		OpListShapes,
		OpGetSlideTextStates,
		OpGetShapeTextState,
		OpGetShapeRuns,
		OpSetShapeRuns,
		OpSetSlideShapeRuns,
		OpUpdateDeckRunTexts,
		OpUpdateSlideRunTexts,
		OpUpdateShapeRunText,
		OpAppendShapeRun,
		OpAddShape,
		OpAddTextbox,
		OpAddTextboxes,
		OpAddConnectors,
		OpReserveShapeIDs,
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
		OpSetNotesShapeText,
		OpSetNotesShapeProps,
		OpSetModifyPassword,
		OpSetMarkAsFinal,
		OpAddTable,
		OpGetTable,
		OpMergeTableCells,
		OpSplitTableCell,
		OpUpdateTableFlags,
		OpUpdateTableCell,
		OpSetTableStyle,
		OpDefineTableStyle,
		OpListTableStyles,
		OpSetTableRowHeight,
		OpSetTableColumnWidth,
		OpAddCustomXML,
		OpListCustomXML,
		OpRemoveCustomXML,
		OpAddVba,
		OpExportPDF,
		OpExportHTML,
		OpListPlaceholders,
		OpSetPlaceholderContent,
		OpGetImageMetadata,
		OpAddVideo,
		OpAddAudio,
		OpAddOLEObject,
		OpMarkdownToSlides,
		OpURLFetchToSlides,
		OpAddMermaidShape,
		OpAddSmartArt,
		OpAddAnimation,
		OpSetSlideTransition,
		OpMoveShapeToIndex,
		OpListSlideImages,
		OpSwapImageByIndex,
		OpSwapImageByRelID,
		OpListNotesShapes,
		OpListNotesPlaceholders,
		OpUpdateNotesMaster,
		OpMergeFromEditor,
		OpUpdateSmartArt,
		OpSetSlideBackground,
		OpSetSlideHeaderFooter,
		OpGetHandoutMaster,
		OpUpdateHandoutMaster,
		OpHasDigitalSignature,
	}
}

// SupportedOps returns the full list of operation codes handled by the bridge.
func SupportedOps() []string {
	ops := supportedSlideAndMetaOps()
	ops = append(ops, supportedContentOps()...)
	return ops
}
