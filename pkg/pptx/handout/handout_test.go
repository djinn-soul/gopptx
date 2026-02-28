package handout_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
)

func TestHandoutLayoutSlidesPerPage(t *testing.T) {
	cases := []struct {
		layout handout.HandoutLayout
		want   int
	}{
		{handout.Layout1Up, 1},
		{handout.Layout2Up, 2},
		{handout.Layout3Up, 3},
		{handout.Layout4Up, 4},
		{handout.Layout6Up, 6},
		{handout.Layout9Up, 9},
		{handout.LayoutOutline, 0},
	}
	for _, tc := range cases {
		if got := tc.layout.SlidesPerPage(); got != tc.want {
			t.Errorf("HandoutLayout(%d).SlidesPerPage() = %d, want %d", tc.layout, got, tc.want)
		}
	}
}

func TestHandoutMasterDefaults(t *testing.T) {
	m := handout.New()
	if m.Layout != handout.Layout1Up {
		t.Errorf("expected default layout Layout1Up, got %v", m.Layout)
	}
	if !m.ShowHeader || !m.ShowFooter || !m.ShowDate || !m.ShowPageNumber {
		t.Error("expected all header/footer/date/page-number to be visible by default")
	}
}

func TestHandoutMasterBuilder(t *testing.T) {
	m := handout.New().
		WithLayout(handout.Layout4Up).
		WithHeader("Acme Corp").
		WithFooter("Confidential").
		HideDate()
	if m.Layout != handout.Layout4Up {
		t.Errorf("expected Layout4Up, got %v", m.Layout)
	}
	if m.HeaderText != "Acme Corp" {
		t.Errorf("expected header text 'Acme Corp', got %q", m.HeaderText)
	}
	if m.ShowDate {
		t.Error("expected ShowDate to be false after HideDate()")
	}
	if !m.ShowPageNumber {
		t.Error("expected ShowPageNumber to remain true")
	}
}

func TestHandoutMasterGenerateXML(t *testing.T) {
	m := handout.New()
	xml := m.GenerateXML()
	if !strings.Contains(xml, "p:handoutMaster") {
		t.Error("expected p:handoutMaster element")
	}
	if !strings.Contains(xml, "p:hf") {
		t.Error("expected p:hf element")
	}
	if !strings.Contains(xml, `hdr="1"`) {
		t.Errorf("expected hdr=1 in XML, got:\n%s", xml)
	}
}

func TestHandoutMasterHideAll(t *testing.T) {
	m := handout.New().HideHeader().HideFooter().HideDate().HidePageNumber()
	xml := m.GenerateXML()
	if strings.Contains(xml, `hdr="1"`) || strings.Contains(xml, `ftr="1"`) ||
		strings.Contains(xml, `dt="1"`) || strings.Contains(xml, `sldNum="1"`) {
		t.Errorf("expected all hf attributes to be 0, got:\n%s", xml)
	}
}

func TestRelationshipsXML(t *testing.T) {
	rels := handout.RelationshipsXML(3)
	if !strings.Contains(rels, "theme3.xml") {
		t.Errorf("expected theme3.xml in rels, got:\n%s", rels)
	}
}
