package editor

import "testing"

const latencyBatchSize = 20

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
