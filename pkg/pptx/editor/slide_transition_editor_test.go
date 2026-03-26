package editor

import (
	"strings"
	"testing"
)

func TestBuildTransitionXMLDisableAdvanceOnClick(t *testing.T) {
	xml := buildTransitionXML("fade", 0, 20000, true)
	if !strings.Contains(xml, `advTm="20000"`) {
		t.Fatalf("expected advTm in transition XML, got %s", xml)
	}
	if !strings.Contains(xml, `advClick="0"`) {
		t.Fatalf("expected advClick=0 in transition XML, got %s", xml)
	}
}

func TestBuildTransitionXMLKeepsClickAdvanceByDefault(t *testing.T) {
	xml := buildTransitionXML("fade", 0, 20000, false)
	if !strings.Contains(xml, `advTm="20000"`) {
		t.Fatalf("expected advTm in transition XML, got %s", xml)
	}
	if strings.Contains(xml, `advClick="0"`) {
		t.Fatalf("did not expect advClick=0 when disable flag is false, got %s", xml)
	}
}
