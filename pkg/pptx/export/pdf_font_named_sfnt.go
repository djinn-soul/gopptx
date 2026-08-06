package export

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

// Reading a font file well enough to index it: its family name, its subfamily
// (which gives the bold/italic style), and nothing else. Only the header, the
// table directory and the name table are read, so indexing a font directory of
// several hundred faces stays cheap.

// sfnt name-table layout.
const (
	nameTableCountOffset        = 2
	nameTableStringOffsetOffset = 4
	nameTableRecordsOffset      = 6
	nameRecordSize              = 12

	nameIDFamily    = 1
	nameIDSubfamily = 2

	namePlatformMac     = 1
	namePlatformWindows = 3

	// Cap on how much of a font file is read while indexing. A name table is a
	// few kilobytes; anything larger is malformed or not worth indexing.
	maxNameTableBytes = 1 << 16

	// fontHeadBufferBytes covers the sfnt header and table directory in one
	// read. A directory of 64 tables needs 1036 bytes, so this has ample room.
	fontHeadBufferBytes = 4096
)

var errFontNameUnreadable = errors.New("font name table unreadable")

// readTTFFamilyAndStyle returns a font file's family name and the gopdf style
// bitmask its subfamily describes.
func readTTFFamilyAndStyle(path string) (string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = file.Close() }()

	// One read covers the header and the table directory for every font in
	// practice; indexing a font directory of several hundred files is otherwise
	// dominated by per-table syscalls.
	head := make([]byte, fontHeadBufferBytes)
	n, err := file.ReadAt(head, 0)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", 0, err
	}
	head = head[:n]

	base, err := sfntFontOffset(head)
	if err != nil {
		return "", 0, err
	}
	offset, length, err := sfntTableExtent(head, base, "name")
	if err != nil {
		return "", 0, err
	}
	if length > maxNameTableBytes {
		length = maxNameTableBytes
	}
	table := make([]byte, length)
	if _, err := file.ReadAt(table, int64(offset)); err != nil && !errors.Is(err, io.EOF) {
		return "", 0, err
	}

	family := sfntName(table, nameIDFamily)
	subfamily := sfntName(table, nameIDSubfamily)
	return family, styleFromSubfamily(subfamily), nil
}

// sfntFontOffset returns the byte offset of the font whose tables should be
// read. For a TrueType Collection that is its first font, which is also the face
// gopdf embeds.
func sfntFontOffset(head []byte) (uint32, error) {
	if len(head) < sfntHeaderSize {
		return 0, errFontNameUnreadable
	}
	if string(head[0:4]) != "ttcf" {
		return 0, nil
	}
	if len(head) < ttcFirstFontOffset+4 {
		return 0, errFontNameUnreadable
	}
	return binary.BigEndian.Uint32(head[ttcFirstFontOffset:]), nil
}

// sfntTableExtent locates one table's offset and length in the font's table
// directory, which the caller has already read into head.
func sfntTableExtent(head []byte, base uint32, tag string) (uint32, uint32, error) {
	if uint64(base)+sfntHeaderSize > uint64(len(head)) {
		return 0, 0, errFontNameUnreadable
	}
	numTables := binary.BigEndian.Uint16(head[base+4:])
	for i := range uint32(numTables) {
		at := uint64(base) + sfntHeaderSize + uint64(i)*sfntTableRecordSize
		if at+sfntTableRecordSize > uint64(len(head)) {
			return 0, 0, errFontNameUnreadable
		}
		record := head[at : at+sfntTableRecordSize]
		if string(record[0:4]) != tag {
			continue
		}
		return binary.BigEndian.Uint32(record[8:12]), binary.BigEndian.Uint32(record[12:16]), nil
	}
	return 0, 0, errFontNameUnreadable
}

// sfntName pulls one name-table string, preferring the Windows platform's
// UTF-16 record and falling back to the Mac platform's single-byte one.
func sfntName(table []byte, nameID uint16) string {
	if len(table) < nameTableRecordsOffset {
		return ""
	}
	count := int(binary.BigEndian.Uint16(table[nameTableCountOffset:]))
	stringBase := int(binary.BigEndian.Uint16(table[nameTableStringOffsetOffset:]))

	best := ""
	for i := range count {
		record := nameTableRecordsOffset + i*nameRecordSize
		if record+nameRecordSize > len(table) {
			break
		}
		platform := binary.BigEndian.Uint16(table[record:])
		if binary.BigEndian.Uint16(table[record+6:]) != nameID {
			continue
		}
		length := int(binary.BigEndian.Uint16(table[record+8:]))
		offset := stringBase + int(binary.BigEndian.Uint16(table[record+10:]))
		if offset < 0 || offset+length > len(table) {
			continue
		}
		value := decodeSFNTNameValue(table[offset:offset+length], platform)
		if value == "" {
			continue
		}
		if platform == namePlatformWindows {
			return value
		}
		if best == "" {
			best = value
		}
	}
	return best
}

func decodeSFNTNameValue(raw []byte, platform uint16) string {
	if platform != namePlatformWindows {
		if platform != namePlatformMac {
			return ""
		}
		return strings.TrimSpace(string(raw))
	}
	// Windows name records are UTF-16BE. Only the Basic Multilingual Plane is
	// decoded, which covers every family name in practice.
	var b strings.Builder
	for i := 0; i+1 < len(raw); i += 2 {
		b.WriteRune(rune(binary.BigEndian.Uint16(raw[i:])))
	}
	return strings.TrimSpace(b.String())
}

// styleFromSubfamily maps an sfnt subfamily ("Bold Italic", "Oblique", …) to the
// gopdf style bitmask.
func styleFromSubfamily(subfamily string) int {
	normalized := strings.ToLower(subfamily)
	style := gopdf.Regular
	if strings.Contains(normalized, "bold") {
		style |= gopdf.Bold
	}
	if strings.Contains(normalized, "italic") || strings.Contains(normalized, "oblique") {
		style |= gopdf.Italic
	}
	return style
}
