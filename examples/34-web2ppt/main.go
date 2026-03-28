// examples/34-web2ppt/main.go documents the Web2PPT feature — converting
// HTML content extracted from web pages into PPTX presentations.
//
// The urlfetch package provides HTMLToPPTX and HTMLToPPTXWithOptions for
// converting raw HTML strings to []byte PPTX data. A real fetcher would
// download the HTML from a URL first; here we use static HTML strings so
// the example is self-contained and works without network access.
//
// Run with: go run ./examples/34-web2ppt/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "34_web2ppt.pptx"
)

// sampleHTML is a minimal article page used to demonstrate content extraction.
const sampleHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Web2PPT Feature Overview</title>
  <meta name="description" content="How gopptx converts web pages to presentations">
</head>
<body>
  <main>
    <h1>Web2PPT: Webpage to Presentation</h1>
    <p>The urlfetch package converts HTML content into slide decks automatically, extracting headings, paragraphs, and lists into structured slides.</p>

    <h2>Content Selectors</h2>
    <p>gopptx tries a prioritised list of CSS selectors to locate the main content area:</p>
    <ul>
      <li>main — primary content landmark</li>
      <li>article — article content element</li>
      <li>.entry-content — common blog post class</li>
      <li>Custom selectors configurable via WithContentSelectors()</li>
    </ul>

    <h2>Conversion Options</h2>
    <ul>
      <li>WithMaxSlides(n) — cap the number of generated slides</li>
      <li>WithMaxBullets(n) — cap bullet points per slide</li>
      <li>WithCode(true) — include code block slides</li>
      <li>WithDownloadImages(true) — fetch and embed referenced images</li>
      <li>WithExcludeSelectors([]) — strip nav, footer, ads, etc.</li>
    </ul>

    <h2>Usage Pattern</h2>
    <pre><code>pptxBytes, err := urlfetch.HTMLToPPTX(html, "https://example.com")
// or with full control:
cfg := urlfetch.DefaultConfig().WithMaxSlides(8).WithCode(true)
opts := urlfetch.DefaultConversionOptions().WithTitle("My Deck").WithSourceURL(true)
pptxBytes, err = urlfetch.HTMLToPPTXWithOptions(html, url, cfg, opts)</code></pre>
  </main>
</body>
</html>`

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// --- Example 1: Convert HTML with defaults ---
	log.Printf("Converting sample HTML with default config...")
	pptxBytes, err := urlfetch.HTMLToPPTX(sampleHTML, "https://example.com/web2ppt")
	if err != nil {
		return fmt.Errorf("HTMLToPPTX: %w", err)
	}
	defaultPath := filepath.Join(outputDir, "34_web2ppt_default.pptx")
	if err = os.WriteFile(defaultPath, pptxBytes, 0o600); err != nil {
		return fmt.Errorf("write default: %w", err)
	}
	log.Printf("Generated %s (%d bytes)\n", defaultPath, len(pptxBytes))

	// --- Example 2: Custom config — limit slides and include code blocks ---
	log.Printf("Converting sample HTML with custom config...")
	cfg := urlfetch.DefaultConfig().
		WithMaxSlides(5).
		WithMaxBullets(4).
		WithCode(true)

	opts := urlfetch.DefaultConversionOptions().
		WithTitle("Web2PPT Feature Overview").
		WithAuthor("gopptx").
		WithSourceURL(true)

	pptxBytes2, err := urlfetch.HTMLToPPTXWithOptions(
		sampleHTML,
		"https://example.com/web2ppt",
		cfg,
		opts,
	)
	if err != nil {
		return fmt.Errorf("HTMLToPPTXWithOptions: %w", err)
	}
	customPath := filepath.Join(outputDir, "34_web2ppt_custom.pptx")
	if err = os.WriteFile(customPath, pptxBytes2, 0o600); err != nil {
		return fmt.Errorf("write custom: %w", err)
	}
	log.Printf("Generated %s (%d bytes)\n", customPath, len(pptxBytes2))

	// --- Example 3: Build a reference presentation documenting Web2PPT ---
	docSlides := []pptx.SlideContent{
		pptx.NewSlide("Web2PPT - Webpage to Presentation").
			AddBullet("Convert web pages directly to PPTX").
			AddBullet("Extracts title, headings, paragraphs, and lists").
			AddBullet("Uses goquery for HTML parsing and content selection"),

		pptx.NewSlide("Supported Content Selectors").
			AddNumbered("main — primary content area").
			AddNumbered("article — article content element").
			AddNumbered(".entry-content — common blog post class").
			AddNumbered("Custom CSS selectors via WithContentSelectors()"),

		pptx.NewSlide("urlfetch.Config Options").
			AddBullet("WithMaxSlides(n)          — cap generated slide count").
			AddBullet("WithMaxBullets(n)         — cap bullets per slide").
			AddBullet("WithCode(true)            — include code block slides").
			AddBullet("WithDownloadImages(true)  — embed images from the page").
			AddBullet("WithExcludeSelectors([])  — strip nav, footer, ads"),

		pptx.NewSlide("Quick Usage").
			AddBullet("import \"github.com/djinn-soul/gopptx/pkg/pptx/urlfetch\"").
			AddBullet("bytes, err := urlfetch.HTMLToPPTX(html, url)").
			AddBullet("os.WriteFile(\"out.pptx\", bytes, 0o600)"),
	}

	docData, err := pptx.CreateWithSlides("Web2PPT Reference", docSlides)
	if err != nil {
		return fmt.Errorf("create doc: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, docData, 0o600); err != nil {
		return fmt.Errorf("write doc: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
