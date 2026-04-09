package pptxxml

import (
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const notesMasterImageRelIDOffset = 2
const defaultNotesFontSizePt = 12

// NotesSlide renders one notes slide XML part.
func NotesSlide(paragraphs []text.Paragraph) string {
	paragraphsXML := notesParagraphsXML(paragraphs)
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:notes xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="0" cy="0"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="0" cy="0"/>
</a:xfrm>
</p:grpSpPr>
<p:sp>
<p:nvSpPr>
<p:cNvPr id="2" name="Slide Image Placeholder 1"/>
<p:cNvSpPr><a:spLocks noGrp="1" noRot="1" noChangeAspect="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="sldImg"/></p:nvPr>
</p:nvSpPr>
<p:spPr/>
</p:sp>
<p:sp>
<p:nvSpPr>
<p:cNvPr id="3" name="Notes Placeholder 2"/>
<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="body" idx="1"/></p:nvPr>
</p:nvSpPr>
<p:spPr/>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>` + paragraphsXML + `
</p:txBody>
</p:sp>
</p:spTree>
</p:cSld>
<p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:notes>`
}

// NotesSlideRelationships renders notesSlideN.xml.rels.
func NotesSlideRelationships(slideNumber int) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="../slides/slide` + strconv.Itoa(slideNumber) + `.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="../notesMasters/notesMaster1.xml"/>
</Relationships>`
}

// NotesMasterSpec defines the content for notesMaster1.xml.
type NotesMasterSpec struct {
	HeaderText   string
	FooterText   string
	ShowDateTime bool
	ShowSlideNum bool
	Background   *SlideBackgroundSpec
	NotesStyle   []TextLevelStyle
}

// NotesMaster renders a notes master part.
func NotesMaster(spec *NotesMasterSpec) string {
	if spec == nil {
		spec = &NotesMasterSpec{
			ShowDateTime: true,
			ShowSlideNum: true,
		}
	}

	notesStyleXML := ""
	if len(spec.NotesStyle) > 0 {
		notesStyleXML = `
<p:notesStyle>` + textLevelStylesXML(spec.NotesStyle) + `
</p:notesStyle>`
	}

	dateXML := ""
	if spec.ShowDateTime {
		dateXML = notesMasterDate()
	}

	slideNumXML := ""
	if spec.ShowSlideNum {
		slideNumXML = notesMasterSlideNum()
	}

	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:notesMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>` +
		backgroundXML(spec.Background) + `
<p:spTree>` +
		notesMasterCommonRootGrp() +
		notesMasterHeader(spec.HeaderText) +
		dateXML +
		notesMasterSlideImage() +
		notesMasterBody() +
		notesMasterFooter(spec.FooterText) +
		slideNumXML + `
</p:spTree>
</p:cSld>
<p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" ` +
		`accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>` + notesStyleXML + `
</p:notesMaster>`
}

func notesMasterCommonRootGrp() string {
	return `
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="0" cy="0"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="0" cy="0"/>
</a:xfrm>
</p:grpSpPr>`
}

func notesMasterHeader(text string) string {
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="2" name="Header Placeholder 1"/>
<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="hdr"/></p:nvPr>
</p:nvSpPr>
<p:spPr/>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
<a:r>
<a:rPr lang="en-US"/>
<a:t>` + Escape(text) + `</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`
}

func notesMasterDate() string {
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="3" name="Date Placeholder 2"/>
<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="dt"/></p:nvPr>
</p:nvSpPr>
<p:spPr/>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p><a:fld id="{8583B92D-B326-4076-96F9-126CB471A9B6}" type="datetime1">` +
		`<a:rPr lang="en-US"/><a:pPr/><a:t></a:t></a:fld><a:endParaRPr lang="en-US"/></a:p>
</p:txBody>
</p:sp>`
}

func notesMasterSlideImage() string {
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="4" name="Slide Image Placeholder 3"/>
<p:cNvSpPr><a:spLocks noGrp="1" noRot="1" noChangeAspect="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="sldImg"/></p:nvPr>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="1143000" y="685800"/>
<a:ext cx="4572000" cy="3429000"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
</p:spPr>
</p:sp>`
}

func notesMasterBody() string {
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="5" name="Notes Placeholder 4"/>
<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="body" idx="1"/></p:nvPr>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="1143000" y="4572000"/>
<a:ext cx="4572000" cy="3886200"/>
</a:xfrm>
</p:spPr>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p><a:endParaRPr lang="en-US"/></a:p>
</p:txBody>
</p:sp>`
}

func notesMasterFooter(text string) string {
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="6" name="Footer Placeholder 5"/>
<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="ftr" idx="3"/></p:nvPr>
</p:nvSpPr>
<p:spPr/>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
<a:r>
<a:rPr lang="en-US"/>
<a:t>` + Escape(text) + `</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`
}

func notesMasterSlideNum() string {
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="7" name="Slide Number Placeholder 6"/>
<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
<p:nvPr><p:ph type="sldNum" idx="4"/></p:nvPr>
</p:nvSpPr>
<p:spPr/>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p><a:fld id="{1E4E639B-83C3-4404-B32B-98A843E836FB}" type="slidenum">` +
		`<a:rPr lang="en-US"/><a:t>‹#›</a:t></a:fld><a:endParaRPr lang="en-US"/></a:p>
</p:txBody>
</p:sp>`
}

// NotesMasterRelationships renders notesMaster1.xml.rels.
func NotesMasterRelationships(themeIndex int, imageTargets []string) string {
	if themeIndex <= 0 {
		themeIndex = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme` + strconv.Itoa(themeIndex) + `.xml"/>`)

	for i, target := range imageTargets {
		rid := i + notesMasterImageRelIDOffset
		b.WriteString(`
<Relationship Id="rId` + strconv.Itoa(rid) + `" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="` + Escape(target) + `"/>`)
	}
	b.WriteString(`
</Relationships>`)
	return b.String()
}

func notesParagraphsXML(paragraphs []text.Paragraph) string {
	if len(paragraphs) == 0 {
		return `
<a:p><a:endParaRPr lang="en-US"/></a:p>`
	}
	var b strings.Builder
	defaultStyle := ContentStyleSpec{SizePt: defaultNotesFontSizePt}

	for _, p := range paragraphs {
		styleSpec := convertNotesStyle(p.Style)
		runSpecs := make([]TextRunSpec, len(p.Runs))
		for i, r := range p.Runs {
			runSpecs[i] = convertNotesRun(r)
		}
		b.WriteString(bulletParagraphRuns(runSpecs, styleSpec, defaultStyle))
	}
	return b.String()
}

func convertNotesStyle(s text.ParagraphStyle) BulletParagraphSpec {
	tabStops := make([]int64, 0, len(s.TabStops))
	for _, stop := range s.TabStops {
		tabStops = append(tabStops, stop.Emu())
	}
	return BulletParagraphSpec{
		Align:          s.Align,
		SpaceBeforePt:  s.SpaceBeforePt,
		SpaceAfterPt:   s.SpaceAfterPt,
		LineSpacingPct: s.LineSpacingPct,
		LineSpacingPts: s.LineSpacingPts,
		BulletStyle:    s.BulletStyle,
		BulletChar:     s.BulletChar,
		BulletColor:    s.BulletColor,
		BulletSize:     s.BulletSize,
		TabStops:       tabStops,
		Level:          s.Level,
		LeftIndent:     int64(s.LeftIndent),
		RightIndent:    int64(s.RightIndent),
		HangingIndent:  int64(s.HangingIndent),
	}
}

func convertNotesRun(r text.Run) TextRunSpec {
	return TextRunSpec{
		Text:          r.Text,
		Bold:          r.Bold,
		Italic:        r.Italic,
		Underline:     r.Underline,
		Strikethrough: r.Strikethrough,
		Subscript:     r.Subscript,
		Superscript:   r.Superscript,
		Color:         r.Color,
		Highlight:     r.Highlight,
		Font:          r.Font,
		SizePt:        r.SizePt,
		Code:          r.Code,
		AllCaps:       r.AllCaps,
		SmallCaps:     r.SmallCaps,
		Lang:          r.Lang,
	}
}
