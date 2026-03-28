package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func TestRenderSmartArtPartsParallel_DeterministicPartOrder(t *testing.T) {
	parts := []SmartArtPart{
		{
			partNumber: 2,
			spec: pptxxml.SmartArtSpec{
				LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/process1",
				Nodes: []pptxxml.SmartArtNodeSpec{
					{Text: "A"},
					{Text: "B"},
				},
			},
		},
		{
			partNumber: 1,
			spec: pptxxml.SmartArtSpec{
				LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/default",
				Nodes: []pptxxml.SmartArtNodeSpec{
					{Text: "X"},
					{Text: "Y"},
				},
			},
		},
	}

	rendered, err := renderSmartArtPartsParallel(parts)
	if err != nil {
		t.Fatalf("renderSmartArtPartsParallel() error = %v", err)
	}
	if len(rendered) != 10 {
		t.Fatalf("expected 10 rendered parts, got %d", len(rendered))
	}

	expectedOrder := []string{
		"ppt/diagrams/data1.xml",
		"ppt/diagrams/layout1.xml",
		"ppt/diagrams/colors1.xml",
		"ppt/diagrams/quickStyle1.xml",
		"ppt/diagrams/drawing1.xml",
		"ppt/diagrams/data2.xml",
		"ppt/diagrams/layout2.xml",
		"ppt/diagrams/colors2.xml",
		"ppt/diagrams/quickStyle2.xml",
		"ppt/diagrams/drawing2.xml",
	}
	for i, wantPath := range expectedOrder {
		if rendered[i].path != wantPath {
			t.Fatalf("rendered[%d].path = %q, want %q", i, rendered[i].path, wantPath)
		}
		if strings.TrimSpace(rendered[i].content) == "" {
			t.Fatalf("rendered[%d] content is empty for %q", i, rendered[i].path)
		}
	}
}
