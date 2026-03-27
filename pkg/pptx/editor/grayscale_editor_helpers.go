package editor

import (
	"bytes"
	"regexp"

	editorgrayscale "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/grayscale"
)

var (
	bgSectionPattern = regexp.MustCompile(`(?s)<p:bg\b.*?</p:bg>`)
	bgColorPattern   = regexp.MustCompile(`(?i)(<a:srgbClr[^>]*\bval=")([0-9a-f]{6})(")`)
	bgEmbedPattern   = regexp.MustCompile(`(?s)<p:bg\b.*?r:embed="([^"]+)"`)
	jpegExtPattern   = regexp.MustCompile(`(?i)\.jpe?g$`)
)

const backgroundColorSubmatchCount = 4

func grayscaleBackgroundXML(content []byte) ([]byte, bool) {
	match := bgSectionPattern.Find(content)
	if len(match) == 0 {
		return content, false
	}
	changed := false
	replaced := bgColorPattern.ReplaceAllFunc(match, func(segment []byte) []byte {
		submatches := bgColorPattern.FindSubmatch(segment)
		if len(submatches) != backgroundColorSubmatchCount {
			return segment
		}
		gray, err := editorgrayscale.HexColor(string(submatches[2]))
		if err != nil {
			return segment
		}
		changed = true
		return bytes.Replace(segment, submatches[2], []byte(gray), 1)
	})
	if !changed {
		return content, false
	}
	return bytes.Replace(content, match, replaced, 1), true
}

func imageFormatFromTarget(target string) string {
	if jpegExtPattern.MatchString(target) {
		return formatJPEG
	}
	return formatPNG
}
