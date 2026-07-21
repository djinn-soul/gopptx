package export

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	osWindows = "windows"
	osDarwin  = "darwin"
)

// PDF exports the presentation to a PDF file.
func PDF(title string, slides []elements.SlideContent, outputPath string) error {
	return PDFWithOptions(title, slides, outputPath, defaultPDFOptions())
}

// PDFWithOptions exports the presentation to a PDF file using the requested driver.
func PDFWithOptions(title string, slides []elements.SlideContent, outputPath string, opts PDFOptions) error {
	driver, err := normalizePDFDriver(opts)
	if err != nil {
		return err
	}

	switch driver {
	case PDFDriverAuto:
		return pdfWithAutoDriver(title, slides, outputPath, opts)
	case PDFDriverNative:
		return pdfViaNative(title, slides, outputPath, opts, optionsPageSize(opts))
	case PDFDriverLibreOffice:
		return pdfViaLibreOffice(title, slides, outputPath, opts)
	case PDFDriverPowerPoint:
		return pdfViaPowerPointFromSlides(title, slides, outputPath, opts)
	default:
		return fmt.Errorf("unsupported PDF driver %q", driver)
	}
}

// PDFFromFile converts an existing PPTX file on disk to PDF.
func PDFFromFile(pptxPath, pdfPath string) error {
	return PDFFromFileWithOptions(pptxPath, pdfPath, defaultPDFOptions())
}

// PDFFromFileWithOptions converts an existing PPTX file on disk to PDF using the requested driver.
func PDFFromFileWithOptions(pptxPath, pdfPath string, opts PDFOptions) error {
	driver, err := normalizePDFDriver(opts)
	if err != nil {
		return err
	}

	absPPTX, err := filepath.Abs(pptxPath)
	if err != nil {
		return fmt.Errorf("invalid input path: %w", err)
	}
	absPDF, err := filepath.Abs(pdfPath)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	switch driver {
	case PDFDriverAuto:
		return pdfFromFileWithAutoDriver(absPPTX, absPDF, opts)
	case PDFDriverNative:
		presTitle, slides, slideSize, readErr := slidesFromPPTXWithSize(absPPTX)
		if readErr != nil {
			return fmt.Errorf("native (PPTX reader): %w", readErr)
		}
		return pdfViaNative(presTitle, slides, absPDF, opts, pageSizeFromEMU(slideSize.Width, slideSize.Height))
	case PDFDriverLibreOffice:
		return pdfFromFileViaLibreOffice(absPPTX, absPDF, opts)
	case PDFDriverPowerPoint:
		if runtime.GOOS != osWindows {
			return errors.New("PowerPoint driver is only available on Windows")
		}
		return exportWithPowerPoint(absPPTX, absPDF, opts)
	default:
		return fmt.Errorf("unsupported PDF driver %q", driver)
	}
}

func pdfWithAutoDriver(title string, slides []elements.SlideContent, outputPath string, opts PDFOptions) error {
	// Attempt 1: LibreOffice
	libreErr := pdfViaLibreOffice(title, slides, outputPath, opts)
	if libreErr == nil {
		return nil
	}

	// Attempt 2: PowerPoint COM (Windows only)
	var pptErr error
	if runtime.GOOS == osWindows {
		pptErr = pdfViaPowerPointFromSlides(title, slides, outputPath, opts)
		if pptErr == nil {
			return nil
		}
	}

	// Attempt 3: Native gopdf engine (experimental fallback)
	nativeErr := pdfViaNative(title, slides, outputPath, opts, optionsPageSize(opts))
	if nativeErr == nil {
		return nil
	}

	compositeErr := fmt.Errorf("LibreOffice driver: %w", libreErr)
	if pptErr != nil {
		compositeErr = fmt.Errorf("%w\nPowerPoint driver: %w", compositeErr, pptErr)
	}
	compositeErr = fmt.Errorf("%w\nnative PDF driver (experimental): %w", compositeErr, nativeErr)
	return compositeErr
}

func pdfFromFileWithAutoDriver(absPPTX, absPDF string, opts PDFOptions) error {
	// Attempt 1: LibreOffice
	libreErr := pdfFromFileViaLibreOffice(absPPTX, absPDF, opts)
	if libreErr == nil {
		return nil
	}

	// Attempt 2: PowerPoint COM (Windows only)
	var pptErr error
	if runtime.GOOS == osWindows {
		pptErr = exportWithPowerPoint(absPPTX, absPDF, opts)
		if pptErr == nil {
			return nil
		}
	}

	// Attempt 3: Native gopdf via PPTX reader (experimental fallback)
	presTitle, slides, slideSize, readErr := slidesFromPPTXWithSize(absPPTX)
	var nativeErr error
	if readErr == nil {
		nativeErr = pdfViaNative(presTitle, slides, absPDF, opts, pageSizeFromEMU(slideSize.Width, slideSize.Height))
	}
	if readErr == nil && nativeErr == nil {
		return nil
	}

	// All failed
	compositeErr := fmt.Errorf("LibreOffice: %w", libreErr)
	if pptErr != nil {
		compositeErr = fmt.Errorf("%w\nPowerPoint: %w", compositeErr, pptErr)
	}
	compositeErr = fmt.Errorf("%w\nnative (experimental): %w", compositeErr, nativeErr)
	if readErr != nil {
		compositeErr = fmt.Errorf("%w\nnative (PPTX reader): %w", compositeErr, readErr)
	}
	return compositeErr
}
