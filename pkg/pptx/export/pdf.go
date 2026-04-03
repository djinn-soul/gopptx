package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
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
		return pdfViaNative(title, slides, outputPath, opts)
	case PDFDriverLibreOffice:
		return pdfViaLibreOffice(title, slides, outputPath)
	case PDFDriverPowerPoint:
		return pdfViaPowerPointFromSlides(title, slides, outputPath)
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
		presTitle, slides, readErr := SlidesFromPPTX(absPPTX)
		if readErr != nil {
			return fmt.Errorf("native (PPTX reader): %w", readErr)
		}
		return pdfViaNative(presTitle, slides, absPDF, opts)
	case PDFDriverLibreOffice:
		return pdfFromFileViaLibreOffice(absPPTX, absPDF)
	case PDFDriverPowerPoint:
		if runtime.GOOS != osWindows {
			return errors.New("PowerPoint driver is only available on Windows")
		}
		return exportWithPowerPoint(absPPTX, absPDF)
	default:
		return fmt.Errorf("unsupported PDF driver %q", driver)
	}
}

func pdfWithAutoDriver(title string, slides []elements.SlideContent, outputPath string, opts PDFOptions) error {
	// Attempt 1: Native gopdf engine
	nativeErr := pdfViaNative(title, slides, outputPath, opts)
	if nativeErr == nil {
		return nil
	}

	// Attempt 2: LibreOffice
	libreErr := pdfViaLibreOffice(title, slides, outputPath)
	if libreErr == nil {
		return nil
	}

	// Attempt 3: PowerPoint COM (Windows only)
	pptErr := pdfViaPowerPointFromSlides(title, slides, outputPath)
	if runtime.GOOS != osWindows {
		return fmt.Errorf("native PDF driver: %w\nLibreOffice driver: %w", nativeErr, libreErr)
	}
	if pptErr == nil {
		return nil
	}

	compositeErr := fmt.Errorf("native PDF driver: %w", nativeErr)
	compositeErr = fmt.Errorf("%w\nLibreOffice driver: %w", compositeErr, libreErr)
	compositeErr = fmt.Errorf("%w\nPowerPoint driver: %w", compositeErr, pptErr)
	return compositeErr
}

func pdfFromFileWithAutoDriver(absPPTX, absPDF string, opts PDFOptions) error {
	// Attempt 1: Native gopdf via PPTX reader
	presTitle, slides, readErr := SlidesFromPPTX(absPPTX)
	var nativeErr error
	if readErr == nil {
		nativeErr = pdfViaNative(presTitle, slides, absPDF, opts)
	}
	if readErr == nil && nativeErr == nil {
		return nil
	}

	// Attempt 2: LibreOffice
	libreErr := pdfFromFileViaLibreOffice(absPPTX, absPDF)
	if libreErr == nil {
		return nil
	}

	// Attempt 3: PowerPoint COM (Windows only)
	var pptErr error
	if runtime.GOOS == osWindows {
		pptErr = exportWithPowerPoint(absPPTX, absPDF)
		if pptErr == nil {
			return nil
		}
	}

	// All failed
	compositeErr := fmt.Errorf("native: %w", nativeErr)
	if readErr != nil {
		compositeErr = fmt.Errorf("%w\nnative (PPTX reader): %w", compositeErr, readErr)
	}
	compositeErr = fmt.Errorf("%w\nLibreOffice: %w", compositeErr, libreErr)
	if pptErr != nil {
		compositeErr = fmt.Errorf("%w\nPowerPoint: %w", compositeErr, pptErr)
	}
	return compositeErr
}

// pdfFromFileViaLibreOffice converts an existing PPTX file via LibreOffice.
func pdfFromFileViaLibreOffice(pptxPath, pdfPath string) error {
	sofficeCmd := "soffice"
	if runtime.GOOS == osDarwin {
		macPath := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		if _, err := os.Stat(macPath); err == nil {
			sofficeCmd = macPath
		}
	}

	if _, err := exec.LookPath(sofficeCmd); err != nil {
		return errors.New("LibreOffice ('soffice') not found in PATH")
	}

	outputDir := filepath.Dir(pdfPath)
	cmd := exec.CommandContext(
		context.Background(),
		sofficeCmd,
		"--headless",
		"--convert-to",
		"pdf",
		pptxPath,
		"--outdir",
		outputDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreOffice conversion failed: %w\nOutput: %s", err, string(output))
	}

	baseName := filepath.Base(pptxPath)
	pdfName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".pdf"
	generatedPDF := filepath.Join(outputDir, pdfName)

	if generatedPDF != pdfPath {
		return moveGeneratedPDF(generatedPDF, pdfPath)
	}
	return nil
}

func sanitizeTitle(title string) string {
	safeTitle := "presentation"
	if title != "" {
		var safeTitleBuilder strings.Builder
		for _, c := range title {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
				safeTitleBuilder.WriteRune(c)
			} else {
				safeTitleBuilder.WriteByte('_')
			}
		}
		if safeTitleBuilder.Len() > 0 {
			safeTitle = safeTitleBuilder.String()
		}
	}
	return safeTitle
}

// pdfViaLibreOffice attempts to use a local LibreOffice installation to perform the conversion.
func pdfViaLibreOffice(title string, slides []elements.SlideContent, outputPath string) error {
	safeTitle := sanitizeTitle(title)

	tmpDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(tmpDir, 0o750); err != nil {
		tmpDir = os.TempDir()
	}
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("gopptx_%s_temp.pptx", safeTitle))

	pptxBytes, err := pptx.CreateWithSlides(title, slides)
	if err != nil {
		return fmt.Errorf("failed to generate PPTX for LibreOffice: %w", err)
	}

	if err := os.WriteFile(tmpFile, pptxBytes, 0o600); err != nil {
		return fmt.Errorf("failed to write temp PPTX: %w", err)
	}
	defer os.Remove(tmpFile)

	sofficeCmd := "soffice"
	if runtime.GOOS == osDarwin {
		macPath := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		if _, err := os.Stat(macPath); err == nil {
			sofficeCmd = macPath
		}
	}

	if _, err := exec.LookPath(sofficeCmd); err != nil {
		return errors.New("LibreOffice ('soffice') not found in PATH")
	}

	outputDir := filepath.Dir(outputPath)
	cmd := exec.CommandContext(
		context.Background(),
		sofficeCmd,
		"--headless",
		"--convert-to",
		"pdf",
		tmpFile,
		"--outdir",
		outputDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreOffice conversion failed: %w\nOutput: %s", err, string(output))
	}

	baseName := filepath.Base(tmpFile)
	pdfName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".pdf"
	generatedPDF := filepath.Join(outputDir, pdfName)

	if generatedPDF != outputPath {
		if err := moveGeneratedPDF(generatedPDF, outputPath); err != nil {
			return err
		}
	}

	return nil
}

func moveGeneratedPDF(generatedPDF, outputPath string) error {
	if _, err := os.Stat(generatedPDF); err != nil {
		return fmt.Errorf("expected generated PDF not found at %s", generatedPDF)
	}
	if err := os.Rename(generatedPDF, outputPath); err == nil {
		return nil
	}
	input, err := os.ReadFile(generatedPDF)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, input, 0o600); err != nil {
		return err
	}
	if err := os.Remove(generatedPDF); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
