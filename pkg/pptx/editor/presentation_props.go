package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const themePartPath = "ppt/theme/theme1.xml"
const corePropertiesBuilderCap = 700
const slideSizeTypeCustom = "custom"

var sldSzPattern = regexp.MustCompile(`(?s)<p:sldSz\b[^>]*/>`)

// ApplyTheme replaces the package theme with the provided palette and fonts.
func (e *PresentationEditor) ApplyTheme(theme styling.Theme) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if !e.parts.Has(themePartPath) {
		return fmt.Errorf("missing required package part %q", themePartPath)
	}
	e.parts.Set(themePartPath, []byte(pptxxml.Theme(mapEditorThemeToSpec(&theme))))
	return nil
}

// SetSlideSize updates presentation dimensions used for existing and future slides.
func (e *PresentationEditor) SetSlideSize(size common.SlideSize) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if size.Width <= 0 || size.Height <= 0 {
		return errors.New("slide size must have positive dimensions")
	}

	rewritten, err := rewritePresentationSlideSize(e.presentationXML, size)
	if err != nil {
		return err
	}

	e.presentationXML = rewritten
	e.metadata.SlideSize = size
	return nil
}

// GetCoreProperties returns the presentation's core properties (Dublin Core metadata).
func (e *PresentationEditor) GetCoreProperties() common.CoreProperties {
	if e == nil {
		return common.CoreProperties{}
	}
	return e.metadata.CoreProperties
}

// SetCoreProperties updates the presentation's core properties.
// These changes are applied to the in-memory metadata and will be written to docProps/core.xml on Save.
func (e *PresentationEditor) SetCoreProperties(p common.CoreProperties) {
	if e == nil {
		return
	}
	e.metadata.CoreProperties = p
	// Sync legacy title field for backward compatibility if needed, though we should likely deprecate it.
	e.metadata.Title = p.Title
}

// SetShowSettings injects or replaces the <p:showPr> element in the presentation XML.
// Use this to configure loop, kiosk mode, browse mode, and animation/timing settings.
func (e *PresentationEditor) SetShowSettings(s common.ShowSettings) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	xmlSpec := pptxxml.ShowSettings{
		Loop:           s.Loop,
		Mode:           pptxxml.ShowMode(s.Mode),
		ShowScrollbar:  s.ShowScrollbar,
		DisableTimings: s.DisableTimings,
		HideAnimation:  s.HideAnimation,
	}
	newElement := pptxxml.ShowPrXML(xmlSpec)
	current := e.presentationXML

	// Remove any existing <p:showPr>...</p:showPr> or <p:showPr ... />
	showPrPattern := regexp.MustCompile(`(?s)<p:showPr\b[^>]*(?:/>|>.*?</p:showPr>)`)
	current = showPrPattern.ReplaceAllString(current, "")

	if newElement == "" {
		e.presentationXML = current
		return nil
	}

	// Inject before </p:presentation>
	if idx := strings.LastIndex(current, "</p:presentation>"); idx >= 0 {
		e.presentationXML = current[:idx] + "\n" + newElement + "\n" + current[idx:]
		return nil
	}
	return errors.New("presentation XML does not contain </p:presentation>")
}

func parsePresentationSlideSize(content []byte) (common.SlideSize, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return common.SlideSize{}, nil
			}
			return common.SlideSize{}, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sldSz" {
			continue
		}
		return parseSlideSizeAttrs(start.Attr)
	}
}

func parseSlideSizeAttrs(attrs []xml.Attr) (common.SlideSize, error) {
	size := common.SlideSize{}
	for _, attr := range attrs {
		switch attr.Name.Local {
		case "cx":
			width, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
			if err != nil {
				return common.SlideSize{}, fmt.Errorf("invalid slide width %q", attr.Value)
			}
			size.Width = width
		case "cy":
			height, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
			if err != nil {
				return common.SlideSize{}, fmt.Errorf("invalid slide height %q", attr.Value)
			}
			size.Height = height
		}
	}
	if size.Width == 0 || size.Height == 0 {
		return common.SlideSize{}, errors.New("slide size entry missing cx/cy")
	}
	return size, nil
}

func rewritePresentationSlideSize(current string, size common.SlideSize) (string, error) {
	if strings.TrimSpace(current) == "" {
		return "", errors.New("missing presentation XML content")
	}
	entry := fmt.Sprintf(`<p:sldSz cx="%d" cy="%d" type="%s"/>`, size.Width, size.Height, slideSizeType(size))

	if sldSzPattern.MatchString(current) {
		replaced := false
		return sldSzPattern.ReplaceAllStringFunc(current, func(match string) string {
			if replaced {
				return match
			}
			replaced = true
			return entry
		}), nil
	}

	if idx := strings.Index(current, "<p:notesSz"); idx >= 0 {
		return current[:idx] + entry + "\n" + current[idx:], nil
	}
	if idx := strings.Index(current, "</p:presentation>"); idx >= 0 {
		return current[:idx] + entry + "\n" + current[idx:], nil
	}
	return "", errors.New("presentation XML does not contain <p:notesSz> or </p:presentation>")
}

func slideSizeType(size common.SlideSize) string {
	if size.Width == common.SlideSize4x3().Width && size.Height == common.SlideSize4x3().Height {
		return "screen4x3"
	} else if size.Width == common.SlideSize16x9().Width && size.Height == common.SlideSize16x9().Height {
		return "screen16x9"
	}
	return slideSizeTypeCustom
}

func mapEditorThemeToSpec(theme *styling.Theme) *pptxxml.ThemeSpec {
	if theme == nil {
		return nil
	}
	return &pptxxml.ThemeSpec{
		Name: theme.Name,
		Colors: pptxxml.ColorSchemeSpec{
			Name:     theme.Colors.Name,
			Dk1:      theme.Colors.Dk1,
			Lt1:      theme.Colors.Lt1,
			Dk2:      theme.Colors.Dk2,
			Lt2:      theme.Colors.Lt2,
			Accent1:  theme.Colors.Accent1,
			Accent2:  theme.Colors.Accent2,
			Accent3:  theme.Colors.Accent3,
			Accent4:  theme.Colors.Accent4,
			Accent5:  theme.Colors.Accent5,
			Accent6:  theme.Colors.Accent6,
			Hlink:    theme.Colors.Hlink,
			FolHlink: theme.Colors.FolHlink,
		},
		Fonts: pptxxml.FontSchemeSpec{
			Name:      theme.Fonts.Name,
			MajorFont: theme.Fonts.MajorFont,
			MinorFont: theme.Fonts.MinorFont,
		},
	}
}

// parseCoreProperties parses docProps/core.xml into CoreProperties.
func parseCoreProperties(content []byte) (common.CoreProperties, error) {
	if len(content) == 0 {
		return common.CoreProperties{}, nil
	}
	var props common.CoreProperties
	if err := xml.Unmarshal(content, &props); err != nil {
		return common.CoreProperties{}, err
	}
	return props, nil
}

func renderCoreProperties(props common.CoreProperties) []byte {
	created := strings.TrimSpace(props.Created)
	if created == "" {
		created = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	}
	modified := strings.TrimSpace(props.Modified)
	if modified == "" {
		modified = created
	}

	lastModifiedBy := strings.TrimSpace(props.LastModifiedBy)
	if lastModifiedBy == "" {
		lastModifiedBy = strings.TrimSpace(props.Creator)
	}

	// Hand-rolled renderer replaces xml.Marshal to eliminate ~73 allocs/save
	// from xml.(*printer).marshalAttr (5 namespace attrs + reflection overhead).
	var b strings.Builder
	// Header ~65 bytes + root open ~350 bytes + up to 10 optional fields ~150 bytes + 2 dates ~100 bytes + close ~20 bytes.
	b.Grow(corePropertiesBuilderCap)
	b.WriteString(xml.Header)
	b.WriteString(`<cp:coreProperties`)
	b.WriteString(` xmlns:cp="` + common.CPNamespace + `"`)
	b.WriteString(` xmlns:dc="` + common.DCNamespace + `"`)
	b.WriteString(` xmlns:dcterms="` + common.DCTermsNamespace + `"`)
	b.WriteString(` xmlns:dcmitype="` + common.DCMITypeNamespace + `"`)
	b.WriteString(` xmlns:xsi="` + common.XSINamespace + `"`)
	b.WriteByte('>')

	writeCorePropsField := func(tag, value string) {
		if value == "" {
			return
		}
		b.WriteByte('<')
		b.WriteString(tag)
		b.WriteByte('>')
		b.WriteString(common.XMLEscape(value))
		b.WriteString("</")
		b.WriteString(tag)
		b.WriteByte('>')
	}

	writeCorePropsField("dc:title", strings.TrimSpace(props.Title))
	writeCorePropsField("dc:subject", strings.TrimSpace(props.Subject))
	writeCorePropsField("dc:creator", strings.TrimSpace(props.Creator))
	writeCorePropsField("cp:keywords", strings.TrimSpace(props.Keywords))
	writeCorePropsField("dc:description", strings.TrimSpace(props.Description))
	writeCorePropsField("cp:lastModifiedBy", lastModifiedBy)
	writeCorePropsField("cp:revision", strings.TrimSpace(props.Revision))

	b.WriteString(`<dcterms:created xsi:type="dcterms:W3CDTF">`)
	b.WriteString(common.XMLEscape(created))
	b.WriteString(`</dcterms:created>`)
	b.WriteString(`<dcterms:modified xsi:type="dcterms:W3CDTF">`)
	b.WriteString(common.XMLEscape(modified))
	b.WriteString(`</dcterms:modified>`)

	writeCorePropsField("cp:category", strings.TrimSpace(props.Category))
	writeCorePropsField("cp:contentStatus", strings.TrimSpace(props.ContentStatus))

	b.WriteString(`</cp:coreProperties>`)

	return []byte(b.String())
}
