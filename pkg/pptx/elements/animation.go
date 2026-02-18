package elements

import (
	"fmt"
	"strings"
)

const (
	twoColumnDivisor = 2
	animIDStride     = 2
	animIDOffset     = 3
)

// CalculateShapeIDs replicates the ID generation order in pptxxml.SlideWithLayout.
func CalculateShapeIDs(s SlideContent) []int {
	// Shape IDs start at 2 (Title = 2).
	nextID := 2

	// Title
	if s.Layout != SlideLayoutBlank {
		nextID++
	}

	// Table or Content
	if s.Table != nil {
		nextID++
	} else if len(s.Bullets) > 0 || len(s.BulletRuns) > 0 {
		nextID++
		if s.Layout == SlideLayoutTwoColumn {
			leftCount := (len(s.Bullets) + 1) / twoColumnDivisor
			if len(s.Bullets[leftCount:]) > 0 {
				nextID++
			}
		}
	}

	// Primary chart object occupies one shape slot.
	if hasPrimaryChart(s) {
		nextID++
	}

	// Calculate base IDs for each pool
	imageStartID := nextID
	nextID += len(s.Images)

	shapeStartID := nextID
	nextID += len(s.Shapes)

	connectorStartID := nextID
	nextID += len(s.Connectors)

	placeholderStartID := nextID

	// Combine IDs into a flat slice for indexing by Animation.ShapeIndex.
	// Order: Shapes (1..N), Connectors (N+1..M), Images (M+1..P), Placeholders (P+1..Q).
	totalCount := len(s.Shapes) + len(s.Connectors) + len(s.Images) + len(s.PlaceholderOverrides)
	allIDs := make([]int, 0, totalCount)

	for i := range s.Shapes {
		allIDs = append(allIDs, shapeStartID+i)
	}
	for i := range s.Connectors {
		allIDs = append(allIDs, connectorStartID+i)
	}
	for i := range s.Images {
		allIDs = append(allIDs, imageStartID+i)
	}
	for i := range s.PlaceholderOverrides {
		allIDs = append(allIDs, placeholderStartID+i)
	}

	return allIDs
}

func hasPrimaryChart(s SlideContent) bool {
	return s.Chart != nil ||
		s.BarHorizontal != nil ||
		s.BarStacked != nil ||
		s.BarStacked100 != nil ||
		s.Line != nil ||
		s.LineMarkers != nil ||
		s.LineStacked != nil ||
		s.Scatter != nil ||
		s.Area != nil ||
		s.AreaStacked != nil ||
		s.AreaStacked100 != nil ||
		s.Pie != nil ||
		s.Doughnut != nil ||
		s.Bubble != nil ||
		s.Radar != nil ||
		s.RadarFilled != nil ||
		s.StockHLC != nil ||
		s.StockOHLC != nil ||
		s.Combo != nil
}

func SlideAnimationsXML(s SlideContent, shapeIDs []int) string {
	if len(s.Animations) == 0 {
		return ""
	}

	animationsXML := make([]string, len(s.Animations))
	for i, anim := range s.Animations {
		actualID := 0
		if anim.ShapeIndex > 0 && anim.ShapeIndex <= len(shapeIDs) {
			actualID = shapeIDs[anim.ShapeIndex-1]
		}
		if actualID == 0 {
			continue
		}
		animationsXML[i] = anim.XML(i*animIDStride+animIDOffset, actualID)
	}

	var finalXML []string
	for _, xml := range animationsXML {
		if xml != "" {
			finalXML = append(finalXML, xml)
		}
	}

	if len(finalXML) == 0 {
		return ""
	}

	return fmt.Sprintf(`
<p:timing>
  <p:tnLst>
    <p:par>
      <p:cTn id="1" dur="indefinite" restart="never" nodeType="tmRoot">
        <p:childTnLst>
          <p:seq concurrent="1" nextAc="seek">
            <p:cTn id="2" dur="indefinite" nodeType="mainSeq">
              <p:childTnLst>
                %s
              </p:childTnLst>
            </p:cTn>
          </p:seq>
        </p:childTnLst>
      </p:cTn>
    </p:par>
  </p:tnLst>
</p:timing>`,
		strings.Join(finalXML, "\n"),
	)
}
