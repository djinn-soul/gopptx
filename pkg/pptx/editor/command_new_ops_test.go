package editor

import (
	"encoding/json"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestCommandUpdateChartData(t *testing.T) {
	e := newChartUpdateEditorFixture()
	e.parts.Set("ppt/charts/chart1.xml", []byte(categoryChartXML()))

	req := `{"api_version":1,"request_id":"r1","op":"update_chart_data","payload":{"slide_index":0,"chart_selector":{"index":0},"data":{"categories":["A"],"series":[{"values":[2]}]}}}`
	resp := ExecuteCommand(e, req)
	var out map[string]any
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
	var out map[string]any
	if err := json.Unmarshal([]byte(listResp), &out); err != nil {
		t.Fatalf("invalid list response: %v", err)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("expected list success: %s", listResp)
	}

	cloneResp := ExecuteCommand(
		e,
		`{"api_version":1,"request_id":"r2","op":"clone_layout_master_family","payload":{"layout_part":"ppt/slideLayouts/slideLayout1.xml"}}`,
	)
	if err := json.Unmarshal([]byte(cloneResp), &out); err != nil {
		t.Fatalf("invalid clone response: %v", err)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("expected clone success: %s", cloneResp)
	}
}

func TestCommandSectionOps(t *testing.T) {
	e := &PresentationEditor{
		parts:  NewPartStore(),
		slides: []common.EditorSlideRef{{SlideID: 256}},
	}

	addReq := `{"api_version":1,"op":"add_section","payload":{"name":"Intro","slide_indices":[0]}}`
	resp := ExecuteCommand(e, addReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("add_section failed: %s", resp)
	}

	renameReq := `{"api_version":1,"op":"rename_section","payload":{"old_name":"Intro","new_name":"Introduction"}}`
	resp = ExecuteCommand(e, renameReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("rename_section failed: %s", resp)
	}

	removeReq := `{"api_version":1,"op":"remove_section","payload":{"name":"Introduction"}}`
	resp = ExecuteCommand(e, removeReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("remove_section failed: %s", resp)
	}
}

func TestCommandPropsOps(t *testing.T) {
	e := &PresentationEditor{
		parts: NewPartStore(),
	}

	setReq := `{"api_version":1,"op":"set_core_properties","payload":{"title":"New Title","creator":"Test"}}`
	resp := ExecuteCommand(e, setReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("set_core_properties failed: %s", resp)
	}

	getReq := `{"api_version":1,"op":"get_core_properties","payload":{}}`
	resp = ExecuteCommand(e, getReq)
	if !strings.Contains(resp, `"title":"New Title"`) {
		t.Fatalf("get_core_properties failed or returned wrong title: %s", resp)
	}
}

func TestCommandThemeAndSizeOps(t *testing.T) {
	ps := NewPartStore()
	ps.Set("ppt/theme/theme1.xml", []byte("<xml/>"))
	e := &PresentationEditor{
		parts:           ps,
		presentationXML: `<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:sldSz cx="9144000" cy="6858000"/></p:presentation>`,
	}

	themeReq := `{"api_version":1,"op":"apply_theme","payload":{"theme_name":"Modern"}}`
	resp := ExecuteCommand(e, themeReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("apply_theme failed: %s", resp)
	}

	sizeReq := `{"api_version":1,"op":"set_slide_size","payload":{"width":12192000,"height":6858000}}`
	resp = ExecuteCommand(e, sizeReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("set_slide_size failed: %s", resp)
	}
}
