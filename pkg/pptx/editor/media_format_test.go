package editor

import "testing"

func TestMediaFormatHelpers(t *testing.T) {
	t.Parallel()

	format, ok := ParseMediaFormat(".WMV")
	if !ok {
		t.Fatal("expected wmv format to parse")
	}
	if format != MediaFormatWMV {
		t.Fatalf("expected MediaFormatWMV, got %q", format)
	}
	if !format.IsVideo() || format.IsAudio() {
		t.Fatalf("expected wmv to be video-only, got isVideo=%v isAudio=%v", format.IsVideo(), format.IsAudio())
	}
	if got := format.MIMEType(); got != "video/x-ms-wmv" {
		t.Fatalf("expected wmv mime, got %q", got)
	}

	audio := MediaFormatOGG
	if audio.IsVideo() || !audio.IsAudio() {
		t.Fatalf("expected ogg to be audio-only, got isVideo=%v isAudio=%v", audio.IsVideo(), audio.IsAudio())
	}
	if got := audio.MIMEType(); got != "audio/ogg" {
		t.Fatalf("expected ogg mime, got %q", got)
	}

	videoMKV, ok := ParseMediaFormat("mkv")
	if !ok || videoMKV != MediaFormatMKV {
		t.Fatalf("expected mkv parse support, got format=%q ok=%v", videoMKV, ok)
	}
	if got := videoMKV.MIMEType(); got != "video/x-matroska" {
		t.Fatalf("expected mkv mime, got %q", got)
	}

	audioFLAC, ok := ParseMediaFormat("flac")
	if !ok || audioFLAC != MediaFormatFLAC {
		t.Fatalf("expected flac parse support, got format=%q ok=%v", audioFLAC, ok)
	}
	if got := audioFLAC.MIMEType(); got != "audio/flac" {
		t.Fatalf("expected flac mime, got %q", got)
	}
}

func TestAddVideoWithFormatRejectsAudioFormat(t *testing.T) {
	e := newMediaEditorFixture()
	if _, err := e.AddVideoWithFormat(0, []byte("video"), []byte("poster"), MediaFormatMP3, 1, 2, 3, 4); err == nil {
		t.Fatal("expected AddVideoWithFormat to reject audio format")
	}
}

func TestAddAudioWithFormatRejectsVideoFormat(t *testing.T) {
	e := newMediaEditorFixture()
	if _, err := e.AddAudioWithFormat(0, []byte("audio"), MediaFormatMOV, 1, 2, 3, 4); err == nil {
		t.Fatal("expected AddAudioWithFormat to reject video format")
	}
}
