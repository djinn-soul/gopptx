package command

type ParseRequestFn[T any] func(map[string]any) (T, bool)
type ExecuteRequestFn[T any] func(T) (any, error)
type ExecuteRequestWithPayloadFn[T any] func(T, map[string]any) (any, error)
type ValidationErrFn func() error

func HandleParsedRequest[T any](
	rawPayload []byte,
	parseRaw ParseRawPayloadFn,
	parseReq ParseRequestFn[T],
	validationErr ValidationErrFn,
	execute ExecuteRequestFn[T],
) (any, error) {
	payload, err := parseRaw(rawPayload)
	if err != nil {
		return nil, err
	}
	req, ok := parseReq(payload)
	if !ok {
		return nil, validationErr()
	}
	return execute(req)
}

func HandleParsedRequestWithPayload[T any](
	rawPayload []byte,
	parseRaw ParseRawPayloadFn,
	parseReq ParseRequestFn[T],
	validationErr ValidationErrFn,
	execute ExecuteRequestWithPayloadFn[T],
) (any, error) {
	payload, err := parseRaw(rawPayload)
	if err != nil {
		return nil, err
	}
	req, ok := parseReq(payload)
	if !ok {
		return nil, validationErr()
	}
	return execute(req, payload)
}

func HandleSlideIndexRequest(
	rawPayload []byte,
	parseRaw ParseRawPayloadFn,
	parseSlideIndex ParseSlideIndexFn,
	validationErr ValidationErrFn,
	execute ExecuteRequestFn[int],
) (any, error) {
	return HandleParsedRequest(
		rawPayload,
		parseRaw,
		func(payload map[string]any) (int, bool) {
			return ParseSlideIndexRequest(payload, parseSlideIndex)
		},
		validationErr,
		execute,
	)
}

func HandleSlideShapeRequest(
	rawPayload []byte,
	parseRaw ParseRawPayloadFn,
	parseSlideIndex ParseSlideIndexFn,
	parseIntField ParseIntFieldFn,
	validationErr ValidationErrFn,
	execute ExecuteRequestFn[SlideShapeRequest],
) (any, error) {
	return HandleParsedRequest(
		rawPayload,
		parseRaw,
		func(payload map[string]any) (SlideShapeRequest, bool) {
			return ParseSlideShapeRequest(payload, parseSlideIndex, parseIntField)
		},
		validationErr,
		execute,
	)
}

func HandleSlideShapeRequestWithPayload(
	rawPayload []byte,
	parseRaw ParseRawPayloadFn,
	parseSlideIndex ParseSlideIndexFn,
	parseIntField ParseIntFieldFn,
	validationErr ValidationErrFn,
	execute ExecuteRequestWithPayloadFn[SlideShapeRequest],
) (any, error) {
	return HandleParsedRequestWithPayload(
		rawPayload,
		parseRaw,
		func(payload map[string]any) (SlideShapeRequest, bool) {
			return ParseSlideShapeRequest(payload, parseSlideIndex, parseIntField)
		},
		validationErr,
		execute,
	)
}

func HandleSlideShapeIDsRequest(
	rawPayload []byte,
	parseRaw ParseRawPayloadFn,
	parseSlideIndex ParseSlideIndexFn,
	parseIntSliceField ParseIntSliceFieldFn,
	validationErr ValidationErrFn,
	execute ExecuteRequestFn[SlideShapeIDsRequest],
) (any, error) {
	return HandleParsedRequest(
		rawPayload,
		parseRaw,
		func(payload map[string]any) (SlideShapeIDsRequest, bool) {
			return ParseSlideShapeIDsRequest(payload, parseSlideIndex, parseIntSliceField)
		},
		validationErr,
		execute,
	)
}
