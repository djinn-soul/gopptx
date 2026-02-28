package customxml

import "github.com/djinn-soul/gopptx/pkg/pptx/common"

// Store is a collection of custom XML parts to embed in a presentation.
type Store struct {
	items []Part
}

// NewStore creates a new custom XML store.
func NewStore() *Store {
	return &Store{}
}

// Add creates a new CustomXMLPart with the specified root element
// and appends it to this store. Returns a pointer to the added part
// for fluent configuration.
func (s *Store) Add(rootElement string) *Part {
	p := Part{
		part: common.CustomXMLPart{
			RootElement: rootElement,
		},
	}
	s.items = append(s.items, p)
	return &s.items[len(s.items)-1]
}

// AddRaw allows appending a pre-built common.CustomXMLPart to the store.
func (s *Store) AddRaw(part common.CustomXMLPart) {
	s.items = append(s.items, Part{part: part})
}

// Items returns a list of all custom XML parts in the store.
func (s *Store) Items() []Part {
	return s.items
}

// Len returns the number of custom XML parts in the store.
func (s *Store) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the store has no items.
func (s *Store) IsEmpty() bool {
	return len(s.items) == 0
}

// ToCommonParts returns a slice of common.CustomXMLPart suitable for
// assignment to common.Metadata.CustomXML.
func (s *Store) ToCommonParts() []common.CustomXMLPart {
	out := make([]common.CustomXMLPart, len(s.items))
	for i, p := range s.items {
		out[i] = p.part
	}
	return out
}

// Part represents a single custom XML document in a fluent builder format.
type Part struct {
	part common.CustomXMLPart
}

// Namespace sets the namespace URI for this custom XML part.
func (p *Part) Namespace(ns string) *Part {
	p.part.Namespace = ns
	return p
}

// Content sets the raw inner content of the custom XML part.
func (p *Part) Content(xml string) *Part {
	p.part.Content = xml
	return p
}

// Property adds a key-value property to the custom XML part.
// Properties are serialized as simple `<Key>Value</Key>` child elements.
func (p *Part) Property(key, value string) *Part {
	p.part.Properties = append(p.part.Properties, common.CustomXMLKV{
		Key:   key,
		Value: value,
	})
	return p
}
