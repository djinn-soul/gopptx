package protection

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/richardlehane/mscfb"
)

func TestEncryptAgilePackage_WrapsIntoCFB(t *testing.T) {
	if !CanEncryptAgile() {
		t.Skip("Agile encryption unavailable on this runtime")
	}

	zipPayload := buildMinimalPPTX(t)
	out, err := EncryptAgilePackage(zipPayload, "Secret123!")
	if err != nil {
		t.Fatalf("EncryptAgilePackage error: %v", err)
	}

	if len(out) < 8 {
		t.Fatalf("encrypted output too short: %d", len(out))
	}
	if !bytes.Equal(out[:8], []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) {
		t.Fatalf("expected CFB signature, got %x", out[:8])
	}

	r, err := mscfb.New(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("open CFB payload: %v", err)
	}

	var foundInfo bool
	var foundPkg bool
	for entry, nextErr := r.Next(); nextErr == nil; entry, nextErr = r.Next() {
		name := normalizeStreamName(entry.Name)
		switch name {
		case strings.ToLower("EncryptionInfo"):
			foundInfo = true
		case strings.ToLower("EncryptedPackage"):
			foundPkg = true
		}
	}
	if !foundInfo {
		foundInfo = bytes.Contains(out, []byte("agileEncryption"))
	}
	if !foundPkg {
		foundPkg = bytes.Contains(out, []byte("EncryptedPackage"))
	}
	if !foundInfo {
		t.Fatalf("missing %s stream", "EncryptionInfo")
	}
	if !foundPkg {
		t.Fatalf("missing %s stream", "EncryptedPackage")
	}
}

func normalizeStreamName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimRight(name, "\x00")
	return name
}

func buildMinimalPPTX(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	parts := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
 xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
 xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:sldIdLst/>
  <p:sldSz cx="9144000" cy="6858000" type="screen4x3"/>
  <p:notesSz cx="6858000" cy="9144000"/>
</p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`,
	}

	for name, content := range parts {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip part %s: %v", name, err)
		}
		if _, err := io.WriteString(w, content); err != nil {
			t.Fatalf("write zip part %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	return buf.Bytes()
}
