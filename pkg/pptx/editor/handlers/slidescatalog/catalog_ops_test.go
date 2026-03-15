package slidescatalog

import (
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestResponseBuilders(t *testing.T) {
	charts := []common.SlideChartRef{{ChartPart: "ppt/charts/chart1.xml"}}
	if got := BuildChartsResponse(charts); len(got["charts"].([]common.SlideChartRef)) != 1 {
		t.Fatalf("BuildChartsResponse unexpected payload: %+v", got)
	}

	layoutRef := BuildLayoutRefResponse(
		"ppt/slideLayouts/slideLayout1.xml",
		"ppt/slideMasters/slideMaster1.xml",
	)
	if layoutRef["layout_part"] != "ppt/slideLayouts/slideLayout1.xml" {
		t.Fatalf("BuildLayoutRefResponse unexpected layout part: %+v", layoutRef)
	}
	if layoutRef["master_part"] != "ppt/slideMasters/slideMaster1.xml" {
		t.Fatalf("BuildLayoutRefResponse unexpected master part: %+v", layoutRef)
	}

	layouts := []common.SlideLayoutInfo{{Part: "ppt/slideLayouts/slideLayout1.xml", Name: "Title"}}
	if got := BuildLayoutsResponse(layouts); len(got["layouts"].([]common.SlideLayoutInfo)) != 1 {
		t.Fatalf("BuildLayoutsResponse unexpected payload: %+v", got)
	}

	masters := []common.SlideMasterInfo{{Part: "ppt/slideMasters/slideMaster1.xml"}}
	if got := BuildMastersResponse(masters); len(got["masters"].([]common.SlideMasterInfo)) != 1 {
		t.Fatalf("BuildMastersResponse unexpected payload: %+v", got)
	}

	state := common.ChartState{
		CategoryAx: common.ChartAxisState{Present: true, TickLabelPos: "nextTo"},
	}
	if got := BuildChartStateResponse(state); !got["state"].(common.ChartState).CategoryAx.Present {
		t.Fatalf("BuildChartStateResponse unexpected payload: %+v", got)
	}

	clone := common.SlideMasterCloneResult{
		MasterPart: "ppt/slideMasters/slideMaster2.xml",
		ThemePart:  "ppt/theme/theme2.xml",
		LayoutMap: map[string]string{
			"ppt/slideLayouts/slideLayout1.xml": "ppt/slideLayouts/slideLayout9.xml",
		},
	}
	cloneResp := BuildCloneFamilyResponse(clone)
	if cloneResp["master_part"] != clone.MasterPart || cloneResp["theme_part"] != clone.ThemePart {
		t.Fatalf("BuildCloneFamilyResponse unexpected payload: %+v", cloneResp)
	}
	if len(cloneResp["layout_map"].(map[string]string)) != 1 {
		t.Fatalf("BuildCloneFamilyResponse missing layout map: %+v", cloneResp)
	}

	if got := BuildAddedMasterResponse("ppt/slideMasters/slideMaster3.xml"); got["master_part"] != "ppt/slideMasters/slideMaster3.xml" {
		t.Fatalf("BuildAddedMasterResponse unexpected payload: %+v", got)
	}
	if got := BuildAddedLayoutResponse("ppt/slideLayouts/slideLayout2.xml"); got["layout_part"] != "ppt/slideLayouts/slideLayout2.xml" {
		t.Fatalf("BuildAddedLayoutResponse unexpected payload: %+v", got)
	}

	if got := BuildUpdatedResponse(); got["updated"] != true {
		t.Fatalf("BuildUpdatedResponse unexpected payload: %+v", got)
	}
	if got := BuildReboundResponse(); got["rebound"] != true {
		t.Fatalf("BuildReboundResponse unexpected payload: %+v", got)
	}
	if got := BuildRemovedResponse(); got["removed"] != true {
		t.Fatalf("BuildRemovedResponse unexpected payload: %+v", got)
	}
}
