package media

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var cTnIDPattern = regexp.MustCompile(`<p:cTn[^>]*\sid="([0-9]+)"`)

const mediaShapePatternTemplate = `(?s)<p:pic\b[^>]*>.*?<p:cNvPr\b[^>]*\bid="%d"[^>]*>.*?</p:cNvPr>.*?<p14:media\b[^>]*\br:embed="([^"]+)"[^>]*/>.*?</p:pic>`
const (
	minRegexSubmatchCount = 2
	maxVolumePercent      = 100
	volumeScaleThousand   = 1000
)

//nolint:revive // Exported name kept for API compatibility across editor and bindings.
type MediaTimingOptions struct {
	AutoPlay         bool
	LoopPlayback     bool
	Muted            bool
	Volume           uint32
	ShowWhenStopped  bool
	PlayAcrossSlides bool
	SlideIndex       int
	SlideCount       int
}

func ApplyMediaTiming(
	content []byte,
	mediaKind string,
	shapeID int,
	options MediaTimingOptions,
) ([]byte, error) {
	slideXML := string(content)
	if !strings.Contains(slideXML, "<p:timing>") {
		withTiming, err := addDefaultTimingBlock(slideXML)
		if err != nil {
			return nil, err
		}
		slideXML = withTiming
	}

	timingStart := strings.Index(slideXML, "<p:timing>")
	timingEnd := strings.Index(slideXML, "</p:timing>")
	if timingStart < 0 || timingEnd < 0 {
		return nil, errors.New("invalid slide xml: timing block not found after insertion")
	}
	timingEnd += len("</p:timing>")
	timingXML := slideXML[timingStart:timingEnd]

	nextID := nextTimingCTnID(timingXML)
	mediaRelID := mediaRelIDForShape(slideXML, shapeID)
	mediaNode := buildMediaTimingNode(mediaKind, shapeID, nextID, mediaRelID, options)

	updatedTiming, err := insertMediaNodeIntoMainSeq(timingXML, mediaNode)
	if err != nil {
		return nil, err
	}

	updatedSlide := slideXML[:timingStart] + updatedTiming + slideXML[timingEnd:]
	return []byte(updatedSlide), nil
}

func addDefaultTimingBlock(slideXML string) (string, error) {
	const timingBlock = `
<p:timing>
  <p:tnLst>
    <p:par>
      <p:cTn id="1" dur="indefinite" restart="never" nodeType="tmRoot">
        <p:childTnLst>
          <p:seq concurrent="1" nextAc="seek">
            <p:cTn id="2" dur="indefinite" nodeType="mainSeq">
              <p:childTnLst>
              </p:childTnLst>
            </p:cTn>
          </p:seq>
        </p:childTnLst>
      </p:cTn>
    </p:par>
  </p:tnLst>
</p:timing>`
	const closeSlide = "</p:sld>"
	idx := strings.LastIndex(slideXML, closeSlide)
	if idx < 0 {
		return "", errors.New("invalid slide xml: missing </p:sld>")
	}
	return slideXML[:idx] + timingBlock + slideXML[idx:], nil
}

func nextTimingCTnID(timingXML string) int {
	matches := cTnIDPattern.FindAllStringSubmatch(timingXML, -1)
	maxID := 2
	for _, match := range matches {
		if len(match) < minRegexSubmatchCount {
			continue
		}
		idValue, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		if idValue > maxID {
			maxID = idValue
		}
	}
	return maxID + 1
}

func insertMediaNodeIntoMainSeq(timingXML, mediaNode string) (string, error) {
	mainSeqIdx := strings.Index(timingXML, `nodeType="mainSeq"`)
	if mainSeqIdx < 0 {
		return "", errors.New("invalid timing xml: mainSeq cTn not found")
	}
	childStart := strings.Index(timingXML[mainSeqIdx:], "<p:childTnLst>")
	if childStart < 0 {
		return "", errors.New("invalid timing xml: mainSeq childTnLst not found")
	}
	childStart += mainSeqIdx
	childEnd := strings.Index(timingXML[childStart:], "</p:childTnLst>")
	if childEnd < 0 {
		return "", errors.New("invalid timing xml: mainSeq childTnLst end not found")
	}
	childEnd += childStart
	return timingXML[:childEnd] + mediaNode + timingXML[childEnd:], nil
}

func buildMediaTimingNode(
	mediaKind string,
	shapeID, cTnID int,
	mediaRelID string,
	options MediaTimingOptions,
) string {
	startDelay := "indefinite"
	if options.AutoPlay {
		startDelay = "0"
	}
	repeatAttr := ""
	if options.LoopPlayback {
		repeatAttr = ` repeatCount="indefinite"`
	}
	muteAttr := ""
	if options.Muted {
		muteAttr = ` mute="1"`
	}
	volAttr := fmt.Sprintf(` vol="%d"`, normalizeMediaVolume(options.Volume))
	showWhenStopped := "1"
	if !options.ShowWhenStopped {
		showWhenStopped = "0"
	}
	numSldAttr := ` numSld="1"`
	if options.PlayAcrossSlides {
		numSlides := max(1, options.SlideCount-options.SlideIndex)
		numSldAttr = fmt.Sprintf(` numSld="%d"`, numSlides)
	}

	// The slide shape already carries media embed references through p14:media.
	// Emitting extLst under p:cMediaNode is invalid per PresentationML schema.
	extXML := ""
	_ = mediaRelID

	return fmt.Sprintf(`
              <p:%s>
                <p:cMediaNode%s%s%s showWhenStopped="%s">
                  <p:cTn id="%d" fill="hold" display="0"%s>
                    <p:stCondLst>
                      <p:cond delay="%s"/>
                    </p:stCondLst>
                  </p:cTn>
                  <p:tgtEl>
                    <p:spTgt spid="%d"/>
                  </p:tgtEl>
%s
                </p:cMediaNode>
              </p:%s>`,
		mediaKind,
		volAttr,
		muteAttr,
		numSldAttr,
		showWhenStopped,
		cTnID,
		repeatAttr,
		startDelay,
		shapeID,
		extXML,
		mediaKind,
	)
}

func mediaRelIDForShape(slideXML string, shapeID int) string {
	pattern := fmt.Sprintf(mediaShapePatternTemplate, shapeID)
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(slideXML)
	if len(match) < minRegexSubmatchCount {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func normalizeMediaVolume(volume uint32) uint32 {
	if volume > maxVolumePercent {
		volume = maxVolumePercent
	}
	return volume * volumeScaleThousand
}
