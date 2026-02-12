package shapes

// Flowchart shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeFlowChartProcess renders a flowchart process shape.
	ShapeTypeFlowChartProcess = "flowChartProcess"
	// ShapeTypeFlowChartDecision renders a flowchart decision shape.
	ShapeTypeFlowChartDecision = "flowChartDecision"
	// ShapeTypeFlowChartTerminator renders a flowchart terminator shape.
	ShapeTypeFlowChartTerminator = "flowChartTerminator"
	// ShapeTypeFlowChartDocument renders a flowchart document shape.
	ShapeTypeFlowChartDocument = "flowChartDocument"
	// ShapeTypeFlowChartData renders a flowchart data shape (parallelogram).
	ShapeTypeFlowChartData = "flowChartInputOutput"

	// ShapeTypeFlowChartAlternateProcess renders a flowchart alternate process.
	ShapeTypeFlowChartAlternateProcess = "flowChartAlternateProcess"
	// ShapeTypeFlowChartCard renders a flowchart card.
	ShapeTypeFlowChartCard = "flowChartCard"
	// ShapeTypeFlowChartCollate renders a flowchart collate.
	ShapeTypeFlowChartCollate = "flowChartCollate"
	// ShapeTypeFlowChartConnector renders a flowchart connector.
	ShapeTypeFlowChartConnector = "flowChartConnector"
	// ShapeTypeFlowChartDelay renders a flowchart delay.
	ShapeTypeFlowChartDelay = "flowChartDelay"
	// ShapeTypeFlowChartDirectAccessStorage renders a flowchart direct-access storage.
	ShapeTypeFlowChartDirectAccessStorage = "flowChartDirectAccessStorage"
	// ShapeTypeFlowChartDisplay renders a flowchart display.
	ShapeTypeFlowChartDisplay = "flowChartDisplay"
	// ShapeTypeFlowChartExtract renders a flowchart extract.
	ShapeTypeFlowChartExtract = "flowChartExtract"
	// ShapeTypeFlowChartInputOutput renders a flowchart input/output.
	ShapeTypeFlowChartInputOutput = "flowChartInputOutput"
	// ShapeTypeFlowChartInternalStorage renders a flowchart internal storage.
	ShapeTypeFlowChartInternalStorage = "flowChartInternalStorage"
	// ShapeTypeFlowChartMagneticDisk renders a flowchart magnetic disk.
	ShapeTypeFlowChartMagneticDisk = "flowChartMagneticDisk"
	// ShapeTypeFlowChartManualInput renders a flowchart manual input.
	ShapeTypeFlowChartManualInput = "flowChartManualInput"
	// ShapeTypeFlowChartManualOperation renders a flowchart manual operation.
	ShapeTypeFlowChartManualOperation = "flowChartManualOperation"
	// ShapeTypeFlowChartMerge renders a flowchart merge.
	ShapeTypeFlowChartMerge = "flowChartMerge"
	// ShapeTypeFlowChartMultidocument renders a flowchart multidocument.
	ShapeTypeFlowChartMultidocument = "flowChartMultidocument"
	// ShapeTypeFlowChartOffpageConnector renders a flowchart off-page connector.
	ShapeTypeFlowChartOffpageConnector = "flowChartOffpageConnector"
	// ShapeTypeFlowChartOr renders a flowchart OR.
	ShapeTypeFlowChartOr = "flowChartOr"
	// ShapeTypeFlowChartPredefinedProcess renders a flowchart predefined process.
	ShapeTypeFlowChartPredefinedProcess = "flowChartPredefinedProcess"
	// ShapeTypeFlowChartPreparation renders a flowchart preparation.
	ShapeTypeFlowChartPreparation = "flowChartPreparation"
	// ShapeTypeFlowChartPunchedTape renders a flowchart punched tape.
	ShapeTypeFlowChartPunchedTape = "flowChartPunchedTape"
	// ShapeTypeFlowChartSequentialAccessStorage renders a flowchart sequential-access storage.
	ShapeTypeFlowChartSequentialAccessStorage = "flowChartSequentialAccessStorage"
	// ShapeTypeFlowChartSort renders a flowchart sort.
	ShapeTypeFlowChartSort = "flowChartSort"
	// ShapeTypeFlowChartStoredData renders a flowchart stored data.
	ShapeTypeFlowChartStoredData = "flowChartStoredData"
	// ShapeTypeFlowChartSummingJunction renders a flowchart summing junction.
	ShapeTypeFlowChartSummingJunction = "flowChartSummingJunction"
)

func init() {
	for _, t := range []string{
		ShapeTypeFlowChartProcess, ShapeTypeFlowChartDecision,
		ShapeTypeFlowChartTerminator, ShapeTypeFlowChartDocument,
		ShapeTypeFlowChartData,
		ShapeTypeFlowChartAlternateProcess, ShapeTypeFlowChartCard,
		ShapeTypeFlowChartCollate, ShapeTypeFlowChartConnector,
		ShapeTypeFlowChartDelay, ShapeTypeFlowChartDirectAccessStorage,
		ShapeTypeFlowChartDisplay, ShapeTypeFlowChartExtract,
		ShapeTypeFlowChartInputOutput, ShapeTypeFlowChartInternalStorage,
		ShapeTypeFlowChartMagneticDisk, ShapeTypeFlowChartManualInput,
		ShapeTypeFlowChartManualOperation, ShapeTypeFlowChartMerge,
		ShapeTypeFlowChartMultidocument, ShapeTypeFlowChartOffpageConnector,
		ShapeTypeFlowChartOr, ShapeTypeFlowChartPredefinedProcess,
		ShapeTypeFlowChartPreparation, ShapeTypeFlowChartPunchedTape,
		ShapeTypeFlowChartSequentialAccessStorage, ShapeTypeFlowChartSort,
		ShapeTypeFlowChartStoredData, ShapeTypeFlowChartSummingJunction,
	} {
		registerShapeType(t)
	}

	// Flowchart aliases.
	registerShapeAlias("flowchartprocess", ShapeTypeFlowChartProcess)
	registerShapeAlias("flowchart-process", ShapeTypeFlowChartProcess)
	registerShapeAlias("flowchart_process", ShapeTypeFlowChartProcess)
	registerShapeAlias("flowchartdecision", ShapeTypeFlowChartDecision)
	registerShapeAlias("flowchart-decision", ShapeTypeFlowChartDecision)
	registerShapeAlias("flowchart_decision", ShapeTypeFlowChartDecision)
	registerShapeAlias("flowchartterminator", ShapeTypeFlowChartTerminator)
	registerShapeAlias("flowchart-terminator", ShapeTypeFlowChartTerminator)
	registerShapeAlias("flowchart_terminator", ShapeTypeFlowChartTerminator)
	registerShapeAlias("document", ShapeTypeFlowChartDocument)
	registerShapeAlias("flowchartdocument", ShapeTypeFlowChartDocument)
	registerShapeAlias("data", ShapeTypeFlowChartData)
	registerShapeAlias("flowchartdata", ShapeTypeFlowChartData)
	registerShapeAlias("flowchartinputoutput", ShapeTypeFlowChartData)
}
