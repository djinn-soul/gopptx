package pptxxml

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestBackgroundXML(t *testing.T) {
	tests := []struct {
		name     string
		bg       *SlideBackgroundSpec
		contains []string
	}{
		{
			name: "nil background",
			bg:   nil,
		},
		{
			name: "empty type",
			bg:   &SlideBackgroundSpec{Type: ""},
		},
		{
			name: "solid fill",
			bg: &SlideBackgroundSpec{
				Type: "solid",
				SolidFill: &ShapeFillSpec{
					Color: "FF0000",
				},
			},
			contains: []string{"<p:bg>", "<a:solidFill>", "val=\"FF0000\""},
		},
		{
			name: "solid fill with hash",
			bg: &SlideBackgroundSpec{
				Type: "solid",
				SolidFill: &ShapeFillSpec{
					Color: "#00FF00",
				},
			},
			contains: []string{"val=\"00FF00\""},
		},
		{
			name: "picture fill",
			bg: &SlideBackgroundSpec{
				Type: "picture",
				PictureFill: &ImageRef{
					RelID: "rId5",
				},
			},
			contains: []string{"<a:blipFill>", "r:embed=\"rId5\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backgroundXML(tt.bg)
			if tt.contains == nil && got != "" {
				t.Errorf("backgroundXML() = %v, want empty string", got)
			}
			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("backgroundXML() = %v, missing %v", got, s)
				}
			}
		})
	}
}

func TestNormalizeSlideLayoutMode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"titleOnly", slideLayoutTitleOnly},
		{"TITLE_ONLY", slideLayoutTitleOnly},
		{"title-only", slideLayoutTitleOnly},
		{"blank", slideLayoutBlank},
		{"centeredTitle", slideLayoutCenteredTitle},
		{"titleAndBigContent", slideLayoutTitleBigContent},
		{"twoColumn", slideLayoutTwoColumn},
		{"unknown", slideLayoutTitleAndContent},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := normalizeSlideLayoutMode(tt.input); got != tt.want {
				t.Errorf("normalizeSlideLayoutMode(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSlideRelationshipsWithMultiCharts(t *testing.T) {
	tests := []struct {
		name           string
		layoutTarget   string
		imageTargets   []string
		chartRel       *ChartRel
		notesTarget    string
		commentsTarget string
		contains       []string
	}{
		{
			name:         "basic layout",
			layoutTarget: "../slideLayouts/slideLayout1.xml",
			contains:     []string{"Target=\"../slideLayouts/slideLayout1.xml\"", "Id=\"rId1\""},
		},
		{
			name:         "with images",
			layoutTarget: "layout.xml",
			imageTargets: []string{"media/image1.png", "media/audio1.wav"},
			contains: []string{
				"Target=\"media/image1.png\"", "Id=\"rId2\"",
				"Target=\"media/audio1.wav\"", "Id=\"rId3\"",
			},
		},
		{
			name:         "with chart",
			layoutTarget: "layout.xml",
			chartRel:     &ChartRel{RID: "rId10", Target: "charts/chart1.xml"},
			contains:     []string{"Id=\"rId10\"", "Target=\"charts/chart1.xml\""},
		},
		{
			name:           "with notes and comments",
			layoutTarget:   "layout.xml",
			notesTarget:    "notes/notes1.xml",
			commentsTarget: "comments/comment1.xml",
			contains: []string{
				"Target=\"notes/notes1.xml\"",
				"Target=\"comments/comment1.xml\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SlideRelationshipsWithMultiCharts(
				tt.layoutTarget, tt.imageTargets, tt.chartRel, nil, nil, tt.notesTarget, nil, tt.commentsTarget,
			)
			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("SlideRelationshipsWithMultiCharts() = %v, missing %v", got, s)
				}
			}
		})
	}
}

func TestSlideWithLayout(t *testing.T) {
	title := TitleSpec{Text: "Test Title"}
	bullets := []string{"Bullet 1", "Bullet 2"}
	width, height := int64(9144000), int64(6858000)

	got := SlideWithLayout(
		"titleAndContent", title, bullets, nil, nil, ContentStyleSpec{},
		nil, nil, nil, nil, nil, nil, nil, nil,
		"", "", false, "", false, width, height,
	)

	expectedStrings := []string{
		"<p:sld",
		"Test Title",
		"Bullet 1",
		"Bullet 2",
		"cx=\"9144000\"",
		"cy=\"6858000\"",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(got, s) {
			t.Errorf("SlideWithLayout() missing expected string: %v", s)
		}
	}
}

func TestRenderTable(t *testing.T) {
	table := &TableSpec{
		X: 100, Y: 200, CX: 300, CY: 400,
		Rows: [][]string{
			{"A1", "B1"},
			{"A2", "B2"},
		},
		AltText: "Test Table",
	}
	got := RenderTable(table, 1)
	expected := []string{
		"<p:graphicFrame>",
		"id=\"1\"",
		"descr=\"Test Table\"",
		"x=\"100\"", "y=\"200\"",
		"cx=\"300\"", "cy=\"400\"",
		"A1", "B1", "A2", "B2",
	}
	for _, s := range expected {
		if !strings.Contains(got, s) {
			t.Errorf("RenderTable() missing %s", s)
		}
	}
}

func TestRenderTableWithStyledRows(t *testing.T) {
	m10 := int64(10)
	wrap := true
	table := &TableSpec{
		StyledRows: [][]TableCellSpec{
			{
				{
					Text: "Styled", Bold: true, BackgroundColor: "FF0000", Color: "0000FF",
					Align: "ctr", VAlign: "ctr",
					MarginLeft: &m10, WrapText: &wrap,
					BorderLeft: &TableCellBorderSpec{Width: 100, Color: "000000", Dash: "dash"},
				},
			},
		},
	}
	got := RenderTable(table, 1)
	expected := []string{
		"Styled",
		"b=\"1\"",
		"val=\"FF0000\"",
		"val=\"0000FF\"",
		"algn=\"ctr\"",
		"anchor=\"ctr\"",
		"marL=\"10\"",
		"wrap=\"square\"",
		"lnL w=\"100\"",
		"val=\"dash\"",
	}
	for _, s := range expected {
		if !strings.Contains(got, s) {
			t.Errorf("RenderTable() with styles missing %s", s)
		}
	}
}

func TestNotesSlide(t *testing.T) {
	paragraphs := []text.Paragraph{
		{
			Runs: []text.Run{
				{Text: "Note text", Bold: true},
			},
		},
	}
	got := NotesSlide(paragraphs)
	if !strings.Contains(got, "Note text") {
		t.Errorf("NotesSlide() missing 'Note text'")
	}
	if !strings.Contains(got, "b=\"1\"") {
		t.Errorf("NotesSlide() missing bold property")
	}

	// Test empty paragraphs
	gotEmpty := NotesSlide(nil)
	if !strings.Contains(gotEmpty, "<a:p><a:endParaRPr lang=\"en-US\"/></a:p>") {
		t.Errorf("NotesSlide(nil) missing empty paragraph")
	}
}

func TestNotesSlideRelationships(t *testing.T) {
	got := NotesSlideRelationships(5)
	if !strings.Contains(got, "Target=\"../slides/slide5.xml\"") {
		t.Errorf("NotesSlideRelationships() missing slide target")
	}
}

func TestCommentsXMLInternal(t *testing.T) {
	// Just a basic check since it's already covered in pptxxml_test
	// but we want to ensure it's called from within the package too if needed.
	authors := []comments.Author{{ID: 1, Name: "Test"}}
	gotAuthors := CommentAuthorsXML(authors)
	if !strings.Contains(gotAuthors, "Test") {
		t.Errorf("CommentAuthorsXML() missing author name")
	}

	cms := []comments.Comment{{AuthorID: 1, Text: "Comment"}}
	gotCms := CommentsXML(cms)
	if !strings.Contains(gotCms, "Comment") {
		t.Errorf("CommentsXML() missing comment text")
	}
}
