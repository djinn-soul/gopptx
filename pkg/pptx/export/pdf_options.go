package export

import (
	"fmt"
	"strings"
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
}

func defaultPDFOptions() PDFOptions {
	return PDFOptions{Driver: PDFDriverAuto}
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
