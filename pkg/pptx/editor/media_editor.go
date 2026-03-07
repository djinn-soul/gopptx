package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

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
	if err := editormodmedia.ValidateMediaSlideIndex(slideIndex, len(e.slides)); err != nil {
		return 0, err
	}

	posterFramePart, err := editormodmedia.RegisterPartFromDataOrPath(
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

	videoPart, err := editormodmedia.RegisterVideoPart(
		videoData,
		videoPath,
		mimeType,
		e.RegisterMedia,
	)
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

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	newID := maxID + 1
	videoXML := editormodmedia.BuildVideoShapeXML(newID, videoRelID, mediaRelID, posterRelID, x, y, w, h)
	updatedContent, err := editormodmedia.AppendShapeXMLToSlide(content, videoXML)
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
	if err := editormodmedia.ValidateMediaSlideIndex(slideIndex, len(e.slides)); err != nil {
		return 0, err
	}

	iconPart, err := editormodmedia.RegisterPartFromDataOrPath(
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

	embedPart, err := editormodmedia.RegisterEmbeddingPart(
		objectData,
		objectPath,
		e.RegisterEmbedding,
	)
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

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	newID := maxID + 1
	oleXML := editormodmedia.BuildOLEObjectShapeXML(newID, slideIndex, embedRelID, iconRelID, progID, x, y, w, h)
	updatedContent, err := editormodmedia.AppendShapeXMLToSlide(content, oleXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)

	return newID, nil
}

// AddAudio adds an audio shape to the slide.
func (e *PresentationEditor) AddAudio(
	slideIndex int,
	audioData []byte,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addAudioGeneric(slideIndex, audioData, "", mimeType, x, y, w, h)
}

// AddAudioFromFile adds an audio shape from a local file.
func (e *PresentationEditor) AddAudioFromFile(
	slideIndex int,
	audioPath string,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addAudioGeneric(slideIndex, nil, audioPath, mimeType, x, y, w, h)
}

func (e *PresentationEditor) addAudioGeneric(
	slideIndex int,
	audioData []byte,
	audioPath string,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	if err := editormodmedia.ValidateMediaSlideIndex(slideIndex, len(e.slides)); err != nil {
		return 0, err
	}

	audioPart, err := editormodmedia.RegisterAudioPart(
		audioData,
		audioPath,
		mimeType,
		e.RegisterMedia,
	)
	if err != nil {
		return 0, fmt.Errorf("register audio media: %w", err)
	}

	// Create relationships
	// Legacy audio rel
	audioRelID, err := e.getOrCreateSlideRelWithType(slideIndex, audioPart, common.RelTypeAudio)
	if err != nil {
		return 0, err
	}
	// Modern media rel
	mediaRelID, err := e.getOrCreateSlideRelWithType(slideIndex, audioPart, common.RelTypeMedia)
	if err != nil {
		return 0, err
	}

	// Generate XML
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	newID := maxID + 1
	audioXML := editormodmedia.BuildAudioShapeXML(newID, audioRelID, mediaRelID, x, y, w, h)
	updatedContent, err := editormodmedia.AppendShapeXMLToSlide(content, audioXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)

	return newID, nil
}
