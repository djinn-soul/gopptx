package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

// runLegacyFlags preserves the historical flag-only mode:
// -md for markdown conversion, otherwise generates a baseline deck.
func runLegacyFlags(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("legacy", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var outPath string
	var markdownPath string
	var title string
	fs.StringVar(&outPath, "out", "output.pptx", "output PPTX file")
	fs.StringVar(&markdownPath, "md", "", "markdown file to convert to PPTX")
	fs.StringVar(&title, "title", "", "presentation title (required with -md)")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid legacy arguments: %v", err)
		printRootUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printRootUsage(stderr)
		return exitUsage
	}

	if strings.TrimSpace(markdownPath) != "" {
		if strings.TrimSpace(title) == "" {
			printErrorf(stderr, "title is required when using -md")
			return exitUsage
		}
		markdown, err := os.ReadFile(markdownPath)
		if err != nil {
			printErrorf(stderr, "failed to read markdown %q: %v", markdownPath, err)
			return exitIO
		}
		slides, err := pptx.SlidesFromMarkdown(string(markdown))
		if err != nil {
			printErrorf(stderr, "markdown parse failed: %v", err)
			return exitParse
		}
		data, err := pptx.CreateWithSlides(strings.TrimSpace(title), slides)
		if err != nil {
			printErrorf(stderr, "pptx generation failed: %v", err)
			return exitInternal
		}
		if err := writeOutputFile(outPath, data); err != nil {
			printErrorf(stderr, "failed to write %q: %v", outPath, err)
			return exitIO
		}
		_, _ = fmt.Fprintf(
			stdout,
			"OK: wrote %s from %s (%d slide(s))\n",
			strings.TrimSpace(outPath),
			markdownPath,
			len(slides),
		)
		return exitOK
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Welcome").AddBullet("Ported from ppt-rs").AddBullet("Go baseline ready"),
		pptx.NewSlide("Next Steps").AddBullet("Add markdown parser").AddBullet("Add tables/charts/images"),
	}
	data, err := pptx.CreateWithSlides("Go PPTX", slides)
	if err != nil {
		printErrorf(stderr, "pptx generation failed: %v", err)
		return exitInternal
	}
	if err := writeOutputFile(outPath, data); err != nil {
		printErrorf(stderr, "failed to write %q: %v", outPath, err)
		return exitIO
	}
	_, _ = fmt.Fprintf(stdout, "OK: wrote %s (%d slide(s))\n", strings.TrimSpace(outPath), len(slides))
	return exitOK
}
