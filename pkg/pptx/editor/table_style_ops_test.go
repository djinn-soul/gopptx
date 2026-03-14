package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	contentTypesFixtureXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
		`<Default Extension="xml" ContentType="application/xml"/>` +
		`</Types>`
	presentationRelsFixtureXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" ` +
		`Target="slides/slide1.xml"/>` +
		`</Relationships>`
	tableStylesListFixtureXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<a:tblStyleLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`def="{A1111111-1111-1111-1111-111111111111}">` +
		`<a:tblStyle styleId="{A1111111-1111-1111-1111-111111111111}" styleName="A"/>` +
		`<a:tblStyle styleId="{B2222222-2222-2222-2222-222222222222}" styleName="B"/>` +
		`</a:tblStyleLst>`
)

func TestDefineTableStyleCreatesPackageInfrastructure(t *testing.T) {
	e := newTableEditorFixture()
	e.parts.Set(
		"[Content_Types].xml",
		[]byte(contentTypesFixtureXML),
	)
	e.parts.Set(
		"ppt/_rels/presentation.xml.rels",
		[]byte(presentationRelsFixtureXML),
	)

	styleID, err := e.DefineTableStyle(common.TableStyleDefinition{
		Name: "My Style",
	})
	if err != nil {
		t.Fatalf("DefineTableStyle failed: %v", err)
	}
	if styleID == "" {
		t.Fatalf("expected style id")
	}

	tableStyles, ok := e.parts.Get("ppt/tableStyles.xml")
	if !ok {
		t.Fatalf("expected ppt/tableStyles.xml")
	}
	if !strings.Contains(string(tableStyles), `styleName="My Style"`) {
		t.Fatalf("expected style entry in table styles part: %s", string(tableStyles))
	}

	contentTypes, _ := e.parts.Get("[Content_Types].xml")
	if !strings.Contains(
		string(contentTypes),
		`application/vnd.openxmlformats-officedocument.presentationml.tableStyles+xml`,
	) {
		t.Fatalf("expected table styles content type override, got: %s", string(contentTypes))
	}

	presentationRels, _ := e.parts.Get("ppt/_rels/presentation.xml.rels")
	if !strings.Contains(
		string(presentationRels),
		`relationships/tableStyles" Target="tableStyles.xml"`,
	) {
		t.Fatalf("expected presentation table styles relationship, got: %s", string(presentationRels))
	}
}

func TestListTableStyles(t *testing.T) {
	e := newTableEditorFixture()
	e.parts.Set(
		"ppt/tableStyles.xml",
		[]byte(tableStylesListFixtureXML),
	)

	styles, err := e.ListTableStyles()
	if err != nil {
		t.Fatalf("ListTableStyles failed: %v", err)
	}
	if len(styles) != 2 {
		t.Fatalf("expected 2 styles, got %d", len(styles))
	}
	if styles[0].Name != "A" || styles[1].Name != "B" {
		t.Fatalf("unexpected styles: %#v", styles)
	}
}
