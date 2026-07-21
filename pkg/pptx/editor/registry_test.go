package editor

import (
	"errors"
	"sync"
	"testing"
)

func TestRegistry(t *testing.T) {
	registry := NewEditorRegistry()
	e := &PresentationEditor{} // Dummy editor
	h := RegisterEditor(registry, e)
	if h == 0 {
		t.Fatal("expected non-zero handle")
	}

	retrieved, ok := GetEditor(registry, h)
	if !ok || retrieved != e {
		t.Fatal("failed to retrieve editor")
	}

	UnregisterEditor(registry, h)
	_, ok = GetEditor(registry, h)
	if ok {
		t.Fatal("handle should be unregistered")
	}
}

// TestLockEditorSerializesConcurrentAccess asserts that two callers sharing one
// handle never touch the editor simultaneously. cgo drops the Python GIL for the
// duration of a call, so this is the only thing standing between two Python
// threads and a torn write. Run under -race.
func TestLockEditorSerializesConcurrentAccess(t *testing.T) {
	registry := NewEditorRegistry()
	e := &PresentationEditor{}
	h := RegisterEditor(registry, e)
	defer UnregisterEditor(registry, h)

	const goroutines, increments = 8, 200

	var wg sync.WaitGroup
	for range goroutines {
		wg.Go(func() {
			for range increments {
				locked, unlock, ok := LockEditor(registry, h)
				if !ok {
					t.Error("LockEditor failed for a live handle")
					return
				}
				// Deliberately a non-atomic read-modify-write: the lock is what
				// makes it safe, so -race flags any regression here.
				locked.nextSlideNum++
				unlock()
			}
		})
	}
	wg.Wait()

	if got := e.nextSlideNum; got != goroutines*increments {
		t.Fatalf("nextSlideNum = %d, want %d (lost updates)", got, goroutines*increments)
	}
}

func TestLockEditorRejectsUnregisteredHandle(t *testing.T) {
	registry := NewEditorRegistry()
	h := RegisterEditor(registry, &PresentationEditor{})

	if _, unlock, ok := LockEditor(registry, h); !ok {
		t.Fatal("expected live handle to lock")
	} else {
		unlock()
	}

	UnregisterEditor(registry, h)

	if _, _, ok := LockEditor(registry, h); ok {
		t.Fatal("expected unregistered handle to fail to lock")
	}
	if _, _, ok := LockEditor(registry, Handle(999999)); ok {
		t.Fatal("expected unknown handle to fail to lock")
	}
}

func TestRegistryErrorTracking(t *testing.T) {
	registry := NewEditorRegistry()
	e := &PresentationEditor{}
	h := RegisterEditor(registry, e)
	defer UnregisterEditor(registry, h)

	const errMsg = "test error"
	SetHandleError(registry, h, errors.New(errMsg))

	if GetHandleError(registry, h) != errMsg {
		t.Fatalf("expected %q, got %q", errMsg, GetHandleError(registry, h))
	}
}

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
