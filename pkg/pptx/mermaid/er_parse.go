package mermaid

import "strings"

// erCardinalityChars are the glyphs Mermaid builds crow's-foot cardinalities
// from: one (|), zero (o) and many (} or {).
const erCardinalityChars = "|o}{"

// renderER parses and renders a Mermaid ER diagram into PowerPoint elements.
func renderER(code string, theme Theme) DiagramElements {
	diagram := parseER(code)
	return generateERElements(diagram, theme)
}

func parseER(code string) *ERDiagram {
	lines := ParseLines(code)
	entities := make(map[string]*EREntity)
	var relationships []ERRelationship
	var currentEntity *EREntity

	for _, line := range lines {
		currentEntity = parseERLine(line, entities, currentEntity, &relationships)
	}

	entityList := make([]EREntity, 0, len(entities))
	for _, entity := range entities {
		entityList = append(entityList, *entity)
	}

	return &ERDiagram{
		Entities:      entityList,
		Relationships: relationships,
	}
}

func parseERLine(
	line string,
	entities map[string]*EREntity,
	currentEntity *EREntity,
	relationships *[]ERRelationship,
) *EREntity {
	switch {
	case strings.HasPrefix(line, "erDiagram"):
		return currentEntity
	case line == "}":
		return nil
	}

	if entityName, ok := parseEREntityBlockStart(line); ok {
		return ensureEREntity(entities, entityName)
	}
	if currentEntity != nil {
		currentEntity.Attributes = append(currentEntity.Attributes, strings.TrimSpace(line))
		return currentEntity
	}
	if rel, ok := parseERRelationshipLine(line); ok {
		*relationships = append(*relationships, rel)
		ensureEREntity(entities, rel.From)
		ensureEREntity(entities, rel.To)
		return currentEntity
	}
	if entityName, ok := parseERSimpleEntity(line); ok {
		ensureEREntity(entities, entityName)
	}
	return currentEntity
}

func parseEREntityBlockStart(line string) (string, bool) {
	before, ok := strings.CutSuffix(line, "{")
	if !ok {
		return "", false
	}
	name := strings.TrimSpace(before)
	return name, name != ""
}

func ensureEREntity(entities map[string]*EREntity, name string) *EREntity {
	if _, ok := entities[name]; !ok {
		entities[name] = &EREntity{Name: name}
	}
	return entities[name]
}

func parseERRelationshipLine(line string) (ERRelationship, bool) {
	from, relType, to, label, found := splitERRelationship(line)
	if !found {
		return ERRelationship{}, false
	}
	return ERRelationship{From: from, To: to, Type: relType, Label: label}, true
}

func parseERSimpleEntity(line string) (string, bool) {
	if strings.Contains(line, " ") || strings.Contains(line, "-") {
		return "", false
	}
	name := strings.TrimSpace(line)
	return name, name != ""
}

// splitERRelationship pulls `CUSTOMER ||--o{ ORDER : places` apart into its
// entities, relation token and label.
//
// The relation used to be matched against a hand-written list of the solid
// `--` forms, so every non-identifying `..` relation — `}|..|{`, `||..o{` —
// missed and the line was dropped along with any entity only named there.
// Finding the connector and growing outwards over the cardinality glyphs
// accepts each combination of the two notations instead.
func splitERRelationship(line string) (string, string, string, string, bool) {
	// A label follows the colon and may itself contain dots, so only the part
	// before it can hold the relation.
	head, label, hasLabel := strings.Cut(line, ":")
	start, end, ok := erRelationSpan(head)
	if !ok {
		return "", "", "", "", false
	}
	from := strings.TrimSpace(head[:start])
	to := strings.TrimSpace(head[end:])
	if from == "" || to == "" {
		return "", "", "", "", false
	}
	if !hasLabel {
		label = ""
	}
	return from, head[start:end], to, strings.TrimSpace(label), true
}

// erRelationSpan locates the relation token: the `--` or `..` connector plus
// the cardinality glyphs on either side of it.
func erRelationSpan(head string) (int, int, bool) {
	connector := strings.Index(head, "--")
	if connector < 0 {
		connector = strings.Index(head, "..")
	}
	if connector < 0 {
		return 0, 0, false
	}
	start := connector
	for start > 0 && strings.ContainsRune(erCardinalityChars, rune(head[start-1])) {
		start--
	}
	end := connector + 2
	for end < len(head) && strings.ContainsRune(erCardinalityChars, rune(head[end])) {
		end++
	}
	return start, end, true
}
