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

type (
	LayoutMasterResolverFn  func(layoutPart string) (string, error)
	MasterLayoutsResolverFn func(masterPart string) ([]string, error)
)

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

// DefaultSlideMaster returns a basic slide master XML.
func DefaultSlideMaster() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:theme="http://schemas.openxmlformats.org/drawingml/2006/theme">
  <p:cSld>
    <p:bg>
      <p:bgRef idx="10000">
        <a:schemeClr val="bg1"/>
      </p:bgRef>
    </p:bg>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr/>
    </p:spTree>
    <p:txStyles>
      <p:titleStyle>
        <a:lvl1pPr algn="ltr" fontSz="32" kerning="1200">
          <a:defRPr sz="3200" kerning="1200">
            <a:solidFill>
              <a:schemeClr val="tx1"/>
            </a:solidFill>
            <a:latin typeface="+mj-lt"/>
          </a:defRPr>
        </a:lvl1pPr>
      </p:titleStyle>
      <p:bodyStyle>
        <a:lvl1pPr algn="ltr" fontSz="18" kerning="1200">
          <a:defRPr sz="1800" kerning="1200">
            <a:solidFill>
              <a:schemeClr val="tx1"/>
            </a:solidFill>
            <a:latin typeface="+mj-lt"/>
          </a:defRPr>
        </a:lvl1pPr>
      </p:bodyStyle>
    </p:txStyles>
  </p:cSld>
  <p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
  <p:transitions/>
</p:sldMaster>`
}

// DefaultSlideMasterRelationships returns default relationships for a slide master.
func DefaultSlideMasterRelationships() string {
	// Link to first layout (slideLayout1.xml) - the first layout will be created automatically.
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="slideLayout1.xml"/>
</Relationships>`
}

// DefaultSlideLayout returns a basic slide layout XML.
func DefaultSlideLayout(layoutName string, layoutNum, masterNum int) string {
	_ = layoutNum // suppress unused parameter warning
	_ = masterNum // suppress unused parameter warning
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" preserve="1">
  <p:cSld name="%s">
    <p:bg>
      <p:bgRef idx="10000">
        <a:schemeClr val="bg1"/>
      </p:bgRef>
    </p:bg>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr/>
    </p:spTree>
  </p:cSld>
  <p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
  <p:sldLayoutPr preserve="1"/>
</p:sldLayout>`, layoutName)
}

// DefaultSlideLayoutRelationships returns default relationships for a slide layout.
func DefaultSlideLayoutRelationships(masterNum int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster%d.xml"/>
</Relationships>`, masterNum)
}
