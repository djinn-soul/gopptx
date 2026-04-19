package mermaid

import "strings"

func classRelationshipTypes() []string {
	return []string{"<|--", "*--", "o--", "-->", "--", "..>", "..", "<|..", "*..", "o.."}
}

// renderClass parses and renders a Mermaid class diagram into PowerPoint elements.
func renderClass(code string, theme Theme) DiagramElements {
	diagram := parseClass(code)
	return generateClassElements(diagram, theme)
}

func parseClass(code string) *ClassDiagram {
	lines := ParseLines(code)
	classes := make(map[string]*ClassNode)
	var relationships []ClassRelationship
	var currentClass *ClassNode

	for index, line := range lines {
		currentClass = parseClassLine(line, index, classes, currentClass, &relationships)
	}

	classList := make([]ClassNode, 0, len(classes))
	for _, class := range classes {
		classList = append(classList, *class)
	}

	return &ClassDiagram{
		Classes:       classList,
		Relationships: relationships,
	}
}

func parseClassLine(
	line string,
	lineIndex int,
	classes map[string]*ClassNode,
	currentClass *ClassNode,
	relationships *[]ClassRelationship,
) *ClassNode {
	if shouldSkipClassLine(line, lineIndex) {
		return currentClass
	}
	if className, ok := classBlockStart(line); ok {
		return ensureClassNode(classes, className)
	}
	if line == "}" {
		return nil
	}
	if currentClass != nil {
		appendClassMember(currentClass, line)
		return currentClass
	}
	if className, member, ok := parseClassInlineMember(line); ok {
		appendClassMember(ensureClassNode(classes, className), member)
		return currentClass
	}
	if rel, ok := parseClassRelationshipLine(line); ok {
		*relationships = append(*relationships, rel)
		ensureClassNode(classes, rel.From)
		ensureClassNode(classes, rel.To)
		return currentClass
	}
	if className, ok := parseSimpleClassDefinition(line); ok {
		ensureClassNode(classes, className)
	}
	return currentClass
}

func shouldSkipClassLine(line string, lineIndex int) bool {
	if strings.HasPrefix(line, "classDiagram") {
		return true
	}
	return strings.HasPrefix(line, "class ") &&
		!strings.Contains(line, "{") &&
		!strings.Contains(line, ":") &&
		lineIndex == 0
}

func classBlockStart(line string) (string, bool) {
	if !strings.HasPrefix(line, "class ") || !strings.HasSuffix(line, "{") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(line, "{"), "class"))
	return name, name != ""
}

func ensureClassNode(classes map[string]*ClassNode, name string) *ClassNode {
	if _, ok := classes[name]; !ok {
		classes[name] = &ClassNode{ID: name, Name: name}
	}
	return classes[name]
}

func appendClassMember(class *ClassNode, member string) {
	trimmed := strings.TrimSpace(member)
	if strings.Contains(trimmed, "(") {
		class.Methods = append(class.Methods, trimmed)
		return
	}
	class.Attributes = append(class.Attributes, trimmed)
}

func parseClassInlineMember(line string) (string, string, bool) {
	if !strings.Contains(line, ":") {
		return "", "", false
	}
	parts := strings.SplitN(line, ":", 2)
	className := strings.TrimSpace(parts[0])
	member := strings.TrimSpace(parts[1])
	return className, member, className != "" && member != ""
}

func parseClassRelationshipLine(line string) (ClassRelationship, bool) {
	from, relType, to, label, found := splitClassRelationship(line)
	if !found {
		return ClassRelationship{}, false
	}
	return ClassRelationship{From: from, To: to, Type: relType, Label: label}, true
}

func parseSimpleClassDefinition(line string) (string, bool) {
	if !strings.HasPrefix(line, "class ") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(line, "class"))
	return name, name != ""
}

func splitClassRelationship(line string) (string, string, string, string, bool) {
	for _, relType := range classRelationshipTypes() {
		before, after, ok := strings.Cut(line, relType)
		if !ok {
			continue
		}
		from := strings.TrimSpace(before)
		rest := strings.TrimSpace(after)
		to := rest
		label := ""
		if beforeLabel, afterLabel, ok := strings.Cut(rest, ":"); ok {
			to = strings.TrimSpace(beforeLabel)
			label = strings.TrimSpace(afterLabel)
		}
		return from, relType, to, label, true
	}
	return "", "", "", "", false
}
