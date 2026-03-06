package markdown

import (
	"strings"
	"testing"
)

func TestSolarizedColors(t *testing.T) {
	// Verify Solarized Dark color palette constants
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		{"base03", SolarizedBase03, "002b36"},
		{"base02", SolarizedBase02, "073642"},
		{"base01", SolarizedBase01, "586e75"},
		{"base00", SolarizedBase00, "657b83"},
		{"base0", SolarizedBase0, "839496"},
		{"base1", SolarizedBase1, "93a1a1"},
		{"base2", SolarizedBase2, "eee8d5"},
		{"base3", SolarizedBase3, "fdf6e3"},
		{"yellow", SolarizedYellow, "b58900"},
		{"orange", SolarizedOrange, "cb4b16"},
		{"red", SolarizedRed, "dc322f"},
		{"magenta", SolarizedMagenta, "d33682"},
		{"violet", SolarizedViolet, "6c71c4"},
		{"blue", SolarizedBlue, "268bd2"},
		{"cyan", SolarizedCyan, "2aa198"},
		{"green", SolarizedGreen, "859900"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.color)
			}
		})
	}
}

func TestGetColor(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenKeyword, SolarizedYellow},
		{TokenString, SolarizedCyan},
		{TokenComment, SolarizedBase01},
		{TokenNumber, SolarizedMagenta},
		{TokenOperator, SolarizedYellow},
		{TokenFunction, SolarizedBlue},
		{TokenTypeName, SolarizedViolet},
		{TokenConstant, SolarizedOrange},
		{TokenPunctuation, SolarizedBase0},
		{TokenText, SolarizedBase0},
		{TokenWhitespace, SolarizedBase0},
		{TokenVariable, SolarizedBase0},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := GetColor(tt.tokenType)
			if got != tt.expected {
				t.Errorf("GetColor(%v) = %s, want %s", tt.tokenType, got, tt.expected)
			}
		})
	}
}

func TestTokenizeGo(t *testing.T) {
	code := `package main

import "fmt"

func main() {
	// This is a comment
	message := "Hello, World!"
	fmt.Println(message)
	const MaxCount = 100
}`

	h := newSyntaxHighlighter("go")
	tokens := h.Tokenize(code)

	// Verify we got tokens
	if len(tokens) == 0 {
		t.Fatal("expected tokens, got none")
	}

	// Check for keywords
	foundPackage := false
	foundFunc := false
	foundImport := false
	for _, tok := range tokens {
		if tok.Text == "package" && tok.Type == TokenKeyword {
			foundPackage = true
		}
		if tok.Text == "func" && tok.Type == TokenKeyword {
			foundFunc = true
		}
		if tok.Text == "import" && tok.Type == TokenKeyword {
			foundImport = true
		}
	}
	if !foundPackage {
		t.Error("expected to find 'package' keyword")
	}
	if !foundFunc {
		t.Error("expected to find 'func' keyword")
	}
	if !foundImport {
		t.Error("expected to find 'import' keyword")
	}

	// Check for string
	foundString := false
	for _, tok := range tokens {
		if tok.Text == `"Hello, World!"` && tok.Type == TokenString {
			foundString = true
			break
		}
	}
	if !foundString {
		t.Error("expected to find string literal")
	}

	// Check for comment
	foundComment := false
	for _, tok := range tokens {
		if strings.Contains(tok.Text, "//") && tok.Type == TokenComment {
			foundComment = true
			break
		}
	}
	if !foundComment {
		t.Error("expected to find comment")
	}

	// Check for number
	foundNumber := false
	for _, tok := range tokens {
		if tok.Text == "100" && tok.Type == TokenNumber {
			foundNumber = true
			break
		}
	}
	if !foundNumber {
		t.Error("expected to find number literal")
	}
}

func TestTokenizePython(t *testing.T) {
	code := `def greet(name):
    # A greeting function
    message = f"Hello, {name}!"
    return message

class Person:
    def __init__(self, name):
        self.name = name`

	h := newSyntaxHighlighter("python")
	tokens := h.Tokenize(code)

	// Check for Python keywords
	foundDef := false
	foundClass := false
	foundReturn := false
	for _, tok := range tokens {
		if tok.Text == "def" && tok.Type == TokenKeyword {
			foundDef = true
		}
		if tok.Text == "class" && tok.Type == TokenKeyword {
			foundClass = true
		}
		if tok.Text == "return" && tok.Type == TokenKeyword {
			foundReturn = true
		}
	}
	if !foundDef {
		t.Error("expected to find 'def' keyword")
	}
	if !foundClass {
		t.Error("expected to find 'class' keyword")
	}
	if !foundReturn {
		t.Error("expected to find 'return' keyword")
	}

	// Check for Python comment
	foundComment := false
	for _, tok := range tokens {
		if strings.Contains(tok.Text, "#") && tok.Type == TokenComment {
			foundComment = true
			break
		}
	}
	if !foundComment {
		t.Error("expected to find Python comment")
	}
}

func TestTokenizeJavaScript(t *testing.T) {
	code := `function greet(name) {
    // A greeting function
    const message = "Hello, " + name + "!";
    return message;
}`

	h := newSyntaxHighlighter("javascript")
	tokens := h.Tokenize(code)

	// Check for JS keywords
	foundFunction := false
	foundConst := false
	foundReturn := false
	for _, tok := range tokens {
		if tok.Text == "function" && tok.Type == TokenKeyword {
			foundFunction = true
		}
		if tok.Text == "const" && tok.Type == TokenKeyword {
			foundConst = true
		}
		if tok.Text == "return" && tok.Type == TokenKeyword {
			foundReturn = true
		}
	}
	if !foundFunction {
		t.Error("expected to find 'function' keyword")
	}
	if !foundConst {
		t.Error("expected to find 'const' keyword")
	}
	if !foundReturn {
		t.Error("expected to find 'return' keyword")
	}
}

func TestTokenizeRust(t *testing.T) {
	code := `fn main() {
    // A Rust comment
    let message = "Hello, Rust!";
    println!("{}", message);
}`

	h := newSyntaxHighlighter("rust")
	tokens := h.Tokenize(code)

	// Check for Rust keywords
	foundFn := false
	foundLet := false
	for _, tok := range tokens {
		if tok.Text == "fn" && tok.Type == TokenKeyword {
			foundFn = true
		}
		if tok.Text == "let" && tok.Type == TokenKeyword {
			foundLet = true
		}
	}
	if !foundFn {
		t.Error("expected to find 'fn' keyword")
	}
	if !foundLet {
		t.Error("expected to find 'let' keyword")
	}

	// Check for string (Chroma may tokenize quotes separately)
	foundString := false
	for _, tok := range tokens {
		if strings.Contains(tok.Text, "Hello, Rust") && (tok.Type == TokenString || tok.Color == SolarizedCyan) {
			foundString = true
			break
		}
	}
	if !foundString {
		t.Error("expected to find string literal")
	}
}

func TestTokenizeShell(t *testing.T) {
	code := `#!/bin/bash
# A shell script
echo "Hello, World!"
if [ -f "$FILE" ]; then
    echo "File exists"
fi`

	h := newSyntaxHighlighter("bash")
	tokens := h.Tokenize(code)

	// Check for shell keywords
	foundIf := false
	foundThen := false
	for _, tok := range tokens {
		if tok.Text == "if" && tok.Type == TokenKeyword {
			foundIf = true
		}
		if tok.Text == "then" && tok.Type == TokenKeyword {
			foundThen = true
		}
	}
	if !foundIf {
		t.Error("expected to find 'if' keyword")
	}
	if !foundThen {
		t.Error("expected to find 'then' keyword")
	}

	// Check for string
	foundString := false
	for _, tok := range tokens {
		if tok.Text == `"Hello, World!"` && tok.Type == TokenString {
			foundString = true
			break
		}
	}
	if !foundString {
		t.Error("expected to find string literal")
	}
}

func TestTokenizeJSON(t *testing.T) {
	code := `{
    "name": "test",
    "count": 42,
    "active": true,
    "empty": null
}`

	h := newSyntaxHighlighter("json")
	tokens := h.Tokenize(code)

	// Check for strings (Chroma may have quotes separate or together, and JSON keys may be NameTag)
	foundNameString := false
	for _, tok := range tokens {
		if strings.Contains(tok.Text, "name") {
			foundNameString = true
			break
		}
	}
	if !foundNameString {
		t.Error("expected to find string literal")
	}

	// Check for numbers
	foundNumber := false
	for _, tok := range tokens {
		if tok.Text == "42" && tok.Type == TokenNumber {
			foundNumber = true
			break
		}
	}
	if !foundNumber {
		t.Error("expected to find number literal")
	}

	// Check for keywords
	foundTrue := false
	foundNull := false
	for _, tok := range tokens {
		if tok.Text == "true" && (tok.Type == TokenKeyword || tok.Type == TokenConstant) {
			foundTrue = true
		}
		if tok.Text == "null" && (tok.Type == TokenKeyword || tok.Type == TokenConstant) {
			foundNull = true
		}
	}
	if !foundTrue {
		t.Error("expected to find 'true' keyword")
	}
	if !foundNull {
		t.Error("expected to find 'null' keyword")
	}
}

func TestTokenizeSQL(t *testing.T) {
	code := `SELECT id, name FROM users
WHERE active = true
ORDER BY name ASC;
-- End of query`

	h := newSyntaxHighlighter("sql")
	tokens := h.Tokenize(code)

	// Check for SQL keywords (case insensitive)
	foundSelect := false
	foundFrom := false
	foundWhere := false
	foundOrder := false
	for _, tok := range tokens {
		if strings.ToUpper(tok.Text) == "SELECT" && tok.Type == TokenKeyword {
			foundSelect = true
		}
		if strings.ToUpper(tok.Text) == "FROM" && tok.Type == TokenKeyword {
			foundFrom = true
		}
		if strings.ToUpper(tok.Text) == "WHERE" && tok.Type == TokenKeyword {
			foundWhere = true
		}
		if strings.ToUpper(tok.Text) == "ORDER" && tok.Type == TokenKeyword {
			foundOrder = true
		}
	}
	if !foundSelect {
		t.Error("expected to find 'SELECT' keyword")
	}
	if !foundFrom {
		t.Error("expected to find 'FROM' keyword")
	}
	if !foundWhere {
		t.Error("expected to find 'WHERE' keyword")
	}
	if !foundOrder {
		t.Error("expected to find 'ORDER' keyword")
	}
}

func TestTokenizeBasicFallback(t *testing.T) {
	code := `Some plain text
// A comment
"A string"`

	h := newSyntaxHighlighter("unknown")
	tokens := h.Tokenize(code)

	// Should still tokenize strings and comments in basic mode
	foundString := false
	foundComment := false
	for _, tok := range tokens {
		if tok.Type == TokenString {
			foundString = true
		}
		if tok.Type == TokenComment {
			foundComment = true
		}
	}
	if !foundString {
		t.Error("expected to find string in fallback mode")
	}
	if !foundComment {
		t.Error("expected to find comment in fallback mode")
	}
}

func TestTokenizeBlockComment(t *testing.T) {
	code := `func main() {
	/* This is a
	multiline comment */
	x := 1
}`

	h := newSyntaxHighlighter("go")
	tokens := h.Tokenize(code)

	// Check that block comment is tokenized as single comment
	foundBlockComment := false
	for i, tok := range tokens {
		if !strings.Contains(tok.Text, "/*") || tok.Type != TokenComment || !strings.Contains(tok.Text, "*/") {
			continue
		}
		foundBlockComment = true
		// Make sure we don't have the code inside the comment as separate tokens.
		if i+1 < len(tokens) {
			next := tokens[i+1]
			if strings.Contains(next.Text, "multiline") && next.Type != TokenComment {
				t.Error("multiline comment content should be part of comment token")
			}
		}
	}
	if !foundBlockComment {
		t.Error("expected to find block comment")
	}
}

func TestTokenizeRawString(t *testing.T) {
	code := "query := `SELECT id FROM table WHERE id = 1`"

	h := newSyntaxHighlighter("go")
	tokens := h.Tokenize(code)

	foundRawString := false
	for _, tok := range tokens {
		if strings.HasPrefix(tok.Text, "`") && strings.HasSuffix(tok.Text, "`") && tok.Type == TokenString {
			foundRawString = true
			break
		}
	}
	if !foundRawString {
		t.Error("expected to find raw string literal")
	}
}

func TestTokenizeTripleQuotePython(t *testing.T) {
	code := `docstring = """This is a
multiline
docstring"""`

	h := newSyntaxHighlighter("python")
	tokens := h.Tokenize(code)

	foundTripleQuote := false
	for _, tok := range tokens {
		if strings.HasPrefix(tok.Text, `"""`) && strings.HasSuffix(tok.Text, `"""`) && tok.Type == TokenString {
			foundTripleQuote = true
			break
		}
	}
	if !foundTripleQuote {
		t.Error("expected to find triple-quoted string")
	}
}

func TestTokenizeMultilineGo(t *testing.T) {
	code := `package main

func main() {
	println("line1")
	println("line2")
}`

	h := newSyntaxHighlighter("go")
	tokens := h.Tokenize(code)

	// Count newlines in tokens
	newlineCount := 0
	for _, tok := range tokens {
		if tok.Text == "\n" && tok.Type == TokenWhitespace {
			newlineCount++
		}
	}

	if newlineCount < 4 {
		t.Errorf("expected at least 4 newlines, got %d", newlineCount)
	}
}

func BenchmarkTokenizeGo(b *testing.B) {
	code := `package main

import "fmt"

func main() {
	// A simple program
	for i := 0; i < 100; i++ {
		fmt.Printf("Number: %d\n", i)
	}
}`

	h := newSyntaxHighlighter("go")

	for b.Loop() {
		h.Tokenize(code)
	}
}
