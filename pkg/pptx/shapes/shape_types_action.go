package shapes

// Action button shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeActionButtonBlank renders a blank action button.
	ShapeTypeActionButtonBlank = "actionButtonBlank"
	// ShapeTypeActionButtonHome renders a home action button.
	ShapeTypeActionButtonHome = "actionButtonHome"
	// ShapeTypeActionButtonHelp renders a help action button.
	ShapeTypeActionButtonHelp = "actionButtonHelp"
	// ShapeTypeActionButtonBackPrevious renders a back/previous action button.
	ShapeTypeActionButtonBackPrevious = "actionButtonBackPrevious"
	// ShapeTypeActionButtonForwardNext renders a forward/next action button.
	ShapeTypeActionButtonForwardNext = "actionButtonForwardNext"
	// ShapeTypeActionButtonBeginning renders a beginning action button.
	ShapeTypeActionButtonBeginning = "actionButtonBeginning"
	// ShapeTypeActionButtonEnd renders an end action button.
	ShapeTypeActionButtonEnd = "actionButtonEnd"
	// ShapeTypeActionButtonReturn renders a return action button.
	ShapeTypeActionButtonReturn = "actionButtonReturn"
	// ShapeTypeActionButtonDocument renders a document action button.
	ShapeTypeActionButtonDocument = "actionButtonDocument"
	// ShapeTypeActionButtonSound renders a sound action button.
	ShapeTypeActionButtonSound = "actionButtonSound"
	// ShapeTypeActionButtonMovie renders a movie action button.
	ShapeTypeActionButtonMovie = "actionButtonMovie"
	// ShapeTypeActionButtonInformation renders an information action button.
	ShapeTypeActionButtonInformation = "actionButtonInformation"
)

func initActionShapes() {
	for _, t := range []string{
		ShapeTypeActionButtonBlank, ShapeTypeActionButtonHome,
		ShapeTypeActionButtonHelp, ShapeTypeActionButtonBackPrevious,
		ShapeTypeActionButtonForwardNext, ShapeTypeActionButtonBeginning,
		ShapeTypeActionButtonEnd, ShapeTypeActionButtonReturn,
		ShapeTypeActionButtonDocument, ShapeTypeActionButtonSound,
		ShapeTypeActionButtonMovie, ShapeTypeActionButtonInformation,
	} {
		registerShapeType(t)
	}
}
