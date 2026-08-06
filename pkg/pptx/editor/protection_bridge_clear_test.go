package editor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// The set_modify_password bridge op with an empty password is the documented way
// to remove write protection, so it must reach the handler rather than being
// rejected as a missing/blank field, and it must drop the verifier the package
// was opened with.
func TestExecuteCommand_SetModifyPasswordEmpty_ClearsProtection(t *testing.T) {
	base := writeDeckFixture(t, "protection-bridge.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	protectedPath := filepath.Join(t.TempDir(), "protected.pptx")
	first, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	first.SetModifyPassword("Secret123!")
	if err := first.Save(protectedPath); err != nil {
		t.Fatalf("save protected deck: %v", err)
	}
	_ = first.Close()

	ed, err := OpenPresentationEditor(protectedPath)
	if err != nil {
		t.Fatalf("reopen protected deck: %v", err)
	}
	defer func() { _ = ed.Close() }()

	resp := ExecuteCommand(ed, `{"api_version":1,"op":"set_modify_password","payload":{"password":""}}`)
	if strings.Contains(resp, `"error"`) {
		t.Fatalf("empty password rejected by the bridge: %s", resp)
	}

	if got := savePresentationXML(t, ed); strings.Contains(got, modifyVerifierTagPrefix) {
		t.Fatal("bridge set_modify_password with an empty password left protection in place")
	}
}

// A bridge call that omits the field entirely is a client error, not a request
// to unprotect the deck.
func TestExecuteCommand_SetModifyPasswordMissingField_KeepsProtection(t *testing.T) {
	base := writeDeckFixture(t, "protection-bridge-missing.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	protectedPath := filepath.Join(t.TempDir(), "protected.pptx")
	first, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	first.SetModifyPassword("Secret123!")
	if err := first.Save(protectedPath); err != nil {
		t.Fatalf("save protected deck: %v", err)
	}
	_ = first.Close()

	ed, err := OpenPresentationEditor(protectedPath)
	if err != nil {
		t.Fatalf("reopen protected deck: %v", err)
	}
	defer func() { _ = ed.Close() }()

	resp := ExecuteCommand(ed, `{"api_version":1,"op":"set_modify_password","payload":{}}`)
	if !strings.Contains(resp, "MISSING_FIELD") {
		t.Fatalf("missing password field should be a MISSING_FIELD error, got: %s", resp)
	}

	if got := savePresentationXML(t, ed); !strings.Contains(got, modifyVerifierTagPrefix) {
		t.Fatal("a rejected command dropped the existing protection")
	}
}
