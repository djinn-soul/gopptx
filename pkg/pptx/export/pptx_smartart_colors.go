package export

import (
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// The native renderer draws SmartArt from the semantic node tree, which says
// nothing about colour, so every diagram used to come out in the same accent
// regardless of the colour style it asked for. Resolving the style against the
// deck's theme here gives each node the fill PowerPoint would paint it.

const smartArtAccentSlots = 6

// resolveSmartArtNodeColors fills in the colour each node is drawn with, leaving
// a node that carries its own colour alone.
func resolveSmartArtNodeColors(nodes []smartart.Node, colorStyleID string, theme deckTheme) []smartart.Node {
	slots := smartArtAccentSlotsFor(colorStyleID)
	if len(slots) == 0 {
		return nodes
	}
	next := 0
	var walk func(items []smartart.Node) []smartart.Node
	walk = func(items []smartart.Node) []smartart.Node {
		out := make([]smartart.Node, 0, len(items))
		for _, node := range items {
			if node.Color == "" {
				if hex, ok := theme.themeColor(slots[next%len(slots)]); ok {
					node.Color = hex
				}
			}
			next++
			node.Children = walk(node.Children)
			out = append(out, node)
		}
		return out
	}
	return walk(nodes)
}

// smartArtAccentSlotsFor turns a colour style into the theme slots it cycles
// through: one accent for the accentN families, all six for the colourful ones.
func smartArtAccentSlotsFor(colorStyleID string) []string {
	name := colorStyleID[strings.LastIndex(colorStyleID, "/")+1:]
	switch {
	case strings.HasPrefix(name, "colorful"):
		slots := make([]string, 0, smartArtAccentSlots)
		for i := 1; i <= smartArtAccentSlots; i++ {
			slots = append(slots, "accent"+strconv.Itoa(i))
		}
		return slots
	case strings.HasPrefix(name, "accent"):
		digits, _, _ := strings.Cut(strings.TrimPrefix(name, "accent"), "_")
		n, err := strconv.Atoi(digits)
		if err != nil || n < 1 || n > smartArtAccentSlots {
			return []string{themeColorAccent1}
		}
		return []string{"accent" + strconv.Itoa(n)}
	case name == "":
		return nil
	default:
		return []string{themeColorAccent1}
	}
}
