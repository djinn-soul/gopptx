package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/shapes"

// Shape type re-exports — expanded ECMA-376 presets.
// Core shapes re-exported in shape.go; this file covers the extended catalog.

// Basic geometry (expanded).
const (
	ShapeTypeOctagon               = shapes.ShapeTypeOctagon
	ShapeTypeDecagon               = shapes.ShapeTypeDecagon
	ShapeTypeDodecagon             = shapes.ShapeTypeDodecagon
	ShapeTypeHeptagon              = shapes.ShapeTypeHeptagon
	ShapeTypeTrapezoid             = shapes.ShapeTypeTrapezoid
	ShapeTypeNonIsoscelesTrapezoid = shapes.ShapeTypeNonIsoscelesTrapezoid
	ShapeTypeIsoscelesTrapezoid    = shapes.ShapeTypeIsoscelesTrapezoid
	ShapeTypeCube                  = shapes.ShapeTypeCube
	ShapeTypeCan                   = shapes.ShapeTypeCan
	ShapeTypeCone                  = shapes.ShapeTypeCone
	ShapeTypeCylinder              = shapes.ShapeTypeCylinder
	ShapeTypeFoldedCorner          = shapes.ShapeTypeFoldedCorner
	ShapeTypeFrame                 = shapes.ShapeTypeFrame
	ShapeTypeHalfFrame             = shapes.ShapeTypeHalfFrame
	ShapeTypeCorner                = shapes.ShapeTypeCorner
	ShapeTypeLShape                = shapes.ShapeTypeLShape
	ShapeTypeDiagStripe            = shapes.ShapeTypeDiagStripe
	ShapeTypePlus                  = shapes.ShapeTypePlus
	ShapeTypeDonut                 = shapes.ShapeTypeDonut
	ShapeTypeNoSmoking             = shapes.ShapeTypeNoSmoking
	ShapeTypeBevel                 = shapes.ShapeTypeBevel
	ShapeTypeFunnel                = shapes.ShapeTypeFunnel
	ShapeTypeTeardrop              = shapes.ShapeTypeTeardrop
	ShapeTypePie                   = shapes.ShapeTypePie
	ShapeTypePieWedge              = shapes.ShapeTypePieWedge
	ShapeTypeChord                 = shapes.ShapeTypeChord
	ShapeTypeArc                   = shapes.ShapeTypeArc
	ShapeTypeBlockArc              = shapes.ShapeTypeBlockArc
	ShapeTypeSmileyFace            = shapes.ShapeTypeSmileyFace
	ShapeTypeSun                   = shapes.ShapeTypeSun
	ShapeTypeMoon                  = shapes.ShapeTypeMoon
	ShapeTypeLightningBolt         = shapes.ShapeTypeLightningBolt
	ShapeTypePlaque                = shapes.ShapeTypePlaque
	ShapeTypeMusicNote             = shapes.ShapeTypeMusicNote
	ShapeTypePictureFrame          = shapes.ShapeTypePictureFrame
	ShapeTypeHomePlate             = shapes.ShapeTypeHomePlate
	ShapeTypeChevron               = shapes.ShapeTypeChevron
)

// Arrows (expanded).
const (
	ShapeTypeLeftRightArrow        = shapes.ShapeTypeLeftRightArrow
	ShapeTypeUpDownArrow           = shapes.ShapeTypeUpDownArrow
	ShapeTypeQuadArrow             = shapes.ShapeTypeQuadArrow
	ShapeTypeLeftUpArrow           = shapes.ShapeTypeLeftUpArrow
	ShapeTypeBentArrow             = shapes.ShapeTypeBentArrow
	ShapeTypeBentUpArrow           = shapes.ShapeTypeBentUpArrow
	ShapeTypeUturnArrow            = shapes.ShapeTypeUturnArrow
	ShapeTypeCircularArrow         = shapes.ShapeTypeCircularArrow
	ShapeTypeNotchedRightArrow     = shapes.ShapeTypeNotchedRightArrow
	ShapeTypeStripedRightArrow     = shapes.ShapeTypeStripedRightArrow
	ShapeTypeSwooshArrow           = shapes.ShapeTypeSwooshArrow
	ShapeTypeLeftRightUpArrow      = shapes.ShapeTypeLeftRightUpArrow
	ShapeTypeCurvedRightArrow      = shapes.ShapeTypeCurvedRightArrow
	ShapeTypeCurvedLeftArrow       = shapes.ShapeTypeCurvedLeftArrow
	ShapeTypeCurvedUpArrow         = shapes.ShapeTypeCurvedUpArrow
	ShapeTypeCurvedDownArrow       = shapes.ShapeTypeCurvedDownArrow
	ShapeTypeCurvedLeftRightArrow  = shapes.ShapeTypeCurvedLeftRightArrow
	ShapeTypeCurvedUpDownArrow     = shapes.ShapeTypeCurvedUpDownArrow
	ShapeTypePentagonArrow         = shapes.ShapeTypePentagonArrow
	ShapeTypeChevronArrow          = shapes.ShapeTypeChevronArrow
	ShapeTypeRightArrowCallout     = shapes.ShapeTypeRightArrowCallout
	ShapeTypeLeftArrowCallout      = shapes.ShapeTypeLeftArrowCallout
	ShapeTypeUpArrowCallout        = shapes.ShapeTypeUpArrowCallout
	ShapeTypeDownArrowCallout      = shapes.ShapeTypeDownArrowCallout
	ShapeTypeLeftRightArrowCallout = shapes.ShapeTypeLeftRightArrowCallout
	ShapeTypeQuadArrowCallout      = shapes.ShapeTypeQuadArrowCallout
	ShapeTypeUpDownArrowCallout    = shapes.ShapeTypeUpDownArrowCallout
)

// Stars, banners, scrolls.
const (
	ShapeTypeStar4            = shapes.ShapeTypeStar4
	ShapeTypeStar6            = shapes.ShapeTypeStar6
	ShapeTypeStar7            = shapes.ShapeTypeStar7
	ShapeTypeStar8            = shapes.ShapeTypeStar8
	ShapeTypeStar10           = shapes.ShapeTypeStar10
	ShapeTypeStar12           = shapes.ShapeTypeStar12
	ShapeTypeStar16           = shapes.ShapeTypeStar16
	ShapeTypeStar24           = shapes.ShapeTypeStar24
	ShapeTypeStar32           = shapes.ShapeTypeStar32
	ShapeTypeIrregularSeal1   = shapes.ShapeTypeIrregularSeal1
	ShapeTypeIrregularSeal2   = shapes.ShapeTypeIrregularSeal2
	ShapeTypeRibbon           = shapes.ShapeTypeRibbon
	ShapeTypeRibbon2          = shapes.ShapeTypeRibbon2
	ShapeTypeEllipseRibbon    = shapes.ShapeTypeEllipseRibbon
	ShapeTypeEllipseRibbon2   = shapes.ShapeTypeEllipseRibbon2
	ShapeTypeWave             = shapes.ShapeTypeWave
	ShapeTypeDoubleWave       = shapes.ShapeTypeDoubleWave
	ShapeTypeVerticalScroll   = shapes.ShapeTypeVerticalScroll
	ShapeTypeHorizontalScroll = shapes.ShapeTypeHorizontalScroll
	ShapeTypeSeal             = shapes.ShapeTypeSeal
	ShapeTypeSeal4            = shapes.ShapeTypeSeal4
	ShapeTypeSeal8            = shapes.ShapeTypeSeal8
	ShapeTypeSeal16           = shapes.ShapeTypeSeal16
	ShapeTypeSeal32           = shapes.ShapeTypeSeal32
)

// Callouts.
const (
	ShapeTypeCallout1             = shapes.ShapeTypeCallout1
	ShapeTypeCallout2             = shapes.ShapeTypeCallout2
	ShapeTypeCallout3             = shapes.ShapeTypeCallout3
	ShapeTypeBorderCallout1       = shapes.ShapeTypeBorderCallout1
	ShapeTypeBorderCallout2       = shapes.ShapeTypeBorderCallout2
	ShapeTypeBorderCallout3       = shapes.ShapeTypeBorderCallout3
	ShapeTypeAccentCallout1       = shapes.ShapeTypeAccentCallout1
	ShapeTypeAccentCallout2       = shapes.ShapeTypeAccentCallout2
	ShapeTypeAccentCallout3       = shapes.ShapeTypeAccentCallout3
	ShapeTypeAccentBorderCallout1 = shapes.ShapeTypeAccentBorderCallout1
	ShapeTypeAccentBorderCallout2 = shapes.ShapeTypeAccentBorderCallout2
	ShapeTypeAccentBorderCallout3 = shapes.ShapeTypeAccentBorderCallout3
	ShapeTypeWedgeRectCallout     = shapes.ShapeTypeWedgeRectCallout
	ShapeTypeWedgeRRectCallout    = shapes.ShapeTypeWedgeRRectCallout
	ShapeTypeWedgeEllipseCallout  = shapes.ShapeTypeWedgeEllipseCallout
	ShapeTypeCloudCallout         = shapes.ShapeTypeCloudCallout
)

// Flowcharts (expanded).
const (
	ShapeTypeFlowChartAlternateProcess        = shapes.ShapeTypeFlowChartAlternateProcess
	ShapeTypeFlowChartCard                    = shapes.ShapeTypeFlowChartCard
	ShapeTypeFlowChartCollate                 = shapes.ShapeTypeFlowChartCollate
	ShapeTypeFlowChartConnector               = shapes.ShapeTypeFlowChartConnector
	ShapeTypeFlowChartDelay                   = shapes.ShapeTypeFlowChartDelay
	ShapeTypeFlowChartDirectAccessStorage     = shapes.ShapeTypeFlowChartDirectAccessStorage
	ShapeTypeFlowChartDisplay                 = shapes.ShapeTypeFlowChartDisplay
	ShapeTypeFlowChartExtract                 = shapes.ShapeTypeFlowChartExtract
	ShapeTypeFlowChartInputOutput             = shapes.ShapeTypeFlowChartInputOutput
	ShapeTypeFlowChartInternalStorage         = shapes.ShapeTypeFlowChartInternalStorage
	ShapeTypeFlowChartMagneticDisk            = shapes.ShapeTypeFlowChartMagneticDisk
	ShapeTypeFlowChartManualInput             = shapes.ShapeTypeFlowChartManualInput
	ShapeTypeFlowChartManualOperation         = shapes.ShapeTypeFlowChartManualOperation
	ShapeTypeFlowChartMerge                   = shapes.ShapeTypeFlowChartMerge
	ShapeTypeFlowChartMultidocument           = shapes.ShapeTypeFlowChartMultidocument
	ShapeTypeFlowChartOffpageConnector        = shapes.ShapeTypeFlowChartOffpageConnector
	ShapeTypeFlowChartOr                      = shapes.ShapeTypeFlowChartOr
	ShapeTypeFlowChartPredefinedProcess       = shapes.ShapeTypeFlowChartPredefinedProcess
	ShapeTypeFlowChartPreparation             = shapes.ShapeTypeFlowChartPreparation
	ShapeTypeFlowChartPunchedTape             = shapes.ShapeTypeFlowChartPunchedTape
	ShapeTypeFlowChartSequentialAccessStorage = shapes.ShapeTypeFlowChartSequentialAccessStorage
	ShapeTypeFlowChartSort                    = shapes.ShapeTypeFlowChartSort
	ShapeTypeFlowChartStoredData              = shapes.ShapeTypeFlowChartStoredData
	ShapeTypeFlowChartSummingJunction         = shapes.ShapeTypeFlowChartSummingJunction
	ShapeTypeFlowChartPunchedCard             = shapes.ShapeTypeFlowChartPunchedCard
	ShapeTypeFlowChartOnlineStorage           = shapes.ShapeTypeFlowChartOnlineStorage
	ShapeTypeFlowChartMagneticTape            = shapes.ShapeTypeFlowChartMagneticTape
	ShapeTypeFlowChartMagneticDrum            = shapes.ShapeTypeFlowChartMagneticDrum
)

// Action buttons.
const (
	ShapeTypeActionButtonBlank        = shapes.ShapeTypeActionButtonBlank
	ShapeTypeActionButtonHome         = shapes.ShapeTypeActionButtonHome
	ShapeTypeActionButtonHelp         = shapes.ShapeTypeActionButtonHelp
	ShapeTypeActionButtonBackPrevious = shapes.ShapeTypeActionButtonBackPrevious
	ShapeTypeActionButtonForwardNext  = shapes.ShapeTypeActionButtonForwardNext
	ShapeTypeActionButtonBeginning    = shapes.ShapeTypeActionButtonBeginning
	ShapeTypeActionButtonEnd          = shapes.ShapeTypeActionButtonEnd
	ShapeTypeActionButtonReturn       = shapes.ShapeTypeActionButtonReturn
	ShapeTypeActionButtonDocument     = shapes.ShapeTypeActionButtonDocument
	ShapeTypeActionButtonSound        = shapes.ShapeTypeActionButtonSound
	ShapeTypeActionButtonMovie        = shapes.ShapeTypeActionButtonMovie
	ShapeTypeActionButtonInformation  = shapes.ShapeTypeActionButtonInformation
)

// Math and special shapes.
const (
	ShapeTypeMathPlus      = shapes.ShapeTypeMathPlus
	ShapeTypeMathMinus     = shapes.ShapeTypeMathMinus
	ShapeTypeMathMultiply  = shapes.ShapeTypeMathMultiply
	ShapeTypeMathDivide    = shapes.ShapeTypeMathDivide
	ShapeTypeMathEqual     = shapes.ShapeTypeMathEqual
	ShapeTypeMathNotEqual  = shapes.ShapeTypeMathNotEqual
	ShapeTypeGear6         = shapes.ShapeTypeGear6
	ShapeTypeGear9         = shapes.ShapeTypeGear9
	ShapeTypeChartPlus     = shapes.ShapeTypeChartPlus
	ShapeTypeChartStar     = shapes.ShapeTypeChartStar
	ShapeTypeChartX        = shapes.ShapeTypeChartX
	ShapeTypeBracePair     = shapes.ShapeTypeBracePair
	ShapeTypeBracketPair   = shapes.ShapeTypeBracketPair
	ShapeTypeLeftBrace     = shapes.ShapeTypeLeftBrace
	ShapeTypeRightBrace    = shapes.ShapeTypeRightBrace
	ShapeTypeLeftBracket   = shapes.ShapeTypeLeftBracket
	ShapeTypeRightBracket  = shapes.ShapeTypeRightBracket
	ShapeTypeDoubleBrace   = shapes.ShapeTypeDoubleBrace
	ShapeTypeDoubleBracket = shapes.ShapeTypeDoubleBracket
	ShapeTypeLine          = shapes.ShapeTypeLine
	ShapeTypeLineInv       = shapes.ShapeTypeLineInv
)

// Rectangle variants.
const (
	ShapeTypeSnip1Rect      = shapes.ShapeTypeSnip1Rect
	ShapeTypeSnip2SameRect  = shapes.ShapeTypeSnip2SameRect
	ShapeTypeSnip2DiagRect  = shapes.ShapeTypeSnip2DiagRect
	ShapeTypeRound1Rect     = shapes.ShapeTypeRound1Rect
	ShapeTypeRound2SameRect = shapes.ShapeTypeRound2SameRect
	ShapeTypeRound2DiagRect = shapes.ShapeTypeRound2DiagRect
	ShapeTypeSnipRoundRect  = shapes.ShapeTypeSnipRoundRect
	ShapeTypePlaqueTabs     = shapes.ShapeTypePlaqueTabs
	ShapeTypeSquareTabs     = shapes.ShapeTypeSquareTabs
	ShapeTypeCornerTabs     = shapes.ShapeTypeCornerTabs
)
