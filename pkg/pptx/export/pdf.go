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

// PDF exports the presentation to a PDF file using LibreOffice.
// Requires 'soffice' (or LibreOffice default path on macOS) to be installed.
func PDF(title string, slides []elements.SlideContent, outputPath string) error {
	// 1. Create temporary PPTX
	// Sanitize title for filename
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

	// Use the output directory for the temp file to avoid "Temp" folder permission/trust issues with Office
	// Office sometimes sandboxes files in AppData\Local\Temp
	tmpDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		tmpDir = os.TempDir() // Fallback
	}
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("gopptx_%s_temp.pptx", safeTitle))

	// Create PPTX bytes
	// We verify title/slides here implicitly via CreateWithSlides
	pptxBytes, err := pptx.CreateWithSlides(title, slides)
	if err != nil {
		return fmt.Errorf("failed to generate PPTX: %w", err)
	}

	if err := os.WriteFile(tmpFile, pptxBytes, 0o666); err != nil {
		return fmt.Errorf("failed to write temp PPTX: %w", err)
	}
	// defer os.Remove(tmpFile)
	defer os.Remove(tmpFile)

	// 3. Find soffice command
	sofficeCmd := "soffice"
	if runtime.GOOS == "darwin" {
		// Check standard macOS path
		macPath := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		if _, err := os.Stat(macPath); err == nil {
			sofficeCmd = macPath
		}
	}

	// Check if soffice is available
	if _, err := exec.LookPath(sofficeCmd); err != nil {
		// Soffice not found.
		if runtime.GOOS == "windows" {
			// Try PowerPoint COM fallback
			// We need absolute paths for COM
			absPPTX, _ := filepath.Abs(tmpFile)
			absPDF, _ := filepath.Abs(outputPath)
			if err := exportWithPowerPoint(absPPTX, absPDF); err != nil {
				return fmt.Errorf("LibreOffice not found and PowerPoint fallback failed: %w", err)
			}
			return nil
		}
		return errors.New("LibreOffice ('soffice') not found in PATH")
	}

	// 4. Run conversion with LibreOffice
	// soffice --headless --convert-to pdf <temp_file> --outdir <output_dir>
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

	// 5. Rename result
	// LibreOffice outputs <filename>.pdf in outdir
	baseName := filepath.Base(tmpFile)
	pdfName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".pdf"
	generatedPDF := filepath.Join(outputDir, pdfName)

	if generatedPDF != outputPath {
		// If output path is different, move it.
		// Note: outputPath might be relative, generatedPDF is absolute if outputDir is absolute?
		// outputDir comes from filepath.Dir(outputPath).
		// So generatedPDF is effectively <dir>/<basename>.pdf.

		// If outputPath is just "output.pdf", outputDir is ".".

		// We should verify existence before rename.
		if _, err := os.Stat(generatedPDF); err != nil {
			return fmt.Errorf("expected generated PDF not found at %s", generatedPDF)
		}

		// Move
		if err := os.Rename(generatedPDF, outputPath); err != nil {
			// Fallback: Copy and delete
			input, err := os.ReadFile(generatedPDF)
			if err != nil {
				return err
			}
			if err := os.WriteFile(outputPath, input, 0o644); err != nil {
				return err
			}
			os.Remove(generatedPDF)
		}
	}

	return nil
}

func exportWithPowerPoint(pptxPath, pdfPath string) error {
	// PowerShell command to automate PowerPoint
	// 32 = ppSaveAsPDF
	// Open(FileName) - simplest form for maximum compatibility
	psScript := fmt.Sprintf(`
$ppt = New-Object -ComObject PowerPoint.Application
# PowerPoint might be finicky about visibility
$ppt.Visible = 1
try {
  $pres = $ppt.Presentations.Open('%s')
  $pres.SaveAs('%s', 32)
  $pres.Close()
} catch {
  Write-Error $_
  exit 1
} finally {
  $ppt.Quit()
  [System.Runtime.Interopservices.Marshal]::ReleaseComObject($ppt) | Out-Null
}
`, pptxPath, pdfPath)

	// Run PowerShell
	cmd := exec.CommandContext(
		context.Background(),
		"powershell",
		"-NoProfile",
		"-NonInteractive",
		"-Command",
		psScript,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PowerShell execution failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}
