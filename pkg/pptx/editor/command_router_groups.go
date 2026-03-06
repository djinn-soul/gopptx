package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

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
	case OpAddChart:
		return handleAddChart, true
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
	case OpAddShape:
		return handleAddShape, true
	case OpAddTextbox:
		return handleAddTextbox, true
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
	case OpHasNotesSlide:
		return handleHasNotesSlide, true
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
	default:
		return nil, false
	}
}

func handleSlideCount(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]int{"count": e.SlideCount()}, nil
}

func handleAddSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.AddSlideRequest, bool) {
			return editorcommand.ParseAddSlideRequest(p, v.OptionalString, v.OptionalStringSlice), true
		},
		v.Error,
		func(request editorcommand.AddSlideRequest) (any, error) {
			slide := elements.NewSlide(request.Title)
			if request.Layout != "" {
				slide = slide.WithLayout(request.Layout)
			}
			for _, bullet := range request.Bullets {
				slide = slide.AddBullet(bullet)
			}
			index, err := e.AddSlide(slide)
			if err != nil {
				return nil, err
			}
			return map[string]int{"index": index}, nil
		},
	)
}

func handleRemoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideIndexRequest, bool) {
			return editorcommand.ParseSlideIndexOnlyRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.SlideIndexRequest) (any, error) {
			if !v.IndexBounds(request.Index, 0, e.SlideCount(), "index") {
				return nil, v.Error()
			}
			return nil, e.RemoveSlide(request.Index)
		},
	)
}

func handleMoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.MoveSlideRequest, bool) {
			return editorcommand.ParseMoveSlideRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.MoveSlideRequest) (any, error) {
			slideCount := e.SlideCount()
			if !v.IndexBounds(request.From, 0, slideCount, "from") {
				return nil, v.Error()
			}
			if !v.IndexBounds(request.To, 0, slideCount, "to") {
				return nil, v.Error()
			}
			return nil, e.MoveSlide(request.From, request.To)
		},
	)
}

func handleDuplicateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.DuplicateSlideRequest, bool) {
			return editorcommand.ParseDuplicateSlideRequest(p, v.RequireInt, v.OptionalInt)
		},
		v.Error,
		func(request editorcommand.DuplicateSlideRequest) (any, error) {
			if !v.IndexBounds(request.Index, 0, e.SlideCount(), "index") {
				return nil, v.Error()
			}
			newIdx, err := e.DuplicateSlide(request.Index, request.InsertAt)
			if err != nil {
				return nil, err
			}
			return map[string]int{"new_index": newIdx}, nil
		},
	)
}

func handleGetMetadata(e *PresentationEditor, _ json.RawMessage) (any, error) {
	m := e.Metadata()
	return map[string]any{
		"title":       m.Title,
		"slide_count": m.SlideCount,
		"size": map[string]int64{
			"width":  m.SlideSize.Width,
			"height": m.SlideSize.Height,
		},
	}, nil
}
