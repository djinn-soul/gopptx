package slide

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"path"
	"regexp"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func NextMasterPartPath(nextMasterNum int) string {
	return fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", nextMasterNum)
}

func BuildLayoutCloneMap(layoutFamily []string, nextLayoutNum int) map[string]string {
	layoutMap := make(map[string]string, len(layoutFamily))
	for _, oldLayout := range layoutFamily {
		layoutMap[oldLayout] = fmt.Sprintf("ppt/slideLayouts/slideLayout%d.xml", nextLayoutNum)
		nextLayoutNum++
	}
	return layoutMap
}

func CloneResultTheme(themePart, newThemePart string) string {
	if newThemePart != "" {
		return newThemePart
	}
	return themePart
}

func NextPartNumber(parts []string, pattern *regexp.Regexp, submatchSize int) int {
	maxNum := 0
	for _, part := range parts {
		base := path.Base(part)
		m := pattern.FindStringSubmatch(base)
		if len(m) != submatchSize {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > maxNum {
			maxNum = n
		}
	}
	return maxNum + 1
}

func ParseLayoutName(layoutXML []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(layoutXML))
	for {
		token, err := decoder.Token()
		if err != nil {
			return ""
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Local == "sldLayout" || start.Name.Local == "cSld" {
			for _, attr := range start.Attr {
				if attr.Name.Local == "name" {
					return attr.Value
				}
			}
		}
	}
}

type LayoutMasterResolverFn func(layoutPart string) (string, error)
type MasterLayoutsResolverFn func(masterPart string) ([]string, error)

func CloneFamilyInputs(
	layoutPart string,
	hasPart func(string) bool,
	canonicalPartPath func(string) string,
	resolveLayoutMaster LayoutMasterResolverFn,
	resolveMasterLayouts MasterLayoutsResolverFn,
) (string, []string, error) {
	layoutPart = canonicalPartPath(layoutPart)
	if !hasPart(layoutPart) {
		return "", nil, fmt.Errorf("layout part %s not found", layoutPart)
	}
	sourceMaster, err := resolveLayoutMaster(layoutPart)
	if err != nil {
		return "", nil, err
	}
	layoutFamily, err := resolveMasterLayouts(sourceMaster)
	if err != nil {
		return "", nil, err
	}
	if len(layoutFamily) == 0 {
		return "", nil, fmt.Errorf("no layouts found for master %s", sourceMaster)
	}
	return sourceMaster, layoutFamily, nil
}

func ResolveLayoutMasterPart(
	layoutPart string,
	getPart func(string) ([]byte, bool),
	parseRelationships func([]byte) ([]common.EditorRelationship, error),
) (string, error) {
	relsPath := common.RelsPathFor(layoutPart)
	relsData, ok := getPart(relsPath)
	if !ok {
		return "", fmt.Errorf("layout rels part not found: %s", relsPath)
	}
	rels, err := parseRelationships(relsData)
	if err != nil {
		return "", fmt.Errorf("parse layout rels: %w", err)
	}
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideMaster {
			return common.CanonicalPartPath(path.Join(path.Dir(layoutPart), rel.Target)), nil
		}
	}
	return "", fmt.Errorf("layout %s has no slideMaster relationship", layoutPart)
}
