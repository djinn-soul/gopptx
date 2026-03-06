package digitalsig

import (
	"strings"
	"testing"
)

func TestHashAlgorithmDefault(t *testing.T) {
	algo := HashAlgorithmSha256
	if algo != HashAlgorithmSha256 {
		t.Errorf("Expected default to be Sha256, got %v", algo)
	}
	if !strings.Contains(algo.URI(), "sha256") {
		t.Errorf("Expected URI to contain sha256, got %s", algo.URI())
	}
	if algo.Name() != "SHA-256" {
		t.Errorf("Expected name to be SHA-256, got %s", algo.Name())
	}
}

func TestHashAlgorithmVariants(t *testing.T) {
	if !strings.Contains(HashAlgorithmSha384.URI(), "sha384") {
		t.Errorf("Expected Sha384 URI to contain sha384")
	}
	if !strings.Contains(HashAlgorithmSha512.URI(), "sha512") {
		t.Errorf("Expected Sha512 URI to contain sha512")
	}
	if !strings.Contains(HashAlgorithmSha1.URI(), "sha1") {
		t.Errorf("Expected Sha1 URI to contain sha1")
	}

	if HashAlgorithmSha384.Name() != "SHA-384" {
		t.Errorf("Expected Sha384 name to be SHA-384, got %s", HashAlgorithmSha384.Name())
	}
	if HashAlgorithmSha512.Name() != "SHA-512" {
		t.Errorf("Expected Sha512 name to be SHA-512, got %s", HashAlgorithmSha512.Name())
	}
	if HashAlgorithmSha1.Name() != "SHA-1" {
		t.Errorf("Expected Sha1 name to be SHA-1, got %s", HashAlgorithmSha1.Name())
	}
}

func TestSignerInfoNew(t *testing.T) {
	signer := NewSignerInfo("Alice")
	if signer.Name != "Alice" {
		t.Errorf("Expected name to be Alice, got %s", signer.Name)
	}
	if signer.Email != nil {
		t.Errorf("Expected email to be nil")
	}
	if signer.Organization != nil {
		t.Errorf("Expected organization to be nil")
	}
}

func TestSignerInfoBuilder(t *testing.T) {
	signer := NewSignerInfo("Bob").
		WithEmail("bob@example.com").
		WithOrganization("Acme Corp").
		WithTitle("Engineer")

	if signer.Name != "Bob" {
		t.Errorf("Expected name to be Bob, got %s", signer.Name)
	}
	if signer.Email == nil || *signer.Email != "bob@example.com" {
		t.Errorf("Expected email to be bob@example.com, got %v", signer.Email)
	}
	if signer.Organization == nil || *signer.Organization != "Acme Corp" {
		t.Errorf("Expected organization to be Acme Corp, got %v", signer.Organization)
	}
	if signer.Title == nil || *signer.Title != "Engineer" {
		t.Errorf("Expected title to be Engineer, got %v", signer.Title)
	}
}

func TestSignatureCommitmentVariants(t *testing.T) {
	if !strings.Contains(SignatureCommitmentCreated.URI(), "Creation") {
		t.Errorf("Expected Created URI to contain Creation")
	}
	if !strings.Contains(SignatureCommitmentApproved.URI(), "Approval") {
		t.Errorf("Expected Approved URI to contain Approval")
	}
	if !strings.Contains(SignatureCommitmentReviewed.URI(), "Review") {
		t.Errorf("Expected Reviewed URI to contain Review")
	}

	if SignatureCommitmentCreated.Label() != "Created" {
		t.Errorf("Expected Created label to be Created, got %s", SignatureCommitmentCreated.Label())
	}
	if SignatureCommitmentApproved.Label() != "Approved" {
		t.Errorf("Expected Approved label to be Approved, got %s", SignatureCommitmentApproved.Label())
	}
	if SignatureCommitmentReviewed.Label() != "Reviewed" {
		t.Errorf("Expected Reviewed label to be Reviewed, got %s", SignatureCommitmentReviewed.Label())
	}
}

func TestDigitalSignatureNew(t *testing.T) {
	sig := NewDigitalSignature(NewSignerInfo("Alice"))
	if sig.Signer.Name != "Alice" {
		t.Errorf("Expected signer name to be Alice, got %s", sig.Signer.Name)
	}
	if sig.HashAlgorithm != HashAlgorithmSha256 {
		t.Errorf("Expected default hash algorithm to be Sha256")
	}
	if sig.CommitmentType != SignatureCommitmentCreated {
		t.Errorf("Expected default commitment type to be Created")
	}
}

func TestDigitalSignatureBuilder(t *testing.T) {
	sig := NewDigitalSignature(NewSignerInfo("Bob")).
		WithHashAlgorithm(HashAlgorithmSha512).
		WithSignDate("2025-06-15T10:00:00Z").
		WithCommitment(SignatureCommitmentApproved).
		WithComments("Looks good")

	if sig.HashAlgorithm != HashAlgorithmSha512 {
		t.Errorf("Expected hash algorithm to be Sha512")
	}
	if sig.SignDate == nil || *sig.SignDate != "2025-06-15T10:00:00Z" {
		t.Errorf("Expected sign date to be 2025-06-15T10:00:00Z, got %v", sig.SignDate)
	}
	if sig.CommitmentType != SignatureCommitmentApproved {
		t.Errorf("Expected commitment type to be Approved")
	}
	if sig.Comments == nil || *sig.Comments != "Looks good" {
		t.Errorf("Expected comments to be 'Looks good', got %v", sig.Comments)
	}
}

func TestSignatureXML(t *testing.T) {
	sig := NewDigitalSignature(NewSignerInfo("Alice")).
		WithSignDate("2025-01-01T00:00:00Z")
	xml := sig.ToSignatureXML()

	if !strings.Contains(xml, "<Signature") {
		t.Errorf("Expected XML to contain <Signature")
	}
	if !strings.Contains(xml, "Alice") {
		t.Errorf("Expected XML to contain Alice")
	}
	if !strings.Contains(xml, "sha256") {
		t.Errorf("Expected XML to contain sha256")
	}
	if !strings.Contains(xml, "SigningTime") {
		t.Errorf("Expected XML to contain SigningTime")
	}
}

func TestSignatureXMLWithComments(t *testing.T) {
	sig := NewDigitalSignature(NewSignerInfo("Bob")).
		WithComments("Reviewed & approved")
	xml := sig.ToSignatureXML()

	// Check that the ampersand is properly escaped in the XML (look for &)
	if !strings.Contains(xml, "amp;") {
		t.Errorf("Expected XML to contain escaped ampersand, got: %s", xml)
	}
}

func TestOriginXML(t *testing.T) {
	sig := NewDigitalSignature(NewSignerInfo("X"))
	xml := sig.ToOriginXML()

	if !strings.Contains(xml, "Relationships") {
		t.Errorf("Expected XML to contain Relationships")
	}
}

func TestContentTypeEntry(t *testing.T) {
	sig := NewDigitalSignature(NewSignerInfo("Test"))
	ct := sig.ContentTypeEntry()

	if !strings.Contains(ct, "digital-signature") {
		t.Errorf("Expected content type to contain digital-signature")
	}
	if !strings.Contains(ct, "sig1.xml") {
		t.Errorf("Expected content type to contain sig1.xml")
	}
}

func TestXMLEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"&", "&amp;"},
		{"<", "&lt;"},
		{">", "&gt;"},
		{"'", "&apos;"},
	}

	for _, tt := range tests {
		result := xmlEscape(tt.input)
		if result != tt.expected {
			t.Errorf("xmlEscape(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}

	// Test quote separately to avoid escaping issues
	quoteResult := xmlEscape("\"")
	if quoteResult != "&quot;" {
		t.Errorf("xmlEscape(quote) = %q, expected &quot;", quoteResult)
	}

	// Test combined string
	combinedResult := xmlEscape("a & b < c > d \"e\" f 'g'")
	expectedCombined := "a &amp; b &lt; c &gt; d &quot;e&quot; f &apos;g&apos;"
	if combinedResult != expectedCombined {
		t.Errorf("xmlEscape(combined) = %q, expected %q", combinedResult, expectedCombined)
	}
}
