package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	defaultPythonScript = "examples/python/scripts/python_group_ungroup_validate.py"
	defaultOutputPPTX   = "examples/output/python_generated_validated.pptx"
	commandTimeout      = 5 * time.Minute
)

func main() {
	pythonScript := flag.String("script", defaultPythonScript, "path to python generation script")
	outputPPTX := flag.String("output", defaultOutputPPTX, "generated pptx file path")
	pythonBin := flag.String("python", "", "optional python executable path")
	flag.Parse()

	root, err := repoRoot()
	if err != nil {
		failf("resolve repo root: %v", err)
	}

	pyExec, err := resolvePythonExecutable(root, *pythonBin)
	if err != nil {
		failf("resolve python executable: %v", err)
	}

	if err := runCommand(root, pyExec, *pythonScript); err != nil {
		failf("run python generator: %v", err)
	}

	if err := runCommand(
		root,
		"go",
		"run",
		"./scripts/smoke/validate_smoke_outputs",
		"-file",
		*outputPPTX,
	); err != nil {
		failf("validate generated pptx: %v", err)
	}

	fmt.Printf("PASS generated and validated %s\n", *outputPPTX)
}

func repoRoot() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return root, nil
}

func resolvePythonExecutable(root, provided string) (string, error) {
	if provided != "" {
		return provided, nil
	}

	candidates := pythonCandidates(root)
	for _, candidate := range candidates {
		if filepath.IsAbs(candidate) {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
			continue
		}
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}

	return "", errors.New("no python executable found (.venv or PATH)")
}

func pythonCandidates(root string) []string {
	if runtime.GOOS == "windows" {
		return []string{
			filepath.Join(root, ".venv", "Scripts", "python.exe"),
			"python",
			"py",
		}
	}
	return []string{
		filepath.Join(root, ".venv", "bin", "python3"),
		filepath.Join(root, ".venv", "bin", "python"),
		"python3",
		"python",
	}
}

func runCommand(root, command string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v: %w", command, args, err)
	}
	return nil
}

func failf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
