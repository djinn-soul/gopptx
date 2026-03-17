package editor

import (
	"fmt"
	"strings"
)

const animPresetDefaultFade = 10

// AddSlideAnimation adds an animation effect to a shape on an existing slide.
// shapeID is the shape's numeric ID in the slide XML.
// effect is one of the animation effect constants (e.g. "entr_fade").
// trigger: "onClick", "withPrev", or "afterPrev". durationMS and delayMS in milliseconds.
func (e *PresentationEditor) AddSlideAnimation(
	slideIndex int,
	shapeID int,
	effect string,
	trigger string,
	durationMS int,
	delayMS int,
) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	data, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %s not found", slideRef.Part)
	}

	animXML := buildAnimationXML(shapeID, effect, trigger, durationMS, delayMS)
	slideXML := injectAnimationIntoSlide(string(data), animXML)
	e.parts.Set(slideRef.Part, []byte(slideXML))
	return nil
}

// buildAnimationXML generates a single p:par animation entry.
func buildAnimationXML(shapeID int, effect, trigger string, durationMS, delayMS int) string {
	if durationMS <= 0 {
		durationMS = animPresetDefaultFade * 50 //nolint:mnd // 500ms default duration
	}
	nodeType := animationNodeType(trigger)
	presetClass := animationPresetClass(effect)
	presetID := animationPresetID(effect)
	return fmt.Sprintf(`<p:par>
  <p:cTn presetID="%d" presetClass="%s" fill="hold" nodeType="%s">
    <p:stCondLst><p:cond delay="%d"/></p:stCondLst>
    <p:childTnLst>
      <p:set>
        <p:cBhvr>
          <p:cTn dur="%d" fill="hold"><p:stCondLst><p:cond delay="0"/></p:stCondLst></p:cTn>
          <p:tgtEl><p:spTgt spid="%d"/></p:tgtEl>
        </p:cBhvr>
      </p:set>
    </p:childTnLst>
  </p:cTn>
</p:par>`, presetID, presetClass, nodeType, delayMS, durationMS, shapeID)
}

func animationNodeType(trigger string) string {
	switch trigger {
	case "withPrev":
		return "withEffect"
	case "afterPrev":
		return "afterEffect"
	default:
		return "clickEffect"
	}
}

func animationPresetClass(effect string) string {
	switch {
	case strings.HasPrefix(effect, "entr_"):
		return "entr"
	case strings.HasPrefix(effect, "exit_"):
		return "exit"
	case strings.HasPrefix(effect, "emph_"):
		return "emph"
	case strings.HasPrefix(effect, "path_"):
		return "path"
	default:
		return "entr"
	}
}

//nolint:cyclop,funlen,mnd // OOXML preset ID table; flat switch is the clearest representation.
func animationPresetID(effect string) int {
	switch effect {
	case "entr_appear", "exit_disappear":
		return 1
	case "entr_flyIn", "exit_flyOut":
		return 2
	case "entr_fade", "exit_fade":
		return animPresetDefaultFade
	case "entr_float", "exit_float":
		return 14
	case "entr_randomBars":
		return 15
	case "entr_split":
		return 16
	case "entr_shape":
		return 17
	case "entr_swivel":
		return 19
	case "entr_wheel":
		return 21
	case "entr_wipe":
		return 22
	case "entr_zoom":
		return 23
	case "entr_bounce":
		return 25
	case "entr_growAndTurn":
		return 26
	case "emph_pulse":
		return 31
	case "emph_colorPulse":
		return 32
	case "emph_teeter":
		return 33
	case "emph_spin":
		return 34
	case "emph_growShrink":
		return 35
	case "emph_desaturate":
		return 36
	case "emph_darken":
		return 37
	case "emph_lighten":
		return 38
	case "emph_transparency":
		return 39
	case "emph_objectColor":
		return 40
	case "path_lines":
		return 42
	case "path_arcs":
		return 43
	case "path_turns":
		return 44
	case "path_shapes":
		return 45
	case "path_loops":
		return 46
	case "path_custom":
		return 47
	default:
		return animPresetDefaultFade
	}
}

// injectAnimationIntoSlide inserts the animation XML into the slide's p:timing block.
// If no timing block exists, one is created.
func injectAnimationIntoSlide(slideXML, animXML string) string {
	const seqOpen = `<p:seq concurrent="1" nextAc="seek">`
	const seqClose = `</p:seq>`
	const timingClose = `</p:timing>`

	if idx := strings.Index(slideXML, seqOpen); idx >= 0 {
		const childTnLstClose = `</p:childTnLst>`
		seqEnd := strings.Index(slideXML[idx:], seqClose)
		if seqEnd >= 0 {
			insertPoint := idx + seqEnd
			childEnd := strings.LastIndex(slideXML[:insertPoint], childTnLstClose)
			if childEnd >= 0 {
				return slideXML[:childEnd] + animXML + slideXML[childEnd:]
			}
		}
	}

	if before, after, ok := strings.Cut(slideXML, timingClose); ok {
		return before + buildFullTimingXML(animXML) + after
	}

	const closeSld = "</p:sld>"
	timingXML := buildFullTimingXML(animXML)
	return strings.Replace(slideXML, closeSld, timingXML+closeSld, 1)
}

func buildFullTimingXML(animXML string) string {
	return `<p:timing>
  <p:tnLst>
    <p:par>
      <p:cTn id="1" dur="indefinite" restart="whenNotActive" nodeType="tmRoot">
        <p:childTnLst>
          <p:seq concurrent="1" nextAc="seek">
            <p:cTn dur="indefinite" nodeType="mainSeq">
              <p:childTnLst>` +
		animXML +
		`</p:childTnLst>
            </p:cTn>
            <p:prevCondLst><p:cond evt="onPrevClick" delay="0"><p:tn/></p:cond></p:prevCondLst>
            <p:nextCondLst><p:cond evt="onNextClick" delay="0"><p:tn/></p:cond></p:nextCondLst>
          </p:seq>
        </p:childTnLst>
      </p:cTn>
    </p:par>
  </p:tnLst>
</p:timing>`
}
