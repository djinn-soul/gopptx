package tplx

import "regexp"

// tokenPattern matches a complete Jinja-style token like {{foo}} or {{#each items}}.
var tokenPattern = regexp.MustCompile(`\{\{[^{}]+\}\}`)

// runData holds the parsed data for one <a:r> XML element.
type runData struct {
	rprBytes []byte
	text     string
	extras   [][]byte
}

// paraChild is a polymorphic child inside a paragraph: either a run or raw XML bytes.
type paraChild struct {
	isRun bool
	run   runData
	raw   []byte
}
