package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// benchLargeDeckSlides is the slide count of the fixture built by
// writeLargeBenchDeck. Real decks measured from docs/assets/pptx carry 27-51
// package parts with 66-94% of them under 2KB; 40 slides reproduces that
// distribution at a size where per-entry costs are visible above timer noise.
const benchLargeDeckSlides = 40

// benchLayoutCount matches the twelve layouts PowerPoint's default master ships.
const benchLayoutCount = 12

// writeLargeBenchDeck builds a package shaped like a real deck: a master, a full
// layout set, per-slide notes, and the relationship parts each of those needs.
//
// The pre-existing writeBenchDeck fixture is 8 parts totalling ~3KB, which is
// small enough that per-entry work in the zip writer is invisible. Save
// benchmarks running on it cannot detect a regression in the path they cover.
func writeLargeBenchDeck(tb testing.TB) string {
	tb.Helper()
	path := filepath.Join(tb.TempDir(), "large-bench.pptx")
	files := make(map[string]string, benchLargeDeckSlides*4+benchLayoutCount*2+8)

	files["_rels/.rels"] = relsXML(
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>`,
	)
	files["docProps/core.xml"] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" ` +
		`xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title>Bench</dc:title></cp:coreProperties>`
	files["docProps/app.xml"] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties">` +
		`<Application>gopptx</Application></Properties>`
	files["ppt/theme/theme1.xml"] = benchThemeXML()

	var contentTypes, presRels, sldIDs, masterLayoutRels strings.Builder
	contentTypes.WriteString(
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.` +
			`presentationml.presentation.main+xml"/>`,
	)

	for i := 1; i <= benchLayoutCount; i++ {
		files[fmt.Sprintf("ppt/slideLayouts/slideLayout%d.xml", i)] = benchLayoutXML(i)
		files[fmt.Sprintf("ppt/slideLayouts/_rels/slideLayout%d.xml.rels", i)] = relsXML(
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/` +
				`slideMaster" Target="../slideMasters/slideMaster1.xml"/>`,
		)
		fmt.Fprintf(&contentTypes,
			`<Override PartName="/ppt/slideLayouts/slideLayout%d.xml" ContentType="application/vnd.`+
				`openxmlformats-officedocument.presentationml.slideLayout+xml"/>`, i)
		fmt.Fprintf(&masterLayoutRels,
			`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/`+
				`slideLayout" Target="../slideLayouts/slideLayout%d.xml"/>`, i, i)
	}

	files["ppt/slideMasters/slideMaster1.xml"] = benchMasterXML()
	files["ppt/slideMasters/_rels/slideMaster1.xml.rels"] = relsXML(masterLayoutRels.String() +
		fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/`+
			`relationships/theme" Target="../theme/theme1.xml"/>`, benchLayoutCount+1))

	for i := 1; i <= benchLargeDeckSlides; i++ {
		files[fmt.Sprintf("ppt/slides/slide%d.xml", i)] = benchSlideXML(i)
		files[fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", i)] = relsXML(fmt.Sprintf(
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/`+
				`slideLayout" Target="../slideLayouts/slideLayout%d.xml"/>`+
				`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/`+
				`notesSlide" Target="../notesSlides/notesSlide%d.xml"/>`,
			(i-1)%benchLayoutCount+1, i))
		files[fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", i)] = benchNotesXML(i)
		files[fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", i)] = relsXML(fmt.Sprintf(
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/`+
				`slide" Target="../slides/slide%d.xml"/>`, i))

		fmt.Fprintf(&contentTypes,
			`<Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-`+
				`officedocument.presentationml.slide+xml"/>`, i)
		fmt.Fprintf(&presRels,
			`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/`+
				`slide" Target="slides/slide%d.xml"/>`, i, i)
		fmt.Fprintf(&sldIDs, `<p:sldId id="%d" r:id="rId%d"/>`, 255+i, i)
	}

	files["[Content_Types].xml"] = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
		contentTypes.String() + `</Types>`
	files["ppt/_rels/presentation.xml.rels"] = relsXML(presRels.String())
	files["ppt/presentation.xml"] = xmlDecl + nsPresentation +
		`<p:sldIdLst>` + sldIDs.String() + `</p:sldIdLst>` +
		`<p:sldSz cx="12192000" cy="6858000"/></p:presentation>`

	if err := writeZipFixture(path, files); err != nil {
		tb.Fatalf("write large fixture: %v", err)
	}
	return path
}

const (
	xmlDecl = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
	nsA     = `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`
	nsR     = `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`
	nsP     = `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`

	nsPresentation = `<p:presentation ` + nsA + ` ` + nsR + ` ` + nsP + `>`
)

func relsXML(inner string) string {
	return xmlDecl + `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		inner + `</Relationships>`
}

// benchSlideXML renders a slide carrying a title and a multi-run body, so the
// part lands in the 2-6KB range typical of authored content rather than the
// near-empty stub the small fixture uses.
func benchSlideXML(n int) string {
	var body strings.Builder
	for i := range 14 {
		fmt.Fprintf(&body,
			`<a:p><a:r><a:rPr lang="en-US" sz="1800" dirty="0"/><a:t>Slide %d bullet %d with enough text `+
				`to resemble authored content.</a:t></a:r></a:p>`, n, i+1)
	}
	return xmlDecl + `<p:sld ` + nsA + ` ` + nsR + ` ` + nsP + `><p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title"/><p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>` +
		`<p:nvPr><p:ph type="ctrTitle"/></p:nvPr></p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/><a:lstStyle/>` +
		fmt.Sprintf(`<a:p><a:r><a:rPr lang="en-US"/><a:t>Benchmark slide %d</a:t></a:r></a:p>`, n) +
		`</p:txBody></p:sp>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="3" name="Body"/><p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>` +
		`<p:nvPr><p:ph type="body" idx="1"/></p:nvPr></p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/><a:lstStyle/>` +
		body.String() + `</p:txBody></p:sp>` +
		`</p:spTree></p:cSld></p:sld>`
}

func benchNotesXML(n int) string {
	return xmlDecl + `<p:notes ` + nsA + ` ` + nsR + ` ` + nsP + `><p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Notes"/><p:cNvSpPr/><p:nvPr><p:ph type="body" idx="1"/></p:nvPr>` +
		`</p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/><a:lstStyle/>` +
		fmt.Sprintf(`<a:p><a:r><a:rPr lang="en-US"/><a:t>Speaker notes for slide %d.</a:t></a:r></a:p>`, n) +
		`</p:txBody></p:sp></p:spTree></p:cSld></p:notes>`
}

func benchLayoutXML(n int) string {
	return xmlDecl + `<p:sldLayout ` + nsA + ` ` + nsR + ` ` + nsP + ` type="obj" preserve="1"><p:cSld>` +
		fmt.Sprintf(`<p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name="Layout %d"/><p:cNvGrpSpPr/><p:nvPr/>`, n) +
		`</p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title Placeholder"/><p:cNvSpPr><a:spLocks noGrp="1"/>` +
		`</p:cNvSpPr><p:nvPr><p:ph type="title"/></p:nvPr></p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/>` +
		`<a:lstStyle/><a:p><a:endParaRPr lang="en-US"/></a:p></p:txBody></p:sp>` +
		`</p:spTree></p:cSld><p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr></p:sldLayout>`
}

func benchMasterXML() string {
	return xmlDecl + `<p:sldMaster ` + nsA + ` ` + nsR + ` ` + nsP + `><p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`</p:spTree></p:cSld><p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" ` +
		`accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" ` +
		`hlink="hlink" folHlink="folHlink"/></p:sldMaster>`
}

// benchThemeXML approximates the ~8KB theme part PowerPoint writes: the full
// colour scheme plus enough fill/line style entries to reach a realistic size.
func benchThemeXML() string {
	var styles strings.Builder
	for i := range 8 {
		fmt.Fprintf(&styles,
			`<a:ln w="%d" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill>`+
				`<a:prstDash val="solid"/></a:ln>`, 6350*(i+1))
	}
	return xmlDecl + `<a:theme ` + nsA + ` name="BenchTheme"><a:themeElements><a:clrScheme name="Bench">` +
		`<a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1>` +
		`<a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1>` +
		`<a:dk2><a:srgbClr val="44546A"/></a:dk2><a:lt2><a:srgbClr val="E7E6E6"/></a:lt2>` +
		`<a:accent1><a:srgbClr val="4472C4"/></a:accent1><a:accent2><a:srgbClr val="ED7D31"/></a:accent2>` +
		`<a:accent3><a:srgbClr val="A5A5A5"/></a:accent3><a:accent4><a:srgbClr val="FFC000"/></a:accent4>` +
		`<a:accent5><a:srgbClr val="5B9BD5"/></a:accent5><a:accent6><a:srgbClr val="70AD47"/></a:accent6>` +
		`<a:hlink><a:srgbClr val="0563C1"/></a:hlink><a:folHlink><a:srgbClr val="954F72"/></a:folHlink>` +
		`</a:clrScheme><a:fontScheme name="Bench">` +
		`<a:majorFont><a:latin typeface="Calibri Light"/><a:ea typeface=""/><a:cs typeface=""/></a:majorFont>` +
		`<a:minorFont><a:latin typeface="Calibri"/><a:ea typeface=""/><a:cs typeface=""/></a:minorFont>` +
		`</a:fontScheme><a:fmtScheme name="Bench"><a:fillStyleLst>` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:fillStyleLst>` +
		`<a:lnStyleLst>` + styles.String() + `</a:lnStyleLst>` +
		`<a:effectStyleLst><a:effectStyle><a:effectLst/></a:effectStyle>` +
		`<a:effectStyle><a:effectLst/></a:effectStyle>` +
		`<a:effectStyle><a:effectLst/></a:effectStyle></a:effectStyleLst>` +
		`<a:bgFillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill>` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:bgFillStyleLst>` +
		`</a:fmtScheme></a:themeElements></a:theme>`
}

func openLargeBenchEditor(tb testing.TB) *PresentationEditor {
	tb.Helper()
	editor, err := OpenPresentationEditor(writeLargeBenchDeck(tb))
	if err != nil {
		tb.Fatalf("open large bench editor: %v", err)
	}
	return editor
}
