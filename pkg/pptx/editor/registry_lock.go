package editor

// Per-handle locking for editors held in a [Registry].
//
// The registry's own mutex guards the handle map, not the editors inside it.
// Guarding only the map still allows two callers to mutate one editor at once,
// which matters because cgo releases the Python GIL for the duration of a call:
// two Python threads sharing a handle really do run concurrently.
//
// These locks are acquired only at the C boundary. Acquiring them inside
// [PresentationEditor] methods would deadlock, because command handlers such as
// export call back into e.Save on the same goroutine.

// LockEditor retrieves an editor and acquires its per-handle lock, returning a
// release function the caller must invoke when done.
//
// Callers crossing the C boundary must use this rather than [Registry.GetEditor].
func (r *Registry) LockEditor(h Handle) (*PresentationEditor, func(), bool) {
	return r.lockEditor(h, false)
}

// TryLockEditor behaves like [Registry.LockEditor] but reports failure instead of
// blocking when the handle is busy.
//
// This is what makes cross-handle operations such as merge_from_editor safe: the
// caller already holds the destination lock, so blocking on a second handle could
// deadlock against a thread merging in the opposite direction. Two such threads
// both fail their try-lock and return an error instead of hanging.
func (r *Registry) TryLockEditor(h Handle) (*PresentationEditor, func(), bool) {
	return r.lockEditor(h, true)
}

func (r *Registry) lockEditor(h Handle, try bool) (*PresentationEditor, func(), bool) {
	noop := func() {}
	if r == nil {
		return nil, noop, false
	}
	r.mu.RLock()
	reg, ok := r.editors[h]
	r.mu.RUnlock()

	if !ok {
		return nil, noop, false
	}
	if try {
		if !reg.mu.TryLock() {
			return nil, noop, false
		}
	} else {
		reg.mu.Lock()
	}
	// The handle may have been unregistered while we waited for reg.mu.
	if reg.editor == nil {
		reg.mu.Unlock()
		return nil, noop, false
	}
	return reg.editor, reg.mu.Unlock, true
}

// LockEditor retrieves an editor from a registry and locks it for exclusive use.
// The returned release function must be called when the operation completes.
func LockEditor(registry *Registry, h Handle) (*PresentationEditor, func(), bool) {
	if registry == nil {
		return nil, func() {}, false
	}
	return registry.LockEditor(h)
}

// TryLockEditor locks an editor for exclusive use, failing instead of blocking
// when the handle is busy. The release function must be called on success.
func TryLockEditor(registry *Registry, h Handle) (*PresentationEditor, func(), bool) {
	if registry == nil {
		return nil, func() {}, false
	}
	return registry.TryLockEditor(h)
}
