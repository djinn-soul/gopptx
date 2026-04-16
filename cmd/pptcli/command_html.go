package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func runHTMLCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("html", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inPath string
	var outPath string
	var title string
	var embedImages bool
	var nav bool

	fs.StringVar(&inPath, "in", "", "input file (.pptx/.pptm/.md)")
	fs.StringVar(&outPath, "out", "", "output HTML file (default: built from input name)")
	fs.StringVar(&title, "title", "Presentation", "presentation title")
	fs.BoolVar(&embedImages, "embed", true, "embed images as base64 in HTML")
	fs.BoolVar(&nav, "nav", true, "include navigation shell script in HTML")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid html arguments: %v", err)
		printHTMLUsage(stderr)
		return exitUsage
	}

	inPath = strings.TrimSpace(inPath)
	if inPath == "" {
		printErrorf(stderr, "html requires -in")
		printHTMLUsage(stderr)
		return exitUsage
	}

	inputKind, err := detectPresentationInputKind(inPath)
	if err != nil {
		printErrorf(stderr, "%v", err)
		printHTMLUsage(stderr)
		return exitUsage
	}

	if strings.TrimSpace(outPath) == "" {
		outPath = defaultSiblingFilePath(inPath, "export", ".html")
	}

	resolvedTitle, slides, exitCode := loadSlidesForHTML(inputKind, inPath, title, stderr)
	if exitCode != exitOK {
		return exitCode
	}

	opts := export.DefaultHTMLOptions()
	opts.EmbedImages = embedImages
	opts.IncludeNavigation = nav

	htmlContent := export.HTMLWithOptions(resolvedTitle, slides, opts)
	if err := os.WriteFile(outPath, []byte(htmlContent), 0o600); err != nil {
		printErrorf(stderr, "failed to write HTML to %q: %v", outPath, err)
		return exitIO
	}

	_, _ = fmt.Fprintf(stdout, "OK: Exported to %s\n", outPath)
	return exitOK
}

func loadSlidesForHTML(
	inputKind presentationInputKind,
	inPath string,
	title string,
	stderr io.Writer,
) (string, []pptx.SlideContent, int) {
	switch inputKind {
	case inputKindMarkdown:
		markdown, err := os.ReadFile(inPath)
		if err != nil {
			printErrorf(stderr, "failed to read markdown %q: %v", inPath, err)
			return "", nil, exitIO
		}

		slides, err := pptx.SlidesFromMarkdown(string(markdown))
		if err != nil {
			printErrorf(stderr, "markdown parse failed: %v", err)
			return "", nil, exitParse
		}
		return strings.TrimSpace(title), slides, exitOK
	case inputKindPPTX:
		presTitle, slides, err := export.SlidesFromPPTX(inPath)
		if err != nil {
			printErrorf(stderr, "failed to read PPTX %q: %v", inPath, err)
			return "", nil, exitIO
		}
		return resolveHTMLTitle(title, presTitle), slides, exitOK
	default:
		printErrorf(stderr, "unsupported input type")
		return "", nil, exitUsage
	}
}

func resolveHTMLTitle(requestedTitle, detectedTitle string) string {
	requestedTitle = strings.TrimSpace(requestedTitle)
	detectedTitle = strings.TrimSpace(detectedTitle)
	if requestedTitle != "" && requestedTitle != "Presentation" {
		return requestedTitle
	}
	if detectedTitle != "" {
		return detectedTitle
	}
	if requestedTitle != "" {
		return requestedTitle
	}
	return "Presentation"
}

func printHTMLUsage(w io.Writer) {
	_, _ = fmt.Fprintln(
		w,
		"Usage: pptcli html -in <file.pptx|file.pptm|file.md> [-out deck.html] [-title TITLE] [-embed=true] [-nav=true]",
	)
}
