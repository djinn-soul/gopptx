package protection

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
)

// EncryptAgilePackage wraps a PPTX zip payload into an Office-compatible
// Agile-encrypted package.
func EncryptAgilePackage(zipPayload []byte, password string) ([]byte, error) {
	if len(zipPayload) == 0 {
		return nil, errors.New("zip payload cannot be empty")
	}
	if password == "" {
		return nil, errors.New("encryption password cannot be empty")
	}
	if err := validatePPTXPayload(zipPayload); err != nil {
		return nil, err
	}
	return encryptAgilePackage(zipPayload, password)
}

// CanEncryptAgile reports whether the current runtime can perform compatible
// Agile encryption.
func CanEncryptAgile() bool {
	return canEncryptAgile()
}

func validatePPTXPayload(payload []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if err != nil {
		return fmt.Errorf("invalid pptx zip payload: %w", err)
	}
	required := map[string]bool{
		"[Content_Types].xml":  false,
		"ppt/presentation.xml": false,
	}
	for _, f := range zr.File {
		if _, ok := required[f.Name]; ok {
			required[f.Name] = true
		}
	}
	for name, ok := range required {
		if !ok {
			return fmt.Errorf("invalid pptx payload: missing required part %s", name)
		}
	}
	return nil
}
