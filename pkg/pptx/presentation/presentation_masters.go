package presentation

import "github.com/djinn-soul/gopptx/pkg/pptx/elements"

func getEffectiveMasters(meta Metadata) []*elements.SlideMaster {
	if len(meta.Masters) > 0 {
		return meta.Masters
	}
	if meta.Master != nil {
		return []*elements.SlideMaster{meta.Master}
	}
	return []*elements.SlideMaster{elements.NewMaster()}
}

func getNotesThemeIndex(hasNotes bool) int {
	if !hasNotes {
		return 0
	}
	return minMasterCountWithNativeNotesTheme
}
