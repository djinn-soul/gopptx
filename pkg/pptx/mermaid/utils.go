package mermaid

import (
	"regexp"
	"strings"
)

var (
	themeInitRegex      = regexp.MustCompile(`(?i)%%\{init:\s*\{.*['"]theme['"]:\s*['"]([^'"]+)['"].*\}\s*\}%%`)
	themeDirectiveRegex = regexp.MustCompile(`(?i)^\s*theme:\s*(\w+)`)
)

// DetectTheme identifies the theme from the Mermaid code.
func DetectTheme(code string) string {
	// Check for init block first
	if matches := themeInitRegex.FindStringSubmatch(code); len(matches) > 1 {
		return strings.ToLower(matches[1])
	}

	// Check line by line for theme directive
	lines := strings.SplitSeq(code, "\n")
	for line := range lines {
		trimmed := strings.TrimSpace(line)
		if matches := themeDirectiveRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			return strings.ToLower(matches[1])
		}
	}

	return "default"
}

// ParseLines splits the code into non-empty, trimmed lines, ignoring comments.
func ParseLines(code string) []string {
	lines := strings.Split(code, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		// Remove inline comments
		if idx := strings.Index(trimmed, "%%"); idx != -1 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ExtractDirection parses the layout direction from a Mermaid flowchart header.
func ExtractDirection(header string) FlowDirection {
	upper := strings.ToUpper(header)
	if strings.Contains(upper, "LR") {
		return FlowDirectionLR
	}
	if strings.Contains(upper, "RL") {
		return FlowDirectionRL
	}
	if strings.Contains(upper, "BT") {
		return FlowDirectionBT
	}
	return FlowDirectionTB
}

// SplitConnection splits a line at a connection arrow and returns (from, arrow, rest, found).
// Handles patterns like: "A --> B", "A -->|label| B", "A -- label --> B".
func SplitConnection(line string) (string, string, string, bool) {
	arrows := []string{"==>", "-.->", "-->", "---", "->"}
	for _, arrow := range arrows {
		if before, after, ok := strings.Cut(line, arrow); ok {
			from := strings.TrimSpace(before)
			rest := strings.TrimSpace(after)
			return from, arrow, rest, true
		}
	}
	return "", "", "", false
}

// ExtractArrowLabel extracts a label from an arrow part like "|label| rest" or "label --> rest".
// Returns (label, rest) where label is empty if no label found.
func ExtractArrowLabel(s string) (string, string) {
	// First check for the standard Mermaid syntax: |label|
	if strings.HasPrefix(s, "|") {
		if endIdx := strings.Index(s[1:], "|"); endIdx != -1 {
			label := s[1 : endIdx+1]
			rest := strings.TrimSpace(s[endIdx+2:])
			return label, rest
		}
	}

	// Handle the alternative syntax: label --> target
	// Split by the next arrow to extract label and remaining node
	arrows := []string{"==>", "-.->", "-->", "---", "->"}
	for _, arrow := range arrows {
		if before, after, ok := strings.Cut(s, arrow); ok {
			label := strings.TrimSpace(before)
			rest := strings.TrimSpace(after)
			if label != "" {
				return label, rest
			}
		}
	}

	return "", s
}

// ParseNodeDef parses a node definition like "A[Text]" and returns (id, label, shape).
func ParseNodeDef(s string) (string, string, NodeShape) {
	s = strings.TrimSpace(s)

	bracketTypes := []struct {
		open  string
		close string
		shape NodeShape
	}{
		{"(([", "]))", NodeShapeStadium}, // Stadium (alternative)
		{"((", "))", NodeShapeCircle},
		{"([", "])", NodeShapeStadium},
		{"{{", "}}", NodeShapeHexagon},
		{"[", "]", NodeShapeRectangle},
		{"(", ")", NodeShapeRoundedRect},
		{"{", "}", NodeShapeDiamond},
	}

	for _, bt := range bracketTypes {
		if startIdx := strings.Index(s, bt.open); startIdx != -1 {
			id := strings.TrimSpace(s[:startIdx])
			if endIdx := strings.LastIndex(s, bt.close); endIdx != -1 && endIdx > startIdx {
				label := s[startIdx+len(bt.open) : endIdx]
				return id, label, bt.shape
			}
		}
	}

	// Plain node ID
	return s, s, NodeShapeRectangle
}
