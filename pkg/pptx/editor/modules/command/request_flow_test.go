package command

import (
	"errors"
	"testing"
)

func TestHandleParsedRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := HandleParsedRequest[int](
			[]byte(`{"slide_index":2}`),
			func(_ []byte) (map[string]any, error) {
				return map[string]any{"slide_index": 2}, nil
			},
			func(payload map[string]any) (int, bool) {
				v, ok := payload["slide_index"].(int)
				return v, ok
			},
			func() error { return errors.New("validation failed") },
			func(req int) (any, error) {
				if req != 2 {
					t.Fatalf("expected request=2, got=%d", req)
				}
				return "ok", nil
			},
		)
		if err != nil || got != "ok" {
			t.Fatalf("unexpected result: got=%v err=%v", got, err)
		}
	})

	t.Run("parseRaw error", func(t *testing.T) {
		parseErr := errors.New("parse raw failed")
		got, err := HandleParsedRequest[int](
			[]byte(`{}`),
			func(_ []byte) (map[string]any, error) {
				return nil, parseErr
			},
			func(map[string]any) (int, bool) {
				t.Fatal("parseReq should not be called on parseRaw error")
				return 0, false
			},
			func() error {
				t.Fatal("validationErr should not be called on parseRaw error")
				return nil
			},
			func(int) (any, error) {
				t.Fatal("execute should not be called on parseRaw error")
				return nil, errors.New("execute should not be called")
			},
		)
		if !errors.Is(err, parseErr) || got != nil {
			t.Fatalf("expected parseErr and nil result, got=%v err=%v", got, err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		validationErr := errors.New("validation failed")
		got, err := HandleParsedRequest[int](
			[]byte(`{}`),
			func(_ []byte) (map[string]any, error) {
				return map[string]any{}, nil
			},
			func(map[string]any) (int, bool) {
				return 0, false
			},
			func() error { return validationErr },
			func(int) (any, error) {
				t.Fatal("execute should not be called when request validation fails")
				return nil, errors.New("execute should not be called")
			},
		)
		if !errors.Is(err, validationErr) || got != nil {
			t.Fatalf("expected validationErr and nil result, got=%v err=%v", got, err)
		}
	})
}

func TestHandleParsedRequestWithPayload(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := HandleParsedRequestWithPayload[int](
			[]byte(`{"slide_index":3}`),
			func(_ []byte) (map[string]any, error) {
				return map[string]any{"slide_index": 3}, nil
			},
			func(payload map[string]any) (int, bool) {
				v, ok := payload["slide_index"].(int)
				return v, ok
			},
			func() error { return errors.New("validation failed") },
			func(req int, payload map[string]any) (any, error) {
				if req != 3 {
					t.Fatalf("expected request=3, got=%d", req)
				}
				if payload["slide_index"] != 3 {
					t.Fatalf("expected payload slide_index=3, got=%v", payload["slide_index"])
				}
				return "ok-with-payload", nil
			},
		)
		if err != nil || got != "ok-with-payload" {
			t.Fatalf("unexpected result: got=%v err=%v", got, err)
		}
	})

	t.Run("parseRaw error", func(t *testing.T) {
		parseErr := errors.New("parse raw failed")
		got, err := HandleParsedRequestWithPayload[int](
			[]byte(`{}`),
			func(_ []byte) (map[string]any, error) {
				return nil, parseErr
			},
			func(map[string]any) (int, bool) {
				t.Fatal("parseReq should not be called on parseRaw error")
				return 0, false
			},
			func() error {
				t.Fatal("validationErr should not be called on parseRaw error")
				return nil
			},
			func(int, map[string]any) (any, error) {
				t.Fatal("execute should not be called on parseRaw error")
				return nil, errors.New("execute should not be called")
			},
		)
		if !errors.Is(err, parseErr) || got != nil {
			t.Fatalf("expected parseErr and nil result, got=%v err=%v", got, err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		validationErr := errors.New("validation failed")
		got, err := HandleParsedRequestWithPayload[int](
			[]byte(`{}`),
			func(_ []byte) (map[string]any, error) {
				return map[string]any{}, nil
			},
			func(map[string]any) (int, bool) {
				return 0, false
			},
			func() error { return validationErr },
			func(int, map[string]any) (any, error) {
				t.Fatal("execute should not be called when request validation fails")
				return nil, errors.New("execute should not be called")
			},
		)
		if !errors.Is(err, validationErr) || got != nil {
			t.Fatalf("expected validationErr and nil result, got=%v err=%v", got, err)
		}
	})
}

func TestHandleSlideIndexRequest(t *testing.T) {
	testHandleRequestFlowPaths(
		t,
		func(rawPayload []byte, parseRaw ParseRawPayloadFn, validationErr ValidationErrFn, expect any) (any, error) {
			return HandleSlideIndexRequest(
				rawPayload,
				parseRaw,
				testParseSlideIndex,
				validationErr,
				func(req int) (any, error) {
					if req != expect.(int) {
						t.Fatalf("expected slide_index=%d, got=%d", expect.(int), req)
					}
					return "ok", nil
				},
			)
		},
		map[string]any{"slide_index": 1},
		1,
		map[string]any{},
	)
}

func TestHandleSlideShapeRequest(t *testing.T) {
	testHandleRequestFlowPaths(
		t,
		func(rawPayload []byte, parseRaw ParseRawPayloadFn, validationErr ValidationErrFn, expect any) (any, error) {
			return HandleSlideShapeRequest(
				rawPayload,
				parseRaw,
				testParseSlideIndex,
				testParseIntField,
				validationErr,
				func(req SlideShapeRequest) (any, error) {
					expected := expect.(SlideShapeRequest)
					if req != expected {
						t.Fatalf("expected request=%+v, got=%+v", expected, req)
					}
					return "ok", nil
				},
			)
		},
		map[string]any{"slide_index": 2, "shape_id": 11},
		SlideShapeRequest{SlideIndex: 2, ShapeID: 11},
		map[string]any{"slide_index": 2},
	)
}

func TestHandleSlideShapeRequestWithPayload(t *testing.T) {
	testHandleRequestFlowPaths(
		t,
		func(rawPayload []byte, parseRaw ParseRawPayloadFn, validationErr ValidationErrFn, expect any) (any, error) {
			return HandleSlideShapeRequestWithPayload(
				rawPayload,
				parseRaw,
				testParseSlideIndex,
				testParseIntField,
				validationErr,
				func(req SlideShapeRequest, payload map[string]any) (any, error) {
					expected := expect.(SlideShapeRequest)
					if req != expected {
						t.Fatalf("expected request=%+v, got=%+v", expected, req)
					}
					if payload["shape_id"] != expected.ShapeID {
						t.Fatalf(
							"expected shape_id in payload=%d, got=%v",
							expected.ShapeID,
							payload["shape_id"],
						)
					}
					return "ok", nil
				},
			)
		},
		map[string]any{"slide_index": 3, "shape_id": 9},
		SlideShapeRequest{SlideIndex: 3, ShapeID: 9},
		map[string]any{"slide_index": 3},
	)
}

func TestHandleSlideShapeIDsRequest(t *testing.T) {
	testHandleRequestFlowPaths(
		t,
		func(rawPayload []byte, parseRaw ParseRawPayloadFn, validationErr ValidationErrFn, expect any) (any, error) {
			return HandleSlideShapeIDsRequest(
				rawPayload,
				parseRaw,
				testParseSlideIndex,
				testParseIntSliceField,
				validationErr,
				func(req SlideShapeIDsRequest) (any, error) {
					expected := expect.(SlideShapeIDsRequest)
					if req.SlideIndex != expected.SlideIndex {
						t.Fatalf(
							"expected slide_index=%d, got=%d",
							expected.SlideIndex,
							req.SlideIndex,
						)
					}
					if len(req.ShapeIDs) != len(expected.ShapeIDs) ||
						req.ShapeIDs[0] != expected.ShapeIDs[0] ||
						req.ShapeIDs[1] != expected.ShapeIDs[1] {
						t.Fatalf("expected shape_ids=%v, got=%v", expected.ShapeIDs, req.ShapeIDs)
					}
					return "ok", nil
				},
			)
		},
		map[string]any{"slide_index": 4, "shape_ids": []int{5, 6}},
		SlideShapeIDsRequest{SlideIndex: 4, ShapeIDs: []int{5, 6}},
		map[string]any{"slide_index": 4},
	)
}

type invokeFlowFn func(rawPayload []byte, parseRaw ParseRawPayloadFn, validationErr ValidationErrFn, expect any) (any, error)

func testHandleRequestFlowPaths(
	t *testing.T,
	invoke invokeFlowFn,
	successPayload map[string]any,
	expectedReq any,
	invalidPayload map[string]any,
) {
	t.Helper()

	t.Run("success", func(t *testing.T) {
		got, err := invoke([]byte(`ok`), func(_ []byte) (map[string]any, error) {
			return successPayload, nil
		}, func() error {
			return errors.New("validation failed")
		}, expectedReq)
		if err != nil || got != "ok" {
			t.Fatalf("unexpected result: got=%v err=%v", got, err)
		}
	})

	t.Run("parseRaw error", func(t *testing.T) {
		parseErr := errors.New("parse raw failed")
		got, err := invoke([]byte(`bad`), func(_ []byte) (map[string]any, error) {
			return nil, parseErr
		}, func() error {
			t.Fatal("validationErr should not be called on parseRaw error")
			return nil
		}, expectedReq)
		if !errors.Is(err, parseErr) || got != nil {
			t.Fatalf("expected parseErr and nil result, got=%v err=%v", got, err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		validationErr := errors.New("validation failed")
		got, err := invoke([]byte(`invalid`), func(_ []byte) (map[string]any, error) {
			return invalidPayload, nil
		}, func() error {
			return validationErr
		}, expectedReq)
		if !errors.Is(err, validationErr) || got != nil {
			t.Fatalf("expected validationErr and nil result, got=%v err=%v", got, err)
		}
	})
}
