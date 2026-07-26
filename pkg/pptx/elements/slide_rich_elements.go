package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// AddImage adds an image to the slide.
func (s SlideContent) AddImage(img shapes.Image) SlideContent {
	s.Images = append(s.Images, img)
	return s
}

// WithSlideNumber enables or disables slide number display.
func (s SlideContent) WithSlideNumber(show bool) SlideContent {
	s.ShowSlideNumber = show
	return s
}

// AddAnimation adds an animation to the slide.
func (s SlideContent) AddAnimation(anim animations.AnimationDefinition) SlideContent {
	s.Animations = append(s.Animations, anim.ToAnimation())
	return s
}

// AddSmartArt adds a SmartArt diagram to the slide.
func (s SlideContent) AddSmartArt(sa smartart.SmartArt) SlideContent {
	s.SmartArtDiagrams = append(s.SmartArtDiagrams, sa)
	return s
}

// AddConnector adds a connector to the slide.
func (s SlideContent) AddConnector(connector shapes.Connector) SlideContent {
	s.Connectors = append(s.Connectors, connector)
	return s
}

// AutoRerouteConnectors recalculates connector sites from current shape positions.
func (s SlideContent) AutoRerouteConnectors() SlideContent {
	rerouted := make([]shapes.Connector, 0, len(s.Connectors))
	for _, connector := range s.Connectors {
		rerouted = append(rerouted, connector.AutoReroute(s.Shapes))
	}
	s.Connectors = rerouted
	return s
}
