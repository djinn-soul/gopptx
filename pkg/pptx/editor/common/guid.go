package common

import pptxcommon "github.com/djinn-soul/gopptx/pkg/pptx/common"

// NewGUID generates a fresh GUID string in {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX} format.
func NewGUID() (string, error) {
	return pptxcommon.NewGUID()
}
