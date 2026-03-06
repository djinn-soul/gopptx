package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestHandleSetPlaceholderContentSupportsTablePayload(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-table-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{
		"slide_index": 0,
		"ph_index": 1,
		"table": {
			"rows": [
				["Quarter", "Revenue"],
				["Q1", "125"]
			]
		}
	}`)
	if _, err := handleSetPlaceholderContent(editor, payload); err != nil {
		t.Fatalf("handleSetPlaceholderContent: %v", err)
	}

	content, ok := editor.parts.Get(editor.slides[0].Part)
	if !ok {
		t.Fatal("expected updated slide part")
	}
	slideXML := string(content)
	if !strings.Contains(slideXML, `<p:graphicFrame`) || !strings.Contains(slideXML, `<a:tbl>`) {
		t.Fatalf("expected placeholder replacement to render a table, got: %s", slideXML)
	}
	if !strings.Contains(slideXML, `<p:ph idx="1" type="body"/>`) {
		t.Fatalf("expected table placeholder to preserve index/type, got: %s", slideXML)
	}
	if !strings.Contains(slideXML, "<a:t>Quarter</a:t>") || !strings.Contains(slideXML, "<a:t>125</a:t>") {
		t.Fatalf("expected table cell text in output xml, got: %s", slideXML)
	}
}

func TestHandleSetPlaceholderContentSupportsChartPayload(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-chart-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{
		"slide_index": 0,
		"ph_index": 1,
		"chart": {
			"chart_type": "bar",
			"title": "Quarterly Revenue",
			"categories": ["Q1", "Q2"],
			"values": [10, 20]
		}
	}`)
	if _, err := handleSetPlaceholderContent(editor, payload); err != nil {
		t.Fatalf("handleSetPlaceholderContent: %v", err)
	}

	content, ok := editor.parts.Get(editor.slides[0].Part)
	if !ok {
		t.Fatal("expected updated slide part")
	}
	slideXML := string(content)
	if !strings.Contains(slideXML, `<c:chart`) || !strings.Contains(slideXML, `<p:graphicFrame`) {
		t.Fatalf("expected placeholder replacement to render chart frame, got: %s", slideXML)
	}
	if !strings.Contains(slideXML, `<p:ph idx="1" type="body"/>`) {
		t.Fatalf("expected chart placeholder to preserve index/type, got: %s", slideXML)
	}

	slideRelsPath := common.SlideRelsPartName(editor.slides[0].Part)
	relsData, ok := editor.parts.Get(slideRelsPath)
	if !ok {
		t.Fatalf("expected slide relationships part %s", slideRelsPath)
	}
	relsXML := string(relsData)
	if !strings.Contains(relsXML, "/relationships/chart") || !strings.Contains(relsXML, `../charts/chart`) {
		t.Fatalf("expected slide chart relationship, got: %s", relsXML)
	}

	chartParts := editor.parts.KeysWithPrefix("ppt/charts/chart")
	if len(chartParts) != 1 {
		t.Fatalf("expected one chart part, got %d (%v)", len(chartParts), chartParts)
	}
	chartXML, ok := editor.parts.Get(chartParts[0])
	if !ok {
		t.Fatalf("expected chart part %s", chartParts[0])
	}
	if !strings.Contains(string(chartXML), "<c:barChart>") {
		t.Fatalf("expected bar chart xml, got: %s", string(chartXML))
	}

	chartRelsPath := common.RelsPathFor(chartParts[0])
	chartRelsData, ok := editor.parts.Get(chartRelsPath)
	if !ok {
		t.Fatalf("expected chart relationships part %s", chartRelsPath)
	}
	if !strings.Contains(string(chartRelsData), "/relationships/package") ||
		!strings.Contains(string(chartRelsData), "../embeddings/") {
		t.Fatalf("expected chart embedding relationship, got: %s", string(chartRelsData))
	}

	embeddingParts := editor.parts.KeysWithPrefix("ppt/embeddings/")
	if len(embeddingParts) != 1 {
		t.Fatalf("expected one embedding part, got %d (%v)", len(embeddingParts), embeddingParts)
	}
}

func TestHandleSetPlaceholderContentSupportsDoughnutChartPayload(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-doughnut-chart-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{
		"slide_index": 0,
		"ph_index": 1,
		"chart": {
			"chart_type": "doughnut",
			"title": "Mix",
			"categories": ["A", "B"],
			"values": [30, 70]
		}
	}`)
	if _, err := handleSetPlaceholderContent(editor, payload); err != nil {
		t.Fatalf("handleSetPlaceholderContent: %v", err)
	}

	chartParts := editor.parts.KeysWithPrefix("ppt/charts/chart")
	if len(chartParts) != 1 {
		t.Fatalf("expected one chart part, got %d (%v)", len(chartParts), chartParts)
	}
	chartXML, ok := editor.parts.Get(chartParts[0])
	if !ok {
		t.Fatalf("expected chart part %s", chartParts[0])
	}
	if !strings.Contains(string(chartXML), "<c:doughnutChart>") {
		t.Fatalf("expected doughnut chart xml, got: %s", string(chartXML))
	}
}

func TestHandleSetPlaceholderContentSupportsScatterChartPayload(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-scatter-chart-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{
		"slide_index": 0,
		"ph_index": 1,
		"chart": {
			"chart_type": "scatter",
			"title": "Points",
			"x_values": [1, 2, 3],
			"y_values": [2, 4, 8]
		}
	}`)
	if _, err := handleSetPlaceholderContent(editor, payload); err != nil {
		t.Fatalf("handleSetPlaceholderContent: %v", err)
	}

	chartParts := editor.parts.KeysWithPrefix("ppt/charts/chart")
	if len(chartParts) != 1 {
		t.Fatalf("expected one chart part, got %d (%v)", len(chartParts), chartParts)
	}
	chartXML, ok := editor.parts.Get(chartParts[0])
	if !ok {
		t.Fatalf("expected chart part %s", chartParts[0])
	}
	if !strings.Contains(string(chartXML), "<c:scatterChart>") {
		t.Fatalf("expected scatter chart xml, got: %s", string(chartXML))
	}
}

func TestHandleSetPlaceholderContentSupportsBubbleChartPayload(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-bubble-chart-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{
		"slide_index": 0,
		"ph_index": 1,
		"chart": {
			"chart_type": "bubble",
			"title": "Bubbles",
			"x_values": [1, 2],
			"y_values": [3, 4],
			"sizes": [5, 10]
		}
	}`)
	if _, err := handleSetPlaceholderContent(editor, payload); err != nil {
		t.Fatalf("handleSetPlaceholderContent: %v", err)
	}

	chartParts := editor.parts.KeysWithPrefix("ppt/charts/chart")
	if len(chartParts) != 1 {
		t.Fatalf("expected one chart part, got %d (%v)", len(chartParts), chartParts)
	}
	chartXML, ok := editor.parts.Get(chartParts[0])
	if !ok {
		t.Fatalf("expected chart part %s", chartParts[0])
	}
	if !strings.Contains(string(chartXML), "<c:bubbleChart>") {
		t.Fatalf("expected bubble chart xml, got: %s", string(chartXML))
	}
}

func TestHandleSetPlaceholderContentRejectsMultipleContentKinds(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-mixed-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{
		"slide_index": 0,
		"ph_index": 1,
		"text": "Hello",
		"image_path": "logo.png"
	}`)
	_, err = handleSetPlaceholderContent(editor, payload)
	if err == nil {
		t.Fatal("expected validation error for mixed content kinds")
	}
	if !strings.Contains(err.Error(), "Only one placeholder content kind is allowed") {
		t.Fatalf("expected single-kind validation error, got: %v", err)
	}
}

func TestHandleSetPlaceholderContentRejectsMissingContentKind(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-bridge-missing-content.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	payload := []byte(`{"slide_index":0,"ph_index":1}`)
	_, err = handleSetPlaceholderContent(editor, payload)
	if err == nil {
		t.Fatal("expected validation error for missing content kind")
	}
	if !strings.Contains(err.Error(), "Must provide exactly one") {
		t.Fatalf("expected missing-kind validation error, got: %v", err)
	}
}
