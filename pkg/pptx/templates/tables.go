package templates

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

// PricingTier defines a structured pricing item.
type PricingTier struct {
	Name     string
	Price    string
	Features []string
}

// Milestone defines a structured timeline item.
type Milestone struct {
	Date   string
	Task   string
	Status string
}

func renderPricingTable(tiers []PricingTier) tables.Table {
	rows := make([][]string, 0, len(tiers))
	for _, tier := range tiers {
		features := strings.Join(tier.Features, ", ")
		rows = append(rows, []string{tier.Name, tier.Price, features})
	}

	return renderTable(
		[]styling.Length{styling.Inches(2.0), styling.Inches(1.5), styling.Inches(4.5)},
		[]string{"Tier", "Price", "Features"},
		rows,
	)
}

func renderTimelineTable(milestones []Milestone) tables.Table {
	rows := make([][]string, 0, len(milestones))
	for _, milestone := range milestones {
		rows = append(rows, []string{milestone.Date, milestone.Task, milestone.Status})
	}

	return renderTable(
		[]styling.Length{styling.Inches(2.0), styling.Inches(4.0), styling.Inches(2.0)},
		[]string{"Date", "Milestone", "Status"},
		rows,
	)
}

func renderTable(colWidths []styling.Length, headers []string, rows [][]string) tables.Table {
	table := tables.NewTable(colWidths)
	headerCells := make([]tables.TableCell, 0, len(headers))
	for _, header := range headers {
		headerCells = append(headerCells, tables.TableCell{
			Text:            header,
			Bold:            true,
			BackgroundColor: "1F497D",
			Color:           "FFFFFF",
			RowSpan:         1,
			ColSpan:         1,
		})
	}

	table = table.AddStyledRow(headerCells)
	for _, row := range rows {
		table = table.AddRow(row)
	}

	return table
}
