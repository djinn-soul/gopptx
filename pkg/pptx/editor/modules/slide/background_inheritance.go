package slide

import (
	"bytes"
	"errors"
	"regexp"
)

var explicitSlideBackgroundPattern = regexp.MustCompile(
	`(?s)<p:bg(?:\s[^>]*)?(?:/>|>.*?</p:bg>)`,
)

var commonSlideDataStartPattern = regexp.MustCompile(`<p:cSld(?:\s[^>]*)?>`)

const noFillSlideBackground = `<p:bg><p:bgPr><a:noFill/><a:effectLst/></p:bgPr></p:bg>`

// ParseFollowMasterBackground reports whether the slide inherits its
// background. A slide follows its master exactly when p:cSld has no p:bg.
func ParseFollowMasterBackground(content []byte) bool {
	return !explicitSlideBackgroundPattern.Match(content)
}

// RewriteFollowMasterBackground restores inheritance by removing p:bg, or
// interrupts it by adding a schema-valid no-fill background when none exists.
func RewriteFollowMasterBackground(content []byte, follow bool) ([]byte, error) {
	if follow {
		return explicitSlideBackgroundPattern.ReplaceAll(content, nil), nil
	}
	if !ParseFollowMasterBackground(content) {
		return bytes.Clone(content), nil
	}
	location := commonSlideDataStartPattern.FindIndex(content)
	if location == nil {
		return nil, errors.New("slide XML does not contain <p:cSld>")
	}
	insertAt := location[1]
	out := make([]byte, 0, len(content)+len(noFillSlideBackground))
	out = append(out, content[:insertAt]...)
	out = append(out, noFillSlideBackground...)
	out = append(out, content[insertAt:]...)
	return out, nil
}
