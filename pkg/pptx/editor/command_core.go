package editor

import (
	"encoding/json"
	"errors"
	"fmt"
)

// RequestEnvelope is the standard wrapper for all incoming commands.
type RequestEnvelope struct {
	APIVersion int             `json:"api_version"`
	Op         string          `json:"op"`
	Payload    json.RawMessage `json:"payload"`
	RequestID  string          `json:"request_id"`
}

// ResponseEnvelope is the standard wrapper for all outgoing responses.
type ResponseEnvelope struct {
	OK        bool         `json:"ok"`
	Result    any          `json:"result,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type commandHandler func(*PresentationEditor, json.RawMessage) (any, error)

func commandHandlerFor(op string) (commandHandler, bool) {
	for _, lookup := range []func(string) (commandHandler, bool){
		commandHandlerForSlides,
		commandHandlerForLayoutMetadata,
		commandHandlerForContent,
		commandHandlerForCommentsShapes,
		commandHandlerForNotesTables,
	} {
		if h, ok := lookup(op); ok {
			return h, true
		}
	}
	return nil, false
}

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
	case OpAddImage:
		return handleAddImage, true
	case OpRemoveShape:
		return handleRemoveShape, true
	case OpUpdateShape:
		return handleUpdateShape, true
	case OpMoveShapeToFront:
		return handleMoveShapeToFront, true
	case OpMoveShapeToBack:
		return handleMoveShapeToBack, true
	default:
		return nil, false
	}
}

func commandHandlerForNotesTables(op string) (commandHandler, bool) {
	switch op {
	case OpGetNotes:
		return handleGetNotes, true
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
	default:
		return nil, false
	}
}

// ExecuteCommand dispatches a JSON command to the appropriate editor method.
func ExecuteCommand(e *PresentationEditor, jsonInput string) string {
	var req RequestEnvelope
	if err := json.Unmarshal([]byte(jsonInput), &req); err != nil {
		return errorResponse(ErrCodeInvalidJSON, err.Error(), "")
	}

	if req.APIVersion != 1 {
		return errorResponse(
			ErrCodeUnsupportedVer,
			fmt.Sprintf("API version %d not supported", req.APIVersion),
			req.RequestID,
		)
	}

	handler, ok := commandHandlerFor(req.Op)
	if !ok {
		return errorResponse(ErrCodeUnknownOp, fmt.Sprintf("Operation %q not recognized", req.Op), req.RequestID)
	}

	result, err := handler(e, req.Payload)
	if err != nil {
		// Check if error is a BridgeError with specific code
		var bridgeErr *BridgeError
		if errors.As(err, &bridgeErr) {
			return errorResponseWithDetails(bridgeErr.Code, bridgeErr.Message, bridgeErr.Details, req.RequestID)
		}
		return errorResponse(ErrCodeOpFailed, err.Error(), req.RequestID)
	}

	resp := ResponseEnvelope{
		OK:        true,
		Result:    result,
		RequestID: req.RequestID,
	}
	out, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return errorResponse(ErrCodeMarshalError, marshalErr.Error(), req.RequestID)
	}
	return string(out)
}

func errorResponse(code, message, reqID string) string {
	resp := ResponseEnvelope{
		OK:        false,
		RequestID: reqID,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	out, err := json.Marshal(resp)
	if err != nil {
		return `{"ok": false, "error": {"code": "INTERNAL_ERROR", "message": "Failed to marshal error response"}}`
	}
	return string(out)
}

func errorResponseWithDetails(code, message string, details any, reqID string) string {
	resp := ResponseEnvelope{
		OK:        false,
		RequestID: reqID,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	out, err := json.Marshal(resp)
	if err != nil {
		return `{"ok": false, "error": {"code": "INTERNAL_ERROR", "message": "Failed to marshal error response"}}`
	}
	return string(out)
}
