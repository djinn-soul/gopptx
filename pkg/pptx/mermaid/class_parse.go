package mermaid

import "strings"

// classDecorationChars are the glyphs Mermaid decorates a class relation's
// ends with: inheritance (<, |, >), composition (*) and aggregation (o).
const classDecorationChars = "<|>*o"

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
	// Relations are tested before inline members. Both can carry a colon, and
	// checking the member form first turned `Animal <|-- Duck : label` into a
	// member of a class literally named `Animal <|-- Duck`.
	if rel, ok := parseClassRelationshipLine(line); ok {
		*relationships = append(*relationships, rel)
		ensureClassNode(classes, rel.From)
		ensureClassNode(classes, rel.To)
		return currentClass
	}
	if className, member, ok := parseClassInlineMember(line); ok {
		appendClassMember(ensureClassNode(classes, className), member)
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

// splitClassRelationship pulls `Animal <|-- Duck : label` apart into its
// classes, relation token and label.
//
// The relation was matched against an ordered list of tokens, and the order was
// wrong: `..` sat before `<|..`, so `Animal <|.. Duck` split at the dots and
// left a class called `Animal <|`. Finding the connector and growing outwards
// over the decoration glyphs accepts every combination without an order to get
// wrong.
func splitClassRelationship(line string) (string, string, string, string, bool) {
	head, label, hasLabel := strings.Cut(line, ":")
	start, end, ok := classRelationSpan(head)
	if !ok {
		return "", "", "", "", false
	}
	from := trimClassCardinality(head[:start])
	to := trimClassCardinality(head[end:])
	if from == "" || to == "" {
		return "", "", "", "", false
	}
	if !hasLabel {
		label = ""
	}
	return from, head[start:end], to, strings.TrimSpace(label), true
}

// classRelationSpan locates the relation token: the `--` or `..` connector plus
// the decoration glyphs on either side of it.
func classRelationSpan(head string) (int, int, bool) {
	connector := strings.Index(head, "--")
	if connector < 0 {
		connector = strings.Index(head, "..")
	}
	if connector < 0 {
		return 0, 0, false
	}
	start := connector
	for start > 0 && isClassDecoration(head, start-1, start-2) {
		start--
	}
	end := connector + 2
	for end < len(head) && isClassDecoration(head, end, end+1) {
		end++
	}
	return start, end, true
}

// isClassDecoration reports whether head[at] decorates the relation rather than
// belonging to a class name. `o` is the ambiguous one — it ends the aggregation
// marker but also plenty of identifiers — so it only counts when the character
// beyond it is not part of a name.
func isClassDecoration(head string, at, beyond int) bool {
	c := head[at]
	if !strings.ContainsRune(classDecorationChars, rune(c)) {
		return false
	}
	if c != 'o' {
		return true
	}
	if beyond < 0 || beyond >= len(head) {
		return true
	}
	return !isIdentifierByte(head[beyond])
}

func isIdentifierByte(c byte) bool {
	switch {
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '_':
		return true
	default:
		return false
	}
}

// trimClassCardinality drops the quoted multiplicity Mermaid allows beside a
// class, so `Customer "1"` names the class Customer.
func trimClassCardinality(part string) string {
	trimmed := strings.TrimSpace(part)
	if before, _, ok := strings.Cut(trimmed, `"`); ok {
		// A leading quote means the cardinality came first: "1" Order.
		if strings.TrimSpace(before) == "" {
			if _, after, ok := strings.Cut(strings.TrimPrefix(trimmed, `"`), `"`); ok {
				return strings.TrimSpace(after)
			}
		}
		return strings.TrimSpace(before)
	}
	return trimmed
}
