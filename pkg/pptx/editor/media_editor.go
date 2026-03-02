package editor

import (
	"bytes"
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
func (e *PresentationEditor) AddVideo(slideIndex int, videoData []byte, posterFrameData []byte, mimeType string, x, y, w, h float64) (int, error) {
	return e.addVideoGeneric(slideIndex, videoData, "", posterFrameData, "", mimeType, x, y, w, h)
}

// AddVideoFromFile adds a video shape from local files.
func (e *PresentationEditor) AddVideoFromFile(slideIndex int, videoPath string, posterFramePath string, mimeType string, x, y, w, h float64) (int, error) {
	return e.addVideoGeneric(slideIndex, nil, videoPath, nil, posterFramePath, mimeType, x, y, w, h)
}

func (e *PresentationEditor) addVideoGeneric(slideIndex int, videoData []byte, videoPath string, posterData []byte, posterPath string, mimeType string, x, y, w, h float64) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	// 1. Register poster frame image
	var posterFramePart string
	var err error
	if len(posterData) > 0 {
		posterFramePart, err = e.RegisterImage(posterData, "png")
	} else if posterPath != "" {
		posterFramePart, err = e.RegisterImageFromFile(posterPath)
	} else {
		return 0, errors.New("poster frame data or path is required for video")
	}
	if err != nil {
		return 0, fmt.Errorf("register poster frame: %w", err)
	}
	posterRelID, err := e.getOrCreateSlideRel(slideIndex, posterFramePart)
	if err != nil {
		return 0, fmt.Errorf("create poster rel: %w", err)
	}

	// 2. Register video media
	var videoPart string
	videoExt := "mp4"
	if mimeType == "video/quicktime" {
		videoExt = "mov"
	}

	if len(videoData) > 0 {
		videoPart, err = e.RegisterMedia(videoData, videoExt)
	} else if videoPath != "" {
		data, err2 := os.ReadFile(videoPath)
		if err2 != nil {
			return 0, err2
		}
		videoPart, err = e.RegisterMedia(data, strings.TrimPrefix(path.Ext(videoPath), "."))
	} else {
		return 0, errors.New("video data or path is required")
	}
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
	name := fmt.Sprintf("Video %d", newID)

	videoXML := fmt.Sprintf(`
<p:pic>
  <p:nvPicPr>
    <p:cNvPr id="%d" name="%s">
      <a:extLst>
        <a:ext uri="{FF2B5EF4-FFF2-40B4-BE49-F238E27FC236}">
          <a16:creationId xmlns:a16="http://schemas.microsoft.com/office/drawing/2014/main" id="{00000000-0000-0000-0000-000000000000}"/>
        </a:ext>
      </a:extLst>
    </p:cNvPr>
    <p:cNvPicPr>
      <a:picLocks noChangeAspect="1"/>
    </p:cNvPicPr>
    <p:nvPr>
      <a:videoFile r:link="%s"/>
      <p:extLst>
        <p:ext uri="{DAA4B4D4-6D71-4841-9C94-3DE7FCFB9230}">
          <p14:media xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" r:embed="%s"/>
        </p:ext>
      </p:extLst>
    </p:nvPr>
  </p:nvPicPr>
  <p:blipFill>
    <a:blip r:embed="%s"/>
    <a:stretch><a:fillRect/></a:stretch>
  </p:blipFill>
  <p:spPr>
    <a:xfrm>
      <a:off x="%d" y="%d"/>
      <a:ext cx="%d" cy="%d"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>
</p:pic>`, newID, name, videoRelID, mediaRelID, posterRelID, int64(x), int64(y), int64(w), int64(h))

	// 5. Update slide XML
	endTree := []byte("</p:spTree>")
	idx := bytes.LastIndex(content, endTree)
	if idx == -1 {
		return 0, errors.New("invalid slide xml: missing spTree end")
	}

	var buf bytes.Buffer
	buf.Write(content[:idx])
	buf.WriteString(videoXML)
	buf.Write(content[idx:])

	e.parts.Set(slideRef.Part, buf.Bytes())

	return newID, nil
}

// AddOLEObject adds an OLE object with an icon to the slide.
func (e *PresentationEditor) AddOLEObject(slideIndex int, objectData []byte, iconData []byte, progID string, x, y, w, h float64) (int, error) {
	return e.addOLEObjectGeneric(slideIndex, objectData, "", iconData, "", progID, x, y, w, h)
}

// AddOLEObjectFromFile adds an OLE object from local files.
func (e *PresentationEditor) AddOLEObjectFromFile(slideIndex int, objectPath string, iconPath string, progID string, x, y, w, h float64) (int, error) {
	return e.addOLEObjectGeneric(slideIndex, nil, objectPath, nil, iconPath, progID, x, y, w, h)
}

func (e *PresentationEditor) addOLEObjectGeneric(slideIndex int, objectData []byte, objectPath string, iconData []byte, iconPath string, progID string, x, y, w, h float64) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	// 1. Register icon image
	var iconPart string
	var err error
	if len(iconData) > 0 {
		iconPart, err = e.RegisterImage(iconData, "png")
	} else if iconPath != "" {
		iconPart, err = e.RegisterImageFromFile(iconPath)
	} else {
		return 0, errors.New("icon data or path is required for OLE object")
	}
	if err != nil {
		return 0, fmt.Errorf("register icon: %w", err)
	}
	iconRelID, err := e.getOrCreateSlideRel(slideIndex, iconPart)
	if err != nil {
		return 0, fmt.Errorf("create icon rel: %w", err)
	}

	// 2. Register OLE embedding
	var embedPart string
	embedExt := "bin"
	if len(objectData) > 0 {
		embedPart, err = e.RegisterEmbedding(objectData, embedExt)
	} else if objectPath != "" {
		data, err2 := os.ReadFile(objectPath)
		if err2 != nil {
			return 0, err2
		}
		embedPart, err = e.RegisterEmbedding(data, strings.TrimPrefix(path.Ext(objectPath), "."))
	} else {
		return 0, errors.New("object data or path is required")
	}
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
	name := fmt.Sprintf("Object %d", newID)

	safeProgID := escapeXMLAttr(progID)
	oleXML := fmt.Sprintf(`
<p:graphicFrame>
  <p:nvGraphicFramePr>
    <p:cNvPr id="%d" name="%s"/>
    <p:cNvGraphicFramePr>
      <a:graphicFrameLocks noChangeAspect="1"/>
    </p:cNvGraphicFramePr>
    <p:nvPr/>
  </p:nvGraphicFramePr>
  <p:xfrm>
    <a:off x="%d" y="%d"/>
    <a:ext cx="%d" cy="%d"/>
  </p:xfrm>
  <a:graphic>
    <a:graphicData uri="http://schemas.openxmlformats.org/presentationml/2006/ole">
      <mc:AlternateContent xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">
        <mc:Choice xmlns:v="urn:schemas-microsoft-com:vml" Requires="v">
          <p:oleObj spid="_x0000_s%d" name="Object" showAsIcon="1" r:id="%s" imgW="0" imgH="0" progId="%s">
            <p:embed/>
          </p:oleObj>
        </mc:Choice>
        <mc:Fallback>
          <p:oleObj name="Object" showAsIcon="1" r:id="%s" imgW="0" imgH="0" progId="%s">
            <p:embed/>
            <p:pic>
              <p:nvPicPr>
                <p:cNvPr id="0" name=""/>
                <p:cNvPicPr/>
                <p:nvPr/>
              </p:nvPicPr>
              <p:blipFill>
                <a:blip r:embed="%s"/>
                <a:stretch><a:fillRect/></a:stretch>
              </p:blipFill>
              <p:spPr>
                <a:xfrm>
                  <a:off x="%d" y="%d"/>
                  <a:ext cx="%d" cy="%d"/>
                </a:xfrm>
                <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
              </p:spPr>
            </p:pic>
          </p:oleObj>
        </mc:Fallback>
      </mc:AlternateContent>
    </a:graphicData>
  </a:graphic>
</p:graphicFrame>`, newID, name, int64(x), int64(y), int64(w), int64(h), 1024+(slideIndex*10000)+newID, embedRelID, safeProgID, embedRelID, safeProgID, iconRelID, int64(x), int64(y), int64(w), int64(h))

	// 4. Update slide XML
	endTree := []byte("</p:spTree>")
	idx := bytes.LastIndex(content, endTree)
	if idx == -1 {
		return 0, errors.New("invalid slide xml: missing spTree end")
	}

	var buf bytes.Buffer
	buf.Write(content[:idx])
	buf.WriteString(oleXML)
	buf.Write(content[idx:])

	e.parts.Set(slideRef.Part, buf.Bytes())

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

func (e *PresentationEditor) getOrCreateSlideRelWithType(slideIndex int, partPath string, relType string) (string, error) {
	slideRef := e.slides[slideIndex]
	rels, err := e.slideRelationships(slideRef.Part)
	if err != nil {
		return "", err
	}

	var target string
	if strings.HasPrefix(partPath, "ppt/media/") {
		target = "../media/" + path.Base(partPath)
	} else if strings.HasPrefix(partPath, "ppt/embeddings/") {
		target = "../embeddings/" + path.Base(partPath)
	} else {
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
