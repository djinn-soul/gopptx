# Web-to-PPT (urlfetch) Enhancement Plan

## Overview
Implementation plan for adding Image Download & Embedding and CSS Selector Customization to the `pkg/pptx/urlfetch` package.

---

## Feature #1: Image Download & Embedding

### Current State
- `urlfetch` extracts image `src` and `alt` attributes but only adds alt-text as bullets
- `shapes.Image` supports multiple sources: file path, bytes, base64, URL
- Media catalog handles image deduplication and registration

### Implementation Requirements

#### 1.1 New Config Options (`config.go`)
```go
// Add to Config struct:
// DownloadImages fetches and embeds remote images when true (replaces alt-text bullets).
DownloadImages bool
// MaxImageSizeBytes caps individual downloaded image size (default: 5MB).
MaxImageSizeBytes int64
// MaxTotalImageSizeBytes caps total downloaded image size per page (default: 20MB).
MaxTotalImageSizeBytes int64
// MaxImagesPerSlide limits images per slide (default: 3).
MaxImagesPerSlide int
// AllowedImageTypes filters by MIME type (default: ["image/png", "image/jpeg", "image/gif", "image/webp"]).
AllowedImageTypes []string
// ImagePlacement sets positioning strategy: "inline" | "bottom-right" | "full-width" (default: "inline").
ImagePlacement string
```

#### 1.2 HTTP Image Fetcher (`fetcher.go`)
```go
// ImageFetcher handles concurrent image downloads with size limits.
type ImageFetcher struct {
    client      http.Client
    cfg         Config
    totalSize   int64
    mu          sync.Mutex
}

// FetchImage downloads image with validation:
// - MIME type check against AllowedImageTypes
// - Size check against MaxImageSizeBytes
// - Returns (data []byte, format string, err error)
func (f *ImageFetcher) FetchImage(url string) ([]byte, string, error)
```

#### 1.3 Image Converter Extension (`converter.go`)
```go
// Modify KindImage case in appendBlock():
case KindImage:
    if c.cfg.DownloadImages && block.ImageSrc != "" {
        // Download and embed actual image
        img, err := c.fetchAndCreateImage(block.ImageSrc, block.ImageAlt)
        if err == nil {
            slide = slide.AddImage(img)
            bulletCount++
        } else {
            // Fallback to alt-text bullet on error
            slide = slide.AddBullet("[Image: " + block.ImageAlt + "]")
        }
    } else if c.cfg.IncludeImages && block.ImageAlt != "" {
        // Current behavior: alt-text only
        slide = slide.AddBullet("[Image: " + block.ImageAlt + "]")
    }

// New helper:
func (c *Converter) fetchAndCreateImage(src, alt string) (shapes.Image, error) {
    // Resolve relative URLs against base URL
    // Fetch image bytes
    // Detect format from content-type or extension
    // Create shapes.Image with appropriate dimensions
    // Use default sizing (e.g., 3x2 inches at 96 DPI)
}
```

#### 1.4 URL Resolution (`parser.go` enhancement)
```go
// Add to WebParser:
baseURL *url.URL

// ResolveImageURL converts relative URLs to absolute:
func (p *WebParser) ResolveImageURL(src string) string {
    if p.baseURL == nil {
        return src
    }
    u, err := p.baseURL.Parse(src)
    if err != nil {
        return src
    }
    return u.String()
}
```

#### 1.5 Image Dimension Detection
```go
// Detect image dimensions without full decode using:
// - PNG: IHDR chunk parsing
// - JPEG: SOF markers
// - GIF: Header parsing
// - WebP: VP8 header
//
// Package: use standard library or lightweight decoder
// Keep aspect ratio, scale to MaxImageWidthEMU x MaxImageHeightEMU
```

### Testing Requirements
1. **Unit Tests** (`urlfetch_test.go`):
   - Test successful image download and embedding
   - Test MIME type filtering
   - Test size limit enforcement
   - Test relative URL resolution
   - Test fallback to alt-text on fetch failure

2. **Integration Tests**:
   - Test with real HTTP server (httptest)
   - Test concurrent image downloads
   - Test invalid image handling

### Implementation Steps
1. [ ] Add config options to `config.go` with builder methods
2. [ ] Create `image_fetcher.go` for HTTP image fetching
3. [ ] Extend `converter.go` with image embedding logic
4. [ ] Add URL resolution to `parser.go`
5. [ ] Add image dimension detection utility
6. [ ] Write comprehensive tests
7. [ ] Update example 34 to demonstrate image embedding

---

## Feature #2: CSS Selector Customization

### Current State
- Fixed selector list in `parser.go`:
  ```go
  []string{
      "main article", "article", "main", "[role=main]",
      ".content", ".post-content", ".article-content",
      "#content", "#main", "body",
  }
  ```

### Implementation Requirements

#### 2.1 New Config Options (`config.go`)
```go
// Add to Config struct:
// ContentSelectors overrides default main-content CSS selectors.
// If empty, uses built-in defaults.
ContentSelectors []string
// ExcludeSelectors removes matching elements from content.
ExcludeSelectors []string
// RequireSelector enforces at least one match (fail if not found).
RequireSelector string
```

#### 2.2 Parser Extension (`parser.go`)
```go
// Modify findMainContent():
func (p *WebParser) findMainContent(doc *goquery.Document) *goquery.Selection {
    selectors := p.cfg.ContentSelectors
    if len(selectors) == 0 {
        selectors = mainContentSelectors() // default
    }

    for _, sel := range selectors {
        found := doc.Find(sel).First()
        if found.Length() == 0 {
            continue
        }
        // Apply exclude selectors
        for _, exclude := range p.cfg.ExcludeSelectors {
            found.Find(exclude).Remove()
        }
        if len(strings.TrimSpace(found.Text())) >= minMainTextLen {
            return found
        }
    }
    return nil
}

// Add exclusion to walkSelection():
func (p *WebParser) walkSelection(sel *goquery.Selection, wc *WebContent, depth int) {
    // Check exclude selectors before processing
    for _, excludeSel := range p.cfg.ExcludeSelectors {
        if sel.Is(excludeSel) {
            return
        }
    }
    // ... rest of existing logic
}
```

#### 2.3 Validation (`config.go`)
```go
// Add to config validation (if validation exists):
func (c Config) Validate() error {
    // Validate ContentSelectors are non-empty strings
    for _, sel := range c.ContentSelectors {
        if strings.TrimSpace(sel) == "" {
            return fmt.Errorf("content selector cannot be empty")
        }
    }
    // Validate ExcludeSelectors
    for _, sel := range c.ExcludeSelectors {
        if strings.TrimSpace(sel) == "" {
            return fmt.Errorf("exclude selector cannot be empty")
        }
    }
    return nil
}
```

### Testing Requirements
1. **Unit Tests**:
   - Test custom selector extraction
   - Test exclude selector removal
   - Test fallback to default selectors
   - Test empty result handling

### Implementation Steps
1. [ ] Add selector config options to `config.go`
2. [ ] Modify `findMainContent()` to use custom selectors
3. [ ] Add exclude logic to `walkSelection()`
4. [ ] Add config validation
5. [ ] Write unit tests
6. [ ] Update example 34 with custom selector demo

---

## API Usage Examples

### Image Download Example
```go
cfg := urlfetch.DefaultConfig().
    WithDownloadImages(true).
    WithMaxImageSizeBytes(2 * 1024 * 1024). // 2MB per image
    WithMaxImagesPerSlide(2).
    WithAllowedImageTypes([]string{"image/png", "image/jpeg"})

opts := urlfetch.DefaultConversionOptions().
    WithTitle("Web Page with Images")

pptx, err := urlfetch.HTMLToPPTXWithOptions(html, url, cfg, opts)
```

### CSS Selector Example
```go
cfg := urlfetch.DefaultConfig().
    WithContentSelectors([]string{
        "article.post-content",
        ".blog-entry",
        "main",
    }).
    WithExcludeSelectors([]string{
        ".advertisement",
        ".social-share",
        "nav",
    })

pptx, err := urlfetch.HTMLToPPTXWithOptions(html, url, cfg, opts)
```

---

## Files to Modify

| File | Changes |
|------|---------|
| `pkg/pptx/urlfetch/config.go` | Add new config options + builders |
| `pkg/pptx/urlfetch/config.go` | Add validation method |
| `pkg/pptx/urlfetch/image_fetcher.go` | **New file** - HTTP image fetching |
| `pkg/pptx/urlfetch/converter.go` | Add image embedding logic |
| `pkg/pptx/urlfetch/parser.go` | Add URL resolution + selector customization |
| `pkg/pptx/urlfetch/urlfetch_test.go` | Add comprehensive tests |
| `examples/34-urlfetch/main.go` | Update with new feature demos |
| `examples/README.md` | Add example 34 documentation |

---

## Dependencies

- Existing: `github.com/PuerkitoBio/goquery` (HTML parsing)
- Existing: `net/http` (HTTP client)
- Existing: `github.com/djinn-soul/gopptx/pkg/pptx/shapes` (Image type)
- New consideration: Image dimension detection (use `image.DecodeConfig` from stdlib)

---

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| Large images causing memory issues | Implement size limits, streaming downloads with LimitReader |
| Malicious image URLs | Validate MIME types, restrict to image formats only |
| Broken image links | Graceful fallback to alt-text |
| Invalid CSS selectors | Config validation, default fallback |
| Concurrent download overhead | Use bounded goroutine pool |

---

## Success Criteria

1. Images from web pages are downloaded and embedded as actual PPTX images
2. Alt-text fallback works when images fail to download
3. Size limits are enforced per-image and total
4. Custom CSS selectors successfully extract content from non-standard pages
5. Exclude selectors remove unwanted elements
6. All tests pass, including new integration tests
7. Example 34 demonstrates both features
