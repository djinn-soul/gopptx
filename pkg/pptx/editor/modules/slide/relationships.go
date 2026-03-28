package slide

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"sync"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

//nolint:gochecknoglobals // pool reduces bytes.Reader allocs in hot XML parsing paths
var bytesReaderPool = sync.Pool{
	New: func() any { return new(bytes.Reader) },
}

type ParsedSlideIDRef struct {
	SlideID int64
	RelID   string
	Hidden  bool
}

const (
	defaultRelsCapacity     = 8
	defaultSlideIDsCapacity = 8
)

func ParseRelationshipsXML(content []byte) ([]common.EditorRelationship, error) {
	reader, ok := bytesReaderPool.Get().(*bytes.Reader)
	if !ok || reader == nil {
		return nil, errors.New("bytes reader pool returned invalid reader")
	}
	r := reader
	r.Reset(content)
	decoder := xml.NewDecoder(r)
	out := make([]common.EditorRelationship, 0, defaultRelsCapacity)

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			bytesReaderPool.Put(r)
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}

		rel := common.EditorRelationship{}
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "Id":
				rel.ID = strings.TrimSpace(attr.Value)
			case "Type":
				rel.Type = strings.TrimSpace(attr.Value)
			case "Target":
				rel.Target = strings.TrimSpace(attr.Value)
			case "TargetMode":
				rel.TargetMode = strings.TrimSpace(attr.Value)
			}
		}
		if rel.ID == "" || rel.Type == "" || rel.Target == "" {
			bytesReaderPool.Put(r)
			return nil, errors.New("relationship with missing Id/Type/Target")
		}
		out = append(out, rel)
	}
	bytesReaderPool.Put(r)
	return out, nil
}

func ParsePresentationSlideIDs(content []byte) ([]ParsedSlideIDRef, error) {
	reader, ok := bytesReaderPool.Get().(*bytes.Reader)
	if !ok || reader == nil {
		return nil, errors.New("bytes reader pool returned invalid reader")
	}
	r := reader
	r.Reset(content)
	decoder := xml.NewDecoder(r)
	out := make([]ParsedSlideIDRef, 0, defaultSlideIDsCapacity)

	for {
		start, ok, err := nextSlideIDStartElement(decoder)
		if err != nil {
			bytesReaderPool.Put(r)
			return nil, err
		}
		if !ok {
			break
		}
		if !isPresentationSlideIDElement(start) {
			continue
		}

		ref, err := parseSlideIDRef(start)
		if err != nil {
			bytesReaderPool.Put(r)
			return nil, err
		}
		if ref.SlideID == 0 || ref.RelID == "" {
			bytesReaderPool.Put(r)
			return nil, errors.New("slide id entry missing id or r:id")
		}
		out = append(out, ref)
	}
	bytesReaderPool.Put(r)
	return out, nil
}

func NormalizePresentationTarget(target string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	clean = strings.TrimPrefix(clean, "/")
	clean = strings.TrimPrefix(clean, "ppt/")
	return path.Clean(clean)
}

func nextSlideIDStartElement(decoder *xml.Decoder) (xml.StartElement, bool, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return xml.StartElement{}, false, nil
			}
			return xml.StartElement{}, false, err
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		return start, true, nil
	}
}

func isPresentationSlideIDElement(start xml.StartElement) bool {
	if start.Name.Local != "sldId" {
		return false
	}
	if start.Name.Space == "" {
		return true
	}
	return start.Name.Space == "http://schemas.openxmlformats.org/presentationml/2006/main"
}

func parseSlideIDRef(start xml.StartElement) (ParsedSlideIDRef, error) {
	ref := ParsedSlideIDRef{}
	for _, attr := range start.Attr {
		switch {
		case attr.Name.Local == "id" && attr.Name.Space == "":
			slideID, parseErr := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
			if parseErr != nil {
				return ParsedSlideIDRef{}, fmt.Errorf("invalid slide id %q", attr.Value)
			}
			ref.SlideID = slideID
		case attr.Name.Local == "show" && attr.Name.Space == "":
			ref.Hidden = strings.TrimSpace(attr.Value) == "0"
		case attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" ||
			attr.Name.Space == "r":
			ref.RelID = strings.TrimSpace(attr.Value)
		}
	}
	return ref, nil
}
