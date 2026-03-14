package media

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var cTnIDPattern = regexp.MustCompile(`<p:cTn[^>]*\sid="([0-9]+)"`)
var mediaShapePatternTemplate = `(?s)<p:pic\b[^>]*>.*?<p:cNvPr\b[^>]*\bid="%d"[^>]*>.*?</p:cNvPr>.*?<p14:media\b[^>]*\br:embed="([^"]+)"[^>]*/>.*?</p:pic>`

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

func ApplyMediaTiming(content []byte, mediaKind string, shapeID int, options MediaTimingOptions) ([]byte, error) {
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
		return nil, fmt.Errorf("invalid slide xml: timing block not found after insertion")
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
		return "", fmt.Errorf("invalid slide xml: missing </p:sld>")
	}
	return slideXML[:idx] + timingBlock + slideXML[idx:], nil
}

func nextTimingCTnID(timingXML string) int {
	matches := cTnIDPattern.FindAllStringSubmatch(timingXML, -1)
	maxID := 2
	for _, match := range matches {
		if len(match) < 2 {
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
		return "", fmt.Errorf("invalid timing xml: mainSeq cTn not found")
	}
	childStart := strings.Index(timingXML[mainSeqIdx:], "<p:childTnLst>")
	if childStart < 0 {
		return "", fmt.Errorf("invalid timing xml: mainSeq childTnLst not found")
	}
	childStart += mainSeqIdx
	childEnd := strings.Index(timingXML[childStart:], "</p:childTnLst>")
	if childEnd < 0 {
		return "", fmt.Errorf("invalid timing xml: mainSeq childTnLst end not found")
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
		numSlides := options.SlideCount - options.SlideIndex
		if numSlides < 1 {
			numSlides = 1
		}
		numSldAttr = fmt.Sprintf(` numSld="%d"`, numSlides)
	}

	extXML := ""
	if strings.TrimSpace(mediaRelID) != "" {
		extXML = fmt.Sprintf(`
                  <p:extLst>
                    <p:ext uri="{DAA4B4D4-6D71-4841-9C94-3DE7FCFB9230}">
                      <p14:media xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" r:embed="%s"/>
                    </p:ext>
                    <p:ext uri="{EFAFB233-063F-42B5-8137-9DF3F51BA10A}">
                      <p15:media xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2012/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" r:embed="%s"/>
                    </p:ext>
                  </p:extLst>`,
			escapeXMLAttr(mediaRelID),
			escapeXMLAttr(mediaRelID),
		)
	}

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
	if len(match) < 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func normalizeMediaVolume(volume uint32) uint32 {
	if volume > 100 {
		volume = 100
	}
	return volume * 1000
}
