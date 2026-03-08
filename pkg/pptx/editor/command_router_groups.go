package editor

func commandHandlerForSlides(op string) (commandHandler, bool) {
	switch op {
	case OpBatchExecute:
		return handleBatchExecute, true
	case OpSlideCount:
		return handleSlideCount, true
	case OpAddSlide:
		return handleAddSlide, true
	case OpRemoveSlide:
		return handleRemoveSlide, true
	case OpMoveSlide:
		return handleMoveSlide, true
	case OpDuplicateSlide:
		return handleDuplicateSlide, true
	case OpListSlides:
		return handleListSlides, true
	case OpSetSlideTitle:
		return handleSetSlideTitle, true
	case OpUpdateSlide:
		return handleUpdateSlide, true
	default:
		return nil, false
	}
}

func commandHandlerForLayoutMetadata(op string) (commandHandler, bool) {
	switch op {
	case OpGetMetadata:
		return handleGetMetadata, true
	case OpListSlideCharts:
		return handleListSlideCharts, true
	case OpUpdateChartData:
		return handleUpdateChartData, true
	case OpUpdateChartFormatting:
		return handleUpdateChartFormatting, true
	case OpGetChartState:
		return handleGetChartState, true
	case OpAddChart:
		return handleAddChart, true
	case OpGetSlideLayoutRef:
		return handleGetSlideLayoutRef, true
	case OpListSlideLayouts:
		return handleListSlideLayouts, true
	case OpListSlideMasters:
		return handleListSlideMasters, true
	case OpListMasterLayouts:
		return handleListMasterLayouts, true
	case OpRebindSlideLayout:
		return handleRebindSlideLayout, true
	case OpCloneLayoutMasterFamily:
		return handleCloneLayoutMasterFamily, true
	case OpAddSlideMaster:
		return handleAddSlideMaster, true
	case OpRemoveSlideMaster:
		return handleRemoveSlideMaster, true
	case OpAddSlideLayout:
		return handleAddSlideLayout, true
	case OpRemoveSlideLayout:
		return handleRemoveSlideLayout, true
	case OpApplyTheme:
		return handleApplyTheme, true
	case OpSetSlideSize:
		return handleSetSlideSize, true
	case OpMergeFromFile:
		return handleMergeFromFile, true
	case OpGetCoreProperties:
		return handleGetCoreProperties, true
	case OpSetCoreProperties:
		return handleSetCoreProperties, true
	case OpListPlaceholders:
		return handleListPlaceholders, true
	case OpSetPlaceholderContent:
		return handleSetPlaceholderContent, true
	default:
		return nil, false
	}
}

func commandHandlerForContent(op string) (commandHandler, bool) {
	switch op {
	case OpAddSection:
		return handleAddSection, true
	case OpRemoveSection:
		return handleRemoveSection, true
	case OpRenameSection:
		return handleRenameSection, true
	case OpGetSections:
		return handleGetSections, true
	case OpFindAndReplace:
		return handleFindAndReplace, true
	case OpSearchShapes:
		return handleSearchShapes, true
	case OpSetModifyPassword:
		return handleSetModifyPassword, true
	case OpSetMarkAsFinal:
		return handleSetMarkAsFinal, true
	case OpAddCustomXML:
		return handleAddCustomXML, true
	case OpListCustomXML:
		return handleListCustomXML, true
	case OpRemoveCustomXML:
		return handleRemoveCustomXML, true
	case OpAddVba:
		return handleAddVba, true
	default:
		return nil, false
	}
}

func commandHandlerForCommentsShapes(op string) (commandHandler, bool) {
	switch op {
	case OpGetAuthors:
		return handleGetAuthors, true
	case OpAddAuthor:
		return handleAddAuthor, true
	case OpGetComments:
		return handleGetComments, true
	case OpAddComment:
		return handleAddComment, true
	case OpRemoveComment:
		return handleRemoveComment, true
	case OpListShapes:
		return handleListShapes, true
	case OpGetSlideTextStates:
		return handleGetSlideTextStates, true
	case OpGetShapeTextState:
		return handleGetShapeTextState, true
	case OpGetShapeRuns:
		return handleGetShapeRuns, true
	case OpSetShapeRuns:
		return handleSetShapeRuns, true
	case OpUpdateDeckRunTexts:
		return handleUpdateDeckRunTexts, true
	case OpUpdateSlideRunTexts:
		return handleUpdateSlideRunTexts, true
	case OpUpdateShapeRunText:
		return handleUpdateShapeRunText, true
	case OpAppendShapeRun:
		return handleAppendShapeRun, true
	default:
		return commandHandlerForShapeMutations(op)
	}
}

func commandHandlerForShapeMutations(op string) (commandHandler, bool) {
	switch op {
	case OpAddShape:
		return handleAddShape, true
	case OpAddTextbox:
		return handleAddTextbox, true
	case OpAddTextboxes:
		return handleAddTextboxes, true
	case OpReserveShapeIDs:
		return handleReserveShapeIDs, true
	case OpAddConnector:
		return handleAddConnector, true
	case OpAddGroupShape:
		return handleAddGroupShape, true
	case OpBuildFreeform:
		return handleBuildFreeform, true
	case OpAddImage:
		return handleAddImage, true
	case OpRemoveShape:
		return handleRemoveShape, true
	case OpGroupShapes:
		return handleGroupShapes, true
	case OpUngroupShapes:
		return handleUngroupShapes, true
	case OpUpdateShape:
		return handleUpdateShape, true
	case OpMoveShapeToFront:
		return handleMoveShapeToFront, true
	case OpMoveShapeToBack:
		return handleMoveShapeToBack, true
	case OpGetImageMetadata:
		return handleGetImageMetadata, true
	case OpAddVideo:
		return handleAddVideo, true
	case OpAddAudio:
		return handleAddAudio, true
	case OpAddOLEObject:
		return handleAddOLEObject, true
	default:
		return nil, false
	}
}

func commandHandlerForNotesTables(op string) (commandHandler, bool) {
	switch op {
	case OpGetNotes:
		return handleGetNotes, true
	case OpNotesSlideExists:
		return handleNotesSlideExists, true
	case OpSetNotes:
		return handleSetNotes, true
	case OpAddTable:
		return handleAddTable, true
	case OpGetTable:
		return handleGetTable, true
	case OpMergeTableCells:
		return handleMergeTableCells, true
	case OpSplitTableCell:
		return handleSplitTableCell, true
	case OpUpdateTableFlags:
		return handleUpdateTableFlags, true
	case OpUpdateTableCell:
		return handleUpdateTableCell, true
	case OpSetTableStyle:
		return handleSetTableStyle, true
	case OpDefineTableStyle:
		return handleDefineTableStyle, true
	case OpListTableStyles:
		return handleListTableStyles, true
	case OpSetTableRowHeight:
		return handleSetTableRowHeight, true
	case OpSetTableColumnWidth:
		return handleSetTableColumnWidth, true
	default:
		return nil, false
	}
}
