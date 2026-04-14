package mermaid

import "strings"

func erRelationshipTypes() []string {
	return []string{
		"||--o{",
		"||--|{",
		"}|--|{",
		"}|--o{",
		"|o--o{",
		"|o--|{",
		"o{--}o",
		"o{--|{",
		"o{--o{",
		"||--||",
		"||--|o",
		"|o--|o",
		"|o--||",
		"||--",
		"}|--",
		"o{--",
		"--o{",
		"--|{",
		"--}o",
		"--",
	}
}

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

func splitERRelationship(line string) (string, string, string, string, bool) {
	for _, relType := range erRelationshipTypes() {
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
