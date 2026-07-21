package export

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// converterWaitDelay bounds how long a converter's output pipes are drained
// after the process is killed on timeout.
const converterWaitDelay = 5 * time.Second

func converterTimeout(opts PDFOptions) time.Duration {
	if opts.Timeout == 0 {
		return DefaultConverterTimeout
	}
	return opts.Timeout
}

const (
	sofficeName    = "soffice"
	sofficeMacPath = "/Applications/LibreOffice.app/Contents/MacOS/soffice"
)

// sofficeExecutable locates the LibreOffice binary, preferring the macOS
// bundle path when the command is not on PATH. The result is always one of the
// two constants above, never a caller-supplied string.
func sofficeExecutable() (string, error) {
	if runtime.GOOS == osDarwin {
		if _, err := os.Stat(sofficeMacPath); err == nil {
			return sofficeMacPath, nil
		}
	}
	if _, err := exec.LookPath(sofficeName); err != nil {
		return "", errors.New("LibreOffice ('soffice') not found in PATH")
	}
	return sofficeName, nil
}

// convertViaLibreOffice runs soffice on pptxPath and returns the path of the
// PDF it produced inside workDir. Conversion always targets a private
// directory so it can never overwrite files the caller owns.
func convertViaLibreOffice(pptxPath, workDir string, opts PDFOptions) (string, error) {
	sofficeCmd, err := sofficeExecutable()
	if err != nil {
		return "", err
	}

	ctx, cancel := converterContext(opts)
	defer cancel()

	// Use literal strings in exec.CommandContext so static analysis can verify
	// no dynamic/user-controlled executable reaches this call site.
	args := []string{"--headless", "--convert-to", "pdf", pptxPath, "--outdir", workDir}

	// The literals are spelled out here rather than referenced through the
	// constants: the rule requires the argument itself to be a string literal.
	var cmd *exec.Cmd
	switch sofficeCmd {
	case sofficeMacPath:
		cmd = exec.CommandContext(ctx, "/Applications/LibreOffice.app/Contents/MacOS/soffice", args...)
	default: // "soffice"
		cmd = exec.CommandContext(ctx, "soffice", args...)
	}
	// Without WaitDelay, killing soffice on timeout still blocks here until every
	// grandchild that inherited the output pipe exits.
	cmd.WaitDelay = converterWaitDelay
	output, err := cmd.CombinedOutput()
	if ctxErr := ctx.Err(); errors.Is(ctxErr, context.DeadlineExceeded) {
		return "", fmt.Errorf("LibreOffice conversion timed out after %s", converterTimeout(opts))
	}
	if err != nil {
		return "", fmt.Errorf("LibreOffice conversion failed: %w\nOutput: %s", err, string(output))
	}

	baseName := filepath.Base(pptxPath)
	generatedPDF := filepath.Join(workDir, strings.TrimSuffix(baseName, filepath.Ext(baseName))+".pdf")

	// soffice exits 0 without writing anything when another instance holds the
	// user-profile lock, so success is not evidence that a PDF exists.
	if err := verifyPDFProduced(generatedPDF); err != nil {
		return "", fmt.Errorf("LibreOffice reported success but %w\nOutput: %s", err, string(output))
	}
	return generatedPDF, nil
}

// verifyPDFProduced confirms a converter actually wrote a usable PDF.
func verifyPDFProduced(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("no PDF was produced at %s", path)
	}
	if info.Size() == 0 {
		return fmt.Errorf("the PDF produced at %s is empty", path)
	}
	return nil
}

// pdfFromFileViaLibreOffice converts an existing PPTX file via LibreOffice.
func pdfFromFileViaLibreOffice(pptxPath, pdfPath string, opts PDFOptions) error {
	workDir, err := os.MkdirTemp("", "gopptx-libreoffice-")
	if err != nil {
		return fmt.Errorf("failed to create temp work directory: %w", err)
	}
	defer os.RemoveAll(workDir)

	generatedPDF, err := convertViaLibreOffice(pptxPath, workDir, opts)
	if err != nil {
		return err
	}
	return moveGeneratedPDF(generatedPDF, pdfPath)
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

// writeTempDeck renders slides to a PPTX inside a private temp directory and
// returns the directory and file path. The caller owns the directory and must
// remove it. Keeping the intermediate out of the caller's output directory
// avoids clobbering their files and makes concurrent exports safe.
func writeTempDeck(title string, slides []elements.SlideContent, prefix string) (string, string, error) {
	workDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp work directory: %w", err)
	}

	pptxBytes, err := pptx.CreateWithSlides(title, slides)
	if err != nil {
		_ = os.RemoveAll(workDir)
		return "", "", fmt.Errorf("failed to generate PPTX: %w", err)
	}

	tmpFile := filepath.Join(workDir, sanitizeTitle(title)+".pptx")
	if err := os.WriteFile(tmpFile, pptxBytes, 0o600); err != nil {
		_ = os.RemoveAll(workDir)
		return "", "", fmt.Errorf("failed to write temp PPTX: %w", err)
	}
	return workDir, tmpFile, nil
}

// pdfViaLibreOffice attempts to use a local LibreOffice installation to perform the conversion.
func pdfViaLibreOffice(title string, slides []elements.SlideContent, outputPath string, opts PDFOptions) error {
	if _, err := sofficeExecutable(); err != nil {
		return err
	}

	workDir, tmpFile, err := writeTempDeck(title, slides, "gopptx-libreoffice-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	generatedPDF, err := convertViaLibreOffice(tmpFile, workDir, opts)
	if err != nil {
		return err
	}
	return moveGeneratedPDF(generatedPDF, outputPath)
}

func moveGeneratedPDF(generatedPDF, outputPath string) error {
	if err := verifyPDFProduced(generatedPDF); err != nil {
		return err
	}
	// Collapse any ".." segments before the path reaches the filesystem, so the
	// destination is exactly the one the caller named.
	dest := filepath.Clean(outputPath)
	if err := os.MkdirAll(filepath.Dir(dest), 0o750); err != nil {
		return fmt.Errorf("failed to create PDF output directory: %w", err)
	}
	// Rename fails across filesystems, which is the norm now that conversion
	// runs in the system temp directory, so fall back to a copy.
	if err := os.Rename(generatedPDF, dest); err == nil {
		return nil
	}
	if err := copyFile(generatedPDF, dest); err != nil {
		return err
	}
	if err := os.Remove(generatedPDF); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// copyFile streams src to dst. Streaming keeps a large deck off the heap, which
// os.ReadFile would not.
func copyFile(src, dst string) error {
	in, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(filepath.Clean(dst), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
