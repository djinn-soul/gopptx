package presentation

import (
	"archive/zip"
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
)

// buildPackageParts generates a deck and returns its parts by name.
func buildPackageParts(t *testing.T, meta Metadata, slides []elements.SlideContent) map[string]string {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if err := WritePresentationPackage(zw, meta, slides, len(slides)); err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip open failed: %v", err)
	}

	parts := make(map[string]string, len(zr.File))
	for _, f := range zr.File {
		rc, openErr := f.Open()
		if openErr != nil {
			t.Fatalf("open %s failed: %v", f.Name, openErr)
		}
		data, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			t.Fatalf("read %s failed: %v", f.Name, readErr)
		}
		parts[f.Name] = string(data)
	}
	return parts
}

// The handout master part, its content type and its relationship were all
// written, but presentation.xml never declared the ID list — so PowerPoint used
// its own built-in handout instead of the one gopptx generated.
func TestHandoutMasterIsDeclaredInPresentation(t *testing.T) {
	meta := Metadata{
		Metadata:      common.Metadata{Title: "Handout"},
		HandoutMaster: handout.New(),
	}
	parts := buildPackageParts(t, meta, []elements.SlideContent{elements.NewSlide("S1")})

	if _, ok := parts["ppt/handoutMasters/handoutMaster1.xml"]; !ok {
		t.Fatal("handout master part missing")
	}

	presentationXML := parts["ppt/presentation.xml"]
	if !strings.Contains(presentationXML, "<p:handoutMasterIdLst>") {
		t.Fatalf("presentation.xml has no handoutMasterIdLst: %s", presentationXML)
	}

	idMatch := regexp.MustCompile(`<p:handoutMasterId r:id="(rId\d+)"/>`).FindStringSubmatch(presentationXML)
	if idMatch == nil {
		t.Fatalf("no p:handoutMasterId element: %s", presentationXML)
	}

	// The declared rId must be the one the rels file binds to the handout part.
	rels := parts["ppt/_rels/presentation.xml.rels"]
	relPattern := regexp.MustCompile(
		`Id="` + idMatch[1] + `"[^>]*Target="handoutMasters/handoutMaster1\.xml"`,
	)
	if !relPattern.MatchString(rels) {
		t.Fatalf("%s does not target the handout master: %s", idMatch[1], rels)
	}
}

func TestNoHandoutMasterIdListWithoutHandout(t *testing.T) {
	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{elements.NewSlide("S1")})
	if strings.Contains(parts["ppt/presentation.xml"], "handoutMasterIdLst") {
		t.Fatal("handoutMasterIdLst emitted for a deck with no handout master")
	}
}

// TestHandoutMasterIDSurvivesCustomXML pins the relationship arithmetic against
// custom XML. The rels file writes one presentation-level relationship per
// custom XML item, and the ID bookkeeping used to reserve two, so the handout
// master's declared rId pointed one past its own relationship for every item
// present — leaving PowerPoint with a handout declaration it cannot resolve.
func TestHandoutMasterIDSurvivesCustomXML(t *testing.T) {
	meta := Metadata{
		Metadata: common.Metadata{
			Title: "Handout+CustomXML",
			CustomXML: []common.CustomXMLPart{
				{RootElement: "one", Namespace: "urn:gopptx:test:one", Content: "<a/>"},
				{RootElement: "two", Namespace: "urn:gopptx:test:two", Content: "<b/>"},
			},
		},
		HandoutMaster: handout.New(),
	}
	parts := buildPackageParts(t, meta, []elements.SlideContent{elements.NewSlide("S1")})

	presentationXML := parts["ppt/presentation.xml"]
	idMatch := regexp.MustCompile(`<p:handoutMasterId r:id="(rId\d+)"/>`).FindStringSubmatch(presentationXML)
	if idMatch == nil {
		t.Fatalf("no p:handoutMasterId element: %s", presentationXML)
	}

	rels := parts["ppt/_rels/presentation.xml.rels"]
	relPattern := regexp.MustCompile(
		`Id="` + idMatch[1] + `"[^>]*Target="handoutMasters/handoutMaster1\.xml"`,
	)
	if !relPattern.MatchString(rels) {
		t.Fatalf("declared %s does not target the handout master with custom XML present:\n%s", idMatch[1], rels)
	}
}
