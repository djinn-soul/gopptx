package pptx_test

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func TestPresentationBuilder(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "builder_test.pptx")

	builder := pptx.NewPresentationBuilder("Fluent Presentation").
		WithMetadata(pptx.PresentationMetadata{PresentationMetadata: common.PresentationMetadata{Creator: "Test Builder"}}).
		AddSlide(pptx.NewSlide("Slide 1").AddShape(pptx.NewRectangle(1, 1, 2, 2))).
		AddSlide(pptx.NewSlide("Slide 2").AddShape(pptx.NewEllipse(3, 1, 2, 2)))

	data, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Build returned empty data")
	}

	if err := builder.WriteToFile(outPath); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}

	emptyBuilder := pptx.NewPresentationBuilder("Empty")
	_, err = emptyBuilder.Build()
	if err == nil {
		t.Errorf("expected error for empty presentation, got nil")
	}
}

func TestCustomXML(t *testing.T) {
	builder := pptx.NewPresentationBuilder("Custom XML Test").
		AddCustomXML("<test>content 1</test>").
		AddCustomXML("<foo>bar</foo>").
		AddSlide(pptx.NewSlide("Slide 1"))

	data, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Build returned empty data")
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	itemIDPattern := regexp.MustCompile(`ds:itemID="\{[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\}"`)
	itemPropsPathPattern := regexp.MustCompile(`^customXml/itemProps\d+\.xml$`)
	foundIDs := make(map[string]struct{})
	propsCount := 0

	for _, f := range zr.File {
		if !itemPropsPathPattern.MatchString(f.Name) {
			continue
		}
		propsCount++

		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		content, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			t.Fatalf("read %s: %v", f.Name, err)
		}

		matches := itemIDPattern.FindAllString(string(content), -1)
		if len(matches) != 1 {
			t.Fatalf("expected exactly one GUID itemID in %s, got %d", f.Name, len(matches))
		}
		foundIDs[matches[0]] = struct{}{}
	}

	if propsCount != 2 {
		t.Fatalf("expected 2 customXml itemProps files, got %d", propsCount)
	}
	if len(foundIDs) != 2 {
		t.Fatalf("expected unique GUID itemIDs for customXml itemProps")
	}
}

func TestSlideNumberingEnabled(t *testing.T) {
	builder := pptx.NewPresentationBuilder("Numbering Test").
		WithSlideNumbers(true).
		AddSlide(pptx.NewSlide("Slide 1"))

	_, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
}

func TestConnectorGranularArrows(t *testing.T) {
	connector := pptx.NewConnector(pptx.ConnectorTypeStraight, 1, 1, 5, 5).
		WithArrows(pptx.ArrowTypeStealth, pptx.ArrowTypeDiamond).
		WithStartArrowSize(pptx.ArrowSizeSmall, pptx.ArrowSizeLarge).
		WithEndArrowSize(pptx.ArrowSizeLarge, pptx.ArrowSizeSmall)

	if connector.StartArrowWidth != pptx.ArrowSizeSmall {
		t.Errorf("expected start arrow width %s, got %s", pptx.ArrowSizeSmall, connector.StartArrowWidth)
	}
	if connector.EndArrowLen != pptx.ArrowSizeSmall {
		t.Errorf("expected end arrow length %s, got %s", pptx.ArrowSizeSmall, connector.EndArrowLen)
	}

	builder := pptx.NewPresentationBuilder("Arrows Test").
		AddSlide(pptx.NewSlide("Slide 1").AddConnector(connector))

	_, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
}
