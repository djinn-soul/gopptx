package chart

import "testing"

func TestExtractChartState(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<c:style val="10"/>
<c:spPr><a:scene3d><a:camera prst="orthographicFront" fov="45"/><a:lightRig rig="threePt" dir="t" rev="1"/></a:scene3d></c:spPr>
<c:chart><c:plotArea>
<c:barChart>
<c:ser><c:tx><c:strRef><c:strCache><c:pt idx="0"><c:v>North</c:v></c:pt></c:strCache></c:strRef></c:tx><c:val><c:numRef><c:numCache><c:pt idx="0"><c:v>1.5</c:v></c:pt><c:pt idx="1"><c:v>2.5</c:v></c:pt></c:numCache></c:numRef></c:val></c:ser>
<c:axId val="1"/><c:axId val="2"/>
</c:barChart>
<c:catAx><c:axId val="1"/><c:tickLblPos val="low"/><c:crosses val="autoZero"/></c:catAx>
<c:valAx><c:axId val="2"/><c:majorGridlines/><c:minorGridlines/><c:tickLblPos val="nextTo"/><c:crosses val="autoZero"/></c:valAx>
</c:plotArea></c:chart></c:chartSpace>`)

	state := ExtractChartState(xml)
	if state.ChartStyle == nil || *state.ChartStyle != 10 {
		t.Fatalf("expected chart style 10, got %#v", state.ChartStyle)
	}
	if !state.CategoryAx.Present || state.CategoryAx.TickLabelPos != "low" {
		t.Fatalf("unexpected category axis state %#v", state.CategoryAx)
	}
	if state.CategoryAx.Crosses != "autoZero" {
		t.Fatalf("expected category axis crosses autoZero, got %#v", state.CategoryAx.Crosses)
	}
	if !state.ValueAx.Present || state.ValueAx.TickLabelPos != "nextTo" || !state.ValueAx.MajorGridline {
		t.Fatalf("unexpected value axis state %#v", state.ValueAx)
	}
	if !state.ValueAx.MinorGridline {
		t.Fatalf("expected value axis minor gridline true, got %#v", state.ValueAx)
	}
	if state.ValueAx.Crosses != "autoZero" {
		t.Fatalf("expected value axis crosses autoZero, got %#v", state.ValueAx.Crosses)
	}
	if len(state.Series) != 1 || len(state.Series[0].Values) != 2 {
		t.Fatalf("unexpected series state %#v", state.Series)
	}
	if state.Scene3D.CameraPreset != "orthographicFront" || state.Scene3D.LightRig != "threePt" {
		t.Fatalf("unexpected scene3d state %#v", state.Scene3D)
	}
	if state.Scene3D.CameraFieldOfView != 45 ||
		state.Scene3D.LightDirection != "t" ||
		!state.Scene3D.LightRigRevolution {
		t.Fatalf("unexpected scene3d values %#v", state.Scene3D)
	}
}
