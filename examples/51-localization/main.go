package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "51_localization.pptx"
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

	slides := []pptx.SlideContent{
		pptx.NewSlide("Localization Support").
			AddBullet("Set document language via core properties").
			AddBullet("Unicode text rendering across scripts").
			AddBullet("Multi-language content in a single deck"),

		pptx.NewSlide("Unicode Text Samples").
			AddBullet("English: Hello, World!").
			AddBullet("Japanese: こんにちは世界").
			AddBullet("Arabic: مرحبا بالعالم").
			AddBullet("Chinese (Simplified): 你好世界").
			AddBullet("Russian: Привет, мир!").
			AddBullet("Greek: Γεια σου, κόσμε!").
			AddBullet("Korean: 안녕하세요, 세계!"),

		pptx.NewSlide("RTL Language Support").
			AddBullet("Arabic script renders right-to-left").
			AddBullet("Hebrew: שלום עולם").
			AddBullet("Persian: سلام دنیا").
			AddBullet("RTL flag available on Metadata for RTL presentations"),
	}

	meta := pptx.Metadata{
		Metadata: pptx.MetadataFields{
			Title:          "Task 51: Localization",
			Creator:        "gopptx",
			CoreProperties: common.CoreProperties{Language: "en-US"},
		},
	}

	data, err := pptx.CreateWithMetadata(meta, slides)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := outputDir + "/" + outputFile
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
