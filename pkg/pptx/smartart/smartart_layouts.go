package smartart

func processLayoutName(l Layout) (string, bool) {
	switch l {
	case BasicBlockList:
		return "Basic Block List", true
	case VerticalBlockList:
		return "Vertical Block List", true
	case HorizontalBulletLst:
		return "Horizontal Bullet List", true
	case SquareAccentList:
		return "Square Accent List", true
	case PictureAccentList:
		return "Picture Accent List", true
	case BasicProcess:
		return "Basic Process", true
	case AccentProcess:
		return "Accent Process", true
	case AlternatingFlow:
		return "Alternating Flow", true
	case ContinuousBlockProcess:
		return "Continuous Block Process", true
	case BasicCycle:
		return "Basic Cycle", true
	case TextCycle:
		return "Text Cycle", true
	case BlockCycle:
		return "Block Cycle", true
	case OrgChart:
		return "Organization Chart", true
	case Hierarchy:
		return "Hierarchy", true
	case HorizontalHierarchy:
		return "Horizontal Hierarchy", true
	default:
		return "", false
	}
}

func relationshipLayoutName(l Layout) (string, bool) {
	switch l {
	case BasicVenn:
		return "Basic Venn", true
	case LinearVenn:
		return "Linear Venn", true
	case StackedVenn:
		return "Stacked Venn", true
	case BasicRadial:
		return "Basic Radial", true
	case BasicMatrix:
		return "Basic Matrix", true
	case TitledMatrix:
		return "Titled Matrix", true
	case BasicPyramid:
		return "Basic Pyramid", true
	case InvertedPyramid:
		return "Inverted Pyramid", true
	default:
		return "", false
	}
}

func matrixPictureLayoutName(l Layout) (string, bool) {
	switch l {
	case PictureStrips:
		return "Picture Strips", true
	case PictureGrid:
		return "Picture Grid", true
	default:
		return "", false
	}
}
