package editor

import editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"

func cloneBytes(b []byte) []byte {
	return editorslide.CloneBytes(b)
}
