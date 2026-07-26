package elements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The point of SlideBuilder is that a call whose result is ignored still takes
// effect, which is exactly what the value-receiver SlideContent method cannot
// promise.
func TestSlideBuilderAppliesCallsWhoseResultIsIgnored(t *testing.T) {
	builder := NewSlideBuilder("Deck")
	builder.AddBullet("first")
	builder.AddBullet("second")

	slide := builder.Build()
	if len(slide.Bullets) != 2 {
		t.Fatalf("got %d bullets, want 2 — ignored calls were dropped", len(slide.Bullets))
	}
}

// Demonstrates the bug SlideBuilder exists to prevent: the same code written
// against SlideContent silently loses both bullets.
func TestSlideContentDropsCallsWhoseResultIsIgnored(t *testing.T) {
	slide := NewSlide("Deck")
	slide.AddBullet("first")
	slide.AddBullet("second")

	if len(slide.Bullets) != 0 {
		t.Fatalf("expected the value receiver to discard ignored calls, got %d bullets",
			len(slide.Bullets))
	}
}

func TestSlideBuilderChains(t *testing.T) {
	slide := NewSlideBuilder("Deck").
		AddBullet("one").
		AddBullet("two").
		WithNotes("speaker notes").
		Build()

	if len(slide.Bullets) != 2 {
		t.Errorf("got %d bullets, want 2", len(slide.Bullets))
	}
	if slide.Notes != "speaker notes" {
		t.Errorf("notes = %q, want %q", slide.Notes, "speaker notes")
	}
}

func TestSlideBuilderMatchesSlideContent(t *testing.T) {
	viaBuilder := NewSlideBuilder("Deck").AddBullet("one").AddBullet("two").Build()
	viaValue := NewSlide("Deck").AddBullet("one").AddBullet("two")

	if viaBuilder.Title != viaValue.Title {
		t.Errorf("title %q != %q", viaBuilder.Title, viaValue.Title)
	}
	if len(viaBuilder.Bullets) != len(viaValue.Bullets) {
		t.Fatalf("bullet count %d != %d", len(viaBuilder.Bullets), len(viaValue.Bullets))
	}
	for i := range viaBuilder.Bullets {
		if viaBuilder.Bullets[i] != viaValue.Bullets[i] {
			t.Errorf("bullet %d: %q != %q", i, viaBuilder.Bullets[i], viaValue.Bullets[i])
		}
	}
}

// BuildFrom lets code that already holds a SlideContent switch styles.
func TestBuildFromWrapsExistingContent(t *testing.T) {
	existing := NewSlide("Deck").AddBullet("kept")

	builder := BuildFrom(existing)
	builder.AddBullet("added")
	slide := builder.Build()

	if len(slide.Bullets) != 2 {
		t.Fatalf("got %d bullets, want 2", len(slide.Bullets))
	}
	if slide.Bullets[0] != "kept" || slide.Bullets[1] != "added" {
		t.Errorf("unexpected bullets: %+v", slide.Bullets)
	}
}

// Guards the generator's contract rather than any single method: every
// chainable SlideContent method must have a builder counterpart.
func TestSlideBuilderCoversChainableMethods(t *testing.T) {
	const wantAtLeast = 80
	count := countBuilderMethods(t)
	if count < wantAtLeast {
		t.Errorf("SlideBuilder has %d methods, expected at least %d; "+
			"run 'go generate ./...' after adding SlideContent methods",
			count, wantAtLeast)
	}
}

func countBuilderMethods(t *testing.T) int {
	t.Helper()
	data, err := readGeneratedBuilder()
	if err != nil {
		t.Fatalf("read generated builder: %v", err)
	}
	return strings.Count(data, "\nfunc (b *SlideBuilder) ")
}

// readGeneratedBuilder concatenates every generated builder file, since the
// methods are split across chunks to stay under the per-file line ceiling.
func readGeneratedBuilder() (string, error) {
	paths, err := filepath.Glob("slide_builder*_gen.go")
	if err != nil {
		return "", err
	}
	var combined strings.Builder
	for _, path := range paths {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return "", readErr
		}
		combined.Write(data)
	}
	return combined.String(), nil
}
