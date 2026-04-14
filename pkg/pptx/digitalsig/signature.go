// Digital signature support for PPTX presentations
//
// Provides digital signature metadata and XML generation for the
// `_xmlsignatures/` package part per the OOXML digital signature spec.
//
// This is an exact port of the ppt-rs Rust implementation.

package digitalsig

import (
	"fmt"
	"strings"
	"time"
)

// HashAlgorithm used for signing.
type HashAlgorithm int

const (
	// HashAlgorithmSha256 is the default hash algorithm (SHA-256).
	HashAlgorithmSha256 HashAlgorithm = iota
	// HashAlgorithmSha384 uses SHA-384.
	HashAlgorithmSha384
	// HashAlgorithmSha512 uses SHA-512.
	HashAlgorithmSha512
	// HashAlgorithmSha1 uses SHA-1 (legacy).
	HashAlgorithmSha1
)

// URI returns the XML URI for the hash algorithm.
func (h HashAlgorithm) URI() string {
	switch h {
	case HashAlgorithmSha256:
		return "http://www.w3.org/2001/04/xmlenc#sha256"
	case HashAlgorithmSha384:
		return "http://www.w3.org/2001/04/xmldsig-more#sha384"
	case HashAlgorithmSha512:
		return "http://www.w3.org/2001/04/xmlenc#sha512"
	case HashAlgorithmSha1:
		return "http://www.w3.org/2000/09/xmldsig#sha1"
	default:
		return "http://www.w3.org/2001/04/xmlenc#sha256"
	}
}

// Name returns the display name for the hash algorithm.
func (h HashAlgorithm) Name() string {
	switch h {
	case HashAlgorithmSha256:
		return "SHA-256"
	case HashAlgorithmSha384:
		return "SHA-384"
	case HashAlgorithmSha512:
		return "SHA-512"
	case HashAlgorithmSha1:
		return "SHA-1"
	default:
		return "SHA-256"
	}
}

// SignerInfo contains identity information for the signer.
type SignerInfo struct {
	Name         string
	Email        *string
	Organization *string
	Title        *string
}

// NewSignerInfo creates a new SignerInfo with the given name.
func NewSignerInfo(name string) SignerInfo {
	return SignerInfo{
		Name: name,
	}
}

// WithEmail sets the email for the signer (builder pattern).
func (s SignerInfo) WithEmail(email string) SignerInfo {
	s.Email = &email
	return s
}

// WithOrganization sets the organization for the signer (builder pattern).
func (s SignerInfo) WithOrganization(org string) SignerInfo {
	s.Organization = &org
	return s
}

// WithTitle sets the title for the signer (builder pattern).
func (s SignerInfo) WithTitle(title string) SignerInfo {
	s.Title = &title
	return s
}

// SignatureCommitment represents the commitment type for the signature.
type SignatureCommitment int

const (
	// SignatureCommitmentCreated indicates the document was created by the signer.
	SignatureCommitmentCreated SignatureCommitment = iota
	// SignatureCommitmentApproved indicates the document was approved by the signer.
	SignatureCommitmentApproved
	// SignatureCommitmentReviewed indicates the document was reviewed by the signer.
	SignatureCommitmentReviewed
)

// URI returns the ETSI URI for the commitment type.
func (s SignatureCommitment) URI() string {
	switch s {
	case SignatureCommitmentCreated:
		return "http://uri.etsi.org/01903/v1.2.2#ProofOfCreation"
	case SignatureCommitmentApproved:
		return "http://uri.etsi.org/01903/v1.2.2#ProofOfApproval"
	case SignatureCommitmentReviewed:
		return "http://uri.etsi.org/01903/v1.2.2#ProofOfReview"
	default:
		return "http://uri.etsi.org/01903/v1.2.2#ProofOfCreation"
	}
}

// Label returns the display label for the commitment type.
func (s SignatureCommitment) Label() string {
	switch s {
	case SignatureCommitmentCreated:
		return "Created"
	case SignatureCommitmentApproved:
		return "Approved"
	case SignatureCommitmentReviewed:
		return "Reviewed"
	default:
		return "Created"
	}
}

// DigitalSignature contains configuration for a presentation signature.
type DigitalSignature struct {
	Signer         SignerInfo
	HashAlgorithm  HashAlgorithm
	SignDate       *string
	CommitmentType SignatureCommitment
	Comments       *string
}

// NewDigitalSignature creates a new DigitalSignature with the given signer.
func NewDigitalSignature(signer SignerInfo) DigitalSignature {
	return DigitalSignature{
		Signer:         signer,
		HashAlgorithm:  HashAlgorithmSha256,
		CommitmentType: SignatureCommitmentCreated,
	}
}

// WithHashAlgorithm sets the hash algorithm (builder pattern).
func (d DigitalSignature) WithHashAlgorithm(algo HashAlgorithm) DigitalSignature {
	d.HashAlgorithm = algo
	return d
}

// WithSignDate sets the sign date in ISO 8601 format (builder pattern).
func (d DigitalSignature) WithSignDate(date string) DigitalSignature {
	d.SignDate = &date
	return d
}

// WithCommitment sets the commitment type (builder pattern).
func (d DigitalSignature) WithCommitment(commitment SignatureCommitment) DigitalSignature {
	d.CommitmentType = commitment
	return d
}

// WithComments sets the signature comments (builder pattern).
func (d DigitalSignature) WithComments(comments string) DigitalSignature {
	d.Comments = &comments
	return d
}

// ToOriginXML generates the `_xmlsignatures/origin.sigs` relationship XML.
func (d DigitalSignature) ToOriginXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`
}

// ToSignatureXML generates signature info XML for `_xmlsignatures/sig1.xml`.
func (d DigitalSignature) ToSignatureXML() string {
	date := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	if d.SignDate != nil {
		date = *d.SignDate
	}

	var commentsXML string
	if d.Comments != nil {
		commentsXML = fmt.Sprintf("<SignatureComments>%s</SignatureComments>", xmlEscape(*d.Comments))
	}

	var xml strings.Builder
	xml.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	xml.WriteString(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">`)
	xml.WriteString(`<SignedInfo>`)
	xml.WriteString(`<CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>`)
	xml.WriteString(fmt.Sprintf(
		`<SignatureMethod Algorithm="%s"/>`,
		d.HashAlgorithm.URI(),
	))
	xml.WriteString(`</SignedInfo>`)
	xml.WriteString(`<SignatureValue/>`)
	xml.WriteString(`<KeyInfo>`)
	xml.WriteString(fmt.Sprintf(
		`<KeyName>%s</KeyName>`,
		xmlEscape(d.Signer.Name),
	))
	xml.WriteString(`</KeyInfo>`)
	xml.WriteString("<Object>")
	xml.WriteString(fmt.Sprintf(
		"<SignatureProperties><SignatureProperty Target=\"#SignatureInfo\"><SignatureInfoV1 xmlns=\"http://schemas.microsoft.com/office/2006/digsig\"><SetupID/><SignatureText>%s</SignatureText>%s<SignatureType>1</SignatureType><SignatureProviderUrl/><SignatureProviderDetails>9</SignatureProviderDetails><ManifestHashAlgorithm>%s</ManifestHashAlgorithm><SignatureProviderId>{{00000000-0000-0000-0000-000000000000}}</SignatureProviderId><CommitmentTypeId>%s</CommitmentTypeId><CommitmentTypeQualifier>%s</CommitmentTypeQualifier><SigningTime>%s</SigningTime></SignatureInfoV1></SignatureProperty></SignatureProperties>",
		xmlEscape(d.Signer.Name),
		commentsXML,
		d.HashAlgorithm.URI(),
		d.CommitmentType.URI(),
		d.CommitmentType.Label(),
		xmlEscape(date),
	))
	xml.WriteString("</Object>")
	xml.WriteString(`</Signature>`)

	return xml.String()
}

// ContentTypeEntry generates the content type entry for digital signatures.
func (d DigitalSignature) ContentTypeEntry() string {
	return `<Override PartName="/_xmlsignatures/sig1.xml" ContentType="application/vnd.openxmlformats-package.digital-signature-xmlsignature+xml"/>`
}

// xmlEscape escapes special XML characters.
func xmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
