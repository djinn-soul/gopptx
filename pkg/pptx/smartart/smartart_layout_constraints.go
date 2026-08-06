package smartart

type layoutConstraint struct {
	maxItems           int
	requireSingleRoot  bool
	maxRootChildren    int
	maxChildrenPerNode int
	maxDepth           int
}

const (
	layoutMaxOne   = 1
	layoutMaxTwo   = 2
	layoutMaxThree = 3
	layoutMaxFour  = 4
)

// Layouts whose template holds a row of peers — the steps of a process, the
// entries of a list — take any number of items: the data model is generated to
// match and PowerPoint lays the diagram out from it. Only their nesting depth is
// fixed, by whether the layout draws a body under each entry.
func layoutConstraintFor(layout Layout) (layoutConstraint, bool) {
	switch layout {
	case BasicBlockList,
		BasicProcess,
		ContinuousBlockProcess,
		BasicCycle,
		TextCycle,
		BlockCycle,
		BasicVenn,
		LinearVenn,
		StackedVenn,
		BasicMatrix,
		BasicPyramid,
		InvertedPyramid,
		PictureStrips,
		PictureGrid:
		return layoutConstraint{maxDepth: layoutMaxOne}, true
	case VerticalBlockList,
		HorizontalBulletLst,
		SquareAccentList,
		PictureAccentList,
		AccentProcess,
		AlternatingFlow:
		return layoutConstraint{maxDepth: layoutMaxTwo}, true
	case BasicRadial:
		// One hub absorbs the list, so the hub plus its satellites is the limit.
		return layoutConstraint{maxItems: layoutMaxFour, maxDepth: layoutMaxTwo}, true
	case TitledMatrix:
		return layoutConstraint{maxItems: layoutMaxFour, maxDepth: layoutMaxTwo}, true
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
	default:
		return layoutConstraint{}, false
	}
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
