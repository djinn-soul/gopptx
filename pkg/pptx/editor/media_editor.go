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

//nolint:gochecknoglobals // Reused escaper table avoids repeated allocations in hot XML-building paths.
var xmlAttrEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

func escapeXMLAttr(value string) string {
	return xmlAttrEscaper.Replace(value)
}

// AddVideo adds a video shape to the slide with a poster frame.
func (e *PresentationEditor) AddVideo(
	slideIndex int,
	videoData []byte,
	posterFrameData []byte,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addVideoGeneric(slideIndex, videoData, "", posterFrameData, "", mimeType, x, y, w, h)
}

// AddVideoFromFile adds a video shape from local files.
func (e *PresentationEditor) AddVideoFromFile(
	slideIndex int,
	videoPath string,
	posterFramePath string,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addVideoGeneric(slideIndex, nil, videoPath, nil, posterFramePath, mimeType, x, y, w, h)
}

func (e *PresentationEditor) addVideoGeneric(
	slideIndex int,
	videoData []byte,
	videoPath string,
	posterData []byte,
	posterPath string,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	if err := validateMediaSlideIndex(e, slideIndex); err != nil {
		return 0, err
	}

	posterFramePart, err := registerPartFromDataOrPath(
		posterData,
		posterPath,
		"poster frame data or path is required for video",
		func(data []byte) (string, error) { return e.RegisterImage(data, "png") },
		func(filePath string) (string, error) { return e.RegisterImageFromFile(filePath) },
	)
	if err != nil {
		return 0, fmt.Errorf("register poster frame: %w", err)
	}
	posterRelID, err := e.getOrCreateSlideRel(slideIndex, posterFramePart)
	if err != nil {
		return 0, fmt.Errorf("create poster rel: %w", err)
	}

	videoPart, err := registerVideoPart(e, videoData, videoPath, mimeType)
	if err != nil {
		return 0, fmt.Errorf("register video media: %w", err)
	}

	// 3. Create relationships
	// Legacy video rel
	videoRelID, err := e.getOrCreateSlideRelWithType(slideIndex, videoPart, common.RelTypeVideo)
	if err != nil {
		return 0, err
	}
	// Modern media rel
	mediaRelID, err := e.getOrCreateSlideRelWithType(slideIndex, videoPart, common.RelTypeMedia)
	if err != nil {
		return 0, err
	}

	// 4. Generate XML
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	maxID := maxObjectID(content)
	newID := maxID + 1
	videoXML := buildVideoShapeXML(newID, videoRelID, mediaRelID, posterRelID, x, y, w, h)
	updatedContent, err := appendShapeXMLToSlide(content, videoXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)

	return newID, nil
}

// AddOLEObject adds an OLE object with an icon to the slide.
func (e *PresentationEditor) AddOLEObject(
	slideIndex int,
	objectData []byte,
	iconData []byte,
	progID string,
	x, y, w, h float64,
) (int, error) {
	return e.addOLEObjectGeneric(slideIndex, objectData, "", iconData, "", progID, x, y, w, h)
}

// AddOLEObjectFromFile adds an OLE object from local files.
func (e *PresentationEditor) AddOLEObjectFromFile(
	slideIndex int,
	objectPath string,
	iconPath string,
	progID string,
	x, y, w, h float64,
) (int, error) {
	return e.addOLEObjectGeneric(slideIndex, nil, objectPath, nil, iconPath, progID, x, y, w, h)
}

func (e *PresentationEditor) addOLEObjectGeneric(
	slideIndex int,
	objectData []byte,
	objectPath string,
	iconData []byte,
	iconPath string,
	progID string,
	x, y, w, h float64,
) (int, error) {
	if err := validateMediaSlideIndex(e, slideIndex); err != nil {
		return 0, err
	}

	iconPart, err := registerPartFromDataOrPath(
		iconData,
		iconPath,
		"icon data or path is required for OLE object",
		func(data []byte) (string, error) { return e.RegisterImage(data, "png") },
		func(filePath string) (string, error) { return e.RegisterImageFromFile(filePath) },
	)
	if err != nil {
		return 0, fmt.Errorf("register icon: %w", err)
	}
	iconRelID, err := e.getOrCreateSlideRel(slideIndex, iconPart)
	if err != nil {
		return 0, fmt.Errorf("create icon rel: %w", err)
	}

	embedPart, err := registerEmbeddingPart(e, objectData, objectPath)
	if err != nil {
		return 0, fmt.Errorf("register embedding: %w", err)
	}
	embedRelID, err := e.getOrCreateSlideRelWithType(slideIndex, embedPart, common.RelTypePackage)
	if err != nil {
		return 0, err
	}

	// 3. Generate XML (graphicFrame containing oleObj)
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	maxID := maxObjectID(content)
	newID := maxID + 1
	oleXML := buildOLEObjectShapeXML(newID, slideIndex, embedRelID, iconRelID, progID, x, y, w, h)
	updatedContent, err := appendShapeXMLToSlide(content, oleXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)

	return newID, nil
}

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
