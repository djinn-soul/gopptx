package opc

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// A part reaches another part only through a relationship, so reading or
// editing parts by hand is only half an escape hatch without them.

// Relationship is one entry of a .rels part.
type Relationship struct {
	ID     string
	Type   string
	Target string
	// External marks a target outside the package, such as a hyperlink or a
	// hosted video.
	External bool
}

// Relationships is the parsed content of one .rels part.
type Relationships struct {
	items []Relationship
}

const (
	relationshipsNamespace = "http://schemas.openxmlformats.org/package/2006/relationships"
	externalTargetMode     = "External"
)

// relationshipsXML is the wire shape of a .rels part.
type relationshipsXML struct {
	XMLName xml.Name           `xml:"Relationships"`
	Items   []relationshipItem `xml:"Relationship"`
}

type relationshipItem struct {
	ID         string `xml:"Id,attr"`
	Type       string `xml:"Type,attr"`
	Target     string `xml:"Target,attr"`
	TargetMode string `xml:"TargetMode,attr,omitempty"`
}

// ParseRelationships reads a .rels document.
func ParseRelationships(data []byte) (*Relationships, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return &Relationships{}, nil
	}
	var doc relationshipsXML
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse relationships: %w", err)
	}
	rels := &Relationships{items: make([]Relationship, 0, len(doc.Items))}
	for _, item := range doc.Items {
		rels.items = append(rels.items, Relationship{
			ID:       item.ID,
			Type:     item.Type,
			Target:   item.Target,
			External: strings.EqualFold(item.TargetMode, externalTargetMode),
		})
	}
	return rels, nil
}

// Relationships returns the relationships of a part. A part with no .rels file
// has none, which is not an error.
func (p *Package) Relationships(partPath string) (*Relationships, error) {
	relsPath := RelationshipsPath(partPath)
	if !p.Has(relsPath) {
		return &Relationships{}, nil
	}
	data, err := p.Part(relsPath)
	if err != nil {
		return nil, err
	}
	return ParseRelationships(data)
}

// SetRelationships writes a part's relationships back into the package.
func (p *Package) SetRelationships(partPath string, rels *Relationships) {
	p.SetPart(RelationshipsPath(partPath), []byte(rels.XML()))
}

// All returns every relationship, in document order.
func (r *Relationships) All() []Relationship {
	return append([]Relationship(nil), r.items...)
}

// Len is how many relationships the part has.
func (r *Relationships) Len() int {
	return len(r.items)
}

// ByID finds a relationship by its rId.
func (r *Relationships) ByID(id string) (Relationship, bool) {
	for _, item := range r.items {
		if item.ID == id {
			return item, true
		}
	}
	return Relationship{}, false
}

// ByType returns every relationship of one type, such as all the slides.
func (r *Relationships) ByType(relType string) []Relationship {
	var out []Relationship
	for _, item := range r.items {
		if item.Type == relType {
			out = append(out, item)
		}
	}
	return out
}

// Add appends a relationship with the next free rId and returns it.
func (r *Relationships) Add(relType, target string) Relationship {
	rel := Relationship{ID: r.nextID(), Type: relType, Target: target}
	r.items = append(r.items, rel)
	return rel
}

// AddExternal appends a relationship whose target is outside the package.
func (r *Relationships) AddExternal(relType, target string) Relationship {
	rel := Relationship{ID: r.nextID(), Type: relType, Target: target, External: true}
	r.items = append(r.items, rel)
	return rel
}

// AddWithID appends a relationship with a chosen id, replacing any relationship
// that already holds it.
func (r *Relationships) AddWithID(id, relType, target string) Relationship {
	rel := Relationship{ID: id, Type: relType, Target: target}
	for i, item := range r.items {
		if item.ID == id {
			r.items[i] = rel
			return rel
		}
	}
	r.items = append(r.items, rel)
	return rel
}

// Remove deletes a relationship by id, reporting whether it was there.
func (r *Relationships) Remove(id string) bool {
	for i, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return true
		}
	}
	return false
}

// XML renders the .rels document.
func (r *Relationships) XML() string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	b.WriteString(`<Relationships xmlns="` + relationshipsNamespace + `">`)
	for _, item := range r.items {
		b.WriteString("\n<Relationship Id=\"" + escapeAttr(item.ID) +
			"\" Type=\"" + escapeAttr(item.Type) +
			"\" Target=\"" + escapeAttr(item.Target) + "\"")
		if item.External {
			b.WriteString(` TargetMode="` + externalTargetMode + `"`)
		}
		b.WriteString("/>")
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}

// nextID is one past the highest rId in use, so ids stay unique within a part.
func (r *Relationships) nextID() string {
	highest := 0
	for _, item := range r.items {
		number, err := strconv.Atoi(strings.TrimPrefix(item.ID, "rId"))
		if err == nil && number > highest {
			highest = number
		}
	}
	return "rId" + strconv.Itoa(highest+1)
}

// Types lists the distinct relationship types a part uses, sorted.
func (r *Relationships) Types() []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range r.items {
		if seen[item.Type] {
			continue
		}
		seen[item.Type] = true
		out = append(out, item.Type)
	}
	sort.Strings(out)
	return out
}

func escapeAttr(value string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(value))
	return b.String()
}
