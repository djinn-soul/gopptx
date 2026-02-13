package editor

import "testing"

func TestEnsureAuthorsLoadedLocked_DoesNotPoisonCacheOnParseFailure(t *testing.T) {
	ps := NewPartStore()
	ps.Set(commentAuthorsPartName, []byte(`<p:cmAuthorLst><p:cmAuthor`))

	e := &PresentationEditor{
		parts: ps,
	}

	if _, err := e.GetAuthors(); err == nil {
		t.Fatal("expected parse error from malformed commentAuthors.xml")
	}
	if _, err := e.AddAuthor("Alice", "AL"); err == nil {
		t.Fatal("expected add author to fail after malformed commentAuthors.xml")
	}
}
