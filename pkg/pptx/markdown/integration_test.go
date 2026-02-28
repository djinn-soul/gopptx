package markdown_test

import (
	"archive/zip"
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/markdown"
)

func TestSlidesFromMarkdown_Integration_GFMTable(t *testing.T) {
	input := `# Data
| Feature | Status |
|---------|--------|
| Tables | Done |
| Mermaid | Done |
`
	slides, err := markdown.SlidesFromMarkdown(input)
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}

	data, err := pptx.CreateWithSlides("Deck", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, "<a:tbl>") {
		t.Fatalf("expected table XML in slide output")
	}
}

func TestSlidesFromMarkdown_Integration_Md2PptDemoFixture(t *testing.T) {
	content, err := os.ReadFile(testutil.RootTestdataPath("md2ppt_demo.md"))
	if err != nil {
		t.Fatalf("read fixture error: %v", err)
	}

	slides, err := markdown.SlidesFromMarkdown(string(content))
	if err != nil {
		t.Fatalf("SlidesFromMarkdown returned error: %v", err)
	}

	data, err := pptx.CreateWithSlides("Fixture Deck", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	foundTable := false
	foundMermaid := false
	for i := 1; i <= len(slides); i++ {
		name := "ppt/slides/slide" + strconv.Itoa(i) + ".xml"
		slideXML := testutil.ReadZipFile(t, zr, name)
		if strings.Contains(slideXML, "<a:tbl>") {
			foundTable = true
		}
		if strings.Contains(slideXML, "Diagram:") || strings.Contains(slideXML, "Web Browser") {
			foundMermaid = true
		}
		if foundTable && foundMermaid {
			break
		}
	}
	if !foundTable {
		t.Fatalf("expected rendered table XML from fixture deck")
	}
	if !foundMermaid {
		t.Fatalf("expected rendered mermaid placeholder XML from fixture deck")
	}
}
