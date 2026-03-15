package protection

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unicode/utf16"
)

const (
	cfbSectorSize      = 512
	cfbDirEntrySize    = 128
	cfbMinRegularBytes = 4096

	cfbEndOfChain = 0xFFFFFFFE
	cfbFreeSect   = 0xFFFFFFFF
	cfbFatSect    = 0xFFFFFFFD
	cfbNoStream   = 0xFFFFFFFF
)

func ensureMinCFBStreamSize(data []byte) []byte {
	if len(data) < cfbMinRegularBytes {
		return append(data, bytes.Repeat([]byte{0}, cfbMinRegularBytes-len(data))...)
	}
	return data
}

func buildCompoundFile(infoStream, pkgStream []byte, infoName, pkgName string) ([]byte, error) {
	infoSectors := sectorsNeeded(len(infoStream))
	pkgSectors := sectorsNeeded(len(pkgStream))
	baseSectors := 1 + infoSectors + pkgSectors
	fatSectors := calcFatSectors(baseSectors)
	if fatSectors > 109 {
		return nil, fmt.Errorf("compound file too large for header DIFAT: fat sectors=%d", fatSectors)
	}

	totalSectors := baseSectors + fatSectors
	fatEntries := make([]uint32, totalSectors)
	for i := range fatEntries {
		fatEntries[i] = cfbFreeSect
	}

	dirSector := 0
	infoStart := 1
	pkgStart := infoStart + infoSectors
	fatStart := baseSectors

	fatEntries[dirSector] = cfbEndOfChain
	linkChain(fatEntries, infoStart, infoSectors)
	linkChain(fatEntries, pkgStart, pkgSectors)
	for i := range fatSectors {
		fatEntries[fatStart+i] = cfbFatSect
	}

	header := buildCFBHeader(fatSectors, dirSector, fatStart)
	infoStartU32, err := checkedUint32FromInt(infoStart)
	if err != nil {
		return nil, err
	}
	pkgStartU32, err := checkedUint32FromInt(pkgStart)
	if err != nil {
		return nil, err
	}
	dirBytes, err := buildDirectorySector(
		infoName,
		pkgName,
		infoStartU32,
		pkgStartU32,
		uint64(len(infoStream)),
		uint64(len(pkgStream)),
	)
	if err != nil {
		return nil, err
	}
	fatBytes := buildFATSectors(fatEntries, fatSectors)

	var out bytes.Buffer
	out.Write(header)
	out.Write(padToSector(dirBytes))
	out.Write(padToSector(infoStream))
	out.Write(padToSector(pkgStream))
	out.Write(fatBytes)
	return out.Bytes(), nil
}

func sectorsNeeded(size int) int {
	if size <= 0 {
		return 1
	}
	return (size + cfbSectorSize - 1) / cfbSectorSize
}

func calcFatSectors(base int) int {
	fat := 1
	for {
		total := base + fat
		need := (total + 127) / 128
		if need == fat {
			return fat
		}
		fat = need
	}
}

func linkChain(fat []uint32, start, count int) {
	if count <= 0 {
		return
	}
	for i := range count - 1 {
		nextSector, err := checkedUint32FromInt(start + i + 1)
		if err != nil {
			panic(err)
		}
		fat[start+i] = nextSector
	}
	fat[start+count-1] = cfbEndOfChain
}

func padToSector(data []byte) []byte {
	if len(data)%cfbSectorSize == 0 {
		return data
	}
	pad := cfbSectorSize - (len(data) % cfbSectorSize)
	return append(data, bytes.Repeat([]byte{0}, pad)...)
}

func buildCFBHeader(numFatSectors, firstDirSector, fatStart int) []byte {
	h := make([]byte, cfbSectorSize)
	copy(h[0:8], []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1})
	binary.LittleEndian.PutUint16(h[24:26], 0x003E)
	binary.LittleEndian.PutUint16(h[26:28], 0x0003)
	binary.LittleEndian.PutUint16(h[28:30], 0xFFFE)
	binary.LittleEndian.PutUint16(h[30:32], 0x0009)
	binary.LittleEndian.PutUint16(h[32:34], 0x0006)
	binary.LittleEndian.PutUint32(h[40:44], mustUint32FromInt(numFatSectors))
	binary.LittleEndian.PutUint32(h[44:48], mustUint32FromInt(firstDirSector))
	binary.LittleEndian.PutUint32(h[56:60], cfbMinRegularBytes)
	binary.LittleEndian.PutUint32(h[60:64], cfbEndOfChain)
	binary.LittleEndian.PutUint32(h[64:68], 0)
	binary.LittleEndian.PutUint32(h[68:72], cfbEndOfChain)
	binary.LittleEndian.PutUint32(h[72:76], 0)
	for i := range 109 {
		binary.LittleEndian.PutUint32(h[76+i*4:80+i*4], cfbFreeSect)
	}
	for i := range min(numFatSectors, 109) {
		binary.LittleEndian.PutUint32(h[76+i*4:80+i*4], mustUint32FromInt(fatStart+i))
	}
	return h
}

func buildDirectorySector(
	infoName, pkgName string,
	infoStart, pkgStart uint32,
	infoSize, pkgSize uint64,
) ([]byte, error) {
	dir := make([]byte, cfbSectorSize)
	// Directory children are represented as a red-black tree.
	// Use EncryptionInfo (entry 2) as root child and EncryptedPackage (entry 1)
	// as its left sibling so both streams are discoverable.
	root, err := newDirEntry("Root Entry", 5, cfbNoStream, cfbNoStream, 2, cfbEndOfChain, 0)
	if err != nil {
		return nil, err
	}
	pkg, err := newDirEntry(pkgName, 2, cfbNoStream, cfbNoStream, cfbNoStream, pkgStart, pkgSize)
	if err != nil {
		return nil, err
	}
	info, err := newDirEntry(infoName, 2, 1, cfbNoStream, cfbNoStream, infoStart, infoSize)
	if err != nil {
		return nil, err
	}
	copy(dir[0:cfbDirEntrySize], root)
	copy(dir[cfbDirEntrySize:cfbDirEntrySize*2], pkg)
	copy(dir[cfbDirEntrySize*2:cfbDirEntrySize*3], info)
	return dir, nil
}

func newDirEntry(
	name string,
	objType byte,
	leftSibling, rightSibling, child uint32,
	startSector uint32,
	size uint64,
) ([]byte, error) {
	if len(name) == 0 {
		return nil, errors.New("directory entry name cannot be empty")
	}
	entry := make([]byte, cfbDirEntrySize)
	u16 := utf16.Encode([]rune(name + "\x00"))
	nameLenBytes := len(u16) * 2
	if nameLenBytes > 64 {
		return nil, fmt.Errorf("directory entry name too long: %q", name)
	}
	for i, v := range u16 {
		binary.LittleEndian.PutUint16(entry[i*2:], v)
	}
	binary.LittleEndian.PutUint16(entry[64:66], mustUint16FromInt(nameLenBytes))
	entry[66] = objType
	entry[67] = 1
	binary.LittleEndian.PutUint32(entry[68:72], leftSibling)
	binary.LittleEndian.PutUint32(entry[72:76], rightSibling)
	binary.LittleEndian.PutUint32(entry[76:80], child)
	binary.LittleEndian.PutUint32(entry[116:120], startSector)
	if size > math.MaxUint32 {
		binary.LittleEndian.PutUint64(entry[120:128], size)
	} else {
		binary.LittleEndian.PutUint32(entry[120:124], uint32(size))
	}
	return entry, nil
}

func buildFATSectors(entries []uint32, fatSectors int) []byte {
	out := make([]byte, fatSectors*cfbSectorSize)
	for i := range fatSectors * 128 {
		v := uint32(cfbFreeSect)
		if i < len(entries) {
			v = entries[i]
		}
		binary.LittleEndian.PutUint32(out[i*4:], v)
	}
	return out
}

func checkedUint32FromInt(v int) (uint32, error) {
	if v < 0 || v > math.MaxUint32 {
		return 0, fmt.Errorf("int value out of uint32 range: %d", v)
	}
	return uint32(v), nil
}

func mustUint32FromInt(v int) uint32 {
	out, err := checkedUint32FromInt(v)
	if err != nil {
		panic(err)
	}
	return out
}

func mustUint16FromInt(v int) uint16 {
	if v < 0 || v > math.MaxUint16 {
		panic(fmt.Sprintf("int value out of uint16 range: %d", v))
	}
	return uint16(v)
}
