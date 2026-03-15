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
	case string(MediaFormatMP4):
		return "video/mp4"
	case string(MediaFormatWEBM):
		return "video/webm"
	case string(MediaFormatAVI):
		return "video/x-msvideo"
	case string(MediaFormatWMV):
		return "video/x-ms-wmv"
	case string(MediaFormatMOV):
		return "video/quicktime"
	case string(MediaFormatMKV):
		return "video/x-matroska"
	case string(MediaFormatM4V):
		return "video/x-m4v"
	case string(MediaFormatMP3):
		return "audio/mpeg"
	case string(MediaFormatWAV):
		return "audio/wav"
	case string(MediaFormatWMA):
		return "audio/x-ms-wma"
	case string(MediaFormatM4A):
		return "audio/mp4"
	case string(MediaFormatOGG):
		return "audio/ogg"
	case string(MediaFormatFLAC):
		return "audio/flac"
	case string(MediaFormatAAC):
		return "audio/aac"
	default:
		return ""
	}
}

// IsVideo reports whether the format is a video type.
func (f MediaFormat) IsVideo() bool {
	switch strings.ToLower(string(f)) {
	case string(MediaFormatMP4),
		string(MediaFormatWEBM),
		string(MediaFormatAVI),
		string(MediaFormatWMV),
		string(MediaFormatMOV),
		string(MediaFormatMKV),
		string(MediaFormatM4V):
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
	case string(MediaFormatMP4):
		return MediaFormatMP4, true
	case string(MediaFormatWEBM):
		return MediaFormatWEBM, true
	case string(MediaFormatAVI):
		return MediaFormatAVI, true
	case string(MediaFormatWMV):
		return MediaFormatWMV, true
	case string(MediaFormatMOV):
		return MediaFormatMOV, true
	case string(MediaFormatMKV):
		return MediaFormatMKV, true
	case string(MediaFormatM4V):
		return MediaFormatM4V, true
	case string(MediaFormatMP3):
		return MediaFormatMP3, true
	case string(MediaFormatWAV):
		return MediaFormatWAV, true
	case string(MediaFormatWMA):
		return MediaFormatWMA, true
	case string(MediaFormatM4A):
		return MediaFormatM4A, true
	case string(MediaFormatOGG):
		return MediaFormatOGG, true
	case string(MediaFormatFLAC):
		return MediaFormatFLAC, true
	case string(MediaFormatAAC):
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
