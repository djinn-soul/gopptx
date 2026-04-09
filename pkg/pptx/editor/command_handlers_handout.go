package editor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------
// Feature 4 – Handout master
// ---------------------------------------------------------------------------

const handoutMasterPath = "ppt/handoutMasters/handoutMaster1.xml"

// GetHandoutMaster returns basic info about the handout master.
func (e *PresentationEditor) GetHandoutMaster() (map[string]any, error) {
	data, ok := e.parts.Get(handoutMasterPath)
	if !ok {
		return map[string]any{"present": false}, nil
	}
	xmlStr := string(data)
	orientation := "landscape"
	if strings.Contains(xmlStr, `orient="portrait"`) {
		orientation = "portrait"
	}
	spp := extractHandoutSlidesPerPage(xmlStr)
	return map[string]any{
		"present":         true,
		"orientation":     orientation,
		"slides_per_page": spp,
	}, nil
}

// extractHandoutSlidesPerPage reads the slides-per-page from handout master XML.
// Falls back to 6 if not determinable.
func extractHandoutSlidesPerPage(xmlStr string) int {
	const (
		defaultSlidesPerPage  = 6
		maxValidSlidesPerPage = 9
	)
	count := strings.Count(xmlStr, "<p:sp>") + strings.Count(xmlStr, "<p:sp ")
	if count > 0 && count <= maxValidSlidesPerPage {
		return count
	}
	return defaultSlidesPerPage
}

// UpdateHandoutMaster sets orientation and slides-per-page in the handout master.
func (e *PresentationEditor) UpdateHandoutMaster(props map[string]any) error {
	data, ok := e.parts.Get(handoutMasterPath)
	var xmlStr string
	if ok {
		xmlStr = string(data)
	} else {
		xmlStr = defaultHandoutMasterXML()
	}

	if orient, castOK := props["orientation"].(string); castOK {
		xmlStr = applyHandoutOrientation(xmlStr, orient)
	}

	e.parts.Set(handoutMasterPath, []byte(xmlStr))
	e.addContentTypeOverride(handoutMasterPath,
		"application/vnd.openxmlformats-officedocument.presentationml.handoutMaster+xml")
	return nil
}

// applyHandoutOrientation sets the orient attribute in handout master XML.
func applyHandoutOrientation(xmlStr, orientation string) string {
	re := regexp.MustCompile(`orient="[^"]*"`)
	if re.MatchString(xmlStr) {
		return re.ReplaceAllString(xmlStr, fmt.Sprintf(`orient="%s"`, orientation))
	}
	reSz := regexp.MustCompile(`(<p:sldSz\b[^/]*)(/>)`)
	return reSz.ReplaceAllString(xmlStr, fmt.Sprintf(`$1 orient="%s"$2`, orientation))
}

// defaultHandoutMasterXML returns a minimal handout master XML.
func defaultHandoutMasterXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:handoutMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` +
		` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/>` +
		`<p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>` +
		`<p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1"` +
		` accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5"` +
		` accent6="accent6" hlink="hlink" folHlink="folHlink"/>` +
		`</p:handoutMaster>`
}

// handleGetHandoutMaster retrieves handout master info.
//
// Payload: {} (empty).
// Response: {"present": bool, "orientation": "landscape"|"portrait", "slides_per_page": N}.
func handleGetHandoutMaster(e *PresentationEditor, _ json.RawMessage) (any, error) {
	result, err := e.GetHandoutMaster()
	if err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return result, nil
}

// handleUpdateHandoutMaster updates handout master properties.
//
// Payload: {"orientation": "landscape"|"portrait", "slides_per_page": N}.
// Response: {"updated": true}.
func handleUpdateHandoutMaster(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	if updateErr := e.UpdateHandoutMaster(p); updateErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, updateErr.Error())
	}
	return respUpdated, nil
}

// ---------------------------------------------------------------------------
// Feature 5 – Digital signature detection
// ---------------------------------------------------------------------------

// HasDigitalSignature returns true if the package contains a digital signature.
func (e *PresentationEditor) HasDigitalSignature() (bool, error) {
	has := e.parts.Has("_xmlsignatures/origin.sigs") ||
		e.parts.Has("_xmlsignatures/sig1.xml")
	return has, nil
}

// handleHasDigitalSignature checks for a digital signature part.
//
// Payload: {} (empty).
// Response: {"has_digital_signature": bool}.
func handleHasDigitalSignature(e *PresentationEditor, _ json.RawMessage) (any, error) {
	has, err := e.HasDigitalSignature()
	if err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return map[string]bool{"has_digital_signature": has}, nil
}
