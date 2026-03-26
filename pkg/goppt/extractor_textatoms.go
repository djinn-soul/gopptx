package goppt

import (
	"bytes"
	"encoding/binary"
	"strings"

	"golang.org/x/text/encoding"
)

func extractTextFromDrawing(drawing record, out *strings.Builder, utf16Decoder *encoding.Decoder) error {
	const headerRecordTypeOffset = 2
	drawingBytes := drawing.Data()
	from := 0
	for {
		pocketIdx := matchPocket(drawingBytes, from)
		if pocketIdx == -1 {
			break
		}
		if pocketIdx >= 2 && bytes.Equal(drawingBytes[pocketIdx-headerRecordTypeOffset:pocketIdx], []byte{0x00, 0x00}) {
			if err := processTextRecord(drawing, pocketIdx, out, utf16Decoder); err != nil {
				return err
			}
		}
		from = pocketIdx + recordTypeOffset
	}
	return nil
}

func processTextRecord(drawing record, pocketIdx int, out *strings.Builder, utf16Decoder *encoding.Decoder) error {
	const headerRecordTypeOffset = 2
	data := drawing.Data()
	if pocketIdx+headerSize-headerRecordTypeOffset > len(data) {
		return nil
	}
	recType := recordType(
		binary.LittleEndian.Uint16(
			data[pocketIdx-headerRecordTypeOffset+2 : pocketIdx-headerRecordTypeOffset+4],
		),
	)
	recLen := binary.LittleEndian.Uint32(data[pocketIdx-headerRecordTypeOffset+4 : pocketIdx-headerRecordTypeOffset+8])
	if pocketIdx-headerRecordTypeOffset+headerSize+int(recLen) > len(data) {
		return nil
	}
	blockData := data[pocketIdx-headerRecordTypeOffset+headerSize : pocketIdx-headerRecordTypeOffset+headerSize+int(recLen)]

	var header [headerSize]byte
	copy(header[:], data[pocketIdx-headerRecordTypeOffset:pocketIdx-headerRecordTypeOffset+headerSize])
	block := record{recordData: blockData, header: header}
	if recType == recordTypeTextCharsAtom {
		return readTextFromTextCharsAtom(block, out, utf16Decoder)
	}
	if recType == recordTypeTextBytesAtom {
		return readTextFromTextBytesAtom(block, out, utf16Decoder)
	}
	return nil
}

func matchPocket(data []byte, from int) int {
	data = data[from:]
	n := len(data)
	for i := range n {
		switch data[i] {
		case recordTypeTextCharsAtom.LowerPart(), recordTypeTextBytesAtom.LowerPart():
			if i < n-1 && data[i+1] == 0x0F {
				return i + from
			}
		}
	}
	return -1
}

func readTextFromTextCharsAtom(atom record, out *strings.Builder, dec *encoding.Decoder) error {
	dec.Reset()
	transformed, err := dec.Bytes(atom.Data())
	if err != nil {
		return err
	}
	out.Write(transformed)
	out.WriteByte(' ')
	return nil
}

func readTextFromTextBytesAtom(atom record, out *strings.Builder, dec *encoding.Decoder) error {
	dec.Reset()
	transformed, err := decodeTextBytesAtom(atom.Data(), dec)
	if err != nil {
		return err
	}
	out.Write(transformed)
	out.WriteByte(' ')
	return nil
}

func decodeTextBytesAtom(data []byte, dec *encoding.Decoder) ([]byte, error) {
	utf16Data := make([]byte, len(data)*utf16WordBytes)
	for i, b := range data {
		utf16Data[i*utf16WordBytes] = b
		utf16Data[i*utf16WordBytes+1] = 0
	}
	return dec.Bytes(utf16Data)
}
