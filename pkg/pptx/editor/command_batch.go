package editor

import "encoding/json"

type batchCommand struct {
	Op        string          `json:"op"`
	Payload   json.RawMessage `json:"payload"`
	RequestID string          `json:"request_id,omitempty"`
}

type batchPayload struct {
	Commands    []batchCommand `json:"commands"`
	StopOnError bool           `json:"stop_on_error,omitempty"`
}

type batchResult struct {
	OK        bool         `json:"ok"`
	Op        string       `json:"op"`
	RequestID string       `json:"request_id,omitempty"`
	Result    any          `json:"result,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
}

func handleBatchExecute(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p batchPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	results := make([]batchResult, 0, len(p.Commands))
	for i, cmd := range p.Commands {
		if cmd.Op == OpBatchExecute {
			results = append(results, batchResult{
				OK:        false,
				Op:        cmd.Op,
				RequestID: cmd.RequestID,
				Error: &ErrorDetail{
					Code:    "INVALID_BATCH_ITEM",
					Message: "nested batch_execute is not supported",
					Details: map[string]int{"index": i},
				},
			})
			if p.StopOnError {
				break
			}
			continue
		}

		handler, ok := commandHandlerFor(cmd.Op)
		if !ok {
			results = append(results, batchResult{
				OK:        false,
				Op:        cmd.Op,
				RequestID: cmd.RequestID,
				Error: &ErrorDetail{
					Code:    ErrCodeUnknownOp,
					Message: "Operation " + `"` + cmd.Op + `"` + " not recognized",
					Details: map[string]int{"index": i},
				},
			})
			if p.StopOnError {
				break
			}
			continue
		}

		result, err := handler(e, cmd.Payload)
		if err != nil {
			// Check if error is a BridgeError with specific code
			var bridgeErr *BridgeError
			code := ErrCodeOpFailed
			details := any(map[string]any{"index": i})
			if AsBridgeError(err, &bridgeErr) {
				code = bridgeErr.Code
				details = withBatchIndex(i, bridgeErr.Details)
			}
			results = append(results, batchResult{
				OK:        false,
				Op:        cmd.Op,
				RequestID: cmd.RequestID,
				Error: &ErrorDetail{
					Code:    code,
					Message: err.Error(),
					Details: details,
				},
			})
			if p.StopOnError {
				break
			}
			continue
		}

		results = append(results, batchResult{
			OK:        true,
			Op:        cmd.Op,
			RequestID: cmd.RequestID,
			Result:    result,
		})
	}

	return map[string]any{"results": results}, nil
}

func withBatchIndex(index int, details any) map[string]any {
	out := map[string]any{"index": index}
	if details == nil {
		return out
	}
	if m, ok := details.(map[string]any); ok {
		for k, v := range m {
			out[k] = v
		}
		return out
	}
	out["cause"] = details
	return out
}

// AsBridgeError checks if an error is a BridgeError and extracts it.
func AsBridgeError(err error, target **BridgeError) bool {
	if be, ok := err.(*BridgeError); ok {
		*target = be
		return true
	}
	return false
}
