package goppt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/richardlehane/mscfb"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"

	"github.com/djinn-soul/gopptx/internal/ioadapters"
)

const (
	userPersistIDRefOffset = 16
	utf16WordBytes         = 2
	persistIDMask          = 0x000FFFFF
	persistCountMask       = 0xFFF00000
	persistCountShift      = 20
	persistCountValueMask  = 0x00000FFF
	recordTypeOffset       = 2
)

// skippedSlideRecordTypes returns metadata or non-readable records in slide container.
func skippedSlideRecordTypes() []recordType {
	return []recordType{
		recordTypeExternalObjectList,
		recordTypeEnvironment,
		recordTypeSoundCollection,
		recordTypeDrawingGroup,
		recordTypeSlideListWithText,
		recordTypeList,
		recordTypeHeadersFooters,
	}
}

// skippedDrawingRecordTypes returns metadata or non-readable records in drawing container.
func skippedDrawingRecordTypes() []recordType {
	return []recordType{
		recordTypeSlideShowSlideInfoAtom,
		recordTypeHeadersFooters,
		recordTypeRoundTripSlideSyncInfo12,
	}
}

// ExtractText parses PPT file represented by Reader r and extracts text from it.
func ExtractText(r io.Reader) (string, error) {
	ra := ioadapters.ToReaderAt(r)

	d, err := mscfb.New(ra)
	if err != nil {
		return "", err
	}
	currentUser, pptDocument := getCurrentUserAndPPTDoc(d)
	if validErr := isValidPPT(currentUser, pptDocument); validErr != nil {
		return "", validErr
	}
	offsetPersistDirectory, liveRecord, err := getUserEditAtomsData(currentUser, pptDocument)
	if err != nil {
		return "", err
	}
	persistDirEntries, err := getPersistDirectoryEntries(pptDocument, offsetPersistDirectory)
	if err != nil {
		return "", err
	}

	// get DocumentContainer reference
	docPersistIDRef := liveRecord.LongAt(userPersistIDRefOffset)
	documentContainer, err := readRecord(pptDocument, persistDirEntries[docPersistIDRef], recordTypeDocument)
	if err != nil {
		return "", err
	}

	return readSlides(documentContainer, pptDocument, persistDirEntries)
}

// getCurrentUserAndPPTDoc extracts necessary mscfb files from PPT file.
func getCurrentUserAndPPTDoc(r *mscfb.Reader) (*mscfb.File, *mscfb.File) {
	var currentUser *mscfb.File
	var pptDocument *mscfb.File

	for _, f := range r.File {
		switch f.Name {
		case "Current User":
			currentUser = f
		case "PowerPoint Document":
			pptDocument = f
		}
	}
	return currentUser, pptDocument
}

// isValidPPT checks if provided file is valid, meaning
// it has both "Current User" and "PowerPoint Document" files
// and "Current User"'s CurrentUserAtom record has valid header token.
func isValidPPT(currentUser, pptDocument *mscfb.File) error {
	const (
		headerTokenOffset      = 12
		encryptedDocumentToken = 0xF3D1C4DF
		plainDocumentToken     = 0xE391C05F
	)

	if currentUser == nil || pptDocument == nil {
		return errors.New(".ppt file must contain \"Current User\" and \"PowerPoint Document\" streams")
	}
	var b [4]byte
	_, err := currentUser.ReadAt(b[:], headerTokenOffset)
	if err != nil {
		return err
	}
	headerToken := binary.LittleEndian.Uint32(b[:])
	if headerToken != plainDocumentToken && headerToken != encryptedDocumentToken {
		return fmt.Errorf("invalid UserEditAtom header token %X", headerToken)
	}
	return nil
}

// getUserEditAtomsData extracts "live record" and persist directory offsets
// according to section 2.1.2 of specification.
func getUserEditAtomsData(currentUser, pptDocument *mscfb.File) ([]int64, record, error) {
	const (
		offsetLastEditInitialPosition  = 16
		offsetLastEditPosition         = 8
		persistDirectoryOffsetPosition = 12
	)
	var persistDirectoryOffsets []int64
	var liveRecord record

	var b [4]byte
	_, err := currentUser.ReadAt(b[:], offsetLastEditInitialPosition)
	if err != nil {
		return nil, record{}, err
	}
	offsetLastEdit := binary.LittleEndian.Uint32(b[:])

	for {
		liveRecord, err = readRecord(pptDocument, int64(offsetLastEdit), recordTypeUserEditAtom)
		if err != nil {
			if errors.Is(err, errMismatchRecordType) {
				break
			}
			return nil, record{}, err
		}
		persistDirectoryOffsets = append(
			persistDirectoryOffsets,
			int64(liveRecord.LongAt(persistDirectoryOffsetPosition)),
		)
		offsetLastEdit = liveRecord.LongAt(offsetLastEditPosition)
		if offsetLastEdit == 0 {
			break
		}
	}

	return persistDirectoryOffsets, liveRecord, nil
}

// getPersistDirectoryEntries transforms offsets into persists directory identifiers and persist offsets.
func getPersistDirectoryEntries(pptDocument *mscfb.File, offsets []int64) (map[uint32]int64, error) {
	const persistOffsetEntrySize = 4

	persistDirEntries := make(map[uint32]int64)
	for _, v := range slices.Backward(offsets) {
		rgPersistDirEntry, err := readRecord(pptDocument, v, recordTypePersistDirectoryAtom)
		if err != nil {
			return nil, err
		}

		rgPersistDirEntryData := rgPersistDirEntry.recordData

		for j := 0; j < len(rgPersistDirEntryData); {
			persist := rgPersistDirEntryData.LongAt(j)
			persistID := persist & persistIDMask
			cPersist := ((persist & persistCountMask) >> persistCountShift) & persistCountValueMask
			j += 4

			for k := range cPersist {
				persistDirEntries[persistID+k] = int64(rgPersistDirEntryData.LongAt(j + int(k)*persistOffsetEntrySize))
			}
			j += int(cPersist * persistOffsetEntrySize)
		}
	}
	return persistDirEntries, nil
}

// readSlides reads text from slides of given DocumentContainer.
func readSlides(documentContainer, pptDocument io.ReaderAt, persistDirEntries map[uint32]int64) (string, error) {
	const slideSkipInitialOffset = 48
	offset, err := skipRecords(documentContainer, slideSkipInitialOffset, skippedSlideRecordTypes())
	if err != nil {
		return "", err
	}
	slideList, err := readRecord(documentContainer, offset, recordTypeSlideListWithText)
	if err != nil {
		return "", err
	}

	utf16Decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	var out strings.Builder
	data := slideList.Data()
	n := len(data)
	for i := 0; i < n; {
		if i+headerSize > n {
			break
		}
		recType := recordType(binary.LittleEndian.Uint16(data[i+2 : i+4]))
		recLen := binary.LittleEndian.Uint32(data[i+4 : i+8])

		if i+headerSize+int(recLen) > n {
			break
		}
		blockData := data[i+headerSize : i+headerSize+int(recLen)]

		// Create a temporary record view for legacy handlers that expect it.
		// We reuse the header slice if possible or just use a local.
		var header [headerSize]byte
		copy(header[:], data[i:i+headerSize])
		block := record{
			recordData: blockData,
			header:     header,
		}

		switch recType {
		case recordTypeSlidePersistAtom:
			err = readTextFromSlidePersistAtom(block, pptDocument, persistDirEntries, &out, utf16Decoder)
		case recordTypeTextCharsAtom:
			err = readTextFromTextCharsAtom(block, &out, utf16Decoder)
		case recordTypeTextBytesAtom:
			err = readTextFromTextBytesAtom(block, &out, utf16Decoder)
		default:
			// Ignore non-text record types within SlideListWithText.
		}
		if err != nil {
			return "", err
		}

		i += int(recLen) + headerSize
	}

	return out.String(), nil
}

func readTextFromSlidePersistAtom(
	block record,
	pptDocument io.ReaderAt,
	persistDirEntries map[uint32]int64,
	out *strings.Builder,
	utf16Decoder *encoding.Decoder,
) error {
	const slidePersistAtomSkipInitialOffset = 32

	persistDirID := block.LongAt(0)
	slide, err := readRecord(pptDocument, persistDirEntries[persistDirID], recordTypeSlide)
	if err != nil {
		return err
	}
	offset, err := skipRecords(slide, slidePersistAtomSkipInitialOffset, skippedDrawingRecordTypes())
	if err != nil {
		return err
	}

	drawing, err := readRecord(slide, offset, recordTypeDrawing)
	if err != nil {
		return err
	}
	return extractTextFromDrawing(drawing, out, utf16Decoder)
}

// skipRecords reads headers and skips data of records of provided types.
func skipRecords(r io.ReaderAt, initialOffset int64, skippedRecordsTypes []recordType) (int64, error) {
	offset := initialOffset

	for i := range skippedRecordsTypes {
		rec, err := readRecordHeaderOnly(r, offset, skippedRecordsTypes[i])
		if err != nil {
			if errors.Is(err, errMismatchRecordType) {
				continue
			}
			return 0, err
		}
		offset += int64(rec.Length() + headerSize)
	}

	return offset, nil
}
