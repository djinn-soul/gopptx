package editor

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestCommandUpdateChartFormatting_AxisDetails(t *testing.T) {
	path := writeDeckFixture(t, "command-chart-axis-details.pptx", []elements.SlideContent{
		elements.NewSlide("Chart"),
	})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	addResponse := ExecuteCommand(
		editor,
		`{"api_version":1,"request_id":"add","op":"add_chart","payload":{"slide_index":0,"chart_type":"bar","categories":["A","B"],"values":[1,2],"x":0,"y":0,"w":1000,"h":800}}`,
	)
	if !strings.Contains(addResponse, `"ok":true`) {
		t.Fatalf("add chart failed: %s", addResponse)
	}
	updateResponse := ExecuteCommand(
		editor,
		`{"api_version":1,"request_id":"format","op":"update_chart_formatting","payload":{"slide_index":0,"chart_selector":{"index":0},"format":{"category_axis_title":"Quarter","value_axis_minimum_scale":0,"value_axis_maximum_scale":200}}}`,
	)
	if !strings.Contains(updateResponse, `"ok":true`) {
		t.Fatalf("update chart formatting failed: %s", updateResponse)
	}
	stateResponse := ExecuteCommand(
		editor,
		`{"api_version":1,"request_id":"state","op":"get_chart_state","payload":{"slide_index":0,"chart_selector":{"index":0}}}`,
	)
	if !strings.Contains(stateResponse, `"title":"Quarter"`) ||
		!strings.Contains(stateResponse, `"maximum_scale":200`) {
		t.Fatalf("chart state missing axis details: %s", stateResponse)
	}
}
