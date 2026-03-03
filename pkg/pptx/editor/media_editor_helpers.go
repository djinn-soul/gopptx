package editor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

const oleShapeIDStride = 10000

func validateMediaSlideIndex(e *PresentationEditor, slideIndex int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}
	return nil
}

func registerPartFromDataOrPath(
	data []byte,
	filePath string,
	missingErr string,
	fromData func([]byte) (string, error),
	fromPath func(string) (string, error),
) (string, error) {
	if len(data) > 0 {
		return fromData(data)
	}
	if filePath != "" {
		return fromPath(filePath)
	}
	return "", errors.New(missingErr)
}

func registerVideoPart(e *PresentationEditor, videoData []byte, videoPath string, mimeType string) (string, error) {
	videoExt := "mp4"
	if mimeType == "video/quicktime" {
		videoExt = "mov"
	}
	return registerPartFromDataOrPath(
		videoData,
		videoPath,
		"video data or path is required",
		func(data []byte) (string, error) {
			return e.RegisterMedia(data, videoExt)
		},
		func(filePath string) (string, error) {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", err
			}
			return e.RegisterMedia(data, strings.TrimPrefix(path.Ext(filePath), "."))
		},
	)
}

func registerEmbeddingPart(e *PresentationEditor, objectData []byte, objectPath string) (string, error) {
	return registerPartFromDataOrPath(
		objectData,
		objectPath,
		"object data or path is required",
		func(data []byte) (string, error) {
			return e.RegisterEmbedding(data, "bin")
		},
		func(filePath string) (string, error) {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", err
			}
			return e.RegisterEmbedding(data, strings.TrimPrefix(path.Ext(filePath), "."))
		},
	)
}

func appendShapeXMLToSlide(content []byte, shapeXML string) ([]byte, error) {
	endTree := []byte("</p:spTree>")
	idx := bytes.LastIndex(content, endTree)
	if idx == -1 {
		return nil, errors.New("invalid slide xml: missing spTree end")
	}

	var buf bytes.Buffer
	buf.Write(content[:idx])
	buf.WriteString(shapeXML)
	buf.Write(content[idx:])
	return buf.Bytes(), nil
}

func buildVideoShapeXML(
	newID int,
	videoRelID string,
	mediaRelID string,
	posterRelID string,
	x, y, w, h float64,
) string {
	name := fmt.Sprintf("Video %d", newID)
	return fmt.Sprintf(`
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
}

func buildOLEObjectShapeXML(
	newID int,
	slideIndex int,
	embedRelID string,
	iconRelID string,
	progID string,
	x,
	y,
	w,
	h float64,
) string {
	name := fmt.Sprintf("Object %d", newID)
	safeProgID := escapeXMLAttr(progID)
	return fmt.Sprintf(`
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
</p:graphicFrame>`, newID, name, int64(x), int64(y), int64(w), int64(h), 1024+(slideIndex*oleShapeIDStride)+newID, embedRelID, safeProgID, embedRelID, safeProgID, iconRelID, int64(x), int64(y), int64(w), int64(h))
}
