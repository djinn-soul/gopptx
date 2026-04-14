package export

import (
	_ "embed"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

//go:embed html_style.css
var htmlStyle string

// HTML exports the presentation slides to a single HTML document using default options.
func HTML(title string, slides []elements.SlideContent) string {
	return HTMLWithOptions(title, slides, DefaultHTMLOptions())
}

// HTMLWithOptions exports the presentation slides to a single HTML document.
func HTMLWithOptions(title string, slides []elements.SlideContent, opts HTMLOptions) string {
	var b strings.Builder
	_ = HTMLToWriter(&b, title, slides, opts)
	return b.String()
}

// HTMLToWriter exports the presentation slides to an [io.Writer].
//
//nolint:funlen,gocognit,nestif // HTML writer stitches optional theme/navigation sections and intentionally keeps output assembly centralized.
func HTMLToWriter(w io.Writer, title string, slides []elements.SlideContent, opts HTMLOptions) error {
	// Header
	if err := writeString(w, "<!DOCTYPE html>\n<html>\n<head>\n"); err != nil {
		return err
	}
	if err := writeString(w, "<meta charset=\"UTF-8\">\n"); err != nil {
		return err
	}
	if err := writeString(
		w,
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n",
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "<title>%s</title>\n", html.EscapeString(title)); err != nil {
		return err
	}

	// CSS
	if err := writeString(w, "<style>\n"); err != nil {
		return err
	}
	if err := writeString(w, htmlStyle); err != nil {
		return err
	}

	if opts.Theme != nil {
		if err := writeString(w, "\n:root {\n"); err != nil {
			return err
		}
		if opts.Theme.TitleColor != "" {
			if _, err := fmt.Fprintf(w, "  --title-color: %s;\n", opts.Theme.TitleColor); err != nil {
				return err
			}
		}
		if opts.Theme.BodyColor != "" {
			if _, err := fmt.Fprintf(w, "  --body-color: %s;\n", opts.Theme.BodyColor); err != nil {
				return err
			}
		}
		if opts.Theme.AccentColor != "" {
			if _, err := fmt.Fprintf(w, "  --accent-color: %s;\n", opts.Theme.AccentColor); err != nil {
				return err
			}
		}
		if opts.Theme.BackgroundColor != "" {
			if _, err := fmt.Fprintf(w, "  --bg-color: %s;\n", opts.Theme.BackgroundColor); err != nil {
				return err
			}
		}
		if opts.Theme.SlideBackground != "" {
			if _, err := fmt.Fprintf(w, "  --slide-bg: %s;\n", opts.Theme.SlideBackground); err != nil {
				return err
			}
		}
		if err := writeString(w, "}\n"); err != nil {
			return err
		}
	}

	if err := writeString(w, "</style>\n"); err != nil {
		return err
	}

	if err := writeString(w, "</head>\n<body>\n"); err != nil {
		return err
	}

	// Title Slide (Presentation Title)
	if err := writeString(w, "<div class=\"slide title-slide\">\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "<h1>%s</h1>\n", html.EscapeString(title)); err != nil {
		return err
	}
	if err := writeString(w, "</div>\n"); err != nil {
		return err
	}

	// Slides
	for i, slide := range slides {
		if err := renderSlideToWriter(w, slide, i+1, opts); err != nil {
			return err
		}
	}

	if opts.IncludeNavigation {
		if err := writeString(w, `<div class="nav-controls">
  <button id="btn-prev" onclick="prevSlide()">Prev</button>
  <span id="slide-indicator">1 / `+strconv.Itoa(len(slides))+`</span>
  <button id="btn-next" onclick="nextSlide()">Next</button>
</div>
<script>
let currentSlide = 0;
const totalSlides = `+strconv.Itoa(len(slides))+`;
const slides = document.querySelectorAll('.slide:not(.title-slide)');
function updateVisibility() {
	slides.forEach((s, idx) => s.style.display = (idx === currentSlide) ? 'block' : 'none');
	document.getElementById('slide-indicator').innerText = (currentSlide + 1) + ' / ' + totalSlides;
}
function nextSlide() { if(currentSlide < totalSlides - 1) { currentSlide++; updateVisibility(); } }
function prevSlide() { if(currentSlide > 0) { currentSlide--; updateVisibility(); } }
updateVisibility();
</script>
`); err != nil {
			return err
		}
	}

	return writeString(w, "</body>\n</html>")
}

//nolint:gocognit // Slide rendering coordinates many optional blocks (bullets/images/shapes/table) for deterministic ordering.
func renderSlideToWriter(w io.Writer, slide elements.SlideContent, index int, opts HTMLOptions) error {
	if _, err := fmt.Fprintf(w, "<div class=\"slide\" id=\"slide-%d\">\n", index); err != nil {
		return err
	}

	// Slide Number
	if opts.IncludeSlideNumbers {
		if _, err := fmt.Fprintf(w, "<div class=\"slide-number\">%d</div>\n", index); err != nil {
			return err
		}
	}

	// Title
	slideTitle := slide.Title
	if slideTitle == "" {
		slideTitle = "Slide " + strconv.Itoa(index)
	}
	if _, err := fmt.Fprintf(w, "<h2>%s</h2>\n", html.EscapeString(slideTitle)); err != nil {
		return err
	}

	// Content Container
	if err := writeString(w, "<div class=\"content\">\n"); err != nil {
		return err
	}

	// Bullets / Content
	if len(slide.Bullets) > 0 {
		if err := writeString(w, "<ul>\n"); err != nil {
			return err
		}
		for _, bullet := range slide.Bullets {
			if _, err := fmt.Fprintf(w, "<li>%s</li>\n", html.EscapeString(bullet)); err != nil {
				return err
			}
		}
		if err := writeString(w, "</ul>\n"); err != nil {
			return err
		}
	}

	// Images
	for _, img := range slide.Images {
		if err := renderImageToWriter(w, img, opts); err != nil {
			return err
		}
	}

	// Vector Shapes & Text
	if len(slide.Shapes) > 0 {
		if err := writeString(w, renderShapesSVG(slide.Shapes)); err != nil {
			return err
		}
	}

	// Table Support
	if slide.Table != nil {
		if err := renderTableToWriter(w, slide.Table); err != nil {
			return err
		}
	}

	if err := writeString(w, "</div>\n"); err != nil { // content
		return err
	}
	return writeString(w, "</div>\n") // slide
}

func renderTableToWriter(w io.Writer, table *tables.Table) error {
	if err := writeString(w, "<table class=\"slide-table\">\n"); err != nil {
		return err
	}
	hasStyled := len(table.StyledRows) > 0

	rowCount := len(table.Rows)
	if hasStyled {
		rowCount = len(table.StyledRows)
	}

	for r := range rowCount {
		if err := writeString(w, "  <tr>\n"); err != nil {
			return err
		}
		var row []string
		if hasStyled {
			row = make([]string, len(table.StyledRows[r]))
			for c, cell := range table.StyledRows[r] {
				row[c] = cell.Text
			}
		} else {
			row = table.Rows[r]
		}
		for _, text := range row {
			if _, err := fmt.Fprintf(w, "    <td>%s</td>\n", html.EscapeString(text)); err != nil {
				return err
			}
		}
		if err := writeString(w, "  </tr>\n"); err != nil {
			return err
		}
	}
	return writeString(w, "</table>\n")
}
