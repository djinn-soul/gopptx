// examples/04-text-styling/main.go demonstrates text formatting features.
//
// Shows how to apply slide-level text styling (bold, italic, underline, size, color)
// to titles and content, matching the legacy ppt-rs behavior.
//
// Run with: go run ./examples/04-text-styling/main.go
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return err
	}

	// 1. large_title.pptx
	prs1 := pptx.NewPresentationBuilder("Large Title Test")
	s1 := pptx.NewSlide("This is a 60pt Title").
		WithTitleSize(60).
		AddBullet("And this is 16pt content.").
		AddBullet("Second point.").
		WithContentSize(16)
	if err := prs1.AddSlide(s1).WriteToFile(filepath.Join(outDir, "04_large_title.pptx")); err != nil {
		return err
	}
	printLine("Wrote 04_large_title.pptx")

	// 2. bold_content.pptx
	prs2 := pptx.NewPresentationBuilder("Bold Content Test")
	s2 := pptx.NewSlide("Bold Title").
		WithTitleBold(true).
		AddBullet("This content should be bold.").
		AddBullet("This one too.").
		WithContentBold(true)
	if err := prs2.AddSlide(s2).WriteToFile(filepath.Join(outDir, "04_bold_content.pptx")); err != nil {
		return err
	}
	printLine("Wrote 04_bold_content.pptx")

	// 3. combined_styling.pptx
	prs3 := pptx.NewPresentationBuilder("Combined Styling Test")
	s3 := pptx.NewSlide("Combined Styling Demo").
		WithTitleColor("FF0000").
		WithTitleItalic(true).
		AddBullet("Content is underlined and blue.").
		WithContentUnderline(true).
		WithContentColor("0000FF")
	if err := prs3.AddSlide(s3).WriteToFile(filepath.Join(outDir, "04_combined_styling.pptx")); err != nil {
		return err
	}
	printLine("Wrote 04_combined_styling.pptx")

	return nil
}

func printLine(args ...any) {
	_, _ = io.WriteString(os.Stdout, fmt.Sprintln(args...))
}
