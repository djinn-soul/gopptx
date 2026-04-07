package editor

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

const benchBatchSize = 10
const latencyBatchSize = 20

func BenchmarkBridgeExecuteSingleSetSlideTitle(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	req := `{"api_version":1,"request_id":"bench-single","op":"set_slide_title","payload":{"slide_index":0,"title":"Bench"}}`

	for b.Loop() {
		_ = ExecuteCommand(editor, req)
	}
}

func BenchmarkBridgeExecuteSingleSetSlideHidden(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	req := `{"api_version":1,"request_id":"bench-hidden-single","op":"set_slide_hidden","payload":{"slide_index":0,"hidden":true}}`

	for b.Loop() {
		_ = ExecuteCommand(editor, req)
	}
}

func BenchmarkBridgeExecuteBatchSetSlideTitle(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	req := buildBatchSetTitleRequest(benchBatchSize)

	for b.Loop() {
		_ = ExecuteCommand(editor, req)
	}
}

func BenchmarkBridgeLatencySingleOps(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	req := `{"api_version":1,"request_id":"bench-latency-single","op":"set_slide_title","payload":{"slide_index":0,"title":"Bench"}}`
	b.ResetTimer()
	for range b.N {
		for range latencyBatchSize {
			_ = ExecuteCommand(editor, req)
		}
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*latencyBatchSize), "ns/op_single")
}

func BenchmarkBridgeLatencySingleSetSlideHiddenOps(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	req := `{"api_version":1,"request_id":"bench-hidden-latency-single","op":"set_slide_hidden","payload":{"slide_index":0,"hidden":true}}`
	b.ResetTimer()
	for range b.N {
		for range latencyBatchSize {
			_ = ExecuteCommand(editor, req)
		}
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*latencyBatchSize), "ns/op_single")
}

func BenchmarkBridgeLatencyBatchOps(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	req := buildBatchSetTitleRequest(latencyBatchSize)
	b.ResetTimer()
	for range b.N {
		_ = ExecuteCommand(editor, req)
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*latencyBatchSize), "ns/op_effective")
}

func BenchmarkBridgeJSONEnvelopeDecode(b *testing.B) {
	req := []byte(
		`{"api_version":1,"request_id":"bench-json","op":"set_slide_title","payload":{"slide_index":0,"title":"Bench"}}`,
	)
	var envelope RequestEnvelope

	for b.Loop() {
		if err := json.Unmarshal(req, &envelope); err != nil {
			b.Fatalf("unmarshal: %v", err)
		}
	}
}

func BenchmarkBridgeJSONEnvelopeEncode(b *testing.B) {
	resp := ResponseEnvelope{
		OK:        true,
		RequestID: "bench-json",
		Result: map[string]any{
			"slide_count": 2,
			"title":       "Bench",
		},
	}

	for b.Loop() {
		if _, err := json.Marshal(resp); err != nil {
			b.Fatalf("marshal: %v", err)
		}
	}
}

func openBenchEditor(b *testing.B) *PresentationEditor {
	b.Helper()
	path := writeBenchDeck(b)
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		b.Fatalf("open editor: %v", err)
	}
	return editor
}

func writeBenchDeck(b *testing.B) string {
	b.Helper()
	path := filepath.Join(b.TempDir(), "bridge-bench.pptx")
	files := map[string]string{
		"[Content_Types].xml":              `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/><Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/></Types>`,
		"_rels/.rels":                      `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml":             `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels":  `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree><p:title><p:txBody><a:p><a:r><a:t>A</a:t></a:r></a:p></p:txBody></p:title></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`,
	}
	if err := writeZipFixture(path, files); err != nil {
		b.Fatalf("write fixture: %v", err)
	}
	return path
}

func buildBatchSetTitleRequest(size int) string {
	commands := make([]string, 0, size)
	for i := range size {
		commands = append(
			commands,
			fmt.Sprintf(`{"op":"set_slide_title","payload":{"slide_index":0,"title":"Bench %d"}}`, i),
		)
	}
	return `{"api_version":1,"request_id":"bench-batch","op":"batch_execute","payload":{"commands":[` + strings.Join(
		commands,
		",",
	) + `]}}`
}
