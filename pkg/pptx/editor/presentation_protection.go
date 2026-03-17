package editor

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
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

var (
	modifyVerifierPattern = regexp.MustCompile(`(?s)<p:modifyVerifier\b[^>]*/>`)
	notesSizePattern      = regexp.MustCompile(`<p:notesSz\b[^>]*/>`)
)

func rewritePresentationModifyVerifier(current []byte, password string) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}

	source := modifyVerifierPattern.ReplaceAllString(string(current), "")
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

	if loc := notesSizePattern.FindStringIndex(source); loc != nil {
		insertAt := loc[1]
		return source[:insertAt] + "\n" + verifier + source[insertAt:], nil
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
