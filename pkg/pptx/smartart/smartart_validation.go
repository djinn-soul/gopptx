package smartart

import (
	"fmt"
	"regexp"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
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
	if !sa.IsDecorative && len(sa.AltText) > common.MaxAltTextLength {
		return fmt.Errorf("slide %d: SmartArt alt text exceeds %d characters", slideIndex, common.MaxAltTextLength)
	}
	for i, node := range sa.Nodes {
		if err := validateNode(node, slideIndex, i); err != nil {
			return err
		}
	}
	if err := validateLayoutConstraint(sa, slideIndex); err != nil {
		return err
	}
	return nil
}

func validateLayoutConstraint(sa SmartArt, slideIndex int) error {
	constraint, ok := smartArtLayoutConstraints[sa.Layout]
	if !ok {
		return nil
	}

	if constraint.maxItems > 0 {
		if hasNestedChildren(sa.Nodes) {
			return fmt.Errorf(
				"slide %d: SmartArt layout %q only supports flat item lists",
				slideIndex,
				sa.Layout.Name(),
			)
		}
		if len(sa.Nodes) > constraint.maxItems {
			return fmt.Errorf(
				"slide %d: SmartArt layout %q supports at most %d item(s), got %d",
				slideIndex,
				sa.Layout.Name(),
				constraint.maxItems,
				len(sa.Nodes),
			)
		}
		return nil
	}

	if constraint.requireSingleRoot && len(sa.Nodes) != 1 {
		return fmt.Errorf(
			"slide %d: SmartArt layout %q requires exactly one root node, got %d",
			slideIndex,
			sa.Layout.Name(),
			len(sa.Nodes),
		)
	}
	if len(sa.Nodes) == 0 {
		return nil
	}
	root := sa.Nodes[0]
	if constraint.maxRootChildren > 0 && len(root.Children) > constraint.maxRootChildren {
		return fmt.Errorf(
			"slide %d: SmartArt layout %q supports at most %d root child node(s), got %d",
			slideIndex,
			sa.Layout.Name(),
			constraint.maxRootChildren,
			len(root.Children),
		)
	}
	if constraint.maxChildrenPerNode > 0 {
		maxChildren := maxChildrenInTree(sa.Nodes)
		if maxChildren > constraint.maxChildrenPerNode {
			return fmt.Errorf(
				"slide %d: SmartArt layout %q supports at most %d children per node, got %d",
				slideIndex,
				sa.Layout.Name(),
				constraint.maxChildrenPerNode,
				maxChildren,
			)
		}
	}
	if constraint.maxDepth > 0 {
		depth := smartArtTreeDepth(sa.Nodes)
		if depth > constraint.maxDepth {
			return fmt.Errorf(
				"slide %d: SmartArt layout %q supports node depth up to %d, got %d",
				slideIndex,
				sa.Layout.Name(),
				constraint.maxDepth,
				depth,
			)
		}
	}
	return nil
}

func validateNode(n Node, slideIndex, nodeIndex int) error {
	if n.Text == "" {
		return fmt.Errorf("slide %d: SmartArt node %d has empty text", slideIndex, nodeIndex)
	}
	if n.Color != "" && !hexColorRE.MatchString(n.Color) {
		return fmt.Errorf(
			"slide %d: SmartArt node %d has invalid color %q (expected 6-digit hex)",
			slideIndex,
			nodeIndex,
			n.Color,
		)
	}
	for i, child := range n.Children {
		if err := validateNode(child, slideIndex, i); err != nil {
			return err
		}
	}
	return nil
}
