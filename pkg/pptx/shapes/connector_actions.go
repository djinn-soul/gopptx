package shapes

import "github.com/djinn-soul/gopptx/pkg/pptx/action"

// WithClickAction adds a click behavior to the connector.
func (c Connector) WithClickAction(link action.Hyperlink) Connector {
	c.ClickAction = &link
	return c
}

// WithHoverAction adds a hover behavior to the connector.
func (c Connector) WithHoverAction(link action.Hyperlink) Connector {
	c.HoverAction = &link
	return c
}
