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
	slideSize4x3 := common.SlideSize4x3()
	if size.Width == slideSize4x3.Width && size.Height == slideSize4x3.Height {
		return "screen4x3"
	}
	slideSize16x9 := common.SlideSize16x9()
	if size.Width == slideSize16x9.Width && size.Height == slideSize16x9.Height {
		return "screen16x9"
	}
	return "custom"
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

func renderCoreProperties(props common.CoreProperties) ([]byte, error) {
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

	type dctermsDate struct {
		XSIType string `xml:"xsi:type,attr"`
		Value   string `xml:",chardata"`
	}
	type corePropertiesXML struct {
		XMLName xml.Name `xml:"cp:coreProperties"`

		XMLNSCP       string `xml:"xmlns:cp,attr"`
		XMLNSDC       string `xml:"xmlns:dc,attr"`
		XMLNSDCTerms  string `xml:"xmlns:dcterms,attr"`
		XMLNSDCMITYpe string `xml:"xmlns:dcmitype,attr"`
		XMLNSXSI      string `xml:"xmlns:xsi,attr"`

		Title          string      `xml:"dc:title,omitempty"`
		Subject        string      `xml:"dc:subject,omitempty"`
		Creator        string      `xml:"dc:creator,omitempty"`
		Keywords       string      `xml:"cp:keywords,omitempty"`
		Description    string      `xml:"dc:description,omitempty"`
		LastModifiedBy string      `xml:"cp:lastModifiedBy,omitempty"`
		Revision       string      `xml:"cp:revision,omitempty"`
		Created        dctermsDate `xml:"dcterms:created"`
		Modified       dctermsDate `xml:"dcterms:modified"`
		Category       string      `xml:"cp:category,omitempty"`
		ContentStatus  string      `xml:"cp:contentStatus,omitempty"`
	}

	doc := corePropertiesXML{
		XMLNSCP:       common.CPNamespace,
		XMLNSDC:       common.DCNamespace,
		XMLNSDCTerms:  common.DCTermsNamespace,
		XMLNSDCMITYpe: common.DCMITypeNamespace,
		XMLNSXSI:      common.XSINamespace,

		Title:          strings.TrimSpace(props.Title),
		Subject:        strings.TrimSpace(props.Subject),
		Creator:        strings.TrimSpace(props.Creator),
		Keywords:       strings.TrimSpace(props.Keywords),
		Description:    strings.TrimSpace(props.Description),
		LastModifiedBy: lastModifiedBy,
		Revision:       strings.TrimSpace(props.Revision),
		Created: dctermsDate{
			XSIType: "dcterms:W3CDTF",
			Value:   created,
		},
		Modified: dctermsDate{
			XSIType: "dcterms:W3CDTF",
			Value:   modified,
		},
		Category:      strings.TrimSpace(props.Category),
		ContentStatus: strings.TrimSpace(props.ContentStatus),
	}

	data, err := xml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), data...), nil
}
