package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/markdown"
	"github.com/djinn-soul/gopptx/pkg/pptx/mermaid"
)

// editorLookupFn resolves an editor handle (int64) to a *PresentationEditor.
// Registered at startup by the bridge (bridge.go) to avoid an import cycle.
//
//nolint:gochecknoglobals // Required for cross-package init-time registration.
var editorLookupFn func(int64) (*PresentationEditor, bool)

// RegisterEditorLookupFn wires the bridge's registry lookup into the editor package.
// Must be called before any merge_from_editor requests are processed.
func RegisterEditorLookupFn(fn func(int64) (*PresentationEditor, bool)) {
	editorLookupFn = fn
}

// editorTryLockFn acquires the per-handle lock for a source editor without
// blocking, returning a release function. See [RegisterEditorTryLockFn].
//
//nolint:gochecknoglobals // Required for cross-package init-time registration.
var editorTryLockFn func(int64) (func(), bool)

// RegisterEditorTryLockFn wires the bridge's non-blocking per-handle lock into
// the editor package, so cross-handle operations can synchronize against the
// source presentation without risking deadlock. See [Registry.TryLockEditor].
func RegisterEditorTryLockFn(fn func(int64) (func(), bool)) {
	editorTryLockFn = fn
}

// handleMergeFromEditor appends all slides from another open editor into this one.
//
// Payload: {"source_handle": N} where N is the integer handle of the source presentation.
// Response: {"merged": true}.
func handleMergeFromEditor(e *PresentationEditor, payload json.RawMessage) (any, error) {
	if editorLookupFn == nil {
		return nil, errors.New("editor lookup not initialized: call editor.RegisterEditorLookupFn at startup")
	}
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	handleVal, ok := v.RequireInt(p, "source_handle")
	if !ok {
		return nil, v.Error()
	}
	src, found := editorLookupFn(int64(handleVal))
	if !found {
		return nil, NewBridgeError(ErrCodeOpFailed, fmt.Sprintf("source handle %d not found", handleVal))
	}
	// The caller already holds the destination lock. Take the source lock too so
	// the merge cannot read a presentation another thread is mutating — but never
	// block on it, or two threads merging each other's decks would deadlock.
	// A self-merge needs no second acquire: we already hold that lock.
	if src != e && editorTryLockFn != nil {
		release, locked := editorTryLockFn(int64(handleVal))
		if !locked {
			return nil, NewBridgeError(
				ErrCodeOpFailed,
				fmt.Sprintf("source handle %d is busy in another operation", handleVal),
			)
		}
		defer release()
	}
	if mergeErr := e.MergeFromEditor(src); mergeErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, mergeErr.Error())
	}
	return respMerged, nil
}

// urlfetchHandlers holds optional handlers for URL-fetch ops registered by
// external packages to avoid an import cycle:
// editor → urlfetch → pptx → editor.
//
//nolint:gochecknoglobals // Global registry required for cross-package init-time registration.
var urlfetchHandlers = map[string]commandHandler{
	// Placeholder stub: replaced by RegisterURLFetchHandler at startup.
	OpURLFetchToSlides: func(_ *PresentationEditor, _ json.RawMessage) (any, error) {
		return nil, errors.New("url_fetch not initialized: call editorurlfetch.Register() at startup")
	},
}

// RegisterURLFetchHandler registers a handler for the given operation, replacing
// any existing placeholder. Must be called before bridge requests are processed.
func RegisterURLFetchHandler(op string, h func(*PresentationEditor, json.RawMessage) (any, error)) {
	urlfetchHandlers[op] = h
}

// handleMarkdownToSlides converts a Markdown string into one or more slides
// and appends them to the presentation.
//
// Payload: {"markdown": "<string>", "layout": "<string>" (optional)}.
// Response: {"slide_count": N, "first_index": I}.
func handleMarkdownToSlides(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	md, ok := v.RequireString(p, "markdown")
	if !ok {
		return nil, v.Error()
	}
	layout := v.OptionalString(p, "layout")

	slides, parseErr := markdown.SlidesFromMarkdown(md)
	if parseErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, parseErr.Error())
	}

	firstIndex := -1
	for i, slide := range slides {
		if layout != "" {
			slide = slide.WithLayout(layout)
		}
		idx, addErr := e.AddSlide(slide)
		if addErr != nil {
			return nil, NewBridgeError(ErrCodeOpFailed, addErr.Error())
		}
		if i == 0 {
			firstIndex = idx
		}
	}

	return markdownSlidesResponse{
		SlideCount: len(slides),
		FirstIndex: firstIndex,
	}, nil
}

// handleURLFetchToSlides delegates to a lazily-registered handler to avoid
// an import cycle between the editor and urlfetch packages.
func handleURLFetchToSlides(e *PresentationEditor, payload json.RawMessage) (any, error) {
	h, ok := urlfetchHandlers[OpURLFetchToSlides]
	if !ok {
		return nil, errors.New("url_fetch handler not registered")
	}
	return h(e, payload)
}

// handleAddMermaidShape creates a Mermaid diagram on an existing slide.
//
// Payload: {"slide_index": N, "diagram": "<mermaid code>"}.
// Response: {"shape_count": N, "connector_count": N}.
func handleAddMermaidShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	diagram, ok := v.RequireString(p, "diagram")
	if !ok {
		return nil, v.Error()
	}

	diagramElements, diagErr := mermaid.CreateDiagram(diagram)
	if diagErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, diagErr.Error())
	}

	addedShapes := 0
	addedConnectors := 0

	for _, s := range diagramElements.Shapes {
		x := float64(s.X)
		y := float64(s.Y)
		w := float64(s.CX)
		h := float64(s.CY)
		shapeType := s.Type
		if shapeType == "" {
			shapeType = shapeTypeRect
		}
		_, shapeErr := e.AddShape(slideIndex, shapeType, x, y, w, h)
		if shapeErr != nil {
			return nil, NewBridgeError(ErrCodeOpFailed, shapeErr.Error())
		}
		addedShapes++
	}

	// Connectors are represented as line shapes from start to end.
	for _, c := range diagramElements.Connectors {
		x := float64(c.StartX)
		y := float64(c.StartY)
		endX := float64(c.EndX)
		endY := float64(c.EndY)
		w := endX - x
		h := endY - y
		if w < 0 {
			x = endX
			w = -w
		}
		if h < 0 {
			y = endY
			h = -h
		}
		if w < 1 {
			w = 1
		}
		if h < 1 {
			h = 1
		}
		_, connErr := e.AddShape(slideIndex, "line", x, y, w, h)
		if connErr != nil {
			return nil, NewBridgeError(ErrCodeOpFailed, connErr.Error())
		}
		addedConnectors++
	}

	return mermaidAddResponse{
		ShapeCount:     addedShapes,
		ConnectorCount: addedConnectors,
	}, nil
}
