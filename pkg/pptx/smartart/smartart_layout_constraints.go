package smartart

type layoutConstraint struct {
	maxItems           int
	requireSingleRoot  bool
	maxRootChildren    int
	maxChildrenPerNode int
	maxDepth           int
}

const (
	layoutMaxTwo   = 2
	layoutMaxThree = 3
	layoutMaxFour  = 4
	layoutMaxFive  = 5
)

func layoutConstraintFor(layout Layout) (layoutConstraint, bool) {
	switch layout {
	case BasicBlockList:
		return layoutConstraint{maxItems: layoutMaxFive}, true
	case VerticalBlockList:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	case HorizontalBulletLst:
		return layoutConstraint{maxItems: layoutMaxFour}, true
	case SquareAccentList:
		return layoutConstraint{maxItems: layoutMaxFour}, true
	case PictureAccentList:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	case BasicProcess:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	case AccentProcess:
		return layoutConstraint{maxItems: layoutMaxFive}, true
	case AlternatingFlow:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	case ContinuousBlockProcess:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	case BasicCycle, TextCycle, BlockCycle:
		return layoutConstraint{maxItems: layoutMaxFive}, true
	case OrgChart:
		return layoutConstraint{
			requireSingleRoot: true,
			maxRootChildren:   layoutMaxFour,
			maxDepth:          layoutMaxTwo,
		}, true
	case Hierarchy, HorizontalHierarchy:
		return layoutConstraint{
			requireSingleRoot:  true,
			maxRootChildren:    layoutMaxTwo,
			maxChildrenPerNode: layoutMaxTwo,
			maxDepth:           layoutMaxThree,
		}, true
	case BasicVenn:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	case LinearVenn, StackedVenn, BasicRadial, BasicMatrix, TitledMatrix, PictureGrid:
		return layoutConstraint{maxItems: layoutMaxFour}, true
	case BasicPyramid, InvertedPyramid, PictureStrips:
		return layoutConstraint{maxItems: layoutMaxThree}, true
	default:
		return layoutConstraint{}, false
	}
}

func hasNestedChildren(nodes []Node) bool {
	for _, n := range nodes {
		if len(n.Children) > 0 {
			return true
		}
	}
	return false
}

func smartArtTreeDepth(nodes []Node) int {
	maxDepth := 0
	for _, n := range nodes {
		depth := nodeDepth(n)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

func nodeDepth(n Node) int {
	maxChildDepth := 0
	for _, child := range n.Children {
		childDepth := nodeDepth(child)
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}
	return 1 + maxChildDepth
}

func maxChildrenInTree(nodes []Node) int {
	maxChildren := 0
	for _, n := range nodes {
		maxChildren = maxInt(maxChildren, len(n.Children))
		maxChildren = maxInt(maxChildren, maxChildrenInTree(n.Children))
	}
	return maxChildren
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
