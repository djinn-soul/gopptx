package editor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type mediaPlacement struct {
	slideIndex int
	x          float64
	y          float64
	w          float64
	h          float64
}

type mediaInsertSpec struct {
	metaKey          string
	primaryPathKey   string
	primaryDataKey   string
	secondaryPathKey string
	secondaryDataKey string
	primaryMaxLen    int
	secondaryMaxLen  int
	primaryLabel     string
	secondaryLabel   string
	insertBinary     func(mediaPlacement, string, []byte, []byte) (int, error)
	insertPath       func(mediaPlacement, string, string, string) (int, error)
}

func parseMediaPlacement(
	e *PresentationEditor,
	payload map[string]any,
	v *PayloadValidator,
) (mediaPlacement, bool) {
	slideIndex, ok := v.RequireInt(payload, "slide_index")
	if !ok {
		return mediaPlacement{}, false
	}
	x, ok := v.RequireFloat64(payload, "x")
	if !ok {
		return mediaPlacement{}, false
	}
	y, ok := v.RequireFloat64(payload, "y")
	if !ok {
		return mediaPlacement{}, false
	}
	w, ok := v.RequireFloat64(payload, "w")
	if !ok {
		return mediaPlacement{}, false
	}
	h, ok := v.RequireFloat64(payload, "h")
	if !ok {
		return mediaPlacement{}, false
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return mediaPlacement{}, false
	}
	return mediaPlacement{
		slideIndex: slideIndex,
		x:          x,
		y:          y,
		w:          w,
		h:          h,
	}, true
}

func decodeOptionalBase64Field(raw string, maxLen int, label string) ([]byte, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	if len(raw) > maxLen {
		return nil, fmt.Errorf("%s data too large", label)
	}
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s base64: %w", label, err)
	}
	return data, nil
}

func insertShapeFromBinaryOrPath(
	hasBinary bool,
	insertFromBinary func() (int, error),
	insertFromPath func() (int, error),
) (int, error) {
	if hasBinary {
		return insertFromBinary()
	}
	return insertFromPath()
}

func handleMediaInsertCommand(
	e *PresentationEditor,
	payload json.RawMessage,
	spec mediaInsertSpec,
) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	placement, ok := parseMediaPlacement(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	metaValue := v.OptionalString(p, spec.metaKey)
	primaryPath := v.OptionalString(p, spec.primaryPathKey)
	primaryBase64 := v.OptionalString(p, spec.primaryDataKey)
	secondaryPath := v.OptionalString(p, spec.secondaryPathKey)
	secondaryBase64 := v.OptionalString(p, spec.secondaryDataKey)

	primaryData, err := decodeOptionalBase64Field(primaryBase64, spec.primaryMaxLen, spec.primaryLabel)
	if err != nil {
		return nil, err
	}
	secondaryData, err := decodeOptionalBase64Field(secondaryBase64, spec.secondaryMaxLen, spec.secondaryLabel)
	if err != nil {
		return nil, err
	}

	shapeID, err := insertShapeFromBinaryOrPath(
		len(primaryData) > 0 || len(secondaryData) > 0,
		func() (int, error) {
			return spec.insertBinary(placement, metaValue, primaryData, secondaryData)
		},
		func() (int, error) {
			return spec.insertPath(placement, metaValue, primaryPath, secondaryPath)
		},
	)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": shapeID}, nil
}
