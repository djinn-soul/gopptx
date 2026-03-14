//go:build windows

package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const pngExportTimeout = 5 * time.Minute

func exportPPTXToPNG(pptxPath, outDir string) error {
	script := buildPowerPointPNGExportScript(pptxPath, outDir)
	ctx, cancel := context.WithTimeout(context.Background(), pngExportTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("powerpoint png export failed: %w; output: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func buildPowerPointPNGExportScript(pptxPath, outDir string) string {
	inQ := quotePowerShellLiteral(pptxPath)
	outQ := quotePowerShellLiteral(outDir)
	return "$ErrorActionPreference='Stop';" +
		"$in=" + inQ + ";" +
		"$out=" + outQ + ";" +
		"$app=New-Object -ComObject PowerPoint.Application;" +
		"try {" +
		"$pres=$app.Presentations.Open($in,$false,$true,$false);" +
		"$pres.SaveAs($out,18);" + // 18 => ppSaveAsPNG
		"$pres.Close();" +
		"} finally {" +
		"if($app -ne $null){$app.Quit() | Out-Null}" +
		"}"
}

func quotePowerShellLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
