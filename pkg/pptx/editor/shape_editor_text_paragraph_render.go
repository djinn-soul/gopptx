package editor

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func renderParagraphPropsXML(paragraph *common.Paragraph) (string, error) {
	if paragraph == nil {
		return "", nil
	}
	spec, err := editorParagraphToSpec(paragraph)
	if err != nil {
		return "", err
	}
	return pptxxml.BulletParagraphPropsXML(spec), nil
}
