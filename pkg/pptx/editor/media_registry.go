package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// RegisterImageFromFile adds an image from a file path.
func (e *PresentationEditor) RegisterImageFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	format := strings.TrimPrefix(path.Ext(filePath), ".")
	return e.RegisterImage(data, format)
}

// RegisterMedia adds a media file to the presentation or reuses an existing one.
func (e *PresentationEditor) RegisterMedia(data []byte, ext string) (string, error) {
	if e == nil {
		return "", errors.New("editor cannot be nil")
	}
	if len(data) == 0 {
		return "", errors.New("media data cannot be empty")
	}

	hash := sha256.Sum256(data)
	hexHash := hex.EncodeToString(hash[:])

	e.mediaMu.Lock()
	defer e.mediaMu.Unlock()

	if part, ok := e.mediaInventory[hexHash]; ok && strings.HasPrefix(part, "ppt/media/") {
		return part, nil
	}

	partPath := "ppt/media/media" + strconv.Itoa(e.nextMediaNum) + "." + ext
	e.nextMediaNum++

	e.parts.Set(partPath, data)
	e.mediaInventory[hexHash] = partPath

	return partPath, nil
}

// RegisterEmbedding adds an OLE embedding to the presentation.
func (e *PresentationEditor) RegisterEmbedding(data []byte, ext string) (string, error) {
	if e == nil {
		return "", errors.New("editor cannot be nil")
	}
	if len(data) == 0 {
		return "", errors.New("embedding data cannot be empty")
	}

	hash := sha256.Sum256(data)
	hexHash := hex.EncodeToString(hash[:])

	e.mediaMu.Lock()
	defer e.mediaMu.Unlock()

	if part, ok := e.mediaInventory[hexHash]; ok && strings.HasPrefix(part, "ppt/embeddings/") {
		return part, nil
	}

	partPath := "ppt/embeddings/oleObject" + strconv.Itoa(e.nextMediaNum) + "." + ext
	e.nextMediaNum++

	e.parts.Set(partPath, data)
	e.mediaInventory[hexHash] = partPath

	return partPath, nil
}

func (e *PresentationEditor) getOrCreateSlideRelWithType(
	slideIndex int,
	partPath string,
	relType string,
) (string, error) {
	slideRef := e.slides[slideIndex]
	rels, err := e.slideRelationships(slideRef.Part)
	if err != nil {
		return "", err
	}

	var target string
	switch {
	case strings.HasPrefix(partPath, "ppt/media/"):
		target = "../media/" + path.Base(partPath)
	case strings.HasPrefix(partPath, "ppt/embeddings/"):
		target = "../embeddings/" + path.Base(partPath)
	default:
		target = path.Base(partPath)
	}

	for _, r := range rels {
		if r.Type == relType && r.Target == target {
			return r.ID, nil
		}
	}

	nextNum := common.NextRelationshipNumber(rels)
	relID := fmt.Sprintf("rId%d", nextNum)

	rels = append(rels, common.EditorRelationship{
		ID:     relID,
		Type:   relType,
		Target: target,
	})

	e.parts.Set(common.SlideRelsPartName(slideRef.Part), []byte(renderRelationshipsXML(rels)))
	return relID, nil
}
