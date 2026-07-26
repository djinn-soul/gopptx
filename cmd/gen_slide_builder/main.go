// Command gen_slide_builder generates SlideBuilder: a pointer-receiver wrapper
// around the value-receiver SlideContent builder methods.
//
// SlideContent's chainable methods take a value receiver and return a new
// SlideContent, so a caller who forgets to reassign the result silently loses
// the change:
//
//	s := NewSlide("t")
//	s.AddBullet("lost")  // compiles, does nothing
//
// SlideBuilder mutates in place instead, so the same call cannot be dropped.
// The wrapper is generated rather than hand-written so it cannot fall behind
// the methods it delegates to; adding a chainable method to SlideContent and
// re-running `go generate ./...` is enough.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const minCLIArgs = 3

func main() {
	if len(os.Args) < minCLIArgs {
		fmt.Fprintln(os.Stderr, "Usage: gen_slide_builder <elements_dir> <output_base_file>")
		os.Exit(1)
	}
	dir, output := os.Args[1], os.Args[2]

	methods, imports, err := collect(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error collecting methods: %v\n", err)
		os.Exit(1)
	}
	if len(methods) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no chainable %s methods found in %s\n", receiverType, dir)
		os.Exit(1)
	}

	if err := writeAll(output, methods, imports); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
}

// writeAll emits the core builder file plus one file per chunk of methods, so
// that no generated file crosses the repository's per-file line ceiling.
func writeAll(output string, methods []method, imports map[string]string) error {
	core, err := renderCore()
	if err != nil {
		return err
	}
	if err := os.WriteFile(output, []byte(core), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", output, err)
	}

	if err := removeStaleChunks(output); err != nil {
		return err
	}

	for index, chunk := range chunkMethods(methods) {
		path := chunkPath(output, index+1)
		content, renderErr := renderChunk(chunk, importsFor(chunk, imports))
		if renderErr != nil {
			return renderErr
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}

// chunkPath names the Nth methods file next to the core file, turning
// slide_builder_gen.go into slide_builder_methods1_gen.go.
func chunkPath(output string, index int) string {
	dir := filepath.Dir(output)
	base := filepath.Base(output)
	stem := base[:len(base)-len(filepath.Ext(base))]
	stem = trimSuffix(stem, genSuffix)
	return filepath.Join(dir, fmt.Sprintf("%s_methods%d%s.go", stem, index, genSuffix))
}

// removeStaleChunks deletes method files left over from a previous run, so that
// shrinking the method set does not leave orphaned generated code behind.
func removeStaleChunks(output string) error {
	dir := filepath.Dir(output)
	base := filepath.Base(output)
	stem := base[:len(base)-len(filepath.Ext(base))]
	stem = trimSuffix(stem, genSuffix)

	matches, err := filepath.Glob(filepath.Join(dir, stem+"_methods*"+genSuffix+".go"))
	if err != nil {
		return fmt.Errorf("glob stale chunks: %w", err)
	}
	for _, path := range matches {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}
	return nil
}

func trimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}
