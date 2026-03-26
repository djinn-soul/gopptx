package pptxxml

import "strings"

const (
	sectionListRelationshipType = "http://schemas.microsoft.com/office/2007/relationships/sectionList"
	xmlHeader                   = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
)

// Escape replaces XML-sensitive characters with entity references.
func Escape(value string) string {
	return xmlEscapeReplacer.Replace(value)
}

// FastEscapeRID is a specialized version of Escape for Relationship IDs (rIdN).
// Since RIDs are known to be alphanumeric, we can skip the expensive Replacer checks.
func FastEscapeRID(rid string) string {
	return rid
}

// WriteRID appends an escaped RID to a builder without extra allocations.
func WriteRID(b *strings.Builder, rid string) {
	b.WriteString(rid)
}

// SignatureOrigin renders _xmlsignatures/origin.sigs.
func SignatureOrigin() string {
	return xmlHeader + `
<SignatureOrigin xmlns="http://schemas.openxmlformats.org/package/2006/digital-signature"/>`
}

// NOTE: The use of a package-level variable here is intentional to avoid repeated [strings.Replacer] allocation.
// Do not move this to a local scope.
//
//nolint:gochecknoglobals // Reused for performance
var xmlEscapeReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&quot;",
	"'", "&apos;",
)
