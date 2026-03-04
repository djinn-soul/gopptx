package command

import (
	"encoding/json"
	"maps"
)

type BatchCommand struct {
	Op        string          `json:"op"`
	Payload   json.RawMessage `json:"payload"`
	RequestID string          `json:"request_id,omitempty"`
}

type BatchPayload struct {
	Commands    []BatchCommand `json:"commands"`
	StopOnError bool           `json:"stop_on_error,omitempty"`
}

type BatchResultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type BatchResult struct {
	OK        bool              `json:"ok"`
	Op        string            `json:"op"`
	RequestID string            `json:"request_id,omitempty"`
	Result    any               `json:"result,omitempty"`
	Error     *BatchResultError `json:"error,omitempty"`
}

type BridgeErrorView struct {
	Code    string
	Message string
	Details any
}

type BatchOptions struct {
	BatchOp       string
	UnknownOpCode string
	OpFailedCode  string
}

func HandleBatchExecute(
	payload json.RawMessage,
	resolve func(op string) (func(json.RawMessage) (any, error), bool),
	asBridgeError func(error) (BridgeErrorView, bool),
	options BatchOptions,
) (map[string]any, error) {
	var p BatchPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	results := make([]BatchResult, 0, len(p.Commands))
	for i, cmd := range p.Commands {
		if cmd.Op == options.BatchOp {
			results = append(results, BatchResult{
				OK:        false,
				Op:        cmd.Op,
				RequestID: cmd.RequestID,
				Error: &BatchResultError{
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

		handler, ok := resolve(cmd.Op)
		if !ok {
			results = append(results, BatchResult{
				OK:        false,
				Op:        cmd.Op,
				RequestID: cmd.RequestID,
				Error: &BatchResultError{
					Code:    options.UnknownOpCode,
					Message: "Operation " + `"` + cmd.Op + `"` + " not recognized",
					Details: map[string]int{"index": i},
				},
			})
			if p.StopOnError {
				break
			}
			continue
		}

		result, err := handler(cmd.Payload)
		if err != nil {
			code := options.OpFailedCode
			details := any(map[string]any{"index": i})
			if be, ok := asBridgeError(err); ok {
				code = be.Code
				details = withBatchIndex(i, be.Details)
			}
			results = append(results, BatchResult{
				OK:        false,
				Op:        cmd.Op,
				RequestID: cmd.RequestID,
				Error: &BatchResultError{
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

		results = append(results, BatchResult{
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
		maps.Copy(out, m)
		return out
	}
	out["cause"] = details
	return out
}
