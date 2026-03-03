package tplx

import (
	"strings"
	"testing"
)

func TestMergeAdjacentRunsPreservesNamespacePrefixes(t *testing.T) {
	xmlIn := `<?xml version="1.0" encoding="UTF-8"?>
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
       xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:sp>
        <p:txBody>
          <a:p>
            <a:r><a:t>{{na</a:t></a:r>
            <a:r><a:t>me}}</a:t></a:r>
          </a:p>
        </p:txBody>
      </p:sp>
    </p:spTree>
  </p:cSld>
</p:sld>`

	out := mergeAdjacentRuns([]byte(xmlIn))
	s := string(out)
	if strings.Contains(s, "<http://") {
		t.Fatalf("invalid namespace URI-as-prefix detected: %s", s)
	}
	if !strings.Contains(s, "<p:sld") {
		t.Fatalf("expected p:sld element, got: %s", s)
	}
	if !strings.Contains(s, "<a:p>") {
		t.Fatalf("expected drawing paragraph tag, got: %s", s)
	}
}
