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

		// Parse node label and shape
		_, label, shape := ParseNodeDef(cleanLine)

		// Handle icons or other mindmap specific syntax (simplified)
		if strings.Contains(label, "::icon") {
			label = strings.Split(label, "::icon")[0]
		}

		node := &MindmapNode{
			Label: strings.TrimSpace(label),
			Level: indent,
			Shape: shape,
		}

		root, stack = appendMindmapNode(root, stack, node)
	}

	return root
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
