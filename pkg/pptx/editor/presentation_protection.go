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
)

func rewritePresentationModifyVerifier(current string, password string) (string, error) {
	if strings.TrimSpace(current) == "" {
		return "", errors.New("missing presentation XML content")
	}

	source := removeSelfClosingTagByPrefix(current, "<p:modifyVerifier")
	password = strings.TrimSpace(password)
	if password == "" {
		return source, nil
	}

	salt := make([]byte, protectionSaltBytesEditor)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate protection salt: %w", err)
	}
	hash := protection.HashModifyPassword(password, salt, protectionSpinCountEditor)
	verifier := buildModifyVerifierXML(base64.StdEncoding.EncodeToString(salt), hash)

	if notesStart := strings.Index(source, "<p:notesSz"); notesStart >= 0 {
		if endRel := strings.Index(source[notesStart:], "/>"); endRel >= 0 {
			insertAt := notesStart + endRel + 2
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

func removeSelfClosingTagByPrefix(source, tagPrefix string) string {
	searchFrom := 0
	for {
		startRel := strings.Index(source[searchFrom:], tagPrefix)
		if startRel < 0 {
			return source
		}
		start := searchFrom + startRel
		endRel := strings.Index(source[start:], "/>")
		if endRel < 0 {
			return source
		}
		end := start + endRel + 2
		source = source[:start] + source[end:]
		searchFrom = start
	}
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
