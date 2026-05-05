package editor

import "maps"

//nolint:gochecknoglobals // package-level dispatch map for O(1) command routing
var staticHandlers map[string]commandHandler

const staticHandlersInitCap = 128

//nolint:gochecknoinits // required to break the init cycle: handleBatchExecute -> commandHandlerFor -> staticHandlers
func init() {
	staticHandlers = make(map[string]commandHandler, staticHandlersInitCap)
	registerHandlers(staticHandlers, slideBasicHandlers())
	registerHandlers(staticHandlers, layoutMetadataHandlers())
	registerHandlers(staticHandlers, themeLayoutHandlers())
	registerHandlers(staticHandlers, contentSectionHandlers())
	registerHandlers(staticHandlers, commentsShapesReadHandlers())
	registerHandlers(staticHandlers, shapeMutationHandlers())
	registerHandlers(staticHandlers, notesTableHandlers())
	registerHandlers(staticHandlers, handoutSignatureHandlers())
	registerHandlers(staticHandlers, templateBuildHandlers())
}

func registerHandlers(dst map[string]commandHandler, src map[string]commandHandler) {
	maps.Copy(dst, src)
}

func slideBasicHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpBatchExecute:        handleBatchExecute,
		OpSlideCount:          handleSlideCount,
		OpAddSlide:            handleAddSlide,
		OpRemoveSlide:         handleRemoveSlide,
		OpMoveSlide:           handleMoveSlide,
		OpDuplicateSlide:      handleDuplicateSlide,
		OpDuplicateSlideAfter: handleDuplicateSlideAfter,
		OpListSlides:          handleListSlides,
		OpSetSlideTitle:       handleSetSlideTitle,
		OpUpdateSlide:         handleUpdateSlide,
		OpValidate:            handleValidate,
		OpRepair:              handleRepair,
		OpSetSlideHidden:      handleSetSlideHidden,
	}
}

func layoutMetadataHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpGetMetadata:             handleGetMetadata,
		OpListSlideCharts:         handleListSlideCharts,
		OpUpdateChartData:         handleUpdateChartData,
		OpUpdateChartDataBatch:    handleUpdateChartDataBatch,
		OpUpdateChartFormatting:   handleUpdateChartFormatting,
		OpGetChartState:           handleGetChartState,
		OpAddChart:                handleAddChart,
		OpGetSlideLayoutRef:       handleGetSlideLayoutRef,
		OpListSlideLayouts:        handleListSlideLayouts,
		OpListSlideMasters:        handleListSlideMasters,
		OpListMasterLayouts:       handleListMasterLayouts,
		OpRebindSlideLayout:       handleRebindSlideLayout,
		OpCloneLayoutMasterFamily: handleCloneLayoutMasterFamily,
		OpAddSlideMaster:          handleAddSlideMaster,
		OpRemoveSlideMaster:       handleRemoveSlideMaster,
		OpAddSlideLayout:          handleAddSlideLayout,
		OpRemoveSlideLayout:       handleRemoveSlideLayout,
		OpApplyTheme:              handleApplyTheme,
		OpSetSlideSize:            handleSetSlideSize,
		OpMergeFromFile:           handleMergeFromFile,
		OpGetCoreProperties:       handleGetCoreProperties,
		OpSetCoreProperties:       handleSetCoreProperties,
		OpListPlaceholders:        handleListPlaceholders,
		OpSetPlaceholderContent:   handleSetPlaceholderContent,
	}
}

func themeLayoutHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpGetLayoutShapes:       handleGetLayoutShapes,
		OpGetMasterShapes:       handleGetMasterShapes,
		OpGetLayoutPlaceholders: handleGetLayoutPlaceholders,
		OpGetMasterPlaceholders: handleGetMasterPlaceholders,
		OpSetGlobalThemePreset:  handleSetGlobalThemePreset,
		OpSetThemeFontScheme:    handleSetThemeFontScheme,
		OpSetThemeColorScheme:   handleSetThemeColorScheme,
		OpGetThemeInventory:     handleGetThemeInventory,
	}
}

func contentSectionHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpAddSection:           handleAddSection,
		OpRemoveSection:        handleRemoveSection,
		OpRenameSection:        handleRenameSection,
		OpGetSections:          handleGetSections,
		OpFindAndReplace:       handleFindAndReplace,
		OpSearchShapes:         handleSearchShapes,
		OpSetModifyPassword:    handleSetModifyPassword,
		OpSetMarkAsFinal:       handleSetMarkAsFinal,
		OpAddCustomXML:         handleAddCustomXML,
		OpListCustomXML:        handleListCustomXML,
		OpRemoveCustomXML:      handleRemoveCustomXML,
		OpAddVba:               handleAddVba,
		OpMarkdownToSlides:     handleMarkdownToSlides,
		OpURLFetchToSlides:     handleURLFetchToSlides,
		OpAddMermaidShape:      handleAddMermaidShape,
		OpAddSmartArt:          handleAddSmartArt,
		OpUpdateSmartArt:       handleUpdateSmartArt,
		OpDeleteSmartArt:       handleDeleteSmartArt,
		OpChangeSmartArtLayout: handleChangeSmartArtLayout,
		OpSetSmartArtStyle:     handleSetSmartArtStyle,
		OpSetSmartArtNodes:     handleSetSmartArtNodes,
		OpSetSlideBackground:   handleSetSlideBackground,
		OpSetSlideHeaderFooter: handleSetSlideHeaderFooter,
		OpGetSlideHeaderFooter: handleGetSlideHeaderFooter,
		OpAddAnimation:         handleAddAnimation,
		OpSetSlideTransition:   handleSetSlideTransition,
		OpMergeFromEditor:      handleMergeFromEditor,
		OpConvertToGrayscale:   handleConvertToGrayscale,
	}
}

func commentsShapesReadHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpGetAuthors:          handleGetAuthors,
		OpAddAuthor:           handleAddAuthor,
		OpGetComments:         handleGetComments,
		OpAddComment:          handleAddComment,
		OpRemoveComment:       handleRemoveComment,
		OpListShapes:          handleListShapes,
		OpGetSlideTextStates:  handleGetSlideTextStates,
		OpGetShapeTextState:   handleGetShapeTextState,
		OpGetShapeRuns:        handleGetShapeRuns,
		OpSetShapeRuns:        handleSetShapeRuns,
		OpSetSlideShapeRuns:   handleSetSlideShapeRuns,
		OpUpdateDeckRunTexts:  handleUpdateDeckRunTexts,
		OpUpdateSlideRunTexts: handleUpdateSlideRunTexts,
		OpUpdateShapeRunText:  handleUpdateShapeRunText,
		OpAppendShapeRun:      handleAppendShapeRun,
	}
}

func shapeMutationHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpAddShape:         handleAddShape,
		OpAddTextbox:       handleAddTextbox,
		OpAddTextboxes:     handleAddTextboxes,
		OpAddConnectors:    handleAddConnectors,
		OpReserveShapeIDs:  handleReserveShapeIDs,
		OpAddConnector:     handleAddConnector,
		OpAddGroupShape:    handleAddGroupShape,
		OpBuildFreeform:    handleBuildFreeform,
		OpAddImage:         handleAddImage,
		OpRemoveShape:      handleRemoveShape,
		OpClearShapes:      handleClearShapes,
		OpGroupShapes:      handleGroupShapes,
		OpUngroupShapes:    handleUngroupShapes,
		OpUpdateShape:      handleUpdateShape,
		OpMoveShapeToFront: handleMoveShapeToFront,
		OpMoveShapeToBack:  handleMoveShapeToBack,
		OpMoveShapeToIndex: handleMoveShapeToIndex,
		OpGetImageMetadata: handleGetImageMetadata,
		OpAddVideo:         handleAddVideo,
		OpAddAudio:         handleAddAudio,
		OpAddOLEObject:     handleAddOLEObject,
		OpListSlideImages:  handleListSlideImages,
		OpSwapImageByIndex: handleSwapImageByIndex,
		OpSwapImageByRelID: handleSwapImageByRelID,
	}
}

func notesTableHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpGetNotes:              handleGetNotes,
		OpNotesSlideExists:      handleNotesSlideExists,
		OpSetNotes:              handleSetNotes,
		OpSetNotesShapeText:     handleSetNotesShapeText,
		OpSetNotesShapeProps:    handleSetNotesShapeProps,
		OpAddTable:              handleAddTable,
		OpGetTable:              handleGetTable,
		OpMergeTableCells:       handleMergeTableCells,
		OpSplitTableCell:        handleSplitTableCell,
		OpUpdateTableFlags:      handleUpdateTableFlags,
		OpUpdateTableCell:       handleUpdateTableCell,
		OpSetTableStyle:         handleSetTableStyle,
		OpDefineTableStyle:      handleDefineTableStyle,
		OpListTableStyles:       handleListTableStyles,
		OpSetTableRowHeight:     handleSetTableRowHeight,
		OpSetTableColumnWidth:   handleSetTableColumnWidth,
		OpAddTableRow:           handleAddTableRow,
		OpAddTableColumn:        handleAddTableColumn,
		OpInsertTableRow:        handleInsertTableRow,
		OpRemoveTableRow:        handleRemoveTableRow,
		OpInsertTableColumn:     handleInsertTableColumn,
		OpRemoveTableColumn:     handleRemoveTableColumn,
		OpUpdateTableCellBorder: handleUpdateTableCellBorder,
		OpListNotesShapes:       handleListNotesShapes,
		OpListNotesPlaceholders: handleListNotesPlaceholders,
		OpUpdateNotesMaster:     handleUpdateNotesMaster,
	}
}

func handoutSignatureHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpGetHandoutMaster:    handleGetHandoutMaster,
		OpUpdateHandoutMaster: handleUpdateHandoutMaster,
		OpIsDigitallySigned:   handleIsDigitallySigned,
	}
}

func templateBuildHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		OpBuildStatusTemplate:    handleBuildStatusTemplate,
		OpBuildSimpleTemplate:    handleBuildSimpleTemplate,
		OpBuildProposalTemplate:  handleBuildProposalTemplate,
		OpBuildTrainingTemplate:  handleBuildTrainingTemplate,
		OpBuildTechnicalTemplate: handleBuildTechnicalTemplate,
		OpRenderTemplate:         handleRenderTemplate,
	}
}
