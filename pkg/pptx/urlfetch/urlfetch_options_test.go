package urlfetch_test

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	urlfetch "github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"
)

func readZipPart(t *testing.T, b []byte, name string) string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatalf("zip reader: %v", err)
	}
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		defer rc.Close()
		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		return string(data)
	}
	t.Fatalf("missing zip part %s", name)
	return ""
}

func TestFetchWithURLReturnsCanonicalRedirectTarget(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			_, _ = w.Write([]byte("<html><body>ok</body></html>"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	f := urlfetch.NewWebFetcherWithConfig(urlfetch.DefaultConfig())
	finalURL, _, err := f.FetchWithURL(srv.URL + "/start")
	if err != nil {
		t.Fatalf("FetchWithURL: %v", err)
	}
	if !strings.HasSuffix(finalURL, "/final") {
		t.Fatalf("expected canonical /final URL, got %s", finalURL)
	}
}

func TestFetchRejectsOversizeBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", 200)))
	}))
	defer srv.Close()

	cfg := urlfetch.DefaultConfig().WithMaxBodyBytes(32)
	f := urlfetch.NewWebFetcherWithConfig(cfg)
	_, err := f.Fetch(srv.URL)
	if err == nil || !strings.Contains(err.Error(), "response too large") {
		t.Fatalf("expected response size error, got %v", err)
	}
}

func TestHTMLToPPTXWithOptionsAppliesAuthorAndPageNumbers(t *testing.T) {
	html := `<!doctype html><html><head><title>A</title></head><body><main><h1>A</h1><p>` +
		strings.Repeat("content ", 30) + `</p></main></body></html>`

	opts := urlfetch.DefaultConversionOptions().
		WithAuthor("urlfetch-author").
		WithPageNumbers(true)

	data, err := urlfetch.HTMLToPPTXWithOptions(html, "https://example.com", urlfetch.DefaultConfig(), opts)
	if err != nil {
		t.Fatalf("HTMLToPPTXWithOptions: %v", err)
	}

	coreXML := readZipPart(t, data, "docProps/core.xml")
	if !strings.Contains(coreXML, "urlfetch-author") {
		t.Fatalf("expected creator in core.xml, got: %s", coreXML)
	}

	slideXML := readZipPart(t, data, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `type="sldNum"`) {
		t.Fatalf("expected slide number placeholder in slide XML")
	}
}

func TestURLToPPTXWithOptionsFollowsRedirects(t *testing.T) {
	const body = `<!doctype html><html><head><title>X</title></head><body><main><h1>X</h1><p>` +
		`This paragraph is deliberately long enough to pass parser thresholds and create a slide body. ` +
		`It keeps going so the extracted main text length is comfortably above one hundred characters. ` +
		`That avoids ErrNoContent from the main-content detector.</p></main></body></html>`

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			_, _ = w.Write([]byte(body))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	data, err := urlfetch.URLToPPTXWithOptions(
		srv.URL+"/start",
		urlfetch.DefaultConfig(),
		urlfetch.DefaultConversionOptions(),
	)
	if err != nil {
		t.Fatalf("URLToPPTXWithOptions: %v", err)
	}
	if !isPPTX(data) {
		t.Fatalf("expected PPTX output for redirected URL")
	}
}

func TestParserExtractLinksToggle(t *testing.T) {
	html := `<!doctype html><html><head><title>L</title></head><body><main><h1>L</h1>` +
		`<p>` + strings.Repeat("text ", 30) + `</p><a href="https://example.com/docs">Docs</a></main></body></html>`

	pOn := urlfetch.NewWebParserWithConfig(urlfetch.DefaultConfig().WithLinks(true))
	wcOn, err := pOn.Parse(html, "https://example.com")
	if err != nil {
		t.Fatalf("parse links on: %v", err)
	}
	found := false
	for _, b := range wcOn.Blocks {
		if b.Kind == urlfetch.KindLink && b.LinkHref == "https://example.com/docs" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected link block when ExtractLinks=true")
	}

	pOff := urlfetch.NewWebParserWithConfig(urlfetch.DefaultConfig().WithLinks(false))
	wcOff, err := pOff.Parse(html, "https://example.com")
	if err != nil {
		t.Fatalf("parse links off: %v", err)
	}
	for _, b := range wcOff.Blocks {
		if b.Kind == urlfetch.KindLink {
			t.Fatalf("did not expect link block when ExtractLinks=false")
		}
	}
}
