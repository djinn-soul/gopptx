package smartart

type layoutConstraint struct {
	maxItems           int
	requireSingleRoot  bool
	maxRootChildren    int
	maxChildrenPerNode int
	maxDepth           int
}

var smartArtLayoutConstraints = map[Layout]layoutConstraint{
	BasicBlockList:      {maxItems: 5},
	VerticalBlockList:   {maxItems: 3},
	HorizontalBulletLst: {maxItems: 4},
	SquareAccentList:    {maxItems: 4},
	PictureAccentList:   {maxItems: 3},

	BasicProcess:           {maxItems: 3},
	AccentProcess:          {maxItems: 5},
	AlternatingFlow:        {maxItems: 3},
	ContinuousBlockProcess: {maxItems: 3},

	BasicCycle: {maxItems: 5},
	TextCycle:  {maxItems: 5},
	BlockCycle: {maxItems: 5},

	OrgChart: {
		requireSingleRoot: true,
		maxRootChildren:   4,
		maxDepth:          2,
	},
	Hierarchy: {
		requireSingleRoot:  true,
		maxRootChildren:    2,
		maxChildrenPerNode: 2,
		maxDepth:           3,
	},
	HorizontalHierarchy: {
		requireSingleRoot:  true,
		maxRootChildren:    2,
		maxChildrenPerNode: 2,
		maxDepth:           3,
	},

	BasicVenn:   {maxItems: 3},
	LinearVenn:  {maxItems: 4},
	StackedVenn: {maxItems: 4},
	BasicRadial: {maxItems: 4},

	BasicMatrix:  {maxItems: 4},
	TitledMatrix: {maxItems: 4},

	BasicPyramid:    {maxItems: 3},
	InvertedPyramid: {maxItems: 3},

	PictureStrips: {maxItems: 3},
	PictureGrid:   {maxItems: 4},
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
