package slide

import (
	"crypto/sha256"
	"encoding/hex"
	"path"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func ParseMediaInventory(ps PartLookup, partKeys []string) (map[string]string, int) {
	inventory := make(map[string]string)
	maxNum := 0
	for _, partPath := range partKeys {
		if !strings.HasPrefix(partPath, "ppt/media/") {
			continue
		}
		data, ok := ps.Get(partPath)
		if !ok {
			continue
		}
		hash := sha256.Sum256(data)
		inventory[hex.EncodeToString(hash[:])] = partPath

		num, ok := parseImagePartNumber(partPath)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return inventory, maxNum + 1
}

func ParseChartInventory(ps PartLookup, partKeys []string) (map[string]string, int, int) {
	inventory := make(map[string]string)
	maxChart := 0
	maxExcel := 0

	for _, p := range partKeys {
		if !isChartPartPath(p) {
			continue
		}
		num, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(p, "ppt/charts/chart"), ".xml"))
		if num > maxChart {
			maxChart = num
		}

		excelPath, nextExcel := findChartEmbedding(ps, p, maxExcel)
		if excelPath == "" {
			continue
		}
		inventory[p] = excelPath
		maxExcel = nextExcel
	}
	return inventory, maxChart + 1, maxExcel + 1
}

func ParseNotesInventory(ps PartLookup, partKeys []string) (map[string]string, int) {
	inventory := make(map[string]string)
	maxNotes := 0

	for _, p := range partKeys {
		if !strings.HasPrefix(p, "ppt/slides/_rels/slide") {
			continue
		}
		if !strings.HasSuffix(p, ".xml.rels") {
			continue
		}
		slidePart := "ppt/slides/" + strings.TrimSuffix(path.Base(p), ".rels")
		relsData, ok := ps.Get(p)
		if !ok {
			continue
		}
		rels, _ := ParseRelationshipsXML(relsData)
		for _, r := range rels {
			if r.Type == "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" {
				notesPath := common.CanonicalPartPath(path.Join("ppt/slides", r.Target))
				inventory[slidePart] = notesPath

				num, _ := strconv.Atoi(
					strings.TrimSuffix(strings.TrimPrefix(path.Base(notesPath), "notesSlide"), ".xml"),
				)
				if num > maxNotes {
					maxNotes = num
				}
			}
		}
	}
	return inventory, maxNotes + 1
}

// ParseDiagramInventory scans part keys for SmartArt data files and returns
// the next available diagram number (max existing + 1, at least 1).
func ParseDiagramInventory(partKeys []string) int {
	maxNum := 0
	for _, p := range partKeys {
		if !strings.HasPrefix(p, "ppt/diagrams/data") || !strings.HasSuffix(p, ".xml") {
			continue
		}
		base := strings.TrimSuffix(strings.TrimPrefix(p, "ppt/diagrams/data"), ".xml")
		num, err := strconv.Atoi(base)
		if err != nil {
			continue
		}
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

func parseImagePartNumber(partPath string) (int, bool) {
	base := path.Base(partPath)
	if !strings.HasPrefix(base, "image") {
		return 0, false
	}
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)
	numStr := strings.TrimPrefix(name, "image")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, false
	}
	return num, true
}

func isChartPartPath(partPath string) bool {
	return strings.HasPrefix(partPath, "ppt/charts/chart") && strings.HasSuffix(partPath, ".xml")
}

func findChartEmbedding(ps PartLookup, chartPart string, currentMaxExcel int) (string, int) {
	relsPath := "ppt/charts/_rels/" + path.Base(chartPart) + ".rels"
	relsData, ok := ps.Get(relsPath)
	if !ok {
		return "", currentMaxExcel
	}

	rels, _ := ParseRelationshipsXML(relsData)
	maxExcel := currentMaxExcel
	for _, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}
		excelPath := common.CanonicalPartPath(path.Join("ppt/charts", rel.Target))
		maxExcel = maxExcelNumber(maxExcel, excelPath)
		return excelPath, maxExcel
	}
	return "", maxExcel
}

func maxExcelNumber(current int, excelPath string) int {
	base := path.Base(excelPath)
	after, ok := strings.CutPrefix(base, "Microsoft_Excel_Worksheet")
	if !ok {
		return current
	}

	enum, _ := strconv.Atoi(strings.TrimSuffix(after, ".xlsx"))
	if enum > current {
		return enum
	}
	return current
}
