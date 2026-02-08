package pptx

import (
	"fmt"
	"strings"
	"sync"
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
// It defines the effect type, trigger condition, direction, duration, and other timing properties.
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

// ToAnimation implements the AnimationDefinition interface.
func (a Animation) ToAnimation() Animation {
	return a
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

func (a Animation) xml(seqID int, actualShapeID int) string {
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

// AnimationDefinition is the interface for types that can be converted to an Animation.
// This allows for extensible animation builders and custom implementations.
type AnimationDefinition interface {
	ToAnimation() Animation
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

func slideAnimationsXML(s SlideContent, shapeIDs []int) string {
	if len(s.Animations) == 0 {
		return ""
	}

	// Calculate concurrently if we have enough animations
	const parallelThreshold = 50
	numAnims := len(s.Animations)
	animationsXML := make([]string, numAnims)

	if numAnims < parallelThreshold {
		for i, anim := range s.Animations {
			actualID := 0
			if anim.ShapeIndex > 0 && anim.ShapeIndex <= len(shapeIDs) {
				actualID = shapeIDs[anim.ShapeIndex-1]
			}
			if actualID == 0 {
				continue
			}
			animationsXML[i] = anim.xml(i*2+3, actualID)
		}
	} else {
		var wg sync.WaitGroup
		numWorkers := 4
		chunkSize := (numAnims + numWorkers - 1) / numWorkers

		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				start := workerID * chunkSize
				end := start + chunkSize
				if end > numAnims {
					end = numAnims
				}

				for i := start; i < end; i++ {
					anim := s.Animations[i]
					actualID := 0
					if anim.ShapeIndex > 0 && anim.ShapeIndex <= len(shapeIDs) {
						actualID = shapeIDs[anim.ShapeIndex-1]
					}
					if actualID == 0 {
						continue
					}
					animationsXML[i] = anim.xml(i*2+3, actualID)
				}
			}(w)
		}
		wg.Wait()
	}

	// Filter out empty strings (from invalid shape IDs)
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

func calculateShapeIDs(s SlideContent) []int {
	// Shape IDs start at 2 (Title = 2).
	// If layout has content placeholder (bullets, table, chart), that takes ID 3.
	// However, gopptx's slide_xml.go renders title (ID 2), then body/content (ID 3), then images/charts/tables if mixed?
	// Actually, looking at slide_xml.go SlideWithLayout:
	// 1. Title (ID 2)
	// 2. Body/Content (ID 3) - this is for bullets, or if table/chart is the *main* content of the layout.
	//
	// In existing gopptx logic, NewSlide creates "titleAndContent" layout.
	// - If we add bullets, they go into the body placeholder (ID 3).
	// - If we add a table/chart via WithTable/WithChart, they might replace the body or be added as separate objects?
	//
	// Let's trace SlideWithLayout in internal/pptxxml/slide_xml.go.
	// It calls `makeTitle(2, ...)`
	// Then `makeBody(3, ...)` if layout supports it.
	// Then `makeTable(..., nextID)`
	//
	// We need to exactly match the ID generation in SlideWithLayout.
	//
	// In SlideWithLayout currently (as of recent view):
	// nextID starts at 2.
	// Title uses nextID (2). nextID++.
	// Body uses nextID (3) if layout has body. nextID++.
	// Table uses nextID. nextID++.
	// Chart uses nextID. nextID++.
	// Images use nextID. nextID++.
	// Shapes use nextID. nextID++.
	//
	// So we need to replicate that sequence.

	nextID := 2 // Title is always 2.
	layoutMode := slideLayoutXMLMode(s.Layout)

	// Title
	nextID++

	// Body / Content Placeholder
	// "titleAndContent", "twoColumn", etc. have a body placeholder.
	// But SlideWithLayout ONLY increments nextID if content is actually added.
	// Check SlideWithLayout logic:
	// if table != nil -> adds table, nextID++
	// else if len(bullets) > 0 -> adds content/bullets, nextID++ (or +2 for twoColumn)
	//
	// So we must check if content is PRESENT.
	hasBodyContent := len(s.Bullets) > 0
	if layoutMode != "blank" && layoutMode != "titleOnly" && s.Table == nil {
		// If table is nil, we MIGHT have bullets.
		if hasBodyContent {
			nextID++
			if layoutMode == "twoColumn" {
				leftCount := (len(s.Bullets) + 1) / 2
				if len(s.Bullets[leftCount:]) > 0 {
					nextID++
				}
			}
		}
	}

	// Table
	if s.Table != nil {
		nextID++
	}

	// Chart
	if chartKindCount(s) > 0 {
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


