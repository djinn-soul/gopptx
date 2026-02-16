package editor

import (
	"errors"
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
