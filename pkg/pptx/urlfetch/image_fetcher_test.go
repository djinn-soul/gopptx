package urlfetch

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageFetcher_FetchImage(t *testing.T) {
	// Create a mock server to serve an image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test.png":
			w.Header().Set("Content-Type", "image/png")
			// Minimal PNG 1x1
			pngData := []byte{
				0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
				0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
				0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
				0x42, 0x60, 0x82,
			}
			w.Write(pngData)
		case "/large.png":
			w.Header().Set("Content-Type", "image/png")
			w.Write(make([]byte, 2000)) // Too large for our test config
		case "/invalid.txt":
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("not an image"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := DefaultConfig().
		WithMaxImageSizeBytes(1000).
		WithMaxTotalImageSizeBytes(5000)

	fetcher := NewImageFetcher(http.DefaultClient, cfg, server.URL)

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid relative", "/test.png", false},
		{"Valid absolute", server.URL + "/test.png", false},
		{"Too large", "/large.png", true},
		{"Invalid type", "/invalid.txt", true},
		{"Server error", "/error", true},
		{"Not found", "/notfound", true},
		{"Data URI (unsupported)", "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := fetcher.FetchImage(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchImage(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if img == nil {
					t.Fatal("expected image, got nil")
				}
				if img.Format != "png" {
					t.Errorf("expected format png, got %s", img.Format)
				}
				if img.Width != 1 || img.Height != 1 {
					t.Errorf("expected 1x1, got %dx%d", img.Width, img.Height)
				}
			}
		})
	}
}

func TestImageFetcher_TotalSizeLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		// Minimal PNG
		pngData := []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
			0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
			0x42, 0x60, 0x82,
		}
		w.Write(pngData)
	}))
	defer server.Close()

	// PNG data is about 67 bytes. Let's set limit to 100.
	cfg := DefaultConfig().WithMaxTotalImageSizeBytes(100)
	fetcher := NewImageFetcher(http.DefaultClient, cfg, server.URL)

	// First fetch should succeed
	_, err := fetcher.FetchImage("/img1.png")
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}

	// Second fetch should fail total limit
	_, err = fetcher.FetchImage("/img2.png")
	if err == nil {
		t.Error("expected total size limit error, got nil")
	}
}

func TestCalculateImageDimensions(t *testing.T) {
	tests := []struct {
		w, h         int
		targetW      int64
		wantW, wantH int64
	}{
		{100, 200, 1000, 1000, 2000},
		{400, 300, 1200, 1200, 900},
		{0, 0, 1000, 1000, 750}, // Default 4:3
	}

	for _, tt := range tests {
		gotW, gotH := CalculateImageDimensions(tt.w, tt.h, tt.targetW)
		if gotW != tt.wantW || gotH != tt.wantH {
			t.Errorf("CalculateImageDimensions(%d, %d, %d) = %d, %d; want %d, %d", tt.w, tt.h, tt.targetW, gotW, gotH, tt.wantW, tt.wantH)
		}
	}
}
