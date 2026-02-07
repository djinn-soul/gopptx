package pptx

// AddShape appends one custom shape and returns the updated slide.
func (s SlideContent) AddShape(shape ShapeDefinition) SlideContent {
	s.Shapes = append(s.Shapes, shape.ToShape())
	return s
}

// AddConnector appends one connector and returns the updated slide.
func (s SlideContent) AddConnector(connector Connector) SlideContent {
	s.Connectors = append(s.Connectors, connector)
	return s
}
