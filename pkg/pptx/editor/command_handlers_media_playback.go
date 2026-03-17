package editor

import (
	"encoding/json"
	"strings"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
)

const maxUint32Value = int(^uint32(0))

func shouldUseMediaPlaybackCommand(payload json.RawMessage) bool {
	raw := strings.ToLower(string(payload))
	for _, token := range []string{
		`"auto_play"`,
		`"loop"`,
		`"muted"`,
		`"volume"`,
		`"hide_when_stopped"`,
		`"hide_during_show"`,
		`"play_across_slides"`,
	} {
		if strings.Contains(raw, token) {
			return true
		}
	}
	return false
}

func handleAddVideoWithPlaybackCommand(
	e *PresentationEditor,
	payload json.RawMessage,
) (any, error) {
	p, placement, v, err := parseMediaInsertPayload(e, payload)
	if err != nil {
		return nil, err
	}

	mimeType := v.OptionalString(p, "mime_type")
	videoPath := v.OptionalString(p, "path")
	videoData, decodeErr := editorcommand.DecodeOptionalBase64Field(
		v.OptionalString(p, "data"),
		maxMediaBase64,
		"video",
	)
	if decodeErr != nil {
		return nil, decodeErr
	}
	posterPath := v.OptionalString(p, "poster_path")
	posterData, decodeErr := editorcommand.DecodeOptionalBase64Field(
		v.OptionalString(p, "poster_data"),
		maxMediaBase64,
		"poster",
	)
	if decodeErr != nil {
		return nil, decodeErr
	}

	opts := parseVideoPlaybackOptionsPayload(p, v)
	if vErr := v.Error(); vErr != nil {
		return nil, vErr
	}

	shapeID, err := editorcommand.InsertShapeFromBinaryOrPath(
		len(videoData) > 0 || len(posterData) > 0,
		func() (int, error) {
			return e.AddVideoWithPlaybackOptions(
				placement.SlideIndex,
				videoData,
				posterData,
				mimeType,
				opts,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
		func() (int, error) {
			return e.AddVideoFromFileWithPlaybackOptions(
				placement.SlideIndex,
				videoPath,
				posterPath,
				mimeType,
				opts,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": shapeID}, nil
}

//nolint:funlen // Handler keeps parse/validate/apply steps together for stable bridge error semantics.
func handleAddAudioWithPlaybackCommand(
	e *PresentationEditor,
	payload json.RawMessage,
) (any, error) {
	p, placement, v, err := parseMediaInsertPayload(e, payload)
	if err != nil {
		return nil, err
	}

	mimeType := v.OptionalString(p, "mime_type")
	audioPath := v.OptionalString(p, "path")
	audioData, decodeErr := editorcommand.DecodeOptionalBase64Field(
		v.OptionalString(p, "data"),
		maxMediaBase64,
		"audio",
	)
	if decodeErr != nil {
		return nil, decodeErr
	}
	iconPath := v.OptionalString(p, "icon_path")
	iconData, decodeErr := editorcommand.DecodeOptionalBase64Field(
		v.OptionalString(p, "icon_data"),
		maxMediaBase64,
		"icon",
	)
	if decodeErr != nil {
		return nil, decodeErr
	}

	opts := parseAudioPlaybackOptionsPayload(p, v)
	if vErr := v.Error(); vErr != nil {
		return nil, vErr
	}

	hasBinary := len(audioData) > 0 || len(iconData) > 0
	hasIcon := len(iconData) > 0 || strings.TrimSpace(iconPath) != ""
	shapeID, err := editorcommand.InsertShapeFromBinaryOrPath(
		hasBinary,
		func() (int, error) {
			if hasIcon {
				return e.AddAudioWithIcon(
					placement.SlideIndex,
					audioData,
					iconData,
					mimeType,
					placement.X,
					placement.Y,
					placement.W,
					placement.H,
				)
			}
			return e.AddAudioWithPlaybackOptions(
				placement.SlideIndex,
				audioData,
				mimeType,
				opts,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
		func() (int, error) {
			if hasIcon {
				return e.AddAudioWithIconFromFile(
					placement.SlideIndex,
					audioPath,
					iconPath,
					mimeType,
					placement.X,
					placement.Y,
					placement.W,
					placement.H,
				)
			}
			return e.AddAudioFromFileWithPlaybackOptions(
				placement.SlideIndex,
				audioPath,
				mimeType,
				opts,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		},
	)
	if err != nil {
		return nil, err
	}

	// Icon-specific insertion APIs do not currently carry playback options, so apply timing now.
	if hasIcon {
		if timingErr := e.applyMediaPlaybackTiming(
			placement.SlideIndex,
			shapeID,
			"audio",
			editormodmedia.MediaTimingOptions{
				AutoPlay:         opts.AutoPlay,
				LoopPlayback:     opts.LoopPlayback,
				Muted:            false,
				Volume:           opts.Volume,
				ShowWhenStopped:  !opts.HideDuringShow,
				PlayAcrossSlides: opts.PlayAcrossSlides,
				SlideIndex:       placement.SlideIndex,
				SlideCount:       len(e.slides),
			},
		); timingErr != nil {
			return nil, timingErr
		}
	}
	return map[string]int{"shape_id": shapeID}, nil
}

func parseMediaInsertPayload(
	e *PresentationEditor,
	payload json.RawMessage,
) (map[string]any, editorcommand.MediaPlacement, *PayloadValidator, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, editorcommand.MediaPlacement{}, nil, err
	}
	v := NewPayloadValidator()
	placement, ok := editorcommand.ParseMediaPlacement(
		p,
		e.SlideCount(),
		v.RequireInt,
		v.RequireFloat64,
		v.IndexBounds,
	)
	if !ok {
		return nil, editorcommand.MediaPlacement{}, v, v.Error()
	}
	return p, placement, v, nil
}

func parseVideoPlaybackOptionsPayload(
	payload map[string]any,
	v *PayloadValidator,
) VideoPlaybackOptions {
	opts := NewVideoPlaybackOptions()
	if val, ok := v.OptionalBool(payload, "auto_play"); ok {
		opts.AutoPlay = val
	}
	if val, ok := v.OptionalBool(payload, "loop"); ok {
		opts.LoopPlayback = val
	}
	if val, ok := v.OptionalBool(payload, "muted"); ok {
		opts.Muted = val
	}
	if val, ok := v.OptionalBool(payload, "hide_when_stopped"); ok {
		opts.HideWhenStopped = val
	}
	if val, ok := v.OptionalInt(payload, "volume"); ok {
		if val < 0 || val > maxUint32Value {
			v.invalidType("volume", "a non-negative integer within uint32 range", val)
		} else {
			opts = opts.WithVolume(uint32(val))
		}
	}
	return opts
}

func parseAudioPlaybackOptionsPayload(
	payload map[string]any,
	v *PayloadValidator,
) AudioPlaybackOptions {
	opts := NewAudioPlaybackOptions()
	if val, ok := v.OptionalBool(payload, "auto_play"); ok {
		opts.AutoPlay = val
	}
	if val, ok := v.OptionalBool(payload, "loop"); ok {
		opts.LoopPlayback = val
	}
	if val, ok := v.OptionalBool(payload, "play_across_slides"); ok {
		opts.PlayAcrossSlides = val
	}
	if val, ok := v.OptionalBool(payload, "hide_during_show"); ok {
		opts.HideDuringShow = val
	}
	if val, ok := v.OptionalInt(payload, "volume"); ok {
		if val < 0 || val > maxUint32Value {
			v.invalidType("volume", "a non-negative integer within uint32 range", val)
		} else {
			opts = opts.WithVolume(uint32(val))
		}
	}
	return opts
}
