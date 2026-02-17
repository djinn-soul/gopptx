package smartart

import (
	"fmt"
	"regexp"
)

var hexColorRE = regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)

// Validate checks the SmartArt diagram for consistency.
func (sa SmartArt) Validate(slideIndex int) error {
	if len(sa.Nodes) == 0 {
		return fmt.Errorf("slide %d: SmartArt requires at least one node", slideIndex)
	}
	if sa.CX <= 0 || sa.CY <= 0 {
		return fmt.Errorf("slide %d: SmartArt size must be positive (cx=%d, cy=%d)", slideIndex, sa.CX, sa.CY)
	}
	if sa.X < 0 || sa.Y < 0 {
		return fmt.Errorf("slide %d: SmartArt position must be non-negative (x=%d, y=%d)", slideIndex, sa.X, sa.Y)
	}
	for i, node := range sa.Nodes {
		if err := validateNode(node, slideIndex, i); err != nil {
			return err
		}
	}
	return nil
}

func validateNode(n Node, slideIndex, nodeIndex int) error {
	if n.Text == "" {
		return fmt.Errorf("slide %d: SmartArt node %d has empty text", slideIndex, nodeIndex)
	}
	if n.Color != "" && !hexColorRE.MatchString(n.Color) {
		return fmt.Errorf("slide %d: SmartArt node %d has invalid color %q (expected 6-digit hex)", slideIndex, nodeIndex, n.Color)
	}
	for i, child := range n.Children {
		if err := validateNode(child, slideIndex, i); err != nil {
			return err
		}
	}
	return nil
}
