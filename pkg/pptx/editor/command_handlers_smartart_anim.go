package editor

import (
	"encoding/json"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// handleAddSmartArt adds a SmartArt diagram to an existing slide.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "layout": "<layout URI or token>",   // required
//	  "items": ["text1", "text2", ...],    // required
//	  "x": <float>, "y": <float>,          // position in EMU (optional, defaults to 1 inch / 2 inches)
//	  "cx": <float>, "cy": <float>         // size in EMU (optional, defaults to 8 × 4 inches)
//	}
//
// Response: {"shape_id": N}.
func handleAddSmartArt(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	layoutURI, ok := v.RequireString(p, "layout")
	if !ok {
		return nil, v.Error()
	}

	// Items list (optional, defaults to empty).
	items, _ := v.OptionalStringSlice(p, "items")

	// Optional bounds in EMU. Use Go default sizes if omitted.
	const (
		defaultX  int64 = 914400  // 1 inch
		defaultY  int64 = 1828800 // 2 inches
		defaultCX int64 = 7315200 // 8 inches
		defaultCY int64 = 3657600 // 4 inches
	)
	x := optionalInt64OrDefault(p, "x", defaultX)
	y := optionalInt64OrDefault(p, "y", defaultY)
	cx := optionalInt64OrDefault(p, "cx", defaultCX)
	cy := optionalInt64OrDefault(p, "cy", defaultCY)

	layout := smartart.CustomLayout(layoutURI)
	sa := smartart.NewSmartArt(layout).AddItems(items)

	// Override position/size from payload.
	sa.X = styling.Length(x)
	sa.Y = styling.Length(y)
	sa.CX = styling.Length(cx)
	sa.CY = styling.Length(cy)

	shapeID, addErr := e.AddSmartArt(slideIndex, sa)
	if addErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, addErr.Error())
	}
	return map[string]int{"shape_id": shapeID}, nil
}

// handleAddAnimation adds an animation effect to a shape on an existing slide.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "shape_id": N,
//	  "effect": "entr_fade",       // required, e.g. "entr_appear", "exit_flyOut"
//	  "trigger": "onClick",        // optional: "onClick" | "withPrev" | "afterPrev"
//	  "duration_ms": 500,          // optional (default 500)
//	  "delay_ms": 0                // optional (default 0)
//	}
//
// Response: {"added": true}.
func handleAddAnimation(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	effect, ok := v.RequireString(p, "effect")
	if !ok {
		return nil, v.Error()
	}

	trigger := v.OptionalString(p, "trigger")
	if trigger == "" {
		trigger = "onClick"
	}
	durationMS, _ := v.OptionalInt(p, "duration_ms")
	if durationMS <= 0 {
		durationMS = 500
	}
	delayMS, _ := v.OptionalInt(p, "delay_ms")

	if addErr := e.AddSlideAnimation(slideIndex, shapeID, effect, trigger, durationMS, delayMS); addErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, addErr.Error())
	}
	return map[string]bool{"added": true}, nil
}

// handleSetSlideTransition sets the transition for an existing slide.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "transition_type": "fade",    // required; use "none" to remove
//	  "duration_ms": 0,             // optional (0 = default)
//	  "advance_ms": -1              // optional (-1 = click-advance only)
//	}
//
// Response: {"updated": true}.
func handleSetSlideTransition(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	transitionType, ok := v.RequireString(p, "transition_type")
	if !ok {
		return nil, v.Error()
	}

	durationMS, _ := v.OptionalInt(p, "duration_ms")
	advanceMS := optionalIntOrDefault(p, "advance_ms", -1)

	if setErr := e.SetSlideTransition(slideIndex, transitionType, durationMS, advanceMS); setErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, setErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}

// optionalInt64OrDefault returns the int64 value of a map key, or the default if missing/invalid.
func optionalInt64OrDefault(p map[string]any, key string, def int64) int64 {
	v, ok := p[key]
	if !ok {
		return def
	}
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	default:
		return def
	}
}

// optionalIntOrDefault returns the int value of a map key, or the default if missing/invalid.
func optionalIntOrDefault(p map[string]any, key string, def int) int {
	v, ok := p[key]
	if !ok {
		return def
	}
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	default:
		return def
	}
}
