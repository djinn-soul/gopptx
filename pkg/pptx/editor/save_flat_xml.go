package editor

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// Flat XML is PowerPoint's "PowerPoint XML Presentation" (ppSaveAsXMLPresentation):
// the whole package as one XML document, each part a <pkg:part>. It is what
// upstream issue #1059 asks for -- a deck that can be read and diffed without
// unzipping anything.
const (
	flatXMLPackageNS   = "http://schemas.microsoft.com/office/2006/xmlPackage"
	flatXMLProgIDPI    = `<?mso-application progid="PowerPoint.Show"?>`
	flatXMLDeclaration = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
	// base64LineWidth matches the wrapping PowerPoint itself writes, which keeps
	// the output diffable rather than one enormous line.
	base64LineWidth = 76
	// defaultPartContentType is used for a part no content type covers. Nothing
	// in a package written by this library should need it.
	defaultPartContentType = "application/octet-stream"
	contentTypesPartName   = "[Content_Types].xml"
)

// SaveFlatXML writes the presentation as a single PowerPoint XML Presentation
// file instead of a zipped package.
//
// PowerPoint opens the result directly; there is no separate reader for it here,
// so a deck saved this way is meant for inspection and version control, not as
// an input format.
func (e *PresentationEditor) SaveFlatXML(filePath string) error {
	if e == nil {
		return errors.New("nil editor")
	}
	pkg, err := e.SaveToBytes()
	if err != nil {
		return err
	}
	flat, err := FlatXMLFromPackage(pkg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, flat, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}
	return nil
}

// FlatXMLFromPackage converts a PPTX package into its flat XML form.
func FlatXMLFromPackage(pkg []byte) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(pkg), int64(len(pkg)))
	if err != nil {
		return nil, fmt.Errorf("read package: %w", err)
	}

	parts := make(map[string][]byte, len(reader.File))
	order := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		data, readErr := readFlatXMLZipEntry(file)
		if readErr != nil {
			return nil, readErr
		}
		if _, seen := parts[file.Name]; !seen {
			order = append(order, file.Name)
		}
		parts[file.Name] = data
	}

	types, err := parseFlatXMLContentTypes(parts[contentTypesPartName])
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	b.WriteString(flatXMLDeclaration + "\n")
	b.WriteString(flatXMLProgIDPI + "\n")
	b.WriteString(`<pkg:package xmlns:pkg="` + flatXMLPackageNS + `">` + "\n")
	for _, name := range order {
		// The flat form carries each part's content type on the part itself, so
		// [Content_Types].xml has no place in it -- PowerPoint rejects a
		// document that still contains one.
		if name == contentTypesPartName {
			continue
		}
		writeFlatXMLPart(&b, name, parts[name], types.contentTypeFor(name))
	}
	b.WriteString("</pkg:package>\n")
	return []byte(b.String()), nil
}

func readFlatXMLZipEntry(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", file.Name, err)
	}
	defer func() { _ = rc.Close() }()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", file.Name, err)
	}
	return data, nil
}

// writeFlatXMLPart writes one <pkg:part>. XML parts are inlined so they stay
// readable; everything else is base64, which is the only way binary survives an
// XML document.
func writeFlatXMLPart(b *strings.Builder, name string, data []byte, contentType string) {
	b.WriteString(`<pkg:part pkg:name="/` + xmlAttrEscape(name) + `" pkg:contentType="` +
		xmlAttrEscape(contentType) + `"`)
	if !isFlatXMLTextPart(contentType) {
		b.WriteString(` pkg:compression="store"><pkg:binaryData>` + "\n")
		b.WriteString(wrapBase64(data))
		b.WriteString("</pkg:binaryData></pkg:part>\n")
		return
	}
	b.WriteString("><pkg:xmlData>\n")
	b.Write(stripXMLDeclaration(data))
	b.WriteString("\n</pkg:xmlData></pkg:part>\n")
}

// isFlatXMLTextPart reports whether a part is XML, and so can be inlined rather
// than base64 encoded.
func isFlatXMLTextPart(contentType string) bool {
	lower := strings.ToLower(contentType)
	return strings.HasSuffix(lower, "+xml") ||
		lower == "text/xml" ||
		lower == "application/xml"
}

// stripXMLDeclaration removes a part's own XML declaration: the flat document
// already has one and a second is not well-formed.
func stripXMLDeclaration(data []byte) []byte {
	trimmed := bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	trimmed = bytes.TrimLeft(trimmed, " \t\r\n")
	if !bytes.HasPrefix(trimmed, []byte("<?xml")) {
		return trimmed
	}
	if _, after, found := bytes.Cut(trimmed, []byte("?>")); found {
		return bytes.TrimLeft(after, " \t\r\n")
	}
	return trimmed
}

func wrapBase64(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	var b strings.Builder
	for start := 0; start < len(encoded); start += base64LineWidth {
		end := min(start+base64LineWidth, len(encoded))
		b.WriteString(encoded[start:end])
		b.WriteString("\n")
	}
	return b.String()
}

func xmlAttrEscape(value string) string {
	var b bytes.Buffer
	if err := xml.EscapeText(&b, []byte(value)); err != nil {
		return value
	}
	return b.String()
}

// flatXMLContentTypes resolves a part name to its content type the way
// [Content_Types].xml does: an Override wins, otherwise the extension Default.
type flatXMLContentTypes struct {
	defaults  map[string]string
	overrides map[string]string
}

func parseFlatXMLContentTypes(data []byte) (*flatXMLContentTypes, error) {
	types := &flatXMLContentTypes{
		defaults:  map[string]string{},
		overrides: map[string]string{},
	}
	if len(data) == 0 {
		return types, nil
	}
	var doc struct {
		Defaults []struct {
			Extension   string `xml:"Extension,attr"`
			ContentType string `xml:"ContentType,attr"`
		} `xml:"Default"`
		Overrides []struct {
			PartName    string `xml:"PartName,attr"`
			ContentType string `xml:"ContentType,attr"`
		} `xml:"Override"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse content types: %w", err)
	}
	for _, d := range doc.Defaults {
		types.defaults[strings.ToLower(d.Extension)] = d.ContentType
	}
	for _, o := range doc.Overrides {
		types.overrides[strings.TrimPrefix(o.PartName, "/")] = o.ContentType
	}
	return types, nil
}

func (t *flatXMLContentTypes) contentTypeFor(name string) string {
	if contentType, ok := t.overrides[name]; ok {
		return contentType
	}
	ext := strings.TrimPrefix(strings.ToLower(path.Ext(name)), ".")
	if contentType, ok := t.defaults[ext]; ok {
		return contentType
	}
	if name == contentTypesPartName {
		return "text/xml"
	}
	return defaultPartContentType
}
