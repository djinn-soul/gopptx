package export

import (
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

//nolint:gocognit,nestif // Image export must branch by source mode/path embedding and preserve explicit error handling.
func renderImageToWriter(w io.Writer, img shapes.Image, opts HTMLOptions) error {
	if img.Path == "" && len(img.Data) == 0 {
		return nil
	}

	if !opts.EmbedImages && img.Path != "" {
		src := filepath.Base(img.Path)
		if strings.TrimSpace(opts.BaseURL) != "" {
			src = strings.TrimRight(opts.BaseURL, "/") + "/" + src
		}
		if _, err := fmt.Fprintf(
			w,
			"<div class=\"image-container\"><img src=\"%s\" alt=\"%s\" /></div>\n",
			html.EscapeString(src),
			html.EscapeString(filepath.Base(img.Path)),
		); err != nil {
			return err
		}
		return nil
	}

	mimeType := "image/png" // default
	ext := strings.ToLower(filepath.Ext(img.Path))
	if ext == "" {
		ext = "." + strings.ToLower(img.Format)
	}
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".svg":
		mimeType = "image/svg+xml"
	}

	if _, err := fmt.Fprintf(w, "<div class=\"image-container\"><img src=\"data:%s;base64,", mimeType); err != nil {
		return err
	}

	// Stream base64 to writer to reduce memory peak
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	if len(img.Data) > 0 {
		if _, err := encoder.Write(img.Data); err != nil {
			_ = encoder.Close()
			return err
		}
	} else if img.Path != "" {
		// Security: Validate path is within a subtree and doesn't escape via parent references
		if !filepath.IsLocal(img.Path) {
			_ = encoder.Close()
			return fmt.Errorf("invalid image path: %s", img.Path)
		}

		// Normalize the path to handle "./foo" and "foo/./bar" patterns
		cleanPath := filepath.Clean(img.Path)
		f, err := os.Open(cleanPath)
		if err == nil {
			_, copyErr := io.Copy(encoder, f)
			closeErr := f.Close()
			if copyErr != nil {
				_ = encoder.Close()
				return copyErr
			}
			if closeErr != nil {
				_ = encoder.Close()
				return closeErr
			}
		}
	}
	if err := encoder.Close(); err != nil {
		return err
	}

	_, err := fmt.Fprintf(w, "\" alt=\"%s\" /></div>\n", html.EscapeString(filepath.Base(img.Path)))
	return err
}
