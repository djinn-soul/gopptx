package pptxxml

import (
	"fmt"
	"strings"
	"testing"
)

func TestNoPlaceholdersRemainAfterInjection(t *testing.T) {
	type tc struct {
		name      string
		layoutURI string
		nodes     int
	}
	tests := []tc{
		{"BasicProcess", "urn:microsoft.com/office/officeart/2005/8/layout/process1", 2},
		{"AccentProcess", "urn:microsoft.com/office/officeart/2005/8/layout/process3", 2},
		{"AlternatingFlow", "urn:microsoft.com/office/officeart/2005/8/layout/hProcess4", 3},
		{"ContinuousBlockProcess", "urn:microsoft.com/office/officeart/2005/8/layout/hProcess9", 2},
		{"BasicCycle", "urn:microsoft.com/office/officeart/2005/8/layout/cycle2", 3},
		{"TextCycle", "urn:microsoft.com/office/officeart/2005/8/layout/cycle1", 3},
		{"BlockCycle", "urn:microsoft.com/office/officeart/2005/8/layout/cycle5", 3},
		{"BasicBlockList", "urn:microsoft.com/office/officeart/2005/8/layout/default", 3},
		{"VerticalBlockList", "urn:microsoft.com/office/officeart/2005/8/layout/vList5", 4},
		{"HorizontalBulletLst", "urn:microsoft.com/office/officeart/2005/8/layout/hList1", 5},
		{"BasicVenn", "urn:microsoft.com/office/officeart/2005/8/layout/venn1", 2},
		{"LinearVenn", "urn:microsoft.com/office/officeart/2005/8/layout/venn3", 2},
		{"StackedVenn", "urn:microsoft.com/office/officeart/2005/8/layout/venn2", 2},
		{"BasicRadial", "urn:microsoft.com/office/officeart/2005/8/layout/radial1", 2},
		{"BasicMatrix", "urn:microsoft.com/office/officeart/2005/8/layout/matrix3", 2},
		{"TitledMatrix", "urn:microsoft.com/office/officeart/2005/8/layout/matrix1", 2},
		{"BasicPyramid", "urn:microsoft.com/office/officeart/2005/8/layout/pyramid1", 2},
		{"InvertedPyramid", "urn:microsoft.com/office/officeart/2005/8/layout/pyramid3", 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nodes := make([]SmartArtNodeSpec, tc.nodes)
			for i := range nodes {
				nodes[i] = SmartArtNodeSpec{Text: fmt.Sprintf("Item %d", i+1)}
			}
			spec := SmartArtSpec{LayoutURI: tc.layoutURI, Nodes: nodes}
			dataXML := SmartArtDataXML(spec)
			remaining := strings.Count(dataXML, `phldr="1"`)
			if remaining > 0 {
				t.Errorf("layout %s with %d nodes still has %d unfilled phldr=\"1\" nodes in data.xml",
					tc.name, tc.nodes, remaining)
			}
		})
	}
}
