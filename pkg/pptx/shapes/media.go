package shapes

import (
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Video and audio could be inserted into an existing deck through the editor,
// but a deck being generated had nowhere to put them: SlideContent had no media
// field and the media package was images-only. These types are the authoring
// side of what the editor already writes.

// MediaKind distinguishes a video clip from an audio one.
type MediaKind int

const (
	// MediaVideo plays in a picture frame on the slide.
	MediaVideo MediaKind = iota
	// MediaAudio shows as a speaker icon.
	MediaAudio
)

// MediaPlayback is how a clip behaves during the show.
type MediaPlayback struct {
	// AutoPlay starts the clip when the slide appears.
	AutoPlay bool
	// Loop restarts the clip when it finishes.
	Loop bool
	// Muted plays with the volume off.
	Muted bool
	// HideWhenStopped hides the frame while the clip is not playing (video).
	HideWhenStopped bool
	// PlayAcrossSlides keeps audio playing as the show moves on (audio).
	PlayAcrossSlides bool
	// StartMS and EndMS trim the clip, in milliseconds. Nil plays it whole.
	StartMS *uint32
	EndMS   *uint32
	// Volume is 0–100. Zero means the clip's own level.
	Volume uint32
}

// Media is a video or audio clip placed on a slide.
type Media struct {
	Kind MediaKind
	// Path is a file to embed. Data takes precedence when both are set.
	Path string
	// Data is the clip's bytes.
	Data []byte
	// Format is the file extension without the dot ("mp4", "mp3"). Empty is
	// taken from Path.
	Format string
	// SourceURL links an online video instead of embedding it. A clip with a
	// URL carries no media part.
	SourceURL string

	X  styling.Length
	Y  styling.Length
	CX styling.Length
	CY styling.Length

	// PosterData is the still image shown before playback. Without one the
	// frame renders empty in PowerPoint, so a video should normally have one.
	PosterData []byte
	// PosterFormat is the poster's extension without the dot.
	PosterFormat string

	AltText  string
	Playback MediaPlayback
}

// Default media sizes, used when a clip states no extent.
const (
	defaultVideoWidth  = styling.Length(4572000) // 5 inches
	defaultVideoHeight = styling.Length(2571750) // 16:9 at that width
	defaultAudioSide   = styling.Length(457200)  // half an inch, the icon
)

// NewVideo embeds a video file at the given position and size.
func NewVideo(path string, x, y, cx, cy styling.Length) Media {
	return Media{Kind: MediaVideo, Path: path, X: x, Y: y, CX: cx, CY: cy}
}

// NewVideoFromBytes embeds video bytes of the given format.
func NewVideoFromBytes(data []byte, format string, x, y, cx, cy styling.Length) Media {
	return Media{Kind: MediaVideo, Data: data, Format: format, X: x, Y: y, CX: cx, CY: cy}
}

// NewOnlineVideo links a hosted video rather than embedding one.
func NewOnlineVideo(url string, x, y, cx, cy styling.Length) Media {
	return Media{Kind: MediaVideo, SourceURL: url, X: x, Y: y, CX: cx, CY: cy}
}

// NewAudio embeds an audio file, shown as a speaker icon at (x, y).
func NewAudio(path string, x, y styling.Length) Media {
	return Media{Kind: MediaAudio, Path: path, X: x, Y: y, CX: defaultAudioSide, CY: defaultAudioSide}
}

// NewAudioFromBytes embeds audio bytes of the given format.
func NewAudioFromBytes(data []byte, format string, x, y styling.Length) Media {
	return Media{
		Kind: MediaAudio, Data: data, Format: format,
		X: x, Y: y, CX: defaultAudioSide, CY: defaultAudioSide,
	}
}

// WithAutoPlay starts the clip when the slide appears.
func (m Media) WithAutoPlay(enabled bool) Media {
	m.Playback.AutoPlay = enabled
	return m
}

// WithLoop restarts the clip when it finishes.
func (m Media) WithLoop(enabled bool) Media {
	m.Playback.Loop = enabled
	return m
}

// WithMuted plays the clip with the volume off.
func (m Media) WithMuted(enabled bool) Media {
	m.Playback.Muted = enabled
	return m
}

// WithVolume sets the playback level, 0–100.
func (m Media) WithVolume(volume uint32) Media {
	const maxVolume = 100
	m.Playback.Volume = min(volume, maxVolume)
	return m
}

// WithTrim plays only the given span of the clip, in milliseconds.
func (m Media) WithTrim(startMS, endMS uint32) Media {
	m.Playback.StartMS = &startMS
	m.Playback.EndMS = &endMS
	return m
}

// WithHideWhenStopped hides a video frame while it is not playing.
func (m Media) WithHideWhenStopped(enabled bool) Media {
	m.Playback.HideWhenStopped = enabled
	return m
}

// WithPlayAcrossSlides keeps audio playing as the show moves on.
func (m Media) WithPlayAcrossSlides(enabled bool) Media {
	m.Playback.PlayAcrossSlides = enabled
	return m
}

// WithPoster sets the still shown before playback.
func (m Media) WithPoster(data []byte, format string) Media {
	m.PosterData, m.PosterFormat = data, format
	return m
}

// WithAltText describes the clip for screen readers.
func (m Media) WithAltText(text string) Media {
	m.AltText = text
	return m
}

// IsOnline reports whether the clip links a hosted video rather than embedding
// bytes.
func (m Media) IsOnline() bool {
	return strings.TrimSpace(m.SourceURL) != "" && len(m.Data) == 0 && m.Path == ""
}

// FormatOrExtension is the clip's file extension without the dot, taken from
// Format or from the Path.
func (m Media) FormatOrExtension() string {
	if format := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(m.Format)), "."); format != "" {
		return format
	}
	if ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(m.Path)), "."); ext != "" {
		return ext
	}
	if m.Kind == MediaAudio {
		return "mp3"
	}
	return "mp4"
}

// SizeOrDefault is the clip's extent, defaulted per kind when unset.
func (m Media) SizeOrDefault() (styling.Length, styling.Length) {
	if m.CX > 0 && m.CY > 0 {
		return m.CX, m.CY
	}
	if m.Kind == MediaAudio {
		return defaultAudioSide, defaultAudioSide
	}
	return defaultVideoWidth, defaultVideoHeight
}
