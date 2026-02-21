package main

import "testing"

func TestUnquoteEscapes(t *testing.T) {
	got, err := unquote(`"line\nnext"`)
	if err != nil {
		t.Fatalf("unquote error: %v", err)
	}
	if got != "line\nnext" {
		t.Fatalf("unexpected unquote value: %q", got)
	}
}
