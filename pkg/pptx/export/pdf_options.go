package export

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

// PDFDriver identifies which PDF backend to use.
type PDFDriver string

const (
	// PDFDriverAuto prefers external converters first and uses native as fallback.
	PDFDriverAuto PDFDriver = "auto"
	// PDFDriverNative uses the built-in renderer (experimental visual fidelity).
	PDFDriverNative      PDFDriver = "native"
	PDFDriverLibreOffice PDFDriver = "libreoffice"
	PDFDriverPowerPoint  PDFDriver = "powerpoint"
)

// PDFOptions configures PDF export behavior.
//
// WARNING:
//   - PDFDriverNative is experimental and may not match PowerPoint rendering
//     fidelity for all decks/layouts.
//   - PDFDriverAuto prefers LibreOffice/PowerPoint first and only falls back to
//     native when those drivers are unavailable or fail.
type PDFOptions struct {
	Driver          PDFDriver
	NativeFontPaths []string

	// Timeout bounds how long an external converter (LibreOffice, PowerPoint)
	// may run before it is killed and the export fails. Zero means
	// DefaultConverterTimeout; a negative value disables the timeout.
	Timeout time.Duration

	// SlideSize sets the page geometry for the native renderer when exporting
	// in-memory slides, which carry no size of their own. Zero means 4:3.
	// Exports that read a PPTX take the size from the file and ignore this.
	SlideSize common.SlideSize
}

// DefaultConverterTimeout bounds external converter runs. LibreOffice and the
// PowerPoint COM automation can both block indefinitely (profile locks, modal
// repair dialogs), which would otherwise hang the caller forever.
const DefaultConverterTimeout = 3 * time.Minute

func defaultPDFOptions() PDFOptions {
	return PDFOptions{Driver: PDFDriverAuto}
}

// converterContext returns a context bounded by the configured timeout, plus
// its cancel func. The caller must always call cancel.
func converterContext(opts PDFOptions) (context.Context, context.CancelFunc) {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = DefaultConverterTimeout
	}
	if timeout < 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), timeout)
}

// ParsePDFDriver validates and normalizes a driver name.
func ParsePDFDriver(value string) (PDFDriver, error) {
	driver := PDFDriver(strings.ToLower(strings.TrimSpace(value)))
	if driver == "" {
		return PDFDriverAuto, nil
	}

	switch driver {
	case PDFDriverAuto, PDFDriverNative, PDFDriverLibreOffice, PDFDriverPowerPoint:
		return driver, nil
	default:
		return "", fmt.Errorf("invalid PDF driver %q (allowed: auto|native|libreoffice|powerpoint)", value)
	}
}

func normalizePDFDriver(opts PDFOptions) (PDFDriver, error) {
	driver := string(opts.Driver)
	if strings.TrimSpace(driver) == "" {
		return PDFDriverAuto, nil
	}
	return ParsePDFDriver(driver)
}
