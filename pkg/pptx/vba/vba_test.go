package vba

import (
	"bytes"
	"testing"
)

func TestVBAModuleType_String(t *testing.T) {
	tests := []struct {
		name string
		t    VBAModuleType
		want string
	}{
		{"Standard", ModuleTypeStandard, "Standard"},
		{"Class", ModuleTypeClass, "Class"},
		{"Form", ModuleTypeForm, "Form"},
		{"Document", ModuleTypeDocument, "Document"},
		{"Unknown", VBAModuleType(999), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("VBAModuleType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewModule(t *testing.T) {
	m := NewModule("Module1", "Sub Test()\nEnd Sub")
	if m.Name != "Module1" {
		t.Errorf("NewModule() Name = %v, want Module1", m.Name)
	}
	if m.Type != ModuleTypeStandard {
		t.Errorf("NewModule() Type = %v, want ModuleTypeStandard", m.Type)
	}
}

func TestNewClassModule(t *testing.T) {
	m := NewClassModule("MyClass", "Private x As Integer")
	if m.Name != "MyClass" {
		t.Errorf("NewClassModule() Name = %v, want MyClass", m.Name)
	}
	if m.Type != ModuleTypeClass {
		t.Errorf("NewClassModule() Type = %v, want ModuleTypeClass", m.Type)
	}
}

func TestVBAModule_WithType(t *testing.T) {
	m := NewModule("Module1", "").WithType(ModuleTypeDocument)
	if m.Type != ModuleTypeDocument {
		t.Errorf("WithType() Type = %v, want ModuleTypeDocument", m.Type)
	}
}

func TestVBAProject_New(t *testing.T) {
	p := New()
	if p.IsMacroEnabled() {
		t.Errorf("New() IsMacroEnabled = %v, want false", p.IsMacroEnabled())
	}
}

func TestVBAProject_AddModule(t *testing.T) {
	p := New().
		AddModule(NewModule("Module1", "Sub Test()\nEnd Sub")).
		AddModule(NewClassModule("Class1", ""))

	if len(p.Modules) != 2 {
		t.Errorf("AddModule() len int = %v, want 2", len(p.Modules))
	}
	if p.IsMacroEnabled() {
		t.Errorf("IsMacroEnabled on modules = %v, want false", p.IsMacroEnabled())
	}
}

func TestVBAProject_FromData(t *testing.T) {
	blob := []byte{0x00, 0x01, 0x02}
	p := FromData(blob)
	if !p.IsMacroEnabled() {
		t.Errorf("IsMacroEnabled on data = %v, want true", p.IsMacroEnabled())
	}
	if !bytes.Equal(p.Data, blob) {
		t.Errorf("FromData() Data = %v, want %v", p.Data, blob)
	}

	// Test SetData
	blob2 := []byte{0xff}
	p.SetData(blob2)
	if !bytes.Equal(p.Data, blob2) {
		t.Errorf("SetData() Data = %v, want %v", p.Data, blob2)
	}
}

func TestVBAProject_Validate(t *testing.T) {
	p := New().AddModule(NewModule("Valid", ""))
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}

	p.AddModule(NewModule("", ""))
	if err := p.Validate(); err == nil {
		t.Error("Validate() error = nil, want error on empty module name")
	}
}
