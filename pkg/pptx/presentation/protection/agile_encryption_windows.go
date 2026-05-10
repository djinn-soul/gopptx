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

//nolint:gochecknoglobals // COM availability probe result is cached process-wide via sync.Once for deterministic behavior.
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

	script := buildPowerPointEncryptScript(inPath, outPath)
	// Pass the password via an environment variable so it never appears in the
	// process command line or in script text visible to other processes.
	// Use a byte slice so the buffer can be zeroed after use, minimising the
	// window during which the plaintext credential remains in memory.
	envLine := []byte("GOPPTX_PW=" + password)
	defer clear(envLine)
	output, err := runPowerShellScript(script, powerPointEncryptTimeout, string(envLine))
	if err != nil {
		// The output is safe to include: the password is only in the env var,
		// not in the script string or in any PowerShell output.
		return nil, fmt.Errorf(
			"powerpoint agile encryption failed: %w; output: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	encrypted, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("read agile encrypted output: %w", err)
	}
	return encrypted, nil
}

func buildPowerPointEncryptScript(inPath, outPath string) string {
	inQ := quotePowerShellLiteral(inPath)
	outQ := quotePowerShellLiteral(outPath)
	// Password is read from the GOPPTX_PW environment variable so it never
	// appears in the script string or in the process command line.
	return "$ErrorActionPreference='Stop';" +
		"$in=" + inQ + ";" +
		"$out=" + outQ + ";" +
		"$pw=$env:GOPPTX_PW;" +
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

func runPowerShellScript(script string, timeout time.Duration, extraEnv ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}
	output, err := cmd.CombinedOutput()
	if err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return output, fmt.Errorf("powershell command timed out after %s", timeout)
	}
	return output, err
}

func errorsAgileUnavailable() error {
	return errAgileUnavailable
}
