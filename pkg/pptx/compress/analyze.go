package compress

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// PartSize is one part's uncompressed size.
type PartSize struct {
	Name  string `json:"name"`
	Bytes int64  `json:"bytes"`
}

// Analysis reports where a package's bytes live, so a caller can tell whether
// stripping notes is worth anything or whether the file is all media.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Analysis struct {
	TotalBytes  int64      `json:"totalBytes"`
	MediaBytes  int64      `json:"mediaBytes"`
	SlideBytes  int64      `json:"slideBytes"`
	ChartBytes  int64      `json:"chartBytes"`
	ThemeBytes  int64      `json:"themeBytes"`
	NotesBytes  int64      `json:"notesBytes"`
	OtherBytes  int64      `json:"otherBytes"`
	PartCount   int        `json:"partCount"`
	ImageCount  int        `json:"imageCount"`
	SlideCount  int        `json:"slideCount"`
	LargestPart []PartSize `json:"largestPart"`
}

const largestPartsReported = 10

// AnalyzeFile reads a PPTX from disk and reports its size breakdown.
func AnalyzeFile(path string) (Analysis, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Analysis{}, fmt.Errorf("read %s: %w", path, err)
	}
	return AnalyzeBytes(data)
}

// AnalyzeBytes reports the size breakdown of an in-memory PPTX package. Sizes
// are uncompressed part sizes; they sum to more than the file on disk.
func AnalyzeBytes(data []byte) (Analysis, error) {
	parts, err := readParts(data)
	if err != nil {
		return Analysis{}, err
	}

	analysis := Analysis{PartCount: len(parts)}
	sizes := make([]PartSize, 0, len(parts))
	for _, p := range parts {
		size := int64(len(p.Data))
		name := normalizeName(p.Name)
		analysis.TotalBytes += size
		sizes = append(sizes, PartSize{Name: name, Bytes: size})

		switch {
		case strings.HasPrefix(name, "ppt/media/"):
			analysis.MediaBytes += size
			analysis.ImageCount++
		case strings.HasPrefix(name, "ppt/slides/") && strings.HasSuffix(name, ".xml"):
			analysis.SlideBytes += size
			analysis.SlideCount++
		case strings.HasPrefix(name, "ppt/charts/") || strings.HasPrefix(name, "ppt/embeddings/"):
			analysis.ChartBytes += size
		case strings.HasPrefix(name, "ppt/theme/"):
			analysis.ThemeBytes += size
		case strings.HasPrefix(name, "ppt/notesSlides/"):
			analysis.NotesBytes += size
		default:
			analysis.OtherBytes += size
		}
	}

	sort.Slice(sizes, func(i, j int) bool { return sizes[i].Bytes > sizes[j].Bytes })
	if len(sizes) > largestPartsReported {
		sizes = sizes[:largestPartsReported]
	}
	analysis.LargestPart = sizes
	return analysis, nil
}

// MediaPercentage returns the share of the package taken by media parts.
func (a Analysis) MediaPercentage() float64 {
	return a.percentage(a.MediaBytes)
}

// ChartPercentage returns the share of the package taken by charts and their
// embedded workbooks.
func (a Analysis) ChartPercentage() float64 {
	return a.percentage(a.ChartBytes)
}

func (a Analysis) percentage(value int64) float64 {
	if a.TotalBytes == 0 {
		return 0
	}
	const toPercent = 100
	return float64(value) / float64(a.TotalBytes) * toPercent
}

// Summary renders a short human-readable breakdown.
func (a Analysis) Summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "parts: %d, slides: %d, images: %d\n", a.PartCount, a.SlideCount, a.ImageCount)
	fmt.Fprintf(&b, "total:  %s\n", humanBytes(a.TotalBytes))
	fmt.Fprintf(&b, "media:  %s (%.1f%%)\n", humanBytes(a.MediaBytes), a.MediaPercentage())
	fmt.Fprintf(&b, "slides: %s\n", humanBytes(a.SlideBytes))
	fmt.Fprintf(&b, "charts: %s (%.1f%%)\n", humanBytes(a.ChartBytes), a.ChartPercentage())
	fmt.Fprintf(&b, "themes: %s\n", humanBytes(a.ThemeBytes))
	fmt.Fprintf(&b, "notes:  %s\n", humanBytes(a.NotesBytes))
	fmt.Fprintf(&b, "other:  %s\n", humanBytes(a.OtherBytes))
	if len(a.LargestPart) > 0 {
		b.WriteString("largest parts:\n")
		for _, p := range a.LargestPart {
			fmt.Fprintf(&b, "  %-48s %s\n", p.Name, humanBytes(p.Bytes))
		}
	}
	return b.String()
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return fmt.Sprintf("%.1f %s", float64(n)/float64(div), units[exp])
}
