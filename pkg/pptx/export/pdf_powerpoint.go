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

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func pdfViaPowerPointFromSlides(
	title string,
	slides []elements.SlideContent,
	outputPath string,
	opts PDFOptions,
) error {
	if runtime.GOOS != osWindows {
		return errors.New("PowerPoint driver is only available on Windows")
	}

	workDir, tmpFile, err := writeTempDeck(title, slides, "gopptx-powerpoint-")
	if err != nil {
		return fmt.Errorf("PowerPoint driver: %w", err)
	}
	defer os.RemoveAll(workDir)

	absPPTX, err := filepath.Abs(tmpFile)
	if err != nil {
		return fmt.Errorf("PowerPoint driver (PPTX path): %w", err)
	}
	absPDF, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("PowerPoint driver (PDF path): %w", err)
	}
	return exportWithPowerPoint(absPPTX, absPDF, opts)
}

func exportWithPowerPoint(pptxPath, pdfPath string, opts PDFOptions) error {
	if err := os.MkdirAll(filepath.Dir(pdfPath), 0o750); err != nil {
		return fmt.Errorf("failed to create PDF output directory: %w", err)
	}

	// Use a temporary script file so path arguments are bound predictably.
	psScript := `
param(
  [Parameter(Mandatory = $true)]
  [string]$pptxPath,
  [Parameter(Mandatory = $true)]
  [string]$pdfPath
)
$ppt = New-Object -ComObject PowerPoint.Application
$ppt.Visible = 1
try {
  $pres = $ppt.Presentations.Open($pptxPath, $false, $true, $false)
  $pres.SaveAs($pdfPath, 32)
  $pres.Close()
} catch {
  Write-Error $_
  exit 1
} finally {
  try { $ppt.Quit() } catch {}
  try { [System.Runtime.Interopservices.Marshal]::ReleaseComObject($ppt) | Out-Null } catch {}
}
`
	scriptFile, err := os.CreateTemp("", "gopptx-export-ppt-*.ps1")
	if err != nil {
		return fmt.Errorf("failed to create PowerShell temp script: %w", err)
	}
	scriptPath := scriptFile.Name()
	if _, err := scriptFile.WriteString(psScript); err != nil {
		_ = scriptFile.Close()
		_ = os.Remove(scriptPath)
		return fmt.Errorf("failed to write PowerShell temp script: %w", err)
	}
	if err := scriptFile.Close(); err != nil {
		_ = os.Remove(scriptPath)
		return fmt.Errorf("failed to close PowerShell temp script: %w", err)
	}
	defer os.Remove(scriptPath)

	psExe, err := findPowerShellExecutable()
	if err != nil {
		return err
	}

	// Use literal strings in exec.CommandContext so static analysis can verify
	// no dynamic/user-controlled executable reaches this call site.
	args := []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-File", scriptPath, "-pptxPath", pptxPath, "-pdfPath", pdfPath}

	ctx, cancel := converterContext(opts)
	defer cancel()

	var cmd *exec.Cmd
	switch psExe {
	case "pwsh":
		cmd = exec.CommandContext(ctx, "pwsh", args...)
	default: // "powershell"
		cmd = exec.CommandContext(ctx, "powershell", args...)
	}
	cmd.WaitDelay = converterWaitDelay
	output, err := cmd.CombinedOutput()
	if ctxErr := ctx.Err(); errors.Is(ctxErr, context.DeadlineExceeded) {
		return fmt.Errorf("PowerPoint export timed out after %s", converterTimeout(opts))
	}
	if err != nil {
		return fmt.Errorf("PowerShell execution failed: %w\nOutput: %s", err, normalizePowerShellOutput(string(output)))
	}
	if err := verifyPDFProduced(pdfPath); err != nil {
		return fmt.Errorf("PowerPoint export completed but %w", err)
	}
	return nil
}

func findPowerShellExecutable() (string, error) {
	for _, candidate := range []string{"powershell", "pwsh"} {
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", errors.New(
		"PowerPoint driver requires PowerShell ('powershell' or 'pwsh') in PATH",
	)
}

func normalizePowerShellOutput(out string) string {
	return strings.TrimSpace(out)
}
