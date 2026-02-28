package presentation

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

func writeVBAParts(pw *pptxxml.PackageWriter, project *vba.VBAProject) error {
	if !project.IsMacroEnabled() {
		return nil
	}
	pw.AddBinaryPart("ppt/vbaProject.bin", project.Data)
	return nil
}
