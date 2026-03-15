package media

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecodeBase64ImagePayloadAndDataURI(t *testing.T) {
	rawBytes := []byte{0x89, 0x50, 0x4E, 0x47}
	rawB64 := base64.StdEncoding.EncodeToString(rawBytes)

	data, format, err := DecodeBase64ImagePayload(rawB64)
	if err != nil {
		t.Fatalf("DecodeBase64ImagePayload(raw) failed: %v", err)
	}
	if format != "" || string(data) != string(rawBytes) {
		t.Fatalf("unexpected raw decode result: format=%q data=%v", format, data)
	}

	uri := "data:image/png;base64," + rawB64
	data, format, err = DecodeBase64ImagePayload(uri)
	if err != nil {
		t.Fatalf("DecodeBase64ImagePayload(data-uri) failed: %v", err)
	}
	if format != "png" || string(data) != string(rawBytes) {
		t.Fatalf("unexpected data-uri decode result: format=%q data=%v", format, data)
	}

	if _, _, err = DecodeBase64ImagePayload(""); err == nil {
		t.Fatal("expected empty payload error")
	}
	if _, _, err = DecodeBase64ImagePayload("data:image/png,abc"); err == nil {
		t.Fatal("expected non-base64 data URI error")
	}
	if _, _, err = DecodeBase64ImagePayload("%%%"); err == nil {
		t.Fatal("expected invalid base64 error")
	}
}

func TestFetchImageFromURLAndFormatDetection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte{1, 2, 3})
		case "/ext.jpg":
			_, _ = w.Write([]byte{4, 5, 6})
		case "/empty":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte{})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	data, format, err := FetchImageFromURL(server.URL + "/ok")
	if err != nil {
		t.Fatalf("FetchImageFromURL(/ok) failed: %v", err)
	}
	if format != "png" || len(data) != 3 {
		t.Fatalf("unexpected fetch result: format=%q len=%d", format, len(data))
	}

	data, format, err = FetchImageFromURL(server.URL + "/ext.jpg")
	if err != nil {
		t.Fatalf("FetchImageFromURL(/ext.jpg) failed: %v", err)
	}
	if format != "jpg" || len(data) != 3 {
		t.Fatalf("unexpected extension-derived format: format=%q len=%d", format, len(data))
	}

	if _, _, err = FetchImageFromURL(server.URL + "/missing"); err == nil {
		t.Fatal("expected non-200 status error")
	}
	if _, _, err = FetchImageFromURL(server.URL + "/empty"); err == nil {
		t.Fatal("expected empty body error")
	}
	if _, _, err = FetchImageFromURL("://bad-url"); err == nil {
		t.Fatal("expected invalid url error")
	}
}

func TestInsertAndSlideHelpers(t *testing.T) {
	if err := ValidateMediaSlideIndex(1, 1); err == nil {
		t.Fatal("expected out-of-range slide index error")
	}
	if err := ValidateMediaSlideIndex(0, 1); err != nil {
		t.Fatalf("unexpected slide index error: %v", err)
	}

	calledData := false
	_, err := RegisterPartFromDataOrPath(
		[]byte{1},
		"",
		"missing",
		func(data []byte) (string, error) {
			calledData = true
			return "ppt/media/x.bin", nil
		},
		func(string) (string, error) {
			t.Fatal("fromPath should not be called when data exists")
			return "", nil
		},
	)
	if err != nil || !calledData {
		t.Fatalf("RegisterPartFromDataOrPath data path failed: err=%v calledData=%v", err, calledData)
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "sample.bin")
	if writeErr := os.WriteFile(tmpFile, []byte("abc"), 0o600); writeErr != nil {
		t.Fatalf("failed to write temp file: %v", writeErr)
	}
	embedPath, err := RegisterEmbeddingPart(nil, tmpFile, func(data []byte, ext string) (string, error) {
		if ext != "bin" || string(data) != "abc" {
			t.Fatalf("unexpected embedding payload ext=%q data=%q", ext, string(data))
		}
		return "ppt/embeddings/ole1.bin", nil
	})
	if err != nil || embedPath == "" {
		t.Fatalf("RegisterEmbeddingPart(path) failed: path=%q err=%v", embedPath, err)
	}

	slideXML := []byte(`<p:sld><p:cSld><p:spTree></p:spTree></p:cSld></p:sld>`)
	appended, err := AppendShapeXMLToSlide(slideXML, `<p:pic id="1"/>`)
	if err != nil {
		t.Fatalf("AppendShapeXMLToSlide failed: %v", err)
	}
	if !strings.Contains(string(appended), `<p:pic id="1"/>`) {
		t.Fatalf("expected shape xml inserted before spTree close: %s", string(appended))
	}
	if _, err = AppendShapeXMLToSlide([]byte("<p:sld/>"), "<p:pic/>"); err == nil {
		t.Fatal("expected missing spTree end error")
	}
}

func TestTimingHelpersAndShapeXMLBuilders(t *testing.T) {
	if len(DefaultVideoPosterPNG()) == 0 {
		t.Fatal("DefaultVideoPosterPNG should return non-empty bytes")
	}

	slideXML := `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
		`<p:cSld><p:spTree>` +
		`<p:pic><p:nvPicPr><p:cNvPr id="7" name="Video"></p:cNvPr></p:nvPicPr><p:nvPr><p:extLst><p:ext><p14:media r:embed="rIdMedia1"/></p:ext></p:extLst></p:nvPr></p:pic>` +
		`</p:spTree></p:cSld></p:sld>`
	withTiming, err := addDefaultTimingBlock(slideXML)
	if err != nil {
		t.Fatalf("addDefaultTimingBlock failed: %v", err)
	}
	if !strings.Contains(withTiming, "<p:timing>") {
		t.Fatalf("expected timing block insertion: %s", withTiming)
	}
	if _, err = addDefaultTimingBlock("<p:sld>"); err == nil {
		t.Fatal("expected invalid slide xml close-tag error")
	}

	if got := nextTimingCTnID(`<p:cTn id="4"/><p:cTn id="22"/>`); got != 23 {
		t.Fatalf("nextTimingCTnID=%d, want 23", got)
	}
	inserted, err := insertMediaNodeIntoMainSeq(
		`<p:timing><p:cTn nodeType="mainSeq"><p:childTnLst></p:childTnLst></p:cTn></p:timing>`,
		`<p:audio/>`,
	)
	if err != nil || !strings.Contains(inserted, "<p:audio/>") {
		t.Fatalf("insertMediaNodeIntoMainSeq failed: inserted=%q err=%v", inserted, err)
	}
	if _, err = insertMediaNodeIntoMainSeq(`<p:timing></p:timing>`, `<p:audio/>`); err == nil {
		t.Fatal("expected missing mainSeq error")
	}

	relID := mediaRelIDForShape(withTiming, 7)
	if relID != "rIdMedia1" {
		t.Fatalf("mediaRelIDForShape=%q, want rIdMedia1", relID)
	}
	if relID = mediaRelIDForShape(withTiming, 9); relID != "" {
		t.Fatalf("expected missing shape rel id to be empty, got %q", relID)
	}

	options := MediaTimingOptions{
		AutoPlay:         true,
		LoopPlayback:     true,
		Muted:            true,
		Volume:           200,
		ShowWhenStopped:  false,
		PlayAcrossSlides: true,
		SlideIndex:       1,
		SlideCount:       4,
	}
	updated, err := ApplyMediaTiming([]byte(withTiming), "audio", 7, options)
	if err != nil {
		t.Fatalf("ApplyMediaTiming failed: %v", err)
	}
	updatedXML := string(updated)
	if !strings.Contains(updatedXML, `spid="7"`) || !strings.Contains(updatedXML, `vol="100000"`) {
		t.Fatalf("expected media timing node with normalized volume and spid: %s", updatedXML)
	}
	if !strings.Contains(updatedXML, `numSld="3"`) {
		t.Fatalf("expected cross-slide numSld=3, got: %s", updatedXML)
	}

	videoXML := BuildVideoShapeXML(5, "rIdV", "rIdM", "rIdP", `A "quoted" & value`, 1, 2, 3, 4)
	if !strings.Contains(videoXML, `descr="A &quot;quoted&quot; &amp; value"`) {
		t.Fatalf("expected escaped video alt text: %s", videoXML)
	}
	audioXML := BuildAudioShapeXML(6, "rIdA", "rIdM", "", `Audio's`, 1, 2, 3, 4)
	if strings.Contains(audioXML, `r:embed=""`) {
		t.Fatalf("empty icon rel should not emit empty embed attribute: %s", audioXML)
	}
	oleXML := BuildOLEObjectShapeXML(7, 2, "rIdE", "rIdI", `Word.Document`, 1, 2, 3, 4)
	if !strings.Contains(oleXML, `progId="Word.Document"`) || !strings.Contains(oleXML, `spid="_x0000_s21031"`) {
		t.Fatalf("unexpected OLE object shape xml: %s", oleXML)
	}

	if got := normalizeMediaVolume(50); got != 50000 {
		t.Fatalf("normalizeMediaVolume(50)=%d, want 50000", got)
	}
	if got := normalizeMediaVolume(101); got != 100000 {
		t.Fatalf("normalizeMediaVolume cap failed: %d", got)
	}
}
