package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"
)

func runURLFetchCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("urlfetch", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var sourceURL string
	var outPath string
	var title string
	var author string
	var includeSourceURL bool
	fs.StringVar(&sourceURL, "url", "", "source URL to convert into a PPTX deck")
	fs.StringVar(&outPath, "out", "", "output PPTX file (default: <host>.pptx)")
	fs.StringVar(&title, "title", "", "optional deck title override")
	fs.StringVar(&author, "author", "", "optional author metadata")
	fs.BoolVar(&includeSourceURL, "source-url", true, "include source URL bullet on title slide")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "invalid urlfetch arguments: %v", err)
		printURLFetchUsage(stderr)
		return exitUsage
	}
	if len(fs.Args()) > 0 {
		printErrorf(stderr, "unexpected trailing arguments: %s", strings.Join(fs.Args(), " "))
		printURLFetchUsage(stderr)
		return exitUsage
	}

	sourceURL = strings.TrimSpace(sourceURL)
	if sourceURL == "" {
		printErrorf(stderr, "urlfetch requires -url")
		printURLFetchUsage(stderr)
		return exitUsage
	}

	opts := urlfetch.DefaultConversionOptions().WithSourceURL(includeSourceURL)
	if trimmedTitle := strings.TrimSpace(title); trimmedTitle != "" {
		opts = opts.WithTitle(trimmedTitle)
	}
	if trimmedAuthor := strings.TrimSpace(author); trimmedAuthor != "" {
		opts = opts.WithAuthor(trimmedAuthor)
	}

	data, err := urlfetch.URLToPPTXWithOptions(sourceURL, urlfetch.DefaultConfig(), opts)
	if err != nil {
		printErrorf(stderr, "urlfetch conversion failed: %v", err)
		return exitInternal
	}

	if strings.TrimSpace(outPath) == "" {
		outPath = defaultOutputPathFromURL(sourceURL)
	}
	if err := writeOutputFile(outPath, data); err != nil {
		printErrorf(stderr, "failed to write %q: %v", outPath, err)
		return exitIO
	}

	_, _ = fmt.Fprintf(stdout, "OK: wrote %s from %s\n", strings.TrimSpace(outPath), sourceURL)
	return exitOK
}

func printURLFetchUsage(w io.Writer) {
	_, _ = fmt.Fprintln(
		w,
		"Usage: pptcli urlfetch -url https://example.com "+
			"[-out file.pptx] [-title TITLE] [-author NAME] [-source-url=true|false]",
	)
}

func defaultOutputPathFromURL(rawURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed == nil || strings.TrimSpace(parsed.Hostname()) == "" {
		return "urlfetch.pptx"
	}

	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	replacer := strings.NewReplacer(":", "_", "/", "_", "\\", "_", " ", "_")
	host = replacer.Replace(host)
	if host == "" {
		host = "urlfetch"
	}
	return filepath.Clean(host + ".pptx")
}
