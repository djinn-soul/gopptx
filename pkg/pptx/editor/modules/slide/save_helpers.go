package slide

import (
	"archive/zip"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

var customXMLNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9._-]*$`)

//nolint:gochecknoglobals // Reused escaper table avoids repeated allocations while serializing custom XML values.
var customXMLEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

func GenerateCustomXMLItem(part common.CustomXMLPart) (string, error) {
	if !customXMLNamePattern.MatchString(part.RootElement) {
		return "", fmt.Errorf("invalid root element name %q", part.RootElement)
	}

	nsAttr := ""
	if part.Namespace != "" {
		nsAttr = fmt.Sprintf(` xmlns="%s"`, escapeCustomXML(part.Namespace))
	}

	inner := part.Content
	if inner == "" && len(part.Properties) > 0 {
		var propsSb strings.Builder
		for j, kv := range part.Properties {
			if !customXMLNamePattern.MatchString(kv.Key) {
				return "", fmt.Errorf("invalid property element name %q", kv.Key)
			}
			if j > 0 {
				propsSb.WriteString("\n  ")
			}
			fmt.Fprintf(&propsSb, "<%s>%s</%s>", kv.Key, escapeCustomXML(kv.Value), kv.Key)
		}
		inner = propsSb.String()
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<%s%s>
  %s
</%s>`, part.RootElement, nsAttr, inner, part.RootElement), nil
}

func VbaProjectFromMetadata(raw any) (*vba.VBAProject, bool) {
	project, ok := raw.(*vba.VBAProject)
	if !ok || project == nil {
		return nil, false
	}
	return project, project.IsMacroEnabled()
}

func MapValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func FilterXMLPartPaths(paths []string) []string {
	filtered := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.HasSuffix(p, ".xml") {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func SaveZipMethod(path string) uint16 {
	lowerPath := strings.ToLower(path)
	if strings.HasPrefix(lowerPath, "ppt/notes") {
		return zip.Store
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tif", ".tiff", ".mp3", ".m4a", ".wav", ".mp4", ".avi":
		return zip.Store
	}
	return zip.Deflate
}

func SerializeCommentAuthorsIfPopulated(
	authorCache map[int64]comments.Author,
	getAuthors func() ([]comments.Author, error),
) ([]byte, bool, error) {
	if authorCache == nil {
		return nil, false, nil
	}
	authors, err := getAuthors()
	if err != nil {
		return nil, false, fmt.Errorf("get authors: %w", err)
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].ID < authors[j].ID
	})
	xmlContent := pptxxml.CommentAuthorsXML(authors)
	return []byte(xmlContent), true, nil
}

func EnsureCommentAuthorsRelationship(
	relationships []common.EditorRelationship,
	nextRelIDNum int,
	relType string,
	target string,
) ([]common.EditorRelationship, int) {
	for _, rel := range relationships {
		if rel.Type == relType {
			return relationships, nextRelIDNum
		}
	}
	relID := fmt.Sprintf("rId%d", nextRelIDNum)
	nextRelIDNum++
	relationships = append(relationships, common.EditorRelationship{
		ID:     relID,
		Type:   relType,
		Target: target,
	})
	return relationships, nextRelIDNum
}

func ResolveNotesMasterRelID(
	relationships []common.EditorRelationship,
	hasNotesMaster bool,
	relType string,
) (string, error) {
	if !hasNotesMaster {
		return "", nil
	}
	for _, rel := range relationships {
		if rel.Type == relType {
			return rel.ID, nil
		}
	}
	return "", errors.New("notes master part exists but presentation relationship is missing")
}

func EscapeCustomXML(value string) string {
	return escapeCustomXML(value)
}

func escapeCustomXML(value string) string {
	return customXMLEscaper.Replace(value)
}
