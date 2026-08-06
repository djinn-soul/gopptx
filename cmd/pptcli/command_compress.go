package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/compress"
)

type compressCommandOptions struct {
	inPath     string
	outPath    string
	level      compress.Level
	targetSize int64
	analyze    bool
	jsonMode   bool
	options    compress.Options
}

// runCompressCommand shrinks a PPTX package, or reports its size breakdown
// when -analyze is passed.
func runCompressCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	opts, exitCode, ok := parseCompressCommandOptions(args, stderr)
	if !ok {
		return exitCode
	}

	if opts.analyze {
		return runCompressAnalyze(opts, stdout, stderr)
	}

	result, err := compress.File(opts.inPath, opts.outPath, opts.options)
	if err != nil {
		printErrorf(stderr, "compress failed: %v", err)
		return exitIO
	}

	if opts.jsonMode {
		payload := map[string]any{
			"input":               opts.inPath,
			"output":              opts.outPath,
			"original_bytes":      result.OriginalBytes,
			"compressed_bytes":    result.CompressedBytes,
			"saved_bytes":         result.SavedBytes(),
			"ratio":               result.Ratio(),
			"removed_parts":       result.RemovedParts,
			"recompressed_images": result.RecompressedImages,
			"resized_images":      result.ResizedImages,
			"final_image_quality": result.FinalImageQuality,
		}
		out, marshalErr := json.MarshalIndent(payload, "", "  ")
		if marshalErr != nil {
			printErrorf(stderr, "failed to marshal JSON: %v", marshalErr)
			return exitIO
		}
		_, _ = fmt.Fprintln(stdout, string(out))
		return exitOK
	}

	_, _ = fmt.Fprintf(stdout, "Compressed %s -> %s\n", opts.inPath, opts.outPath)
	_, _ = fmt.Fprintf(
		stdout,
		"  %d -> %d bytes (%.1f%% of original, saved %d)\n",
		result.OriginalBytes,
		result.CompressedBytes,
		result.Ratio()*100, //nolint:mnd // percentage
		result.SavedBytes(),
	)
	_, _ = fmt.Fprintf(
		stdout,
		"  images recompressed: %d, resized: %d, quality: %d\n",
		result.RecompressedImages,
		result.ResizedImages,
		result.FinalImageQuality,
	)
	if len(result.RemovedParts) > 0 {
		_, _ = fmt.Fprintf(stdout, "  removed %d parts\n", len(result.RemovedParts))
	}
	return exitOK
}

func runCompressAnalyze(opts compressCommandOptions, stdout io.Writer, stderr io.Writer) int {
	analysis, err := compress.AnalyzeFile(opts.inPath)
	if err != nil {
		printErrorf(stderr, "analyze failed: %v", err)
		return exitIO
	}
	if opts.jsonMode {
		out, marshalErr := json.MarshalIndent(analysis, "", "  ")
		if marshalErr != nil {
			printErrorf(stderr, "failed to marshal JSON: %v", marshalErr)
			return exitIO
		}
		_, _ = fmt.Fprintln(stdout, string(out))
		return exitOK
	}
	_, _ = fmt.Fprint(stdout, analysis.Summary())
	return exitOK
}

func parseCompressCommandOptions(args []string, stderr io.Writer) (compressCommandOptions, int, bool) {
	opts := compressCommandOptions{}

	fs := flag.NewFlagSet("compress", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		levelName      string
		format         string
		removeMedia    bool
		removeProps    bool
		removeNotes    bool
		removeComments bool
		optimizeXML    bool
	)
	fs.StringVar(&opts.inPath, "in", "", "Input PPTX file path")
	fs.StringVar(&opts.outPath, "out", "", "Output PPTX file path (defaults to overwriting the input)")
	fs.StringVar(&levelName, "level", "balanced", "Compression level: light|balanced|maximum")
	fs.Int64Var(&opts.targetSize, "target-size", 0, "Best-effort maximum output size in bytes")
	fs.BoolVar(&opts.analyze, "analyze", false, "Report the size breakdown instead of compressing")
	fs.StringVar(&format, "format", "text", "Output format: text or json")
	fs.BoolVar(&removeMedia, "remove-unused-media", true, "Drop media parts no relationship points at")
	fs.BoolVar(&removeProps, "remove-properties", false, "Drop custom properties and the package thumbnail")
	fs.BoolVar(&removeNotes, "remove-notes", false, "Drop notes slides")
	fs.BoolVar(&removeComments, "remove-comments", false, "Drop comments and comment authors")
	fs.BoolVar(&optimizeXML, "optimize-xml", true, "Strip insignificant whitespace from XML parts")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid compress arguments: %v", err)
		printCompressUsage(stderr)
		return compressCommandOptions{}, exitUsage, false
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printCompressUsage(stderr)
		return compressCommandOptions{}, exitUsage, false
	}

	opts.inPath = strings.TrimSpace(opts.inPath)
	if opts.inPath == "" {
		printErrorf(stderr, "compress requires -in")
		printCompressUsage(stderr)
		return compressCommandOptions{}, exitUsage, false
	}
	if opts.outPath == "" {
		opts.outPath = opts.inPath
	}

	level, ok := parseCompressLevel(levelName)
	if !ok {
		printErrorf(stderr, "unsupported level %q (allowed: light|balanced|maximum)", levelName)
		printCompressUsage(stderr)
		return compressCommandOptions{}, exitUsage, false
	}
	opts.level = level
	opts.jsonMode = strings.EqualFold(strings.TrimSpace(format), "json")
	opts.options = compress.Options{
		Level:             level,
		RemoveUnusedMedia: removeMedia,
		RemoveProperties:  removeProps,
		RemoveNotes:       removeNotes,
		RemoveComments:    removeComments,
		OptimizeXML:       optimizeXML,
		TargetSizeBytes:   opts.targetSize,
	}
	return opts, exitOK, true
}

func parseCompressLevel(name string) (compress.Level, bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "light":
		return compress.LevelLight, true
	case "balanced", "":
		return compress.LevelBalanced, true
	case "maximum", "max":
		return compress.LevelMaximum, true
	default:
		return compress.LevelBalanced, false
	}
}

func printCompressUsage(w io.Writer) {
	_, _ = fmt.Fprintln(
		w,
		"Usage: pptcli compress -in file.pptx [-out small.pptx] [-level light|balanced|maximum] "+
			"[-target-size BYTES] [-remove-notes] [-remove-comments] [-remove-properties] "+
			"[-analyze] [-format text|json]",
	)
}
