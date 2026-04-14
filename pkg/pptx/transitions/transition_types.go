package transitions

// SlideTransition is the extensibility contract for slide transitions.
type SlideTransition interface {
	Validate() error
	XML() string
}

// TransitionType is the built-in transition enum.
type TransitionType string

const (
	TransitionNone       TransitionType = "none"
	TransitionCut        TransitionType = "cut"
	TransitionFade       TransitionType = "fade"
	TransitionPush       TransitionType = "push"
	TransitionWipe       TransitionType = "wipe"
	TransitionSplit      TransitionType = "split"
	TransitionReveal     TransitionType = "reveal"
	TransitionCover      TransitionType = "cover"
	TransitionZoom       TransitionType = "zoom"
	TransitionRandomBars TransitionType = "randomBar"
	TransitionShape      TransitionType = "circle"
	TransitionUncover    TransitionType = "pull"
	TransitionFlash      TransitionType = "flash"
	TransitionStrips     TransitionType = "strips"
	TransitionBlinds     TransitionType = "blinds"
	TransitionClock      TransitionType = "wheel"
	TransitionRipple     TransitionType = "ripple"
	TransitionHoneycomb  TransitionType = "honeycomb"
	TransitionGlitter    TransitionType = "glitter"
	TransitionVortex     TransitionType = "vortex"
	TransitionShred      TransitionType = "shred"
	TransitionSwitch     TransitionType = "switch"
	TransitionFlip       TransitionType = "flip"
	TransitionGallery    TransitionType = "gallery"
	TransitionCube       TransitionType = "cube"
	TransitionDoors      TransitionType = "doors"
	TransitionBox        TransitionType = "box"
	TransitionRandom     TransitionType = "random"
	TransitionMorph      TransitionType = "morph"
)

// MorphOption defines the granularity of a Morph transition.
type MorphOption string

const (
	MorphOptionObject    MorphOption = "obj"
	MorphOptionWord      MorphOption = "word"
	MorphOptionCharacter MorphOption = "char"
)

// TransitionDirection defines the direction of a transition.
type TransitionDirection string

const (
	TransitionDirIn        TransitionDirection = "in"
	TransitionDirOut       TransitionDirection = "out"
	TransitionDirUp        TransitionDirection = "u"
	TransitionDirDown      TransitionDirection = "d"
	TransitionDirLeft      TransitionDirection = "l"
	TransitionDirRight     TransitionDirection = "r"
	TransitionDirUpLeft    TransitionDirection = "lu"
	TransitionDirUpRight   TransitionDirection = "ru"
	TransitionDirDownLeft  TransitionDirection = "ld"
	TransitionDirDownRight TransitionDirection = "rd"
)

// TransitionOrientation defines the orientation of a transition.
type TransitionOrientation string

const (
	TransitionOrientHorizontal TransitionOrientation = "horz"
	TransitionOrientVertical   TransitionOrientation = "vert"
)

// TransitionSound defines the audio configuration for a transition.
type TransitionSound struct {
	RelID string // Relationship ID for the audio file (required)
	Name  string // Display name (e.g., "Applause")
	Loop  bool   // Whether to loop the sound
}

// TransitionOptions provides advanced configuration for a slide transition.
type TransitionOptions struct {
	Type                  TransitionType
	DurationMS            uint32
	Direction             TransitionDirection
	Orientation           TransitionOrientation
	SpokeCount            uint32
	ThruBlk               bool
	Sound                 *TransitionSound
	DisableAdvanceOnClick bool
	AdvanceAfterMS        uint32
	MorphOption           MorphOption
}
