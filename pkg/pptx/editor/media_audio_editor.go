package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// AddAudio adds an audio shape to the slide without a custom icon.
func (e *PresentationEditor) AddAudio(
	slideIndex int,
	audioData []byte,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addAudioGeneric(slideIndex, audioData, "", nil, "", mimeType, "", x, y, w, h)
}

// AddAudioFromFile adds an audio shape from a local file without a custom icon.
func (e *PresentationEditor) AddAudioFromFile(
	slideIndex int,
	audioPath string,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addAudioGeneric(slideIndex, nil, audioPath, nil, "", mimeType, "", x, y, w, h)
}

// AddAudioWithIcon adds an audio shape with a custom playback icon.
func (e *PresentationEditor) AddAudioWithIcon(
	slideIndex int,
	audioData []byte,
	iconData []byte,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addAudioGeneric(slideIndex, audioData, "", iconData, "", mimeType, "", x, y, w, h)
}

// AddAudioWithIconFromFile adds an audio shape from files with a custom icon.
func (e *PresentationEditor) AddAudioWithIconFromFile(
	slideIndex int,
	audioPath string,
	iconPath string,
	mimeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addAudioGeneric(slideIndex, nil, audioPath, nil, iconPath, mimeType, "", x, y, w, h)
}

func (e *PresentationEditor) addAudioGeneric(
	slideIndex int,
	audioData []byte,
	audioPath string,
	iconData []byte,
	iconPath string,
	mimeType string,
	altText string,
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

	iconRelID, err := e.registerOptionalAudioIcon(slideIndex, iconData, iconPath)
	if err != nil {
		return 0, err
	}

	audioRelID, err := e.getOrCreateSlideRelWithType(slideIndex, audioPart, common.RelTypeAudio)
	if err != nil {
		return 0, err
	}
	mediaRelID, err := e.getOrCreateSlideRelWithType(slideIndex, audioPart, common.RelTypeMedia)
	if err != nil {
		return 0, err
	}

	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	newID := maxID + 1
	audioXML := editormodmedia.BuildAudioShapeXML(newID, audioRelID, mediaRelID, iconRelID, altText, x, y, w, h)
	updatedContent, err := editormodmedia.AppendShapeXMLToSlide(content, audioXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)

	return newID, nil
}

func (e *PresentationEditor) registerOptionalAudioIcon(
	slideIndex int,
	iconData []byte,
	iconPath string,
) (string, error) {
	if len(iconData) == 0 && iconPath == "" {
		return "", nil
	}

	iconPart, err := editormodmedia.RegisterPartFromDataOrPath(
		iconData,
		iconPath,
		"icon data or path is required for audio icon",
		func(data []byte) (string, error) { return e.RegisterImage(data, "png") },
		func(filePath string) (string, error) { return e.RegisterImageFromFile(filePath) },
	)
	if err != nil {
		return "", fmt.Errorf("register audio icon: %w", err)
	}
	iconRelID, err := e.getOrCreateSlideRel(slideIndex, iconPart)
	if err != nil {
		return "", fmt.Errorf("create audio icon rel: %w", err)
	}
	return iconRelID, nil
}
