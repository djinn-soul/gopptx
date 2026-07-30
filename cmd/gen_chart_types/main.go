// Command gen_chart_types generates the Python ChartType enum from the
// canonical XLChartType* constants declared in pkg/pptx/enums/shape_chart.go.
package main

import (
	"fmt"
	"os"
)

const minCLIArgs = 4

func main() {
	if len(os.Args) < minCLIArgs {
		fmt.Fprintln(os.Stderr,
			"Usage: gen_chart_types <enums_go_file> <chartkind_go_file> <output_py_file>")
		os.Exit(1)
	}

	enumsFile := os.Args[1]
	kindFile := os.Args[2]
	outputPy := os.Args[3]

	kinds, err := parseChartKinds(kindFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing chart kinds: %v\n", err)
		os.Exit(1)
	}

	types, err := parseChartTypes(enumsFile, kinds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing chart types: %v\n", err)
		os.Exit(1)
	}
	if len(types) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no %s* constants found in %s\n", chartTypePrefix, enumsFile)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPy, []byte(renderPython(types)), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPy, err)
		os.Exit(1)
	}
}
