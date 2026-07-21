package pptxxml

// OOXML token values repeated across the renderers in this package.
//
// These deliberately duplicate names that also exist in higher-level packages
// (pkg/pptx/shapes, pkg/pptx/enums). pptxxml sits below those packages in the
// dependency graph and must not import them, so the tokens are defined locally.
const (
	imageExtPNG = "png"
	imageExtJPG = "jpg"

	placeholderPicture = "picture"
	placeholderBody    = "body"

	fillTypeSolid = "solid"

	lineJoinMiter = "miter"

	strokeDashDash = "dash"

	axisCrossBetweenMidCat = "midCat"

	langEnUS = "en-US"

	textAutoFitTag     = "spAutoFit"
	textAutoFitElement = "<a:spAutoFit/>"

	smartArtDefaultLayoutURN = "urn:microsoft.com/office/officeart/2005/8/layout/default"
)
