package pptxxml

// SmartArtDataXML renders the dgm:dataModel XML (dataX.xml).
func SmartArtDataXML(spec SmartArtSpec) string {
	return renderSmartArtDataFromTemplate(spec)
}

// SmartArtLayoutXML renders dgm:layoutDef (layoutX.xml) from template.
func SmartArtLayoutXML(layoutURI, _ string) string {
	return renderSmartArtLayoutFromTemplate(layoutURI)
}

// SmartArtColorsXML renders dgm:colorsDef (colorsX.xml) from template.
func SmartArtColorsXML(colorStyleID string) string {
	return renderSmartArtColorsFromTemplate(colorStyleID)
}

// SmartArtStyleXML renders dgm:styleDef (quickStyleX.xml) from template.
func SmartArtStyleXML(quickStyleID string) string {
	return renderSmartArtStyleFromTemplate(quickStyleID)
}

// SmartArtDrawingXML renders dsp:drawing (drawingX.xml) from template.
func SmartArtDrawingXML(spec SmartArtSpec) string {
	return renderSmartArtDrawingFromTemplate(spec)
}

func defaultColorStyleID(id string) string {
	if id != "" {
		return id
	}
	return "urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"
}

func defaultQuickStyleID(id string) string {
	if id != "" {
		return id
	}
	return "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"
}
