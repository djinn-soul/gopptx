package export

import (
	"fmt"
	"strings"
)

// PDFDriver identifies which PDF backend to use.
type PDFDriver string

const (
	PDFDriverAuto        PDFDriver = "auto"
	PDFDriverNative      PDFDriver = "native"
	PDFDriverLibreOffice PDFDriver = "libreoffice"
	PDFDriverPowerPoint  PDFDriver = "powerpoint"
)

// PDFOptions configures PDF export behavior.
type PDFOptions struct {
	Driver PDFDriver
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
