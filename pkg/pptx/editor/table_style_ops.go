package editor

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	tableStylesPartPath    = "ppt/tableStyles.xml"
	tableStylesContentType = "application/vnd.openxmlformats-officedocument.presentationml.tableStyles+xml"
	tableStyleEntryGroups  = 3
	guidByteLength         = 16
)

var (
	reTableStyleEntry = regexp.MustCompile(`(?s)<a:tblStyle\b[^>]*styleId="([^"]+)"[^>]*styleName="([^"]*)"[^>]*/?>`)
	reTableStyleByID  = regexp.MustCompile(`(?s)<a:tblStyle\b[^>]*styleId="%s"[^>]*/?>`)
	reDefaultStyleID  = regexp.MustCompile(`\bdef="[^"]*"`)
	reStyleGUID       = regexp.MustCompile(`^\{[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\}$`)
)

func (e *PresentationEditor) DefineTableStyle(def common.TableStyleDefinition) (string, error) {
	name := strings.TrimSpace(def.Name)
	if name == "" {
		return "", errors.New("name is required")
	}

	styleID, err := normalizeOrGenerateStyleID(def.StyleID)
	if err != nil {
		return "", err
	}

	if err := e.ensureTableStylesInfrastructure(styleID); err != nil {
		return "", err
	}
	part, _ := e.parts.Get(tableStylesPartPath)
	xml := string(part)
	entry := fmt.Sprintf(`<a:tblStyle styleId="%s" styleName="%s"/>`, common.XMLEscape(styleID), common.XMLEscape(name))
	pattern := regexp.MustCompile(fmt.Sprintf(reTableStyleByID.String(), regexp.QuoteMeta(styleID)))
	if pattern.MatchString(xml) {
		xml = pattern.ReplaceAllString(xml, entry)
	} else {
		xml = strings.Replace(xml, "</a:tblStyleLst>", entry+"</a:tblStyleLst>", 1)
	}
	e.parts.Set(tableStylesPartPath, []byte(xml))
	return styleID, nil
}

func (e *PresentationEditor) ListTableStyles() ([]common.TableStyleInfo, error) {
	part, ok := e.parts.Get(tableStylesPartPath)
	if !ok {
		return []common.TableStyleInfo{}, nil
	}
	matches := reTableStyleEntry.FindAllStringSubmatch(string(part), -1)
	styles := make([]common.TableStyleInfo, 0, len(matches))
	for _, match := range matches {
		if len(match) != tableStyleEntryGroups {
			continue
		}
		styles = append(styles, common.TableStyleInfo{StyleID: match[1], Name: match[2]})
	}
	return styles, nil
}

func (e *PresentationEditor) ensureTableStylesInfrastructure(defaultStyleID string) error {
	part, ok := e.parts.Get(tableStylesPartPath)
	if !ok {
		xml := fmt.Sprintf(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><a:tblStyleLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" def="%s"></a:tblStyleLst>`,
			common.XMLEscape(defaultStyleID),
		)
		e.parts.Set(tableStylesPartPath, []byte(xml))
		e.addContentTypeOverride(tableStylesPartPath, tableStylesContentType)
		return e.ensurePresentationTableStylesRelationship()
	}

	xml := string(part)
	if !reDefaultStyleID.MatchString(xml) {
		if strings.Contains(xml, "<a:tblStyleLst ") {
			xml = strings.Replace(
				xml,
				"<a:tblStyleLst ",
				`<a:tblStyleLst def="`+common.XMLEscape(defaultStyleID)+`" `,
				1,
			)
		} else {
			xml = strings.Replace(
				xml,
				"<a:tblStyleLst>",
				`<a:tblStyleLst def="`+common.XMLEscape(defaultStyleID)+`">`,
				1,
			)
		}
		e.parts.Set(tableStylesPartPath, []byte(xml))
	}

	e.addContentTypeOverride(tableStylesPartPath, tableStylesContentType)
	return e.ensurePresentationTableStylesRelationship()
}

func (e *PresentationEditor) ensurePresentationTableStylesRelationship() error {
	relsData, ok := e.parts.Get(common.PresentationRelPath)
	if !ok {
		return fmt.Errorf("presentation relationships part not found: %s", common.PresentationRelPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse %s: %w", common.PresentationRelPath, err)
	}
	for _, rel := range rels {
		if rel.Type == common.RelTypeTableStyles {
			return nil
		}
	}

	nextNum := common.NextRelationshipNumber(rels)
	rels = append(rels, common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", nextNum),
		Type:   common.RelTypeTableStyles,
		Target: "tableStyles.xml",
	})
	e.parts.Set(common.PresentationRelPath, []byte(renderRelationshipsXML(rels)))
	return nil
}

func normalizeOrGenerateStyleID(styleID string) (string, error) {
	raw := strings.ToUpper(strings.TrimSpace(styleID))
	if raw != "" {
		if !strings.HasPrefix(raw, "{") {
			raw = "{" + raw
		}
		if !strings.HasSuffix(raw, "}") {
			raw += "}"
		}
		if !reStyleGUID.MatchString(raw) {
			return "", errors.New("style_id must be a GUID like {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}")
		}
		return raw, nil
	}

	buf := make([]byte, guidByteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate style id: %w", err)
	}
	hexValue := strings.ToUpper(hex.EncodeToString(buf))
	return fmt.Sprintf(
		"{%s-%s-%s-%s-%s}",
		hexValue[0:8],
		hexValue[8:12],
		hexValue[12:16],
		hexValue[16:20],
		hexValue[20:32],
	), nil
}
