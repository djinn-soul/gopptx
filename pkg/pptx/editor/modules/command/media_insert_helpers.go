package command

import (
	"encoding/base64"
	"fmt"
	"strings"
)

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
