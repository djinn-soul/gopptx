package slidescatalog

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

type CatalogReader interface {
	ListSlideCharts(slideIndex int) ([]common.SlideChartRef, error)
	GetSlideLayoutRef(slideIndex int) (string, string, error)
	ListSlideLayouts() ([]common.SlideLayoutInfo, error)
	ListSlideMasters() ([]common.SlideMasterInfo, error)
	ListMasterLayouts(masterPart string) ([]common.SlideLayoutInfo, error)
	GetChartState(slideIndex int, selector common.ChartSelector) (common.ChartState, error)
}

type CatalogMutator interface {
	UpdateChartData(slideIndex int, selector common.ChartSelector, req common.ChartDataUpdate) error
	UpdateChartFormatting(slideIndex int, selector common.ChartSelector, req common.ChartFormatUpdate) error
	RebindSlideLayout(slideIndex int, layoutPart string) error
	CloneLayoutMasterFamily(layoutPart string) (common.SlideMasterCloneResult, error)
	AddSlideMaster() (string, error)
	RemoveSlideMaster(masterPart string) error
	AddSlideLayout(masterPart, layoutName string) (string, error)
	RemoveSlideLayout(layoutPart string) error
}

func BuildChartsResponse(charts []common.SlideChartRef) map[string]any {
	return map[string]any{"charts": charts}
}

func BuildLayoutRefResponse(layoutPart, masterPart string) map[string]any {
	return map[string]any{
		"layout_part": layoutPart,
		"master_part": masterPart,
	}
}

func BuildLayoutsResponse(layouts []common.SlideLayoutInfo) map[string]any {
	return map[string]any{"layouts": layouts}
}

func BuildMastersResponse(masters []common.SlideMasterInfo) map[string]any {
	return map[string]any{"masters": masters}
}

func BuildChartStateResponse(state common.ChartState) map[string]any {
	return map[string]any{"state": state}
}

func BuildCloneFamilyResponse(result common.SlideMasterCloneResult) map[string]any {
	return map[string]any{
		"master_part": result.MasterPart,
		"theme_part":  result.ThemePart,
		"layout_map":  result.LayoutMap,
	}
}

func BuildAddedMasterResponse(masterPart string) map[string]any {
	return map[string]any{"master_part": masterPart}
}

func BuildAddedLayoutResponse(layoutPart string) map[string]any {
	return map[string]any{"layout_part": layoutPart}
}

func BuildUpdatedResponse() map[string]bool {
	return map[string]bool{"updated": true}
}

func BuildReboundResponse() map[string]bool {
	return map[string]bool{"rebound": true}
}

func BuildRemovedResponse() map[string]bool {
	return map[string]bool{"removed": true}
}
