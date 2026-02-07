package main

var chartOrder = []string{
	"bar",
	"barHorizontal",
	"barStacked",
	"barStacked100",
	"line",
	"lineMarkers",
	"lineStacked",
	"area",
	"areaStacked",
	"areaStacked100",
	"pie",
	"doughnut",
	"scatter",
	"scatterLines",
	"scatterSmooth",
	"bubble",
	"radar",
	"radarFilled",
	"stockHLC",
	"stockOHLC",
	"combo",
}

var signatureTokens = []string{
	`<c:barChart`,
	`<c:barDir val="bar"/>`,
	`<c:barDir val="col"/>`,
	`<c:grouping val="clustered"/>`,
	`<c:grouping val="stacked"/>`,
	`<c:grouping val="percentStacked"/>`,
	`<c:lineChart`,
	`<c:grouping val="standard"/>`,
	`<c:marker><c:symbol val="circle"/>`,
	`<c:areaChart`,
	`<c:pieChart`,
	`<c:doughnutChart`,
	`<c:holeSize`,
	`<c:scatterChart`,
	`<c:scatterStyle val="marker"/>`,
	`<c:scatterStyle val="lineMarker"/>`,
	`<c:scatterStyle val="smoothMarker"/>`,
	`<c:bubbleChart`,
	`<c:varyColors val="0"/>`,
	`<c:bubbleScale`,
	`<c:bubbleSize`,
	`<c:xVal`,
	`<c:yVal`,
	`<c:radarChart`,
	`<c:radarStyle val="marker"`,
	`<c:radarStyle val="filled"`,
	`<c:stockChart`,
	`<c:tx><c:v>Open</c:v></c:tx>`,
	`<c:tx><c:v>High</c:v></c:tx>`,
	`<c:tx><c:v>Low</c:v></c:tx>`,
	`<c:tx><c:v>Close</c:v></c:tx>`,
}

var requiredTokenOverrides = map[string]map[string]string{
	"bar": {
		`<c:barDir val="bar"/>`: `<c:barDir val="col"/>`,
	},
	"barStacked": {
		`<c:barDir val="bar"/>`:         `<c:barDir val="col"/>`,
		`<c:grouping val="clustered"/>`: `<c:grouping val="stacked"/>`,
	},
	"barStacked100": {
		`<c:barDir val="bar"/>`:         `<c:barDir val="col"/>`,
		`<c:grouping val="clustered"/>`: `<c:grouping val="percentStacked"/>`,
	},
}

type compareResult struct {
	Chart         string
	RefSeries     int
	OurSeries     int
	Required      []string
	Missing       []string
	Pass          bool
	ReferenceOnly bool
}
