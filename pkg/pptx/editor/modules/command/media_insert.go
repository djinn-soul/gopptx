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

	primaryData, err := decodeOptionalBase64Field(primaryBase64, spec.PrimaryMaxLen, spec.PrimaryLabel)
	if err != nil {
		return 0, err
	}
	secondaryData, err := decodeOptionalBase64Field(secondaryBase64, spec.SecondaryMaxLen, spec.SecondaryLabel)
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
