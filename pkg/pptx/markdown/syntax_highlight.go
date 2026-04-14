package markdown

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Solarized Dark color palette
// https://ethanschoonover.com/solarized/
const (
	// SolarizedBase03 starts the base-tone palette group (background/foreground).
	SolarizedBase03 = "002b36" // background
	SolarizedBase02 = "073642" // background highlights
	SolarizedBase01 = "586e75" // comments/secondary content
	SolarizedBase00 = "657b83" // body text
	SolarizedBase0  = "839496" // body text (default code foreground)
	SolarizedBase1  = "93a1a1" // optional emphasized content
	SolarizedBase2  = "eee8d5" // background highlights
	SolarizedBase3  = "fdf6e3" // background

	// SolarizedYellow starts the accent color palette group.
	SolarizedYellow  = "b58900" // keywords
	SolarizedOrange  = "cb4b16" // special keywords
	SolarizedRed     = "dc322f" // errors
	SolarizedMagenta = "d33682" // numbers
	SolarizedViolet  = "6c71c4" // types
	SolarizedBlue    = "268bd2" // functions
	SolarizedCyan    = "2aa198" // strings
	SolarizedGreen   = "859900" // comments
)

// TokenType represents the type of a syntax token.
type TokenType int

const (
	TokenText TokenType = iota
	TokenKeyword
	TokenString
	TokenComment
	TokenNumber
	TokenOperator
	TokenFunction
	TokenTypeName
	TokenVariable
	TokenConstant
	TokenPunctuation
	TokenWhitespace
)

// Token represents a single syntax token with its type.
type Token struct {
	Text  string
	Type  TokenType
	Color string
}

// syntaxHighlighter provides token-level syntax highlighting using Chroma.
type syntaxHighlighter struct {
	lang string
}

// newSyntaxHighlighter creates a new highlighter for the given language.
func newSyntaxHighlighter(lang string) *syntaxHighlighter {
	return &syntaxHighlighter{lang: strings.ToLower(lang)}
}

// Tokenize breaks code into colored tokens using Chroma lexer.
func (h *syntaxHighlighter) Tokenize(code string) []Token {
	// Get lexer for language
	lexer := h.getLexer(code)
	if lexer == nil {
		// Fallback to basic tokenization
		return h.tokenizeBasic(code)
	}

	// Tokenize using Chroma
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return h.tokenizeBasic(code)
	}

	// Convert Chroma tokens to our Token format
	var tokens []Token
	style := styles.Get("solarized-dark")
	if style == nil {
		style = styles.Fallback
	}

	for _, tok := range iterator.Tokens() {
		tokenType := h.mapChromaType(tok.Type)
		color := h.getColorForToken(tok.Type, style)

		tokens = append(tokens, Token{
			Text:  tok.Value,
			Type:  tokenType,
			Color: color,
		})
	}

	return tokens
}

// getLexer returns the appropriate Chroma lexer for the language.
func (h *syntaxHighlighter) getLexer(code string) chroma.Lexer {
	// Map common language names to Chroma lexers
	lang := h.lang
	if lang == "" || lang == "text" || lang == "plain" {
		return nil
	}

	// Try exact match first
	lexer := lexers.Get(lang)
	if lexer != nil {
		return lexer
	}

	// Try aliases
	if mapped, ok := syntaxHighlightLexerAlias(lang); ok {
		lexer = lexers.Get(mapped)
		if lexer != nil {
			return lexer
		}
	}

	// Try to analyze content
	return lexers.Analyse(code)
}

// mapChromaType maps Chroma token types to our simplified TokenType.
func (h *syntaxHighlighter) mapChromaType(t chroma.TokenType) TokenType {
	switch t {
	case chroma.Keyword, chroma.KeywordConstant, chroma.KeywordDeclaration,
		chroma.KeywordNamespace, chroma.KeywordPseudo, chroma.KeywordReserved,
		chroma.KeywordType:
		return TokenKeyword
	case chroma.String, chroma.StringAffix, chroma.StringBacktick,
		chroma.StringChar, chroma.StringDelimiter, chroma.StringDoc,
		chroma.StringDouble, chroma.StringEscape, chroma.StringHeredoc,
		chroma.StringInterpol, chroma.StringOther, chroma.StringRegex,
		chroma.StringSingle, chroma.StringSymbol:
		return TokenString
	case chroma.Comment, chroma.CommentHashbang, chroma.CommentMultiline,
		chroma.CommentPreproc, chroma.CommentPreprocFile, chroma.CommentSingle,
		chroma.CommentSpecial:
		return TokenComment
	case chroma.Number, chroma.NumberBin, chroma.NumberFloat,
		chroma.NumberHex, chroma.NumberInteger, chroma.NumberIntegerLong,
		chroma.NumberOct:
		return TokenNumber
	case chroma.Operator, chroma.OperatorWord:
		return TokenOperator
	case chroma.NameFunction, chroma.NameFunctionMagic, chroma.NameProperty:
		return TokenFunction
	case chroma.NameClass, chroma.NameException, chroma.NameDecorator,
		chroma.NameBuiltin, chroma.NameBuiltinPseudo, chroma.NameAttribute:
		return TokenTypeName
	case chroma.NameVariable, chroma.NameVariableClass, chroma.NameVariableGlobal,
		chroma.NameVariableInstance, chroma.NameVariableMagic:
		return TokenVariable
	case chroma.NameConstant, chroma.NameTag, chroma.NameLabel:
		return TokenConstant
	case chroma.Punctuation:
		return TokenPunctuation
	case chroma.Text, chroma.TextWhitespace:
		return TokenWhitespace
	default:
		return TokenText
	}
}

// getColorForToken extracts the color from Chroma style.
func (h *syntaxHighlighter) getColorForToken(t chroma.TokenType, style *chroma.Style) string {
	entry := style.Get(t)
	if entry.Colour.IsSet() {
		return entry.Colour.String()
	}
	// Fallback to our Solarized colors
	return GetColor(h.mapChromaType(t))
}

// GetColor returns the Solarized color for a token type.
func GetColor(t TokenType) string {
	switch t {
	case TokenKeyword:
		return SolarizedYellow
	case TokenString:
		return SolarizedCyan
	case TokenComment:
		return SolarizedBase01
	case TokenNumber:
		return SolarizedMagenta
	case TokenOperator:
		return SolarizedYellow
	case TokenFunction:
		return SolarizedBlue
	case TokenTypeName:
		return SolarizedViolet
	case TokenConstant:
		return SolarizedOrange
	case TokenPunctuation:
		return SolarizedBase0
	default:
		return SolarizedBase0
	}
}

// tokenizeBasic provides basic fallback tokenization.
func (h *syntaxHighlighter) tokenizeBasic(code string) []Token {
	var tokens []Token
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		// Check for comments in various formats
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*") {
			tokens = append(tokens, Token{
				Text:  line,
				Type:  TokenComment,
				Color: SolarizedBase01,
			})
		} else {
			// Highlight quoted strings
			tokens = append(tokens, h.tokenizeStrings(line)...)
		}
		if i < len(lines)-1 {
			tokens = append(tokens, Token{Text: "\n", Type: TokenWhitespace, Color: SolarizedBase0})
		}
	}
	return tokens
}

// tokenizeStrings extracts string literals from a line.
func (h *syntaxHighlighter) tokenizeStrings(line string) []Token {
	var tokens []Token
	i := 0
	for i < len(line) {
		if isQuoteChar(line[i]) {
			start := i
			i = scanQuotedStringEnd(line, i)
			tokens = append(tokens, Token{
				Text:  line[start:i],
				Type:  TokenString,
				Color: SolarizedCyan,
			})
		} else {
			// Collect non-string text
			start := i
			for i < len(line) && !isQuoteChar(line[i]) {
				i++
			}
			if i > start {
				tokens = append(tokens, Token{
					Text:  line[start:i],
					Type:  TokenText,
					Color: SolarizedBase0,
				})
			}
		}
	}
	if len(tokens) == 0 {
		tokens = append(tokens, Token{
			Text:  line,
			Type:  TokenText,
			Color: SolarizedBase0,
		})
	}
	return tokens
}

func isQuoteChar(char byte) bool {
	return char == '"' || char == '\'' || char == '`'
}

func scanQuotedStringEnd(line string, start int) int {
	quote := line[start]
	i := start + 1
	for i < len(line) && line[i] != quote {
		if line[i] == '\\' && i+1 < len(line) {
			i += 2
			continue
		}
		i++
	}
	if i < len(line) {
		return i + 1
	}
	return i
}
