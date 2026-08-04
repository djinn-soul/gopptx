package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func contentTypesWithMain(mainType string) map[string][]byte {
	return map[string][]byte{
		common.ContentTypesPath: []byte(
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
				`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
				`<Override PartName="/ppt/presentation.xml" ContentType="` + mainType + `"/>` +
				`</Types>`,
		),
	}
}

// A deck opened from a .potx keeps the template main content type. Saving it as
// a .pptx without rewriting that leaves a file PowerPoint still opens as a
// template (upstream issue #1070).
func TestApplyOutputPresentationKindConvertsTemplateToPresentation(t *testing.T) {
	parts := contentTypesWithMain(templateMainContentType)

	if err := applyOutputPresentationKind(parts, "deck.pptx"); err != nil {
		t.Fatalf("apply: %v", err)
	}

	got := string(parts[common.ContentTypesPath])
	if !strings.Contains(got, presentationMainContentType) {
		t.Fatalf("expected the presentation main content type:\n%s", got)
	}
	if strings.Contains(got, templateMainContentType) {
		t.Fatalf("expected the template content type to be gone:\n%s", got)
	}
}

func TestApplyOutputPresentationKindWritesTemplate(t *testing.T) {
	parts := contentTypesWithMain(presentationMainContentType)

	if err := applyOutputPresentationKind(parts, "theme.potx"); err != nil {
		t.Fatalf("apply: %v", err)
	}

	if !strings.Contains(string(parts[common.ContentTypesPath]), templateMainContentType) {
		t.Fatalf("expected the template main content type:\n%s", parts[common.ContentTypesPath])
	}
}

func TestApplyOutputPresentationKindLeavesMatchingPackageAlone(t *testing.T) {
	parts := contentTypesWithMain(presentationMainContentType)
	before := string(parts[common.ContentTypesPath])

	if err := applyOutputPresentationKind(parts, "deck.pptx"); err != nil {
		t.Fatalf("apply: %v", err)
	}

	if string(parts[common.ContentTypesPath]) != before {
		t.Fatal("expected no rewrite when the package already declares the right kind")
	}
}

func TestApplyOutputPresentationKindRejectsPackageWithNoMainPart(t *testing.T) {
	parts := map[string][]byte{common.ContentTypesPath: []byte(`<Types/>`)}

	if err := applyOutputPresentationKind(parts, "deck.pptx"); err == nil {
		t.Fatal("expected an error when no main content type is present")
	}
}

func TestApplyOutputPresentationKindIgnoresUnknownExtension(t *testing.T) {
	parts := contentTypesWithMain(templateMainContentType)

	if err := applyOutputPresentationKind(parts, "deck.zip"); err != nil {
		t.Fatalf("apply: %v", err)
	}

	if !strings.Contains(string(parts[common.ContentTypesPath]), templateMainContentType) {
		t.Fatal("expected an unlisted extension to leave the package as it is")
	}
}
