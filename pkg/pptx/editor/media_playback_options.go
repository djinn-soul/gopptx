package editor

import (
	"fmt"

	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
)

// VideoPlaybackOptions mirrors ppt-rs generator video playback options.
type VideoPlaybackOptions struct {
	AutoPlay        bool
	LoopPlayback    bool
	HideWhenStopped bool
	Muted           bool
	StartTimeMS     *uint32
	EndTimeMS       *uint32
	Volume          uint32 // 0-100
	AltText         string
}

// NewVideoPlaybackOptions returns default video playback options.
func NewVideoPlaybackOptions() VideoPlaybackOptions {
	return VideoPlaybackOptions{Volume: 100}
}

// NewAutoPlayVideoPlaybackOptions returns default video options with autoplay.
func NewAutoPlayVideoPlaybackOptions() VideoPlaybackOptions {
	opts := NewVideoPlaybackOptions()
	opts.AutoPlay = true
	return opts
}

func (o VideoPlaybackOptions) WithLoop(loop bool) VideoPlaybackOptions {
	o.LoopPlayback = loop
	return o
}

func (o VideoPlaybackOptions) WithMuted(muted bool) VideoPlaybackOptions {
	o.Muted = muted
	return o
}

func (o VideoPlaybackOptions) WithVolume(volume uint32) VideoPlaybackOptions {
	if volume > 100 {
		volume = 100
	}
	o.Volume = volume
	return o
}

func (o VideoPlaybackOptions) WithStartTimeMS(ms uint32) VideoPlaybackOptions {
	o.StartTimeMS = &ms
	return o
}

func (o VideoPlaybackOptions) WithEndTimeMS(ms uint32) VideoPlaybackOptions {
	o.EndTimeMS = &ms
	return o
}

func (o VideoPlaybackOptions) WithAltText(altText string) VideoPlaybackOptions {
	o.AltText = altText
	return o
}

// AudioPlaybackOptions mirrors ppt-rs generator audio playback options.
type AudioPlaybackOptions struct {
	AutoPlay         bool
	LoopPlayback     bool
	HideDuringShow   bool
	PlayAcrossSlides bool
	Volume           uint32 // 0-100
	AltText          string
}

// NewAudioPlaybackOptions returns default audio playback options.
func NewAudioPlaybackOptions() AudioPlaybackOptions {
	return AudioPlaybackOptions{Volume: 100}
}

// NewAutoPlayAudioPlaybackOptions returns default audio options with autoplay.
func NewAutoPlayAudioPlaybackOptions() AudioPlaybackOptions {
	opts := NewAudioPlaybackOptions()
	opts.AutoPlay = true
	return opts
}

func (o AudioPlaybackOptions) WithLoop(loop bool) AudioPlaybackOptions {
	o.LoopPlayback = loop
	return o
}

func (o AudioPlaybackOptions) WithPlayAcrossSlides(play bool) AudioPlaybackOptions {
	o.PlayAcrossSlides = play
	return o
}

func (o AudioPlaybackOptions) WithVolume(volume uint32) AudioPlaybackOptions {
	if volume > 100 {
		volume = 100
	}
	o.Volume = volume
	return o
}

func (o AudioPlaybackOptions) WithAltText(altText string) AudioPlaybackOptions {
	o.AltText = altText
	return o
}

// AddVideoWithPlaybackOptions inserts video media and applies playback timing metadata.
func (e *PresentationEditor) AddVideoWithPlaybackOptions(
	slideIndex int,
	videoData []byte,
	posterFrameData []byte,
	mimeType string,
	options VideoPlaybackOptions,
	x, y, w, h float64,
) (int, error) {
	shapeID, err := e.addVideoGeneric(
		slideIndex,
		videoData,
		"",
		posterFrameData,
		"",
		mimeType,
		options.AltText,
		x,
		y,
		w,
		h,
	)
	if err != nil {
		return 0, err
	}
	if err := e.applyMediaPlaybackTiming(slideIndex, shapeID, "video", editormodmedia.MediaTimingOptions{
		AutoPlay:         options.AutoPlay,
		LoopPlayback:     options.LoopPlayback,
		Muted:            options.Muted,
		Volume:           options.Volume,
		ShowWhenStopped:  !options.HideWhenStopped,
		PlayAcrossSlides: false,
		SlideIndex:       slideIndex,
		SlideCount:       len(e.slides),
	}); err != nil {
		return 0, err
	}
	return shapeID, nil
}

// AddVideoFromFileWithPlaybackOptions inserts video media from a local file and
// applies playback timing metadata.
func (e *PresentationEditor) AddVideoFromFileWithPlaybackOptions(
	slideIndex int,
	videoPath string,
	posterFramePath string,
	mimeType string,
	options VideoPlaybackOptions,
	x, y, w, h float64,
) (int, error) {
	shapeID, err := e.addVideoGeneric(
		slideIndex,
		nil,
		videoPath,
		nil,
		posterFramePath,
		mimeType,
		options.AltText,
		x,
		y,
		w,
		h,
	)
	if err != nil {
		return 0, err
	}
	if err := e.applyMediaPlaybackTiming(slideIndex, shapeID, "video", editormodmedia.MediaTimingOptions{
		AutoPlay:         options.AutoPlay,
		LoopPlayback:     options.LoopPlayback,
		Muted:            options.Muted,
		Volume:           options.Volume,
		ShowWhenStopped:  !options.HideWhenStopped,
		PlayAcrossSlides: false,
		SlideIndex:       slideIndex,
		SlideCount:       len(e.slides),
	}); err != nil {
		return 0, err
	}
	return shapeID, nil
}

// AddAudioWithPlaybackOptions inserts audio media and applies playback timing metadata.
func (e *PresentationEditor) AddAudioWithPlaybackOptions(
	slideIndex int,
	audioData []byte,
	mimeType string,
	options AudioPlaybackOptions,
	x, y, w, h float64,
) (int, error) {
	shapeID, err := e.addAudioGeneric(slideIndex, audioData, "", nil, "", mimeType, options.AltText, x, y, w, h)
	if err != nil {
		return 0, err
	}
	if err := e.applyMediaPlaybackTiming(slideIndex, shapeID, "audio", editormodmedia.MediaTimingOptions{
		AutoPlay:         options.AutoPlay,
		LoopPlayback:     options.LoopPlayback,
		Muted:            false,
		Volume:           options.Volume,
		ShowWhenStopped:  !options.HideDuringShow,
		PlayAcrossSlides: options.PlayAcrossSlides,
		SlideIndex:       slideIndex,
		SlideCount:       len(e.slides),
	}); err != nil {
		return 0, err
	}
	return shapeID, nil
}

// AddAudioFromFileWithPlaybackOptions inserts audio media from a local file and
// applies playback timing metadata.
func (e *PresentationEditor) AddAudioFromFileWithPlaybackOptions(
	slideIndex int,
	audioPath string,
	mimeType string,
	options AudioPlaybackOptions,
	x, y, w, h float64,
) (int, error) {
	shapeID, err := e.addAudioGeneric(slideIndex, nil, audioPath, nil, "", mimeType, options.AltText, x, y, w, h)
	if err != nil {
		return 0, err
	}
	if err := e.applyMediaPlaybackTiming(slideIndex, shapeID, "audio", editormodmedia.MediaTimingOptions{
		AutoPlay:         options.AutoPlay,
		LoopPlayback:     options.LoopPlayback,
		Muted:            false,
		Volume:           options.Volume,
		ShowWhenStopped:  !options.HideDuringShow,
		PlayAcrossSlides: options.PlayAcrossSlides,
		SlideIndex:       slideIndex,
		SlideCount:       len(e.slides),
	}); err != nil {
		return 0, err
	}
	return shapeID, nil
}

func (e *PresentationEditor) applyMediaPlaybackTiming(
	slideIndex int,
	shapeID int,
	mediaKind string,
	options editormodmedia.MediaTimingOptions,
) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index out of range: %d", slideIndex)
	}
	slidePart := e.slides[slideIndex].Part
	content, ok := e.parts.Get(slidePart)
	if !ok {
		return fmt.Errorf("read slide part: %s", slidePart)
	}
	updated, err := editormodmedia.ApplyMediaTiming(content, mediaKind, shapeID, options)
	if err != nil {
		return fmt.Errorf("apply %s playback timing: %w", mediaKind, err)
	}
	e.parts.Set(slidePart, updated)
	return nil
}
