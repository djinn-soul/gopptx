package structural

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var entityPattern = regexp.MustCompile(`&(amp|lt|gt|quot|apos|#\d+|#x[0-9a-fA-F]+);`)

func (r *Repairer) repairInvalidXML(p string) error {
	data, ok := r.modifier.Get(p)
	if !ok {
		return fmt.Errorf("part not found: %s", p)
	}

	content := string(data)
	if !strings.HasPrefix(strings.TrimSpace(content), "<?xml") {
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + content
	}

	repaired := escapeBareAmpersands(content)
	decoder := xml.NewDecoder(strings.NewReader(repaired))
	for {
		_, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("XML repair failed to produce valid XML: %w", err)
		}
	}

	r.modifier.Set(p, []byte(repaired))
	return nil
}

func escapeBareAmpersands(s string) string {
	var result strings.Builder
	last := 0
	for _, match := range entityPattern.FindAllStringIndex(s, -1) {
		result.WriteString(strings.ReplaceAll(s[last:match[0]], "&", "&amp;"))
		result.WriteString(s[match[0]:match[1]])
		last = match[1]
	}
	result.WriteString(strings.ReplaceAll(s[last:], "&", "&amp;"))
	return result.String()
}

func tryUnmarshalRelationships(data []byte, rels *relationshipsXML) bool {
	return xml.Unmarshal(data, rels) == nil
}
