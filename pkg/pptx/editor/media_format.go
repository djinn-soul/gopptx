package editor

import (
	"fmt"
	"strings"
)

// MediaFormat models ppt-rs media format variants for editor media insertion.
type MediaFormat string

const (
	MediaFormatMP4  MediaFormat = "mp4"
	MediaFormatWEBM MediaFormat = "webm"
	MediaFormatAVI  MediaFormat = "avi"
	MediaFormatWMV  MediaFormat = "wmv"
	MediaFormatMOV  MediaFormat = "mov"
	MediaFormatMKV  MediaFormat = "mkv"
	MediaFormatM4V  MediaFormat = "m4v"
	MediaFormatMP3  MediaFormat = "mp3"
	MediaFormatWAV  MediaFormat = "wav"
	MediaFormatWMA  MediaFormat = "wma"
	MediaFormatM4A  MediaFormat = "m4a"
	MediaFormatOGG  MediaFormat = "ogg"
	MediaFormatFLAC MediaFormat = "flac"
	MediaFormatAAC  MediaFormat = "aac"
)

// MIMEType returns the canonical MIME type for the format.
func (f MediaFormat) MIMEType() string {
	switch strings.ToLower(string(f)) {
	case "mp4":
		return "video/mp4"
	case "webm":
		return "video/webm"
	case "avi":
		return "video/x-msvideo"
	case "wmv":
		return "video/x-ms-wmv"
	case "mov":
		return "video/quicktime"
	case "mkv":
		return "video/x-matroska"
	case "m4v":
		return "video/x-m4v"
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "wma":
		return "audio/x-ms-wma"
	case "m4a":
		return "audio/mp4"
	case "ogg":
		return "audio/ogg"
	case "flac":
		return "audio/flac"
	case "aac":
		return "audio/aac"
	default:
		return ""
	}
}

// IsVideo reports whether the format is a video type.
func (f MediaFormat) IsVideo() bool {
	switch strings.ToLower(string(f)) {
	case "mp4", "webm", "avi", "wmv", "mov", "mkv", "m4v":
		return true
	default:
		return false
	}
}

// IsAudio reports whether the format is an audio type.
func (f MediaFormat) IsAudio() bool {
	return f.MIMEType() != "" && !f.IsVideo()
}

// ParseMediaFormat parses a format extension into MediaFormat.
func ParseMediaFormat(ext string) (MediaFormat, bool) {
	switch strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".") {
	case "mp4":
		return MediaFormatMP4, true
	case "webm":
		return MediaFormatWEBM, true
	case "avi":
		return MediaFormatAVI, true
	case "wmv":
		return MediaFormatWMV, true
	case "mov":
		return MediaFormatMOV, true
	case "mkv":
		return MediaFormatMKV, true
	case "m4v":
		return MediaFormatM4V, true
	case "mp3":
		return MediaFormatMP3, true
	case "wav":
		return MediaFormatWAV, true
	case "wma":
		return MediaFormatWMA, true
	case "m4a":
		return MediaFormatM4A, true
	case "ogg":
		return MediaFormatOGG, true
	case "flac":
		return MediaFormatFLAC, true
	case "aac":
		return MediaFormatAAC, true
	default:
		return "", false
	}
}

// AddVideoWithFormat inserts video media using a typed format.
func (e *PresentationEditor) AddVideoWithFormat(
	slideIndex int,
	videoData []byte,
	posterFrameData []byte,
	format MediaFormat,
	x, y, w, h float64,
) (int, error) {
	if !format.IsVideo() {
		return 0, fmt.Errorf("media format %q is not a video format", format)
	}
	return e.AddVideo(slideIndex, videoData, posterFrameData, format.MIMEType(), x, y, w, h)
}

// AddAudioWithFormat inserts audio media using a typed format.
func (e *PresentationEditor) AddAudioWithFormat(
	slideIndex int,
	audioData []byte,
	format MediaFormat,
	x, y, w, h float64,
) (int, error) {
	if !format.IsAudio() {
		return 0, fmt.Errorf("media format %q is not an audio format", format)
	}
	return e.AddAudio(slideIndex, audioData, format.MIMEType(), x, y, w, h)
}
