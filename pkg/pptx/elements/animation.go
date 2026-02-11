package elements

import (
	"fmt"
	"strings"
)

// AnimationEffect defines the type of animation effect.
type AnimationEffect string

const (
	// Entrance effects
	AnimationEntranceAppear      AnimationEffect = "entr_appear"
	AnimationEntranceFade        AnimationEffect = "entr_fade"
	AnimationEntranceFlyIn       AnimationEffect = "entr_flyIn"
	AnimationEntranceFloat       AnimationEffect = "entr_float"
	AnimationEntranceSplit       AnimationEffect = "entr_split"
	AnimationEntranceWipe        AnimationEffect = "entr_wipe"
	AnimationEntranceShape       AnimationEffect = "entr_shape"
	AnimationEntranceWheel       AnimationEffect = "entr_wheel"
	AnimationEntranceRandomBars  AnimationEffect = "entr_randomBars"
	AnimationEntranceGrowAndTurn AnimationEffect = "entr_growAndTurn"
	AnimationEntranceZoom        AnimationEffect = "entr_zoom"
	AnimationEntranceSwivel      AnimationEffect = "entr_swivel"
	AnimationEntranceBounce      AnimationEffect = "entr_bounce"

	// Exit effects
	AnimationExitDisappear AnimationEffect = "exit_disappear"
	AnimationExitFadeOut   AnimationEffect = "exit_fade"
	AnimationExitFlyOut    AnimationEffect = "exit_flyOut"
	AnimationExitFloatOut  AnimationEffect = "exit_float"

	// Emphasis effects
	AnimationEmphasisPulse        AnimationEffect = "emph_pulse"
	AnimationEmphasisColorPulse   AnimationEffect = "emph_colorPulse"
	AnimationEmphasisTeeter       AnimationEffect = "emph_teeter"
	AnimationEmphasisSpin         AnimationEffect = "emph_spin"
	AnimationEmphasisGrowShrink   AnimationEffect = "emph_growShrink"
	AnimationEmphasisDesaturate   AnimationEffect = "emph_desaturate"
	AnimationEmphasisDarken       AnimationEffect = "emph_darken"
	AnimationEmphasisLighten      AnimationEffect = "emph_lighten"
	AnimationEmphasisTransparency AnimationEffect = "emph_transparency"
	AnimationEmphasisObjectColor  AnimationEffect = "emph_objectColor"

	// Motion paths
	AnimationPathLines  AnimationEffect = "path_lines"
	AnimationPathArcs   AnimationEffect = "path_arcs"
	AnimationPathTurns  AnimationEffect = "path_turns"
	AnimationPathShapes AnimationEffect = "path_shapes"
	AnimationPathLoops  AnimationEffect = "path_loops"
	AnimationPathCustom AnimationEffect = "path_custom"
)

// AnimationTrigger defines when the animation starts.
type AnimationTrigger string

const (
	AnimationOnClick       AnimationTrigger = "onClick"
	AnimationWithPrevious  AnimationTrigger = "withPrev"
	AnimationAfterPrevious AnimationTrigger = "afterPrev"
)

// AnimationDirection defines the direction or subtype of the animation.
type AnimationDirection string

const (
	AnimationDirIn        AnimationDirection = "in"
	AnimationDirOut       AnimationDirection = "out"
	AnimationDirUp        AnimationDirection = "u"
	AnimationDirDown      AnimationDirection = "d"
	AnimationDirLeft      AnimationDirection = "l"
	AnimationDirRight     AnimationDirection = "r"
	AnimationDirUpLeft    AnimationDirection = "ul"
	AnimationDirUpRight   AnimationDirection = "ur"
	AnimationDirDownLeft  AnimationDirection = "dl"
	AnimationDirDownRight AnimationDirection = "dr"
)

// Animation represents a single animation effect on a slide object.
type Animation struct {
	ShapeIndex  int // 1-indexed, matches Shape/Connector index
	Effect      AnimationEffect
	Trigger     AnimationTrigger
	Direction   AnimationDirection
	DurationMS  uint32
	DelayMS     uint32
	RepeatCount uint32
	AutoReverse bool
}

// NewAnimation creates a new animation with default settings (500ms duration, OnClick).
func NewAnimation(shapeIndex int, effect AnimationEffect) Animation {
	return Animation{
		ShapeIndex: shapeIndex,
		Effect:     effect,
		Trigger:    AnimationOnClick,
		DurationMS: 500,
	}
}

// WithTrigger sets the trigger for an animation.
func (a Animation) WithTrigger(trigger AnimationTrigger) Animation {
	a.Trigger = trigger
	return a
}

// WithDuration sets the duration of the animation in milliseconds.
func (a Animation) WithDuration(durationMS uint32) Animation {
	a.DurationMS = durationMS
	return a
}

// WithDelay sets the delay before the animation starts in milliseconds.
func (a Animation) WithDelay(delayMS uint32) Animation {
	a.DelayMS = delayMS
	return a
}

// WithRepeat sets the number of times the animation repeats.
func (a Animation) WithRepeat(count uint32) Animation {
	a.RepeatCount = count
	return a
}

// WithAutoReverse enables or disables auto-reverse for the animation.
func (a Animation) WithAutoReverse(autoReverse bool) Animation {
	a.AutoReverse = autoReverse
	return a
}

func (a Animation) Validate() error {
	if err := a.Effect.Validate(); err != nil {
		return err
	}
	if err := a.Trigger.Validate(); err != nil {
		return err
	}
	return nil
}

func (e AnimationEffect) Validate() error {
	switch e {
	case AnimationEntranceAppear, AnimationEntranceFade, AnimationEntranceFlyIn, AnimationEntranceFloat,
		AnimationEntranceSplit, AnimationEntranceWipe, AnimationEntranceShape, AnimationEntranceWheel,
		AnimationEntranceRandomBars, AnimationEntranceGrowAndTurn, AnimationEntranceZoom, AnimationEntranceSwivel,
		AnimationEntranceBounce, AnimationExitDisappear, AnimationExitFadeOut, AnimationExitFlyOut,
		AnimationExitFloatOut, AnimationEmphasisPulse, AnimationEmphasisColorPulse, AnimationEmphasisTeeter,
		AnimationEmphasisSpin, AnimationEmphasisGrowShrink, AnimationEmphasisDesaturate, AnimationEmphasisDarken,
		AnimationEmphasisLighten, AnimationEmphasisTransparency, AnimationEmphasisObjectColor, AnimationPathLines,
		AnimationPathArcs, AnimationPathTurns, AnimationPathShapes, AnimationPathLoops, AnimationPathCustom:
		return nil
	default:
		return fmt.Errorf("unsupported animation effect: %q", string(e))
	}
}

func (t AnimationTrigger) Validate() error {
	switch t {
	case AnimationOnClick, AnimationWithPrevious, AnimationAfterPrevious:
		return nil
	default:
		return fmt.Errorf("unsupported animation trigger: %q", string(t))
	}
}

func (a Animation) presetID() uint32 {
	switch a.Effect {
	case AnimationEntranceAppear, AnimationExitDisappear:
		return 1
	case AnimationEntranceFade, AnimationExitFadeOut:
		return 10
	case AnimationEntranceFlyIn, AnimationExitFlyOut:
		return 2
	case AnimationEntranceFloat, AnimationExitFloatOut:
		return 14
	case AnimationEntranceSplit:
		return 16
	case AnimationEntranceWipe:
		return 22
	case AnimationEntranceShape:
		return 17
	case AnimationEntranceWheel:
		return 21
	case AnimationEntranceRandomBars:
		return 15
	case AnimationEntranceGrowAndTurn:
		return 26
	case AnimationEntranceZoom:
		return 23
	case AnimationEntranceSwivel:
		return 19
	case AnimationEntranceBounce:
		return 25
	case AnimationEmphasisPulse:
		return 31
	case AnimationEmphasisColorPulse:
		return 32
	case AnimationEmphasisTeeter:
		return 33
	case AnimationEmphasisSpin:
		return 34
	case AnimationEmphasisGrowShrink:
		return 35
	case AnimationEmphasisDesaturate:
		return 36
	case AnimationEmphasisDarken:
		return 37
	case AnimationEmphasisLighten:
		return 38
	case AnimationEmphasisTransparency:
		return 39
	case AnimationEmphasisObjectColor:
		return 40
	case AnimationPathLines:
		return 42
	case AnimationPathArcs:
		return 43
	case AnimationPathTurns:
		return 44
	case AnimationPathShapes:
		return 45
	case AnimationPathLoops:
		return 46
	case AnimationPathCustom:
		return 47
	default:
		return 0
	}
}

func (a Animation) presetClass() string {
	switch {
	case strings.HasPrefix(string(a.Effect), "entr_"):
		return "entr"
	case strings.HasPrefix(string(a.Effect), "exit_"):
		return "exit"
	case strings.HasPrefix(string(a.Effect), "emph_"):
		return "emph"
	case strings.HasPrefix(string(a.Effect), "path_"):
		return "path"
	default:
		return "entr"
	}
}

func (a Animation) XML(seqID int, actualShapeID int) string {
	repeatAttr := ""
	if a.RepeatCount > 0 {
		repeatAttr = fmt.Sprintf(` repeatCount="%d000"`, a.RepeatCount)
	}
	reverseAttr := ""
	if a.AutoReverse {
		reverseAttr = ` autoRev="1"`
	}

	return fmt.Sprintf(`
<p:par>
  <p:cTn id="%d" presetID="%d" presetClass="%s" presetSubtype="0" fill="hold" nodeType="%s">
    <p:stCondLst>
      <p:cond delay="%d"/>
    </p:stCondLst>
    <p:childTnLst>
      <p:set>
        <p:cBhvr>
          <p:cTn id="%d" dur="%d" fill="hold"%s%s>
            <p:stCondLst><p:cond delay="0"/></p:stCondLst>
          </p:cTn>
          <p:tgtEl>
            <p:spTgt spid="%d"/>
          </p:tgtEl>
        </p:cBhvr>
      </p:set>
    </p:childTnLst>
  </p:cTn>
</p:par>`,
		seqID,
		a.presetID(),
		a.presetClass(),
		a.nodeType(),
		a.DelayMS,
		seqID+1,
		a.DurationMS,
		repeatAttr,
		reverseAttr,
		actualShapeID,
	)
}

func (a Animation) nodeType() string {
	switch a.Trigger {
	case AnimationOnClick:
		return "clickEffect"
	case AnimationWithPrevious:
		return "withEffect"
	case AnimationAfterPrevious:
		return "afterEffect"
	default:
		return "clickEffect"
	}
}

// AnimationDefinition is the interface for types that can be converted to an Animation.
type AnimationDefinition interface {
	ToAnimation() Animation
}

func (a Animation) ToAnimation() Animation {
	return a
}

// CalculateShapeIDs replicates the ID generation order in slide_xml.go.
func CalculateShapeIDs(s SlideContent) []int {
	// Shape IDs start at 2 (Title = 2).
	nextID := 2

	// Title
	if s.Layout != SlideLayoutBlank {
		nextID++
	}

	// Table or Content
	if s.Table != nil {
		nextID++
	} else if len(s.Bullets) > 0 || len(s.BulletRuns) > 0 {
		nextID++
		if s.Layout == SlideLayoutTwoColumn {
			leftCount := (len(s.Bullets) + 1) / 2
			if len(s.Bullets[leftCount:]) > 0 {
				nextID++
			}
		}
	}

	// Primary chart object occupies one shape slot.
	if hasPrimaryChart(s) {
		nextID++
	}

	// Images
	nextID += len(s.Images)

	// Custom Shapes
	shapeIDs := make([]int, len(s.Shapes))
	for i := range s.Shapes {
		shapeIDs[i] = nextID + i
	}
	return shapeIDs
}

func hasPrimaryChart(s SlideContent) bool {
	return s.Chart != nil ||
		s.BarHorizontal != nil ||
		s.BarStacked != nil ||
		s.BarStacked100 != nil ||
		s.Line != nil ||
		s.LineMarkers != nil ||
		s.LineStacked != nil ||
		s.Scatter != nil ||
		s.Area != nil ||
		s.AreaStacked != nil ||
		s.AreaStacked100 != nil ||
		s.Pie != nil ||
		s.Doughnut != nil ||
		s.Bubble != nil ||
		s.Radar != nil ||
		s.RadarFilled != nil ||
		s.StockHLC != nil ||
		s.StockOHLC != nil ||
		s.Combo != nil
}

func SlideAnimationsXML(s SlideContent, shapeIDs []int) string {
	if len(s.Animations) == 0 {
		return ""
	}

	animationsXML := make([]string, len(s.Animations))
	for i, anim := range s.Animations {
		actualID := 0
		if anim.ShapeIndex > 0 && anim.ShapeIndex <= len(shapeIDs) {
			actualID = shapeIDs[anim.ShapeIndex-1]
		}
		if actualID == 0 {
			continue
		}
		animationsXML[i] = anim.XML(i*2+3, actualID)
	}

	var finalXML []string
	for _, xml := range animationsXML {
		if xml != "" {
			finalXML = append(finalXML, xml)
		}
	}

	if len(finalXML) == 0 {
		return ""
	}

	return fmt.Sprintf(`
<p:timing>
  <p:tnLst>
    <p:par>
      <p:cTn id="1" dur="indefinite" restart="never" nodeType="tmRoot">
        <p:childTnLst>
          <p:seq concurrent="1" nextAc="seek">
            <p:cTn id="2" dur="indefinite" nodeType="mainSeq">
              <p:childTnLst>
                %s
              </p:childTnLst>
            </p:cTn>
          </p:seq>
        </p:childTnLst>
      </p:cTn>
    </p:par>
  </p:tnLst>
</p:timing>`,
		strings.Join(finalXML, "\n"),
	)
}
