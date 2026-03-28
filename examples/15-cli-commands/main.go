// examples/15-cli-commands/main.go builds a reference presentation documenting
// the gopptx CLI subcommands.
//
// No CLI execution is performed — this example creates a PPTX that serves as
// a quick-reference guide for the available commands and their purpose.
//
// Run with: go run ./examples/15-cli-commands/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "15_cli_commands.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	slides := []pptx.SlideContent{
		// Slide 1: title / overview
		pptx.NewSlide("gopptx CLI Reference").
			AddBullet("The gopptx binary exposes several subcommands for common PPTX tasks.").
			AddBullet("Each subsequent slide documents one subcommand.").
			AddBullet("Run `gopptx --help` to see all available commands."),

		// Slide 2: create
		pptx.NewSlide("gopptx create").
			AddBullet("Creates a new blank PPTX presentation.").
			AddBullet("Usage:  gopptx create --title \"My Deck\" --slides 5 -o out.pptx").
			AddBullet("Flags:").
			AddBullet("  --title   Presentation title (required)").
			AddBullet("  --slides  Number of blank slides to generate").
			AddBullet("  -o        Output file path"),

		// Slide 3: md2ppt
		pptx.NewSlide("gopptx md2ppt").
			AddBullet("Converts a Markdown file into a PPTX presentation.").
			AddBullet("Usage:  gopptx md2ppt input.md -o out.pptx").
			AddBullet("Each top-level heading becomes a slide title.").
			AddBullet("Bullet lists under a heading become slide bullets.").
			AddBullet("Code blocks are rendered as monospace text boxes."),

		// Slide 4: info
		pptx.NewSlide("gopptx info").
			AddBullet("Displays metadata and structural information about a PPTX file.").
			AddBullet("Usage:  gopptx info presentation.pptx").
			AddBullet("Output includes:").
			AddBullet("  - Title, author, and creation date").
			AddBullet("  - Slide count and slide sizes").
			AddBullet("  - Embedded media and chart counts"),

		// Slide 5: validate
		pptx.NewSlide("gopptx validate").
			AddBullet("Validates the structural integrity of a PPTX file.").
			AddBullet("Usage:  gopptx validate presentation.pptx").
			AddBullet("Checks for:").
			AddBullet("  - Missing or malformed XML parts").
			AddBullet("  - Broken relationship references").
			AddBullet("  - Invalid media or chart entries").
			AddBullet("Exits with a non-zero code if issues are found."),

		// Slide 6: merge
		pptx.NewSlide("gopptx merge").
			AddBullet("Merges two or more PPTX files into a single presentation.").
			AddBullet("Usage:  gopptx merge a.pptx b.pptx -o merged.pptx").
			AddBullet("Slides are appended in the order the input files are given.").
			AddBullet("Themes and masters from the first file are preserved.").
			AddBullet("Assets (images, charts) are deduplicated automatically."),

		// Slide 7: version
		pptx.NewSlide("gopptx version").
			AddBullet("Prints the current gopptx version and build information.").
			AddBullet("Usage:  gopptx version").
			AddBullet("Output example:").
			AddBullet("  gopptx v1.2.3 (commit abc1234, built 2025-01-01)").
			AddBullet("Useful for confirming the installed binary in CI pipelines."),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(outputPath, "gopptx CLI Reference", slides); err != nil {
		return fmt.Errorf("write presentation: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
