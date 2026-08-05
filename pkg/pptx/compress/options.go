// Package compress shrinks an existing PPTX package: it recompresses raster
// images, drops media no relationship points at, strips optional parts
// (notes, comments, custom properties, thumbnails), and minifies the XML.
//
// It also reports where the bytes went, so a caller can decide what to strip,
// see Analyze.
package compress

// Level is a preset that decides how aggressively images are recompressed.
type Level int

const (
	// LevelLight keeps near-original image quality.
	LevelLight Level = iota
	// LevelBalanced is the default trade-off between size and quality.
	LevelBalanced
	// LevelMaximum favours the smallest file over image fidelity.
	LevelMaximum
)

const (
	lightQuality    = 90
	balancedQuality = 78
	maximumQuality  = 60

	balancedMaxDimension = 1920
	maximumMaxDimension  = 1280

	minRetryQuality  = 30
	retryQualityStep = 10
)

// ImageQuality returns the JPEG quality (1-100) used at this level.
func (l Level) ImageQuality() int {
	switch l {
	case LevelLight:
		return lightQuality
	case LevelMaximum:
		return maximumQuality
	case LevelBalanced:
		return balancedQuality
	default:
		return balancedQuality
	}
}

// MaxImageDimension returns the longest edge images are scaled down to, or 0
// when this level leaves image dimensions alone.
func (l Level) MaxImageDimension() int {
	switch l {
	case LevelLight:
		return 0
	case LevelBalanced:
		return balancedMaxDimension
	case LevelMaximum:
		return maximumMaxDimension
	default:
		return balancedMaxDimension
	}
}

// ResizesImages reports whether this level scales oversized images down.
func (l Level) ResizesImages() bool {
	return l.MaxImageDimension() > 0
}

// Options controls one compression run.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Options struct {
	Level Level
	// RemoveUnusedMedia drops ppt/media parts no relationship targets.
	RemoveUnusedMedia bool
	// RemoveProperties drops docProps/custom.xml and the package thumbnail.
	// The required core and app properties are always kept.
	RemoveProperties bool
	// RemoveNotes drops every notes slide (the notes master is kept).
	RemoveNotes bool
	// RemoveComments drops comment parts and the comment author list.
	RemoveComments bool
	// OptimizeXML strips insignificant whitespace between elements.
	OptimizeXML bool
	// TargetSizeBytes, when > 0, lowers image quality step by step until the
	// package fits, or until quality would drop below a usable floor.
	TargetSizeBytes int64
}

// DefaultOptions returns a safe configuration: balanced image quality, unused
// media dropped, everything else preserved.
func DefaultOptions() Options {
	return Options{
		Level:             LevelBalanced,
		RemoveUnusedMedia: true,
		OptimizeXML:       true,
	}
}

// MaximumOptions strips every optional part and uses the smallest images.
func MaximumOptions() Options {
	return Options{
		Level:             LevelMaximum,
		RemoveUnusedMedia: true,
		RemoveProperties:  true,
		RemoveNotes:       true,
		RemoveComments:    true,
		OptimizeXML:       true,
	}
}

// WebOptions targets fast download: small images, notes and comments kept.
func WebOptions() Options {
	return Options{
		Level:             LevelMaximum,
		RemoveUnusedMedia: true,
		RemoveProperties:  true,
		OptimizeXML:       true,
	}
}

// WithLevel sets the compression level.
func (o Options) WithLevel(level Level) Options {
	o.Level = level
	return o
}

// WithUnusedMediaRemoval toggles dropping unreferenced media.
func (o Options) WithUnusedMediaRemoval(remove bool) Options {
	o.RemoveUnusedMedia = remove
	return o
}

// WithPropertiesRemoval toggles dropping custom properties and thumbnails.
func (o Options) WithPropertiesRemoval(remove bool) Options {
	o.RemoveProperties = remove
	return o
}

// WithNotesRemoval toggles dropping notes slides.
func (o Options) WithNotesRemoval(remove bool) Options {
	o.RemoveNotes = remove
	return o
}

// WithCommentsRemoval toggles dropping comments.
func (o Options) WithCommentsRemoval(remove bool) Options {
	o.RemoveComments = remove
	return o
}

// WithXMLOptimization toggles XML whitespace minification.
func (o Options) WithXMLOptimization(optimize bool) Options {
	o.OptimizeXML = optimize
	return o
}

// WithTargetSize sets a best-effort size ceiling in bytes.
func (o Options) WithTargetSize(bytes int64) Options {
	o.TargetSizeBytes = bytes
	return o
}

// Result reports what one compression run changed.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Result struct {
	OriginalBytes      int64
	CompressedBytes    int64
	RemovedParts       []string
	RecompressedImages int
	ResizedImages      int
	FinalImageQuality  int
}

// SavedBytes returns how many bytes the run removed. It can be negative when a
// package was already smaller than what re-encoding produces.
func (r Result) SavedBytes() int64 {
	return r.OriginalBytes - r.CompressedBytes
}

// Ratio returns compressed size divided by original size. It returns 0 for an
// empty input.
func (r Result) Ratio() float64 {
	if r.OriginalBytes == 0 {
		return 0
	}
	return float64(r.CompressedBytes) / float64(r.OriginalBytes)
}
