package goppt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"golang.org/x/text/encoding/unicode"
)

func encodeRecord(rt recordType, payload []byte) []byte {
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint16(header[2:4], uint16(rt))
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(payload)))
	return append(header, payload...)
}

func TestRecordAndRecordDataHelpers(t *testing.T) {
	if got := recordTypeTextCharsAtom.LowerPart(); got != byte(recordTypeTextCharsAtom&0xFF) {
		t.Fatalf("LowerPart mismatch: got=%d", got)
	}

	payload := []byte{1, 2, 3, 4}
	raw := encodeRecord(recordTypeSlide, payload)
	rr := bytes.NewReader(raw)

	rec, err := readRecordHeaderOnly(rr, 0, recordTypeSlide)
	if err != nil {
		t.Fatalf("readRecordHeaderOnly failed: %v", err)
	}
	if rec.Type() != recordTypeSlide || rec.Length() != 4 {
		t.Fatalf("unexpected record header parse: type=%v len=%d", rec.Type(), rec.Length())
	}

	full, err := readRecord(rr, 0, recordTypeSlide)
	if err != nil {
		t.Fatalf("readRecord failed: %v", err)
	}
	if !bytes.Equal(full.Data(), payload) {
		t.Fatalf("unexpected record payload: %v", full.Data())
	}

	if _, err = readRecordHeaderOnly(rr, 0, recordTypeDocument); !errors.Is(err, errMismatchRecordType) {
		t.Fatalf("expected errMismatchRecordType, got %v", err)
	}

	var out [2]byte
	n, err := full.recordData.ReadAt(out[:], 1)
	if err != nil || n != 2 || out[0] != 2 || out[1] != 3 {
		t.Fatalf("recordData.ReadAt unexpected result n=%d out=%v err=%v", n, out, err)
	}
	if got := full.recordData.LongAt(0); got != 0x04030201 {
		t.Fatalf("recordData.LongAt(0)=%#x, want %#x", got, uint32(0x04030201))
	}
	if got := full.recordData.LongAt(2); got != 0 {
		t.Fatalf("recordData.LongAt out-of-range should return 0, got %#x", got)
	}
}

func TestTextDecodersAndPocketScanner(t *testing.T) {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	// "A" + "B" encoded as TextBytesAtom lower-byte form.
	textBytes, err := decodeTextBytesAtom([]byte{'A', 'B'}, decoder)
	if err != nil {
		t.Fatalf("decodeTextBytesAtom failed: %v", err)
	}
	if string(textBytes) != "AB" {
		t.Fatalf("decodeTextBytesAtom output=%q, want AB", string(textBytes))
	}

	var sb strings.Builder
	atom := record{
		recordData: []byte{0x41, 0x00, 0x42, 0x00}, // UTF-16LE "AB"
		header:     [headerSize]byte{},
	}
	if err = readTextFromTextCharsAtom(atom, &sb, decoder); err != nil {
		t.Fatalf("readTextFromTextCharsAtom failed: %v", err)
	}
	if sb.String() != "AB " {
		t.Fatalf("readTextFromTextCharsAtom output=%q, want %q", sb.String(), "AB ")
	}

	sb.Reset()
	atom = record{recordData: []byte{'C', 'D'}}
	if err = readTextFromTextBytesAtom(atom, &sb, decoder); err != nil {
		t.Fatalf("readTextFromTextBytesAtom failed: %v", err)
	}
	if sb.String() != "CD " {
		t.Fatalf("readTextFromTextBytesAtom output=%q, want %q", sb.String(), "CD ")
	}

	// 0xA0 0x0F and 0xA8 0x0F are pocket signatures (TextChars/TextBytes).
	if idx := matchPocket([]byte{0x00, 0xA0, 0x0F, 0x00}, 0); idx != 1 {
		t.Fatalf("matchPocket chars idx=%d, want 1", idx)
	}
	if idx := matchPocket([]byte{0x00, 0xA8, 0x0F, 0x00}, 0); idx != 1 {
		t.Fatalf("matchPocket bytes idx=%d, want 1", idx)
	}
	if idx := matchPocket([]byte{0x00, 0x01, 0x02}, 0); idx != -1 {
		t.Fatalf("matchPocket no-match idx=%d, want -1", idx)
	}
}

func TestSkipRecords(t *testing.T) {
	// Stream with UserEditAtom then PersistDirectoryAtom then Slide.
	buf := append(
		encodeRecord(recordTypeUserEditAtom, []byte{1}),
		encodeRecord(recordTypePersistDirectoryAtom, []byte{2})...,
	)
	buf = append(buf, encodeRecord(recordTypeSlide, []byte{3, 4})...)

	offset, err := skipRecords(
		bytes.NewReader(buf),
		0,
		[]recordType{recordTypeUserEditAtom, recordTypePersistDirectoryAtom},
	)
	if err != nil {
		t.Fatalf("skipRecords failed: %v", err)
	}
	expected := int64((headerSize + 1) + (headerSize + 1))
	if offset != expected {
		t.Fatalf("skipRecords offset=%d, want %d", offset, expected)
	}

	// If first expected record type mismatches, skipRecords keeps scanning the sequence and continues.
	offset, err = skipRecords(bytes.NewReader(buf), 0, []recordType{recordTypeSlide})
	if err != nil {
		t.Fatalf("skipRecords mismatch case failed: %v", err)
	}
	if offset != 0 {
		t.Fatalf("skipRecords mismatch case should leave offset unchanged, got %d", offset)
	}
}
