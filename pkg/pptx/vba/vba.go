// Package vba provides types and helpers for embedding VBA macro projects in
// macro-enabled PowerPoint presentations (.pptm).
//
// Ported from ppt-rs src/parts/vba_macro.rs.
//
// A VBA project is stored as a Compound File Binary (CFB) blob at
// ppt/vbaProject.bin inside the OOXML package. gopptx treats this blob as
// opaque: it can preserve an existing blob from an opened .pptm file or embed
// a caller-supplied blob verbatim. Synthesising CFB bytecode from VBA source
// text is out of scope — callers must supply a pre-built .bin produced by
// Microsoft Office or a compatible tool. The package exposes CFB inspection
// helpers backed by github.com/richardlehane/mscfb for validation and stream
// inventory.
package vba

import "fmt"

// VBAModuleType categorises a VBA module the same way the VBA IDE does.
//
//nolint:revive // Public type name kept for compatibility with existing consumers.
type VBAModuleType int

const (
	// ModuleTypeStandard is a normal code module (.bas).
	ModuleTypeStandard VBAModuleType = iota
	// ModuleTypeClass is a class module (.cls).
	ModuleTypeClass
	// ModuleTypeForm is a UserForm module (.frm).
	ModuleTypeForm
	// ModuleTypeDocument is the document/presentation module (ThisPresentation).
	ModuleTypeDocument
)

// String returns a human-readable name for the module type.
func (t VBAModuleType) String() string {
	switch t {
	case ModuleTypeStandard:
		return "Standard"
	case ModuleTypeClass:
		return "Class"
	case ModuleTypeForm:
		return "Form"
	case ModuleTypeDocument:
		return "Document"
	default:
		return "Unknown"
	}
}

// VBAModule holds metadata about a single VBA module.
// The code is stored as plain text; it is the caller's responsibility to
// compile or merge it into the vbaProject.bin blob.
//
//nolint:revive // Public type name kept for compatibility with existing consumers.
type VBAModule struct {
	// Name is the module name as it appears in the VBA IDE.
	Name string
	// Code holds the VBA source text (informational only; not compiled by gopptx).
	Code string
	// Type is the module type.
	Type VBAModuleType
}

// NewModule creates a standard (.bas) VBA module.
func NewModule(name, code string) VBAModule {
	return VBAModule{Name: name, Code: code, Type: ModuleTypeStandard}
}

// NewClassModule creates a class (.cls) VBA module.
func NewClassModule(name, code string) VBAModule {
	return VBAModule{Name: name, Code: code, Type: ModuleTypeClass}
}

// WithType returns a copy of the module with the given type.
func (m VBAModule) WithType(t VBAModuleType) VBAModule {
	m.Type = t
	return m
}

// VBAProject represents an embedded VBA project in a presentation.
// Either Data (a pre-built vbaProject.bin blob) or Modules (informational
// metadata) — or both — may be set.
//
//nolint:revive // Public type name kept for compatibility with existing consumers.
type VBAProject struct {
	// Data is the raw CFB binary for ppt/vbaProject.bin.
	// When non-empty it is embedded verbatim in the OOXML package.
	Data []byte
	// Modules holds parsed/informational module metadata.
	// gopptx does NOT compile these into Data automatically.
	Modules []VBAModule
}

// New creates an empty VBAProject ready to receive modules or a binary blob.
func New() *VBAProject {
	return &VBAProject{}
}

// FromData creates a VBAProject that wraps an existing vbaProject.bin blob.
func FromData(data []byte) *VBAProject {
	return &VBAProject{Data: append([]byte(nil), data...)}
}

// AddModule appends a module and returns the project for chaining.
func (p *VBAProject) AddModule(m VBAModule) *VBAProject {
	p.Modules = append(p.Modules, m)
	return p
}

// SetData replaces the binary blob.
func (p *VBAProject) SetData(data []byte) *VBAProject {
	p.Data = append([]byte(nil), data...)
	return p
}

// IsMacroEnabled reports whether the project carries a binary VBA blob that
// can be written to ppt/vbaProject.bin.
func (p *VBAProject) IsMacroEnabled() bool {
	if p == nil {
		return false
	}
	return len(p.Data) > 0
}

// PackagePath is the canonical OPC path of the VBA project part.
const PackagePath = "ppt/vbaProject.bin"

// ContentType is the MIME type registered in [Content_Types].xml.
const ContentType = "application/vnd.ms-office.vbaProject"

// RelationshipType is the relationship type used in ppt/_rels/presentation.xml.rels.
const RelationshipType = "http://schemas.microsoft.com/office/2006/relationships/vbaProject"

// FileExtension is the macro-enabled presentation file extension.
const FileExtension = "pptm"

// Validate returns an error if the project is in an invalid state.
func (p *VBAProject) Validate() error {
	if p == nil {
		return nil
	}
	if len(p.Data) > 0 {
		if _, err := p.InspectCFB(); err != nil {
			return fmt.Errorf("invalid vbaProject.bin CFB data: %w", err)
		}
	}
	for i, m := range p.Modules {
		if m.Name == "" {
			return fmt.Errorf("vba module at index %d has an empty name", i)
		}
	}
	return nil
}
