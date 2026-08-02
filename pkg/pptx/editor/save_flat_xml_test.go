package editor

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func zipPackage(t *testing.T, parts map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	for name, body := range parts {
		w, err := writer.Create(name)
		if err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
		if _, err = w.Write([]byte(body)); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

const flatXMLTestContentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="png" ContentType="image/png"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ` +
	`ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
</Types>`

// The flat form has to be one well-formed document, carry each part's content
// type on the part, and leave [Content_Types].xml out -- PowerPoint refuses a
// document that still contains one (upstream issue #1059).
func TestFlatXMLFromPackage(t *testing.T) {
	pkg := zipPackage(t, map[string]string{
		"[Content_Types].xml":  flatXMLTestContentTypes,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation/>`,
		"ppt/media/image1.png": "\x89PNG\r\n\x1a\n binary bytes",
	})

	flat, err := FlatXMLFromPackage(pkg)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	text := string(flat)

	if !strings.HasPrefix(text, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`+"\n"+flatXMLProgIDPI) {
		t.Fatalf("expected the declaration and the PowerPoint processing instruction:\n%s", text[:120])
	}
	if strings.Contains(text, "[Content_Types].xml") {
		t.Fatalf("expected the content types part to be left out:\n%s", text)
	}
	if !strings.Contains(text, `pkg:name="/ppt/presentation.xml"`) {
		t.Fatalf("expected the presentation part:\n%s", text)
	}
	if !strings.Contains(
		text,
		`pkg:contentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"`,
	) {
		t.Fatalf("expected the override content type on the part:\n%s", text)
	}
	if strings.Count(text, "<?xml") != 1 {
		t.Fatalf("expected exactly one XML declaration in the document:\n%s", text)
	}
	if !strings.Contains(text, `pkg:name="/ppt/media/image1.png"`) ||
		!strings.Contains(text, "<pkg:binaryData>") {
		t.Fatalf("expected the image part as base64 binary data:\n%s", text)
	}

	var doc struct {
		Parts []struct {
			Name        string `xml:"name,attr"`
			ContentType string `xml:"contentType,attr"`
		} `xml:"part"`
	}
	if err := xml.Unmarshal(flat, &doc); err != nil {
		t.Fatalf("flat document is not well-formed: %v", err)
	}
	if len(doc.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(doc.Parts))
	}
}
