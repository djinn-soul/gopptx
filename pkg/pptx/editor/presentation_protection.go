package editor

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
)

const (
	protectionSaltBytesEditor = 16
	protectionSpinCountEditor = 100000
	protectionHashAlgSID      = 14
	selfClosingTagSuffixLen   = 2
	modifyVerifierTagPrefix   = "<p:modifyVerifier"
)

// SetModifyPassword sets the presentation's "password to modify". An empty
// password removes the protection, discarding any verifier the opened package
// carried; assigning to Metadata().Protection.ModifyPassword directly cannot
// express that removal, because an empty value there means "unchanged".
func (e *PresentationEditor) SetModifyPassword(password string) {
	e.metadata.Protection.ModifyPassword = password
	e.existingModifyVerifier = ""
}

// rewritePresentationModifyVerifier emits the p:modifyVerifier element for the
// saved presentation.xml.
//
// preserved is the verifier the package already carried when it was opened. It
// is re-emitted whenever no new password is supplied, because the element holds
// only a salted hash: keeping a deck's "password to modify" across an
// open/edit/save cycle never requires knowing the plaintext. Without it, every
// save of an existing protected deck would silently drop the protection.
// Clearing protection is therefore an explicit act — the caller sets an empty
// password, which also drops preserved.
func rewritePresentationModifyVerifier(current, password, preserved string) (string, error) {
	if strings.TrimSpace(current) == "" {
		return "", errors.New("missing presentation XML content")
	}

	source := removeSelfClosingTagByPrefix(current, modifyVerifierTagPrefix)
	password = strings.TrimSpace(password)

	verifier := strings.TrimSpace(preserved)
	if password != "" {
		salt := make([]byte, protectionSaltBytesEditor)
		defer clear(salt)
		if _, err := rand.Read(salt); err != nil {
			return "", fmt.Errorf("generate protection salt: %w", err)
		}
		hash := protection.HashModifyPassword(password, salt, protectionSpinCountEditor)
		verifier = buildModifyVerifierXML(base64.StdEncoding.EncodeToString(salt), hash)
	}
	if verifier == "" {
		return source, nil
	}

	if notesStart := strings.Index(source, "<p:notesSz"); notesStart >= 0 {
		if endRel := strings.Index(source[notesStart:], "/>"); endRel >= 0 {
			insertAt := notesStart + endRel + selfClosingTagSuffixLen
			return source[:insertAt] + "\n" + verifier + source[insertAt:], nil
		}
	}
	if extIdx := strings.Index(source, "<p:extLst>"); extIdx >= 0 {
		return source[:extIdx] + verifier + "\n" + source[extIdx:], nil
	}
	endIdx := strings.LastIndex(source, "</p:presentation>")
	if endIdx < 0 {
		return "", errors.New("presentation XML does not contain </p:presentation>")
	}
	return source[:endIdx] + verifier + "\n" + source[endIdx:], nil
}

// extractModifyVerifierTag returns the p:modifyVerifier element verbatim, or ""
// when the presentation carries no write protection.
func extractModifyVerifierTag(presentationXML string) string {
	start := strings.Index(presentationXML, modifyVerifierTagPrefix)
	if start < 0 {
		return ""
	}
	endRel := strings.Index(presentationXML[start:], "/>")
	if endRel < 0 {
		return ""
	}
	return presentationXML[start : start+endRel+selfClosingTagSuffixLen]
}

func removeSelfClosingTagByPrefix(source, tagPrefix string) string {
	var b strings.Builder
	b.Grow(len(source))

	lastWrite := 0
	searchFrom := 0
	found := false

	for {
		startRel := strings.Index(source[searchFrom:], tagPrefix)
		if startRel < 0 {
			break
		}

		start := searchFrom + startRel
		endRel := strings.Index(source[start:], "/>")
		if endRel < 0 {
			break
		}

		end := start + endRel + selfClosingTagSuffixLen
		b.WriteString(source[lastWrite:start])

		found = true
		lastWrite = end
		searchFrom = end
	}

	if !found {
		return source
	}

	b.WriteString(source[lastWrite:])
	return b.String()
}

func buildModifyVerifierXML(saltData, hashData string) string {
	var b strings.Builder
	b.WriteString(
		`<p:modifyVerifier cryptProviderType="rsaAES" cryptAlgorithmClass="hash" cryptAlgorithmType="typeAny" cryptAlgorithmSid="`,
	)
	b.WriteString(strconv.Itoa(protectionHashAlgSID))
	b.WriteString(`" spinCount="`)
	b.WriteString(strconv.Itoa(protectionSpinCountEditor))
	b.WriteString(`" saltData="`)
	b.WriteString(common.XMLEscape(saltData))
	b.WriteString(`" hashData="`)
	b.WriteString(common.XMLEscape(hashData))
	b.WriteString(`"/>`)
	return b.String()
}
