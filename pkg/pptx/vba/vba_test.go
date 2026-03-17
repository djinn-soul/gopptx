package vba

import (
	"bytes"
	"encoding/binary"
	"strings"
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

func TestInspectCFB_Invalid(t *testing.T) {
	if _, err := InspectCFB([]byte("not a cfb")); err == nil {
		t.Fatal("InspectCFB() expected error for invalid data")
	}
}

func TestInspectCFB_ValidFixture(t *testing.T) {
	data := minimalValidCFB()
	info, err := InspectCFB(data)
	if err != nil {
		t.Fatalf("InspectCFB(valid fixture) error = %v", err)
	}
	if len(info.Entries) == 0 {
		t.Fatal("expected at least one CFB entry from valid fixture")
	}
}

// minimalValidCFB constructs a minimal Compound File Binary blob (v3, 512-byte sectors)
// containing a root storage and one child stream, which is sufficient for InspectCFB.
func minimalValidCFB() []byte {
	const sectorSize = 512
	const (
		freeSect   = uint32(0xFFFFFFFF)
		endOfChain = uint32(0xFFFFFFFE)
		fatSect    = uint32(0xFFFFFFFD)
		noStream   = uint32(0xFFFFFFFF)
	)

	buf := make([]byte, 4*sectorSize) // header + FAT sector + dir sector + stream sector

	// ---- Header (offset 0-511) ----
	copy(buf[0:], []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) // magic
	binary.LittleEndian.PutUint16(buf[24:], 0x003E)                       // minor version
	binary.LittleEndian.PutUint16(buf[26:], 0x0003)                       // major version 3
	binary.LittleEndian.PutUint16(buf[28:], 0xFFFE)                       // little-endian
	binary.LittleEndian.PutUint16(buf[30:], 9)                            // sector size power (2^9 = 512)
	binary.LittleEndian.PutUint16(buf[32:], 6)                            // mini sector size power (2^6 = 64)
	binary.LittleEndian.PutUint32(buf[44:], 1)                            // 1 FAT sector
	binary.LittleEndian.PutUint32(buf[48:], 1)                            // first directory sector = 1
	binary.LittleEndian.PutUint32(buf[56:], 4096)                         // mini stream cutoff size
	binary.LittleEndian.PutUint32(buf[60:], endOfChain)                   // first mini FAT sector (none)
	binary.LittleEndian.PutUint32(buf[68:], endOfChain)                   // first DIFAT sector (none)
	binary.LittleEndian.PutUint32(buf[76:], 0)                            // DIFAT[0] = sector 0 (FAT)
	for i := 80; i < 512; i += 4 {
		binary.LittleEndian.PutUint32(buf[i:], freeSect)
	}

	// ---- FAT sector (sector 0, offset 512) ----
	fat := buf[sectorSize:]
	binary.LittleEndian.PutUint32(fat[0:], fatSect)    // sector 0 = this FAT sector
	binary.LittleEndian.PutUint32(fat[4:], endOfChain) // sector 1 = dir sector (end)
	binary.LittleEndian.PutUint32(fat[8:], endOfChain) // sector 2 = stream sector (end)
	for i := 12; i < sectorSize; i += 4 {
		binary.LittleEndian.PutUint32(fat[i:], freeSect)
	}

	// ---- Directory sector (sector 1, offset 1024) — 4 entries × 128 bytes ----
	dir := buf[2*sectorSize:]

	// Entry 0: Root Entry
	cfbPutUTF16LE(dir[0:64], "Root Entry")
	binary.LittleEndian.PutUint16(dir[64:], 22)          // name len: 11 chars × 2 bytes
	dir[66] = 5                                          // object type: root storage
	dir[67] = 1                                          // color: black
	binary.LittleEndian.PutUint32(dir[68:], noStream)    // no left sibling
	binary.LittleEndian.PutUint32(dir[72:], noStream)    // no right sibling
	binary.LittleEndian.PutUint32(dir[76:], 1)           // child = entry 1
	binary.LittleEndian.PutUint32(dir[116:], endOfChain) // start sector (no mini stream)

	// Entry 1: "VBA" stream (offset 128 in dir sector)
	e1 := dir[128:]
	cfbPutUTF16LE(e1[0:64], "VBA")
	binary.LittleEndian.PutUint16(e1[64:], 8)           // name len: 4 chars × 2 bytes (V,B,A,\0)
	e1[66] = 2                                          // object type: stream
	e1[67] = 1                                          // color: black
	binary.LittleEndian.PutUint32(e1[68:], noStream)    // no left sibling
	binary.LittleEndian.PutUint32(e1[72:], noStream)    // no right sibling
	binary.LittleEndian.PutUint32(e1[76:], noStream)    // no child
	binary.LittleEndian.PutUint32(e1[116:], 2)          // starting sector = 2
	binary.LittleEndian.PutUint32(e1[120:], sectorSize) // size = 512 bytes

	// Entries 2-3: unallocated (all zeros, type 0x00) — already zeroed by make().

	// Stream data (sector 2, offset 1536): zeros, already zeroed by make().
	return buf
}

// cfbPutUTF16LE writes s as UTF-16-LE into dst, with a null terminator.
func cfbPutUTF16LE(dst []byte, s string) {
	for i, c := range s {
		if i*2+1 >= len(dst) {
			break
		}
		dst[i*2] = byte(c)
		dst[i*2+1] = 0
	}
}

func TestVBAProject_Validate_InvalidCFBData(t *testing.T) {
	p := FromData([]byte("invalid-cfb"))
	err := p.Validate()
	if err == nil {
		t.Fatal("expected Validate() to fail for invalid CFB data")
	}
	if !strings.Contains(err.Error(), "invalid vbaProject.bin CFB data") {
		t.Fatalf("expected CFB validation message, got: %v", err)
	}
}
