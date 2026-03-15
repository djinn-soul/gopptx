package editor

import (
	"crypto/sha256"
	"encoding/hex"
)

func (e *PresentationEditor) ensureExcelEmbeddingHashIndex() map[string]string {
	index := make(map[string]string)
	for _, part := range e.parts.KeysWithPrefix("ppt/embeddings/") {
		data, ok := e.parts.Get(part)
		if !ok {
			continue
		}
		sum := sha256.Sum256(data)
		index[hex.EncodeToString(sum[:])] = part
	}
	e.excelEmbeddingHashIndex = index
	return index
}
