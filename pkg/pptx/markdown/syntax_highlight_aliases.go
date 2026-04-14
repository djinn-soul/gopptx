package markdown

func syntaxHighlightLexerAlias(lang string) (string, bool) {
	if mapped, ok := syntaxHighlightLexerAliasPrimary(lang); ok {
		return mapped, true
	}
	return syntaxHighlightLexerAliasSecondary(lang)
}

func syntaxHighlightLexerAliasPrimary(lang string) (string, bool) {
	switch lang {
	case "go", "golang":
		return "go", true
	case "py", "python":
		return "python", true
	case "js", "javascript":
		return "javascript", true
	case "ts", "typescript":
		return "typescript", true
	case "rs", "rust":
		return "rust", true
	case "java", "c":
		return lang, true
	case "cpp", "c++", "cxx":
		return "cpp", true
	case "cs", "csharp":
		return "csharp", true
	case "sh", "bash", "shell":
		return "bash", true
	case "zsh", "json", "xml", "html", "sql":
		return lang, true
	case "vim", "lua", "php", "swift", "scala":
		return lang, true
	case "r", "matlab", "dart", "groovy", "jsonnet":
		return lang, true
	case "nginx", "toml", "ini":
		return lang, true
	default:
		return "", false
	}
}

func syntaxHighlightLexerAliasSecondary(lang string) (string, bool) {
	switch lang {
	case "svg":
		return "xml", true
	case "yaml", "yml":
		return "yaml", true
	case "md", "markdown":
		return "markdown", true
	case "dockerfile":
		return "docker", true
	case "makefile":
		return "makefile", true
	case "ruby", "rb":
		return "ruby", true
	case "kotlin", "kt":
		return "kotlin", true
	case "perl", "pl":
		return "perl", true
	case "haskell", "hs":
		return "haskell", true
	case "clojure", "clj":
		return "clojure", true
	case "erlang", "erl":
		return "erlang", true
	case "elixir", "ex":
		return "elixir", true
	case "tf", "terraform":
		return "terraform", true
	case "apache":
		return "apacheconf", true
	case "properties":
		return "ini", true
	case "diff", "patch":
		return "diff", true
	default:
		return "", false
	}
}
