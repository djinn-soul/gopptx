//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func main() {
	pptxPath := `E:\Github\gopptx\examples\output\gopptx_mermaid_full_demo_v4.pptx`
	outPath := `E:\Github\gopptx\examples\output\test_native_images.pdf`

	title, slides, err := export.SlidesFromPPTX(pptxPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SlidesFromPPTX error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Title: %q, Slides: %d\n", title, len(slides))
	fmt.Println("Note: Mermaid diagrams are written as vector shapes in PPTX, so slide image counts can be zero.")
	for i, s := range slides {
		fmt.Printf("  Slide %d: title=%q bullets=%d shapes=%d images=%d\n",
			i+1, s.Title, len(s.Bullets), len(s.Shapes), len(s.Images))
		for j, img := range s.Images {
			fmt.Printf("    Image %d: format=%s bytes=%d\n", j+1, img.Format, len(img.Data))
		}
	}

	if err := export.PDFFromFile(pptxPath, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "PDFFromFile error: %v\n", err)
		os.Exit(1)
	}
	fi, _ := os.Stat(outPath)
	fmt.Printf("OK: %s (%d bytes)\n", outPath, fi.Size())
}
