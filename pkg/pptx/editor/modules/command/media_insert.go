package command

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type MediaPlacement struct {
	SlideIndex int
	X          float64
	Y          float64
	W          float64
	H          float64
}

type MediaInsertSpec struct {
	MetaKey          string
	PrimaryPathKey   string
	PrimaryDataKey   string
	SecondaryPathKey string
	SecondaryDataKey string
	PrimaryMaxLen    int
	SecondaryMaxLen  int
	PrimaryLabel     string
	SecondaryLabel   string
	InsertBinary     func(MediaPlacement, string, []byte, []byte) (int, error)
	InsertPath       func(MediaPlacement, string, string, string) (int, error)
}

type VideoBinaryInsertFn func(int, []byte, []byte, string, float64, float64, float64, float64) (int, error)
type VideoPathInsertFn func(int, string, string, string, float64, float64, float64, float64) (int, error)
type OLEBinaryInsertFn func(int, []byte, []byte, string, float64, float64, float64, float64) (int, error)
type OLEPathInsertFn func(int, string, string, string, float64, float64, float64, float64) (int, error)

func AdaptVideoBinaryInsert(insert VideoBinaryInsertFn) func(MediaPlacement, string, []byte, []byte) (int, error) {
	return func(placement MediaPlacement, mimeType string, videoData []byte, posterData []byte) (int, error) {
		return insert(
			placement.SlideIndex,
			videoData,
			posterData,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptVideoPathInsert(insert VideoPathInsertFn) func(MediaPlacement, string, string, string) (int, error) {
	return func(placement MediaPlacement, mimeType string, videoPath string, posterPath string) (int, error) {
		return insert(
			placement.SlideIndex,
			videoPath,
			posterPath,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptOLEBinaryInsert(insert OLEBinaryInsertFn) func(MediaPlacement, string, []byte, []byte) (int, error) {
	return func(placement MediaPlacement, progID string, objectData []byte, iconData []byte) (int, error) {
		return insert(
			placement.SlideIndex,
			objectData,
			iconData,
			progID,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptOLEPathInsert(insert OLEPathInsertFn) func(MediaPlacement, string, string, string) (int, error) {
	return func(placement MediaPlacement, progID string, objectPath string, iconPath string) (int, error) {
		return insert(
			placement.SlideIndex,
			objectPath,
			iconPath,
			progID,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func NewVideoInsertSpec(
	maxLen int,
	insertBinary func(MediaPlacement, string, []byte, []byte) (int, error),
	insertPath func(MediaPlacement, string, string, string) (int, error),
) MediaInsertSpec {
	return MediaInsertSpec{
		MetaKey:          "mime_type",
		PrimaryPathKey:   "path",
		PrimaryDataKey:   "data",
		SecondaryPathKey: "poster_path",
		SecondaryDataKey: "poster_data",
		PrimaryMaxLen:    maxLen,
		SecondaryMaxLen:  maxLen,
		PrimaryLabel:     "video",
		SecondaryLabel:   "poster",
		InsertBinary:     insertBinary,
		InsertPath:       insertPath,
	}
}

func NewOLEInsertSpec(
	maxLen int,
	insertBinary func(MediaPlacement, string, []byte, []byte) (int, error),
	insertPath func(MediaPlacement, string, string, string) (int, error),
) MediaInsertSpec {
	return MediaInsertSpec{
		MetaKey:          "prog_id",
		PrimaryPathKey:   "path",
		PrimaryDataKey:   "data",
		SecondaryPathKey: "icon_path",
		SecondaryDataKey: "icon_data",
		PrimaryMaxLen:    maxLen,
		SecondaryMaxLen:  maxLen,
		PrimaryLabel:     "object",
		SecondaryLabel:   "icon",
		InsertBinary:     insertBinary,
		InsertPath:       insertPath,
	}
}

type RequireIntFn func(map[string]any, string) (int, bool)
type RequireFloatFn func(map[string]any, string) (float64, bool)
type IndexBoundsFn func(int, int, int, string) bool
type ParseRawPayloadFn func([]byte) (map[string]any, error)
type BuildResultFn func(int) any

func ExecuteMediaInsert(
	payload map[string]any,
	placement MediaPlacement,
	optionalString func(map[string]any, string) string,
	spec MediaInsertSpec,
) (int, error) {
	metaValue := optionalString(payload, spec.MetaKey)
	primaryPath := optionalString(payload, spec.PrimaryPathKey)
	primaryBase64 := optionalString(payload, spec.PrimaryDataKey)
	secondaryPath := optionalString(payload, spec.SecondaryPathKey)
	secondaryBase64 := optionalString(payload, spec.SecondaryDataKey)

	primaryData, err := DecodeOptionalBase64Field(primaryBase64, spec.PrimaryMaxLen, spec.PrimaryLabel)
	if err != nil {
		return 0, err
	}
	secondaryData, err := DecodeOptionalBase64Field(secondaryBase64, spec.SecondaryMaxLen, spec.SecondaryLabel)
	if err != nil {
		return 0, err
	}

	return insertShapeFromBinaryOrPath(
		len(primaryData) > 0 || len(secondaryData) > 0,
		func() (int, error) {
			return spec.InsertBinary(placement, metaValue, primaryData, secondaryData)
		},
		func() (int, error) {
			return spec.InsertPath(placement, metaValue, primaryPath, secondaryPath)
		},
	)
}

func HandleMediaInsertCommand(
	rawPayload []byte,
	slideCount int,
	parseRawPayload ParseRawPayloadFn,
	requireInt RequireIntFn,
	requireFloat RequireFloatFn,
	indexBounds IndexBoundsFn,
	optionalString func(map[string]any, string) string,
	validationErr func() error,
	buildResult BuildResultFn,
	spec MediaInsertSpec,
) (any, error) {
	p, err := parseRawPayload(rawPayload)
	if err != nil {
		return nil, err
	}

	placement, ok := ParseMediaPlacement(p, slideCount, requireInt, requireFloat, indexBounds)
	if !ok {
		return nil, validationErr()
	}

	shapeID, err := ExecuteMediaInsert(p, placement, optionalString, spec)
	if err != nil {
		return nil, err
	}
	return buildResult(shapeID), nil
}

func ParseMediaPlacement(
	payload map[string]any,
	slideCount int,
	requireInt RequireIntFn,
	requireFloat RequireFloatFn,
	indexBounds IndexBoundsFn,
) (MediaPlacement, bool) {
	slideIndex, ok := requireInt(payload, "slide_index")
	if !ok {
		return MediaPlacement{}, false
	}
	x, ok := requireFloat(payload, "x")
	if !ok {
		return MediaPlacement{}, false
	}
	y, ok := requireFloat(payload, "y")
	if !ok {
		return MediaPlacement{}, false
	}
	w, ok := requireFloat(payload, "w")
	if !ok {
		return MediaPlacement{}, false
	}
	h, ok := requireFloat(payload, "h")
	if !ok {
		return MediaPlacement{}, false
	}
	if !indexBounds(slideIndex, 0, slideCount, "slide_index") {
		return MediaPlacement{}, false
	}
	return MediaPlacement{
		SlideIndex: slideIndex,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
	}, true
}

func DecodeOptionalBase64Field(raw string, maxLen int, label string) ([]byte, error) {
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
