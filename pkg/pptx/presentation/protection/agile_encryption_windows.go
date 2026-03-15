//go:build windows

package protection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	powerPointProbeTimeout   = 8 * time.Second
	powerPointEncryptTimeout = 90 * time.Second
)

var (
	agileAvailabilityOnce sync.Once
	agileAvailable        bool
)

func canEncryptAgile() bool {
	agileAvailabilityOnce.Do(func() {
		agileAvailable = runPowerShellWithTimeout(
			"$app = New-Object -ComObject PowerPoint.Application; $app.Quit() | Out-Null; 'OK'",
			powerPointProbeTimeout,
		) == nil
	})
	return agileAvailable
}

func encryptAgilePackage(zipPayload []byte, password string) ([]byte, error) {
	if !canEncryptAgile() {
		return nil, errorsAgileUnavailable()
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-agile-*")
	if err != nil {
		return nil, fmt.Errorf("create agile temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inPath := filepath.Join(tmpDir, "input.pptx")
	outPath := filepath.Join(tmpDir, "encrypted.pptx")
	if err := os.WriteFile(inPath, zipPayload, 0o600); err != nil {
		return nil, fmt.Errorf("write agile input payload: %w", err)
	}

	script := buildPowerPointEncryptScript(inPath, outPath, password)
	output, err := runPowerShellScript(script, powerPointEncryptTimeout)
	if err != nil {
		return nil, fmt.Errorf("powerpoint agile encryption failed: %w; output: %s", err, strings.TrimSpace(string(output)))
	}

	encrypted, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read agile encrypted output: %w", err)
	}
	return encrypted, nil
}

func buildPowerPointEncryptScript(inPath, outPath, password string) string {
	inQ := quotePowerShellLiteral(inPath)
	outQ := quotePowerShellLiteral(outPath)
	pwQ := quotePowerShellLiteral(password)
	return "$ErrorActionPreference='Stop';" +
		"$in=" + inQ + ";" +
		"$out=" + outQ + ";" +
		"$pw=" + pwQ + ";" +
		"$app=New-Object -ComObject PowerPoint.Application;" +
		"try {" +
		"$pres=$app.Presentations.Open($in,0,0,0);" +
		"$pres.Password=$pw;" +
		"$pres.SaveAs($out);" +
		"$pres.Close();" +
		"} finally {" +
		"if($app -ne $null){$app.Quit() | Out-Null}" +
		"}"
}

func quotePowerShellLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func runPowerShellWithTimeout(script string, timeout time.Duration) error {
	_, err := runPowerShellScript(script, timeout)
	return err
}

func runPowerShellScript(script string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return output, fmt.Errorf("powershell command timed out after %s", timeout)
	}
	return output, err
}

func errorsAgileUnavailable() error {
	return errAgileUnavailable
}
