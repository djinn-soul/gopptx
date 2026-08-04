package mathml

import (
	"fmt"
	"strings"
	"unicode"
)

// parser walks the LaTeX source once, left to right.
type parser struct {
	src []rune
	pos int
}

// sequenceBuilder accumulates one parsed sequence: finished OMML elements plus
// the literal characters not yet flushed into an <m:r>.
type sequenceBuilder struct {
	out     strings.Builder
	pending strings.Builder
}

// flush turns any buffered literal characters into a run.
func (s *sequenceBuilder) flush() {
	if s.pending.Len() > 0 {
		s.out.WriteString(runXML(s.pending.String()))
		s.pending.Reset()
	}
}

// parseSequence reads runs until the source ends or a closing brace at depth is
// reached, joining them into one OMML fragment.
func (p *parser) parseSequence(depth int) (string, error) {
	var b sequenceBuilder

	for p.pos < len(p.src) {
		ch := p.src[p.pos]
		if ch == '}' {
			if depth == 0 {
				return "", fmt.Errorf("unmatched '}' at position %d", p.pos)
			}
			b.flush()
			return b.out.String(), nil
		}
		if err := p.parseSequenceToken(ch, depth, &b); err != nil {
			return "", err
		}
	}
	if depth != 0 {
		return "", fmt.Errorf("unclosed '{' in %q", string(p.src))
	}
	b.flush()
	return b.out.String(), nil
}

// parseSequenceToken consumes the one construct starting at ch and appends it
// to the builder. Splitting it out of parseSequence keeps that function to the
// loop and the brace bookkeeping.
func (p *parser) parseSequenceToken(ch rune, depth int, b *sequenceBuilder) error {
	switch {
	case ch == '{':
		return p.parseGroupToken(depth, b)
	case ch == '^' || ch == '_':
		return p.parseScriptToken(ch, b)
	case ch == '\\':
		return p.parseCommandToken(b)
	case unicode.IsSpace(ch):
		p.pos++
		return nil
	default:
		b.pending.WriteRune(ch)
		p.pos++
		return nil
	}
}

func (p *parser) parseGroupToken(depth int, b *sequenceBuilder) error {
	p.pos++
	group, err := p.parseSequence(depth + 1)
	if err != nil {
		return err
	}
	if err := p.expect('}'); err != nil {
		return err
	}
	b.flush()
	b.out.WriteString(group)
	return nil
}

func (p *parser) parseScriptToken(marker rune, b *sequenceBuilder) error {
	p.pos++
	operand, err := p.parseOperand()
	if err != nil {
		return err
	}
	base := takeLastElement(&b.out, &b.pending)
	if base == "" {
		return fmt.Errorf("'%c' at position %d has nothing to attach to", marker, p.pos)
	}
	b.flush()
	b.out.WriteString(scriptXML(marker, base, operand))
	return nil
}

func (p *parser) parseCommandToken(b *sequenceBuilder) error {
	b.flush()
	command, err := p.readCommand()
	if err != nil {
		return err
	}
	rendered, err := p.renderCommand(command)
	if err != nil {
		return err
	}
	b.out.WriteString(rendered)
	return nil
}

// parseOperand reads the argument of a script or command: a braced group or a
// single character.
func (p *parser) parseOperand() (string, error) {
	p.skipSpaces()
	if p.pos >= len(p.src) {
		return "", fmt.Errorf("expected an operand at end of %q", string(p.src))
	}
	if p.src[p.pos] == '{' {
		p.pos++
		group, err := p.parseSequence(1)
		if err != nil {
			return "", err
		}
		if err := p.expect('}'); err != nil {
			return "", err
		}
		return group, nil
	}
	if p.src[p.pos] == '\\' {
		command, err := p.readCommand()
		if err != nil {
			return "", err
		}
		return p.renderCommand(command)
	}
	ch := string(p.src[p.pos])
	p.pos++
	return runXML(ch), nil
}

func (p *parser) renderCommand(command string) (string, error) {
	if letter, ok := greekLetters[command]; ok {
		return runXML(letter), nil
	}
	if symbol, ok := symbols[command]; ok {
		return runXML(symbol), nil
	}
	if operator, ok := naryOperators[command]; ok {
		return p.parseNary(operator)
	}
	switch command {
	case "frac":
		return p.parseFraction()
	case "sqrt":
		return p.parseRadical()
	default:
		return "", fmt.Errorf("unsupported LaTeX command \\%s", command)
	}
}

func (p *parser) parseFraction() (string, error) {
	numerator, err := p.parseOperand()
	if err != nil {
		return "", fmt.Errorf("\\frac numerator: %w", err)
	}
	denominator, err := p.parseOperand()
	if err != nil {
		return "", fmt.Errorf("\\frac denominator: %w", err)
	}
	return `<m:f><m:fPr><m:ctrlPr/></m:fPr>` +
		`<m:num>` + numerator + `</m:num>` +
		`<m:den>` + denominator + `</m:den></m:f>`, nil
}

func (p *parser) parseRadical() (string, error) {
	degree := ""
	p.skipSpaces()
	if p.pos < len(p.src) && p.src[p.pos] == '[' {
		p.pos++
		var raw strings.Builder
		for p.pos < len(p.src) && p.src[p.pos] != ']' {
			raw.WriteRune(p.src[p.pos])
			p.pos++
		}
		if err := p.expect(']'); err != nil {
			return "", err
		}
		degree = runXML(raw.String())
	}
	radicand, err := p.parseOperand()
	if err != nil {
		return "", fmt.Errorf("\\sqrt radicand: %w", err)
	}
	hideDegree := "1"
	if degree != "" {
		hideDegree = "0"
	}
	return `<m:rad><m:radPr><m:degHide m:val="` + hideDegree + `"/><m:ctrlPr/></m:radPr>` +
		`<m:deg>` + degree + `</m:deg>` +
		`<m:e>` + radicand + `</m:e></m:rad>`, nil
}

// parseNary reads an n-ary operator with its optional limits and its operand.
func (p *parser) parseNary(operator string) (string, error) {
	sub, sup := "", ""
	for {
		p.skipSpaces()
		if p.pos >= len(p.src) {
			break
		}
		marker := p.src[p.pos]
		if marker != '_' && marker != '^' {
			break
		}
		p.pos++
		operand, err := p.parseOperand()
		if err != nil {
			return "", err
		}
		if marker == '_' {
			sub = operand
		} else {
			sup = operand
		}
	}
	body, err := p.parseOperand()
	if err != nil {
		return "", fmt.Errorf("n-ary operand: %w", err)
	}

	hideSub, hideSup := "1", "1"
	if sub != "" {
		hideSub = "0"
	}
	if sup != "" {
		hideSup = "0"
	}
	return `<m:nary><m:naryPr><m:chr m:val="` + operator + `"/>` +
		`<m:limLoc m:val="undOvr"/>` +
		`<m:subHide m:val="` + hideSub + `"/><m:supHide m:val="` + hideSup + `"/>` +
		`<m:ctrlPr/></m:naryPr>` +
		`<m:sub>` + sub + `</m:sub><m:sup>` + sup + `</m:sup>` +
		`<m:e>` + body + `</m:e></m:nary>`, nil
}

func (p *parser) readCommand() (string, error) {
	p.pos++ // consume the backslash
	start := p.pos
	for p.pos < len(p.src) && unicode.IsLetter(p.src[p.pos]) {
		p.pos++
	}
	if p.pos == start {
		return "", fmt.Errorf("expected a command name after '\\' at position %d", start)
	}
	return string(p.src[start:p.pos]), nil
}

func (p *parser) expect(want rune) error {
	if p.pos >= len(p.src) || p.src[p.pos] != want {
		return fmt.Errorf("expected %q at position %d", string(want), p.pos)
	}
	p.pos++
	return nil
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.src) && unicode.IsSpace(p.src[p.pos]) {
		p.pos++
	}
}

// takeLastElement pulls the base a script attaches to: the last pending literal
// character if there is one, otherwise the last element already written.
func takeLastElement(b, pending *strings.Builder) string {
	if pending.Len() > 0 {
		text := pending.String()
		runes := []rune(text)
		last := string(runes[len(runes)-1])
		pending.Reset()
		pending.WriteString(string(runes[:len(runes)-1]))
		return runXML(last)
	}
	written := b.String()
	start := lastElementStart(written)
	if start < 0 {
		return ""
	}
	base := written[start:]
	b.Reset()
	b.WriteString(written[:start])
	return base
}

// lastElementStart finds where the final top-level element of an OMML fragment
// begins, so a script can wrap it.
func lastElementStart(fragment string) int {
	if fragment == "" {
		return -1
	}
	depth := 0
	for i := len(fragment) - 1; i >= 0; i-- {
		if fragment[i] != '<' {
			continue
		}
		switch {
		case strings.HasPrefix(fragment[i:], "</"):
			depth++
		case strings.HasPrefix(fragment[i:], "<m:") || strings.HasPrefix(fragment[i:], "<a:"):
			// A self-closing tag such as <m:chr m:val="∑"/> opens nothing, so
			// counting it as an opening would unbalance the scan and hand back
			// a fragment that is not a whole element.
			if isSelfClosingTag(fragment[i:]) {
				continue
			}
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// isSelfClosingTag reports whether the tag starting at the front of fragment
// ends with "/>".
func isSelfClosingTag(fragment string) bool {
	end := strings.Index(fragment, ">")
	if end < 1 {
		return false
	}
	return fragment[end-1] == '/'
}

// runXML renders literal text as an OMML run.
func runXML(text string) string {
	if text == "" {
		return ""
	}
	return `<m:r><m:t>` + escapeXML(text) + `</m:t></m:r>`
}

// scriptXML wraps a base in a superscript or subscript element.
func scriptXML(marker rune, base, operand string) string {
	if marker == '^' {
		return `<m:sSup><m:sSupPr><m:ctrlPr/></m:sSupPr>` +
			`<m:e>` + base + `</m:e><m:sup>` + operand + `</m:sup></m:sSup>`
	}
	return `<m:sSub><m:sSubPr><m:ctrlPr/></m:sSubPr>` +
		`<m:e>` + base + `</m:e><m:sub>` + operand + `</m:sub></m:sSub>`
}

func escapeXML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(text)
}
