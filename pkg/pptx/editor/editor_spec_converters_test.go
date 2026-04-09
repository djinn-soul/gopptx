package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestEditorParagraphToSpec_RoundTripsSpaceBeforePts(t *testing.T) {
	spaceBefore := 1200
	p := &common.Paragraph{
		SpaceBeforePts: &spaceBefore,
	}
	xml, err := renderParagraphPropsXML(p)
	if err != nil {
		t.Fatalf("renderParagraphPropsXML failed: %v", err)
	}
	if !strings.Contains(xml, `<a:spcBef><a:spcPts val="1200"/></a:spcBef>`) {
		t.Fatalf("expected spcBef val=1200 round-trip, got: %s", xml)
	}
}
