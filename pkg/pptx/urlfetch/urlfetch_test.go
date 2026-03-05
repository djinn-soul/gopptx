package urlfetch_test

import (
	"archive/zip"
	"bytes"
	"strings"
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
	data, err := urlfetch.NewURLFetchConverterWithConfig(cfg).Convert(wc, &opts)
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
	data, err := urlfetch.NewURLFetchConverterWithConfig(cfg).Convert(wc, &opts)
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
	data, err := urlfetch.NewURLFetchConverterWithConfig(cfg).Convert(wc, &opts)
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
	data, err := urlfetch.NewURLFetchConverter().Convert(wc, &opts)
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

func TestConfig_NewBuilderOptions(t *testing.T) {
	cfg := urlfetch.DefaultConfig().
		WithDownloadImages(true).
		WithMaxImageSizeBytes(2 * 1024 * 1024).
		WithMaxTotalImageSizeBytes(10 * 1024 * 1024).
		WithMaxImagesPerSlide(5).
		WithAllowedImageTypes([]string{"image/png", "image/jpeg"}).
		WithContentSelectors([]string{"article", ".content"}).
		WithExcludeSelectors([]string{"nav", ".ads"})

	if !cfg.DownloadImages {
		t.Error("DownloadImages should be true")
	}
	if cfg.MaxImageSizeBytes != 2*1024*1024 {
		t.Errorf("MaxImageSizeBytes: want %d, got %d", 2*1024*1024, cfg.MaxImageSizeBytes)
	}
	if cfg.MaxTotalImageSizeBytes != 10*1024*1024 {
		t.Errorf("MaxTotalImageSizeBytes: want %d, got %d", 10*1024*1024, cfg.MaxTotalImageSizeBytes)
	}
	if cfg.MaxImagesPerSlide != 5 {
		t.Errorf("MaxImagesPerSlide: want 5, got %d", cfg.MaxImagesPerSlide)
	}
	if len(cfg.AllowedImageTypes) != 2 {
		t.Errorf("AllowedImageTypes length: want 2, got %d", len(cfg.AllowedImageTypes))
	}
	if len(cfg.ContentSelectors) != 2 {
		t.Errorf("ContentSelectors length: want 2, got %d", len(cfg.ContentSelectors))
	}
	if len(cfg.ExcludeSelectors) != 2 {
		t.Errorf("ExcludeSelectors length: want 2, got %d", len(cfg.ExcludeSelectors))
	}
}

// --- Custom selectors tests ---

func TestParser_CustomContentSelectors(t *testing.T) {
	customHTML := `<html><body>
		<div class="ignored">This should be ignored because it does not have enough text to meet the minimum requirements.</div>
		<article class="content">
			<h1>Custom Article</h1>
			<p>This is the main content that should be extracted from the article element. It has plenty of text to ensure it passes the minimum content length check.</p>
		</article>
	</body></html>`

	cfg := urlfetch.DefaultConfig().
		WithContentSelectors([]string{"article.content"})

	p := urlfetch.NewWebParserWithConfig(cfg)
	wc, err := p.Parse(customHTML, "https://example.com")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if wc.Title != "Custom Article" {
		t.Errorf("expected title 'Custom Article', got %q", wc.Title)
	}

	// Should find content from article, not from ignored div
	foundContent := false
	for _, block := range wc.Blocks {
		if block.Kind == urlfetch.KindParagraph && strings.Contains(block.Text, "main content") {
			foundContent = true
			break
		}
	}
	if !foundContent {
		t.Error("expected to find paragraph from article.content")
	}
}

func TestParser_ExcludeSelectors(t *testing.T) {
	htmlWithAds := `<html><body>
		<main>
			<h1>Page Title</h1>
			<p>This is good content that should be kept and has enough text to pass the minimum length requirement.</p>
			<div class="advertisement">Buy our product now! Limited time offer!</div>
			<p>This is also good content with sufficient length to be included.</p>
		</main>
	</body></html>`

	cfg := urlfetch.DefaultConfig().
		WithExcludeSelectors([]string{".advertisement"})

	p := urlfetch.NewWebParserWithConfig(cfg)
	wc, err := p.Parse(htmlWithAds, "https://example.com")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Should not contain advertisement text
	for _, block := range wc.Blocks {
		if strings.Contains(block.Text, "Buy our product") {
			t.Error("advertisement content should have been excluded")
		}
	}

	// Should still have good content
	goodContentCount := 0
	for _, block := range wc.Blocks {
		if block.Kind == urlfetch.KindParagraph && strings.Contains(block.Text, "good content") {
			goodContentCount++
		}
	}
	if goodContentCount != 2 {
		t.Errorf("expected 2 good content paragraphs, got %d", goodContentCount)
	}
}

func TestParser_ImageURLResolution(t *testing.T) {
	htmlWithRelativeImage := `<html><body>
		<main>
			<h1>Test</h1>
			<p>This paragraph has enough text to meet the minimum content length requirement for the parser to find valid content.</p>
			<img src="/images/test.png" alt="Test Image">
			<img src="https://cdn.example.com/absolute.png" alt="Absolute Image">
			<img src="//protocol-relative.jpg" alt="Protocol Relative">
		</main>
	</body></html>`

	p := urlfetch.NewWebParser()
	wc, err := p.Parse(htmlWithRelativeImage, "https://example.com/page/article")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	var foundRelative, foundAbsolute, foundProtocolRelative bool
	for _, img := range wc.Images {
		src := img[0]
		if src == "https://example.com/images/test.png" {
			foundRelative = true
		}
		if src == "https://cdn.example.com/absolute.png" {
			foundAbsolute = true
		}
		if src == "https://protocol-relative.jpg" {
			foundProtocolRelative = true
		}
	}

	if !foundRelative {
		t.Error("expected relative URL to be resolved to https://example.com/images/test.png")
	}
	if !foundAbsolute {
		t.Error("expected absolute URL to remain unchanged")
	}
	if !foundProtocolRelative {
		t.Error("expected protocol-relative URL to be resolved to https://protocol-relative.jpg")
	}
}
