package shape

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// CopyTextRuns returns a detached copy of runs for safe caller-side mutation.
func CopyTextRuns(runs []common.TextRun) []common.TextRun {
	if len(runs) == 0 {
		return nil
	}
	out := make([]common.TextRun, len(runs))
	copy(out, runs)
	return out
}

// UpdateRunText updates one run by index and returns a copied run slice.
func UpdateRunText(runs []common.TextRun, runIndex int, text string) ([]common.TextRun, error) {
	updated := CopyTextRuns(runs)
	if runIndex < 0 || runIndex >= len(updated) {
		return nil, fmt.Errorf("run index %d out of range [0,%d)", runIndex, len(updated))
	}
	updated[runIndex].Text = text
	return updated, nil
}

// AppendRun appends one run and returns a copied run slice.
func AppendRun(runs []common.TextRun, run common.TextRun) []common.TextRun {
	updated := CopyTextRuns(runs)
	return append(updated, run)
}
