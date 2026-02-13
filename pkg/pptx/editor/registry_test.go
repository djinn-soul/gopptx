package editor

import (
	"errors"
	"testing"
)

func TestRegistry(t *testing.T) {
	e := &PresentationEditor{} // Dummy editor
	h := RegisterEditor(e)
	if h == 0 {
		t.Fatal("expected non-zero handle")
	}

	retrieved, ok := GetEditor(h)
	if !ok || retrieved != e {
		t.Fatal("failed to retrieve editor")
	}

	UnregisterEditor(h)
	_, ok = GetEditor(h)
	if ok {
		t.Fatal("handle should be unregistered")
	}
}

func TestRegistryErrorTracking(t *testing.T) {
	e := &PresentationEditor{}
	h := RegisterEditor(e)
	defer UnregisterEditor(h)

	const errMsg = "test error"
	SetHandleError(h, errors.New(errMsg))

	if GetHandleError(h) != errMsg {
		t.Fatalf("expected %q, got %q", errMsg, GetHandleError(h))
	}
}
