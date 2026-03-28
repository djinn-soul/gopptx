// examples/51-language-localization/main.go demonstrates BCP-47 language tags on
// text runs, RTL paragraph direction, and multi-script Unicode content.
//
// Run with: go run ./examples/51-language-localization/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "51_language_localization.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// --- Slide 1: Multi-language runs with BCP-47 language tags ---
	// Each run carries its own language tag so PowerPoint can apply
	// correct spell-checking and hyphenation per script.
	arabicRun := pptx.NewRun("مرحبا بالعالم").
		WithLang("ar-SA")

	frenchRun := pptx.NewRun("Bonjour le monde").
		WithLang("fr-FR")

	japaneseRun := pptx.NewRun("こんにちは世界").
		WithLang("ja-JP")

	chineseRun := pptx.NewRun("你好世界").
		WithLang("zh-CN")

	russianRun := pptx.NewRun("Привет, мир!").
		WithLang("ru-RU")

	slide1 := pptx.NewSlide("Language Tags per Run").
		AddBulletRuns([]pptx.Run{arabicRun}).
		AddBulletRuns([]pptx.Run{frenchRun}).
		AddBulletRuns([]pptx.Run{japaneseRun}).
		AddBulletRuns([]pptx.Run{chineseRun}).
		AddBulletRuns([]pptx.Run{russianRun})

	// --- Slide 2: RTL paragraph direction for right-to-left scripts ---
	rtlStyle := text.NewParagraphStyle().WithRTL(true)

	slide2 := pptx.NewSlide("RTL Paragraph Direction").
		AddBulletRunsWithStyle(
			[]pptx.Run{pptx.NewRun("مرحبا بالعالم").WithLang("ar-SA")},
			rtlStyle,
		).
		AddBulletRunsWithStyle(
			[]pptx.Run{pptx.NewRun("שלום עולם").WithLang("he-IL")},
			rtlStyle,
		).
		AddBulletRunsWithStyle(
			[]pptx.Run{pptx.NewRun("سلام دنیا").WithLang("fa-IR")},
			rtlStyle,
		).
		AddBullet("Left-to-right fallback (English)")

	// --- Slide 3: Document-level language set via core properties ---
	slide3 := pptx.NewSlide("Document Language").
		AddBullet("Document language is set via CoreProperties.Language").
		AddBullet("Individual runs override the document default").
		AddBullet("PowerPoint uses run language for spell-check & hyphenation").
		AddBullet("Omitting WithLang() defaults to en-US")

	meta := pptx.Metadata{
		Metadata: pptx.MetadataFields{
			Title:          "Task 51: Language Localization",
			Creator:        "gopptx",
			CoreProperties: common.CoreProperties{Language: "en-US"},
		},
	}

	data, err := pptx.CreateWithMetadata(meta, []pptx.SlideContent{slide1, slide2, slide3})
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
