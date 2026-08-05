package mermaid

import "strings"

func parseMindmap(code string) *MindmapNode {
	lines := strings.Split(code, "\n")
	var root *MindmapNode
	var stack []*MindmapNode

	for _, line := range lines {
		_, cleanLine, indent, ok := parseMindmapLine(line)
		if !ok {
			continue
		}

		// An "::icon(...)" line decorates the node above it rather than
		// declaring one, so it must not become a node of its own.
		if strings.HasPrefix(cleanLine, "::icon") {
			continue
		}

		// Parse node label and shape
		label, shape := parseMindmapNodeDef(cleanLine)

		node := &MindmapNode{
			Label: strings.TrimSpace(label),
			Level: indent,
			Shape: shape,
		}

		root, stack = appendMindmapNode(root, stack, node)
	}

	return root
}

// parseMindmapNodeDef reads a mindmap node, including the two bracket forms the
// shared flowchart parser has no equivalent for: `id))bang((` and `id)cloud(`.
// Both used to fall through and keep their brackets in the label.
func parseMindmapNodeDef(line string) (string, NodeShape) {
	if label, ok := cutMindmapBrackets(line, "))", "(("); ok {
		return label, NodeShapeMindmapBang
	}
	if label, ok := cutMindmapBrackets(line, ")", "("); ok {
		return label, NodeShapeMindmapCloud
	}
	_, label, shape := ParseNodeDef(line)
	if strings.Contains(label, "::icon") {
		label = strings.Split(label, "::icon")[0]
	}
	return strings.TrimSpace(label), shape
}

// cutMindmapBrackets pulls the label out of an inverted bracket pair, where the
// closing glyph opens the label and the opening glyph closes it.
func cutMindmapBrackets(line, opener, closer string) (string, bool) {
	start := strings.Index(line, opener)
	if start < 0 {
		return "", false
	}
	end := strings.LastIndex(line, closer)
	if end <= start+len(opener)-1 {
		return "", false
	}
	return strings.TrimSpace(line[start+len(opener) : end]), true
}

func parseMindmapLine(line string) (string, string, int, bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" || strings.HasPrefix(strings.TrimSpace(trimmed), "%%") {
		return "", "", 0, false
	}

	cleanLine := strings.TrimSpace(trimmed)
	if strings.EqualFold(cleanLine, "mindmap") {
		return "", "", 0, false
	}
	return trimmed, cleanLine, leadingIndent(line), true
}

func leadingIndent(line string) int {
	indent := 0
	for _, char := range line {
		switch char {
		case ' ':
			indent++
		case '\t':
			indent += 4
		default:
			return indent
		}
	}
	return indent
}

func appendMindmapNode(root *MindmapNode, stack []*MindmapNode, node *MindmapNode) (*MindmapNode, []*MindmapNode) {
	if root == nil {
		return node, []*MindmapNode{node}
	}

	stack = popMindmapParents(stack, node.Level)
	if len(stack) > 0 {
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, node)
	}
	stack = append(stack, node)
	return root, stack
}

func popMindmapParents(stack []*MindmapNode, indent int) []*MindmapNode {
	for len(stack) > 0 && stack[len(stack)-1].Level >= indent {
		stack = stack[:len(stack)-1]
	}
	return stack
}
