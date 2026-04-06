package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func writePresentation(pptxPath, pdfPath string, slides []pptx.SlideContent) error {
	data, err := pptx.CreateWithSlides("Go Chart Types Export Demo", slides)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pptxPath, data, 0o600); err != nil {
		return err
	}
	writeLinef("Saved PPTX: %s (%d slides)", pptxPath, len(slides))

	opts := export.PDFOptions{Driver: export.PDFDriverNative}
	if err := export.PDFWithOptions("Go Chart Types Export Demo", slides, pdfPath, opts); err != nil {
		return err
	}
	writeLinef("Saved PDF:  %s", pdfPath)
	return nil
}

func writeLinef(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format+"\n", args...)
}
