package command

import (
	"encoding/base64"
	"errors"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestParseAddImageRequest_BasicAndValidation(t *testing.T) {
	payload := map[string]any{
		"slide_index": 2,
		"x":           1.0,
		"y":           2.0,
		"w":           3.0,
		"h":           4.0,
		"path":        "img.png",
		"url":         "https://example.com/img.png",
		"data":        "Zm9v",
		"format":      "png",
	}
	req, ok, err := ParseAddImageRequest(
		payload,
		func(m map[string]any) (int, bool) { return int(m["slide_index"].(int)), true },
		func(m map[string]any, key string) (float64, bool) {
			v, ok := m[key].(float64)
			return v, ok
		},
		func(m map[string]any, key string) string {
			s, _ := m[key].(string)
			return s
		},
	)
	if err != nil || !ok {
		t.Fatalf("expected parse success, got ok=%v err=%v", ok, err)
	}
	if req.SlideIndex != 2 || req.X != 1 || req.Y != 2 || req.W != 3 || req.H != 4 {
		t.Fatalf("unexpected geometry parse: %+v", req)
	}
	if req.ImagePath != "img.png" || req.ImageURL == "" || req.Base64Data == "" || req.Format != "png" {
		t.Fatalf("unexpected source fields: %+v", req)
	}

	_, ok, err = ParseAddImageRequest(
		payload,
		func(map[string]any) (int, bool) { return 0, false },
		func(map[string]any, string) (float64, bool) { return 0, true },
		func(map[string]any, string) string { return "" },
	)
	if err != nil || ok {
		t.Fatalf("expected validation miss without error, got ok=%v err=%v", ok, err)
	}
}

func TestParseAddImageRequest_OptionsDecodeError(t *testing.T) {
	payload := map[string]any{
		"slide_index": 1,
		"x":           1.0,
		"y":           2.0,
		"w":           3.0,
		"h":           4.0,
		"options":     []any{"invalid"},
	}
	_, ok, err := ParseAddImageRequest(
		payload,
		func(map[string]any) (int, bool) { return 1, true },
		func(m map[string]any, key string) (float64, bool) {
			v, ok := m[key].(float64)
			return v, ok
		},
		func(map[string]any, string) string { return "" },
	)
	if !ok || err == nil {
		t.Fatalf("expected decode error with ok=true, got ok=%v err=%v", ok, err)
	}
}

func TestDecodeImagePayload_Branches(t *testing.T) {
	decoded, err := DecodeImagePayload("", "", 32)
	if err != nil || decoded != nil {
		t.Fatalf("expected nil,nil for empty payload, got data=%v err=%v", decoded, err)
	}
	if _, err := DecodeImagePayload("Zm9v", "", 32); err == nil {
		t.Fatalf("expected format-required error")
	}
	if _, err := DecodeImagePayload("not-base64", "png", 32); err == nil {
		t.Fatalf("expected invalid base64 error")
	}
}

func TestExecuteAddImageRequest_RoutesSources(t *testing.T) {
	base64Data := base64.StdEncoding.EncodeToString([]byte("hello"))
	req := AddImageRequest{
		SlideIndex: 1,
		X:          10,
		Y:          20,
		W:          30,
		H:          40,
		Base64Data: base64Data,
		Format:     "png",
	}
	called := ""
	id, err := ExecuteAddImageRequest(
		req,
		1024,
		func(slideIndex int, data []byte, format string, x, y, w, h float64, _ *common.ShapeUpdate) (int, error) {
			called = "bytes"
			if slideIndex != 1 || string(data) != "hello" || format != "png" || x != 10 || y != 20 || w != 30 || h != 40 {
				t.Fatalf("unexpected bytes args")
			}
			return 11, nil
		},
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) {
			t.Fatalf("url handler should not be called")
			return 0, nil
		},
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) {
			t.Fatalf("path handler should not be called")
			return 0, nil
		},
	)
	if err != nil || id != 11 || called != "bytes" {
		t.Fatalf("unexpected bytes branch result id=%d called=%s err=%v", id, called, err)
	}

	req = AddImageRequest{SlideIndex: 2, X: 1, Y: 2, W: 3, H: 4, ImageURL: "https://x/y.png"}
	id, err = ExecuteAddImageRequest(
		req,
		1024,
		func(int, []byte, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) {
			t.Fatalf("bytes handler should not be called")
			return 0, nil
		},
		func(slideIndex int, sourceURL string, _, _, _, _ float64, _ *common.ShapeUpdate) (int, error) {
			if slideIndex != 2 || sourceURL == "" {
				t.Fatalf("unexpected url args")
			}
			return 22, nil
		},
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) {
			t.Fatalf("path handler should not be called")
			return 0, nil
		},
	)
	if err != nil || id != 22 {
		t.Fatalf("unexpected url branch result id=%d err=%v", id, err)
	}

	req = AddImageRequest{SlideIndex: 3, X: 5, Y: 6, W: 7, H: 8, ImagePath: "local.png"}
	id, err = ExecuteAddImageRequest(
		req,
		1024,
		func(int, []byte, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) {
			t.Fatalf("bytes handler should not be called")
			return 0, nil
		},
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) {
			t.Fatalf("url handler should not be called")
			return 0, nil
		},
		func(slideIndex int, imagePath string, _, _, _, _ float64, _ *common.ShapeUpdate) (int, error) {
			if slideIndex != 3 || imagePath != "local.png" {
				t.Fatalf("unexpected path args")
			}
			return 33, nil
		},
	)
	if err != nil || id != 33 {
		t.Fatalf("unexpected path branch result id=%d err=%v", id, err)
	}
}

func TestExecuteAddImageRequest_ByteDecodeError(t *testing.T) {
	req := AddImageRequest{Base64Data: "bad", Format: "png"}
	_, err := ExecuteAddImageRequest(
		req,
		1024,
		func(int, []byte, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) { return 0, nil },
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) { return 0, nil },
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) { return 0, nil },
	)
	if err == nil {
		t.Fatalf("expected decode error")
	}

	req = AddImageRequest{Base64Data: "Zm9v", Format: ""}
	_, err = ExecuteAddImageRequest(
		req,
		1024,
		func(int, []byte, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) { return 0, errors.New("should not call") },
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) { return 0, nil },
		func(int, string, float64, float64, float64, float64, *common.ShapeUpdate) (int, error) { return 0, nil },
	)
	if err == nil {
		t.Fatalf("expected format-required error")
	}
}
