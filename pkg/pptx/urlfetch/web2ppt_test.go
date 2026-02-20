package urlfetch_test

import (
	"archive/zip"
	"bytes"
	"testing"

	urlfetch "github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"
)

// simpleHTML is a minimal HTML page used for most tests.
const simpleHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Test Page</title>
  <meta name="description" content="A test description">
</head>
<body>
  <main>
    <h1>Main Title</h1>
    <p>This is a paragraph with enough text to be meaningful content for the parser.</p>
    <h2>Section One</h2>
    <p>Section one content here, providing context for the first section of the page.</p>
    <ul>
      <li>List item alpha</li>
      <li>List item beta</li>
      <li>List item gamma</li>
    </ul>
    <h2>Section Two</h2>
    <blockquote>A notable quotation that captures an important idea.</blockquote>
    <pre><code>fmt.Println("hello")</code></pre>
  </main>
</body>
</html>`

const tableHTML = `<!DOCTYPE html>
<html><head><title>Table Test</title></head>
<body>
  <article>
    <h1>Data Overview</h1>
    <p>The following table shows sample data.</p>
    <table>
      <tr><th>Name</th><th>Value</th><th>Unit</th></tr>
      <tr><td>Alpha</td><td>100</td><td>ms</td></tr>
      <tr><td>Beta</td><td>200</td><td>ms</td></tr>
    </table>
  </article>
</body>
</html>`

// isPPTX reports whether b is a valid zip archive (all PPTX are zip).
func isPPTX(b []byte) bool {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return false
	}
	for _, f := range r.File {
		if f.Name == "[Content_Types].xml" {
			return true
		}
	}
	return false
}

// --- Parser tests ---

func TestParser_Title(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if wc.Title == "" {
		t.Error("expected non-empty title")
	}
	if wc.Title != "Test Page" {
		t.Errorf("expected title 'Test Page', got %q", wc.Title)
	}
}

func TestParser_Description(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if wc.Description != "A test description" {
		t.Errorf("unexpected description: %q", wc.Description)
	}
}

func TestParser_BlocksExtracted(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(wc.Blocks) == 0 {
		t.Fatal("expected content blocks, got none")
	}
}

func TestParser_GroupedByHeadings(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	groups := wc.GroupedByHeadings()
	if len(groups) == 0 {
		t.Fatal("expected heading groups, got none")
	}
}

func TestParser_TableRows(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(tableHTML, "https://example.com/data")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	found := false
	for _, b := range wc.Blocks {
		if b.Kind == urlfetch.KindTable {
			found = true
			if len(b.TableRows) != 3 {
				t.Errorf("expected 3 table rows, got %d", len(b.TableRows))
			}
		}
	}
	if !found {
		t.Error("expected a KindTable block")
	}
}

func TestParser_NoContent(t *testing.T) {
	p := urlfetch.NewWebParser()
	_, err := p.Parse("<html><body></body></html>", "https://example.com")
	if err == nil {
		t.Error("expected ErrNoContent for empty page, got nil")
	}
}

// --- Converter tests ---

func TestConverter_GroupedSlides(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	cfg := urlfetch.DefaultConfig()
	cfg.GroupByHeadings = true

	opts := urlfetch.DefaultConversionOptions()
	data, err := urlfetch.NewWeb2PptWithConfig(cfg).Convert(wc, &opts)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX (zip) archive")
	}
}

func TestConverter_LinearSlides(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	cfg := urlfetch.DefaultConfig()
	cfg.GroupByHeadings = false

	opts := urlfetch.DefaultConversionOptions()
	data, err := urlfetch.NewWeb2PptWithConfig(cfg).Convert(wc, &opts)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX (zip) archive")
	}
}

func TestConverter_MaxSlides(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	cfg := urlfetch.DefaultConfig().WithMaxSlides(2)
	opts := urlfetch.DefaultConversionOptions()
	data, err := urlfetch.NewWeb2PptWithConfig(cfg).Convert(wc, &opts)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX (zip) archive")
	}
}

func TestConverter_CustomTitle(t *testing.T) {
	p := urlfetch.NewWebParser()
	wc, err := p.Parse(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	opts := urlfetch.DefaultConversionOptions().WithTitle("My Custom Title")
	data, err := urlfetch.NewWeb2Ppt().Convert(wc, &opts)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX (zip) archive")
	}
}

// --- Package-level API tests ---

func TestHTMLToPPTX(t *testing.T) {
	data, err := urlfetch.HTMLToPPTX(simpleHTML, "https://example.com")
	if err != nil {
		t.Fatalf("HTMLToPPTX: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX archive")
	}
}

func TestHTMLToPPTXWithOptions(t *testing.T) {
	cfg := urlfetch.DefaultConfig().
		WithMaxSlides(5).
		WithMaxBullets(4).
		WithCode(true)

	opts := urlfetch.DefaultConversionOptions().
		WithTitle("Rust Quick Start").
		WithAuthor("gopptx").
		WithSourceURL(true)

	data, err := urlfetch.HTMLToPPTXWithOptions(simpleHTML, "https://example.com", cfg, opts)
	if err != nil {
		t.Fatalf("HTMLToPPTXWithOptions: %v", err)
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX archive")
	}
}

func TestHTMLToPPTX_TablePage(t *testing.T) {
	data, err := urlfetch.HTMLToPPTX(tableHTML, "https://example.com/data")
	if err != nil {
		t.Fatalf("HTMLToPPTX (table page): %v", err)
	}
	if !isPPTX(data) {
		t.Error("output is not a valid PPTX archive")
	}
}

// --- Config builder tests ---

func TestConfig_Builder(t *testing.T) {
	cfg := urlfetch.DefaultConfig().
		WithMaxSlides(10).
		WithMaxBullets(3).
		WithImages(false).
		WithTables(false).
		WithCode(false).
		WithGroupByHeadings(false).
		WithUserAgent("TestAgent/1.0").
		WithTimeout(10)

	if cfg.MaxSlides != 10 {
		t.Errorf("MaxSlides: want 10, got %d", cfg.MaxSlides)
	}
	if cfg.MaxBulletsPerSlide != 3 {
		t.Errorf("MaxBulletsPerSlide: want 3, got %d", cfg.MaxBulletsPerSlide)
	}
	if cfg.IncludeImages {
		t.Error("IncludeImages should be false")
	}
	if cfg.GroupByHeadings {
		t.Error("GroupByHeadings should be false")
	}
}
