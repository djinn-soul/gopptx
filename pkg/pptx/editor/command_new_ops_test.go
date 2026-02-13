package editor

import (
	"encoding/json"
	"testing"
)

func TestCommandUpdateChartData(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXML()))

	req := `{"api_version":1,"request_id":"r1","op":"update_chart_data","payload":{"slide_index":0,"chart_selector":{"index":0},"data":{"categories":["A"],"series":[{"values":[2]}]}}}`
	resp := ExecuteCommand(e, req)
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &out); err != nil {
		t.Fatalf("invalid response json: %v", err)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("expected ok response: %s", resp)
	}
}

func TestCommandLayoutOps(t *testing.T) {
	e := newLayoutFixtureEditor(t)

	listResp := ExecuteCommand(e, `{"api_version":1,"request_id":"r1","op":"list_slide_layouts","payload":{}}`)
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(listResp), &out); err != nil {
		t.Fatalf("invalid list response: %v", err)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("expected list success: %s", listResp)
	}

	cloneResp := ExecuteCommand(e, `{"api_version":1,"request_id":"r2","op":"clone_layout_master_family","payload":{"layout_part":"ppt/slideLayouts/slideLayout1.xml"}}`)
	if err := json.Unmarshal([]byte(cloneResp), &out); err != nil {
		t.Fatalf("invalid clone response: %v", err)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("expected clone success: %s", cloneResp)
	}
}
