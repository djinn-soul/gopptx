//go:build !windows

package main

import "errors"

func exportPPTXToPNG(_, _ string) error {
	return errors.New("PNG export requires Windows PowerPoint automation")
}
