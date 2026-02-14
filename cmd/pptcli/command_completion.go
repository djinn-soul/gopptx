package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

func runCompletionCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("completion", flag.ContinueOnError)
	fs.SetOutput(stderr)

	shell := fs.String("shell", "bash", "completion shell: bash|zsh")

	if err := fs.Parse(args); err != nil {
		printErrorf(stderr, "completion: %v", err)
		return exitUsage
	}

	switch strings.ToLower(strings.TrimSpace(*shell)) {
	case "bash":
		_, _ = fmt.Fprintln(stdout, bashCompletionScript())
		return exitOK
	case "zsh":
		_, _ = fmt.Fprintln(stdout, zshCompletionScript())
		return exitOK
	default:
		printErrorf(stderr, "unsupported shell %q (expected bash or zsh)", *shell)
		return exitUsage
	}
}

func bashCompletionScript() string {
	return `# bash completion for pptcli
_pptcli_complete() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    local cmds="create md2ppt info validate merge completion version help"
    local global_flags="-h --help -version --version"

    case "${COMP_CWORD}" in
        1)
            COMPREPLY=( $(compgen -W "${cmds} ${global_flags}" -- "${cur}") )
            return 0
            ;;
    esac

    case "${COMP_WORDS[1]}" in
        completion)
            COMPREPLY=( $(compgen -W "-shell bash zsh" -- "${cur}") )
            return 0
            ;;
        create)
            COMPREPLY=( $(compgen -W "-out -title -slides" -- "${cur}") )
            return 0
            ;;
        md2ppt)
            COMPREPLY=( $(compgen -W "-in -out -title" -- "${cur}") )
            return 0
            ;;
        info|validate)
            COMPREPLY=( $(compgen -W "-file" -- "${cur}") )
            return 0
            ;;
        merge)
            COMPREPLY=( $(compgen -W "-out" -- "${cur}") )
            return 0
            ;;
    esac
}
complete -F _pptcli_complete pptcli`
}

func zshCompletionScript() string {
	return `#compdef pptcli

_pptcli() {
  local -a commands
  commands=(
    'create:Create a new presentation'
    'md2ppt:Convert markdown to PPTX'
    'info:Show PPTX metadata'
    'validate:Validate PPTX package'
    'merge:Merge multiple presentations'
    'completion:Generate shell completion script'
    'version:Show version'
    'help:Show help'
  )

  if (( CURRENT == 2 )); then
    _describe 'command' commands
    return
  fi

  case "$words[2]" in
    completion)
      _arguments '-shell[completion shell]:shell:(bash zsh)'
      ;;
    create)
      _arguments '-out[output file]' '-title[presentation title]' '-slides[number of slides]'
      ;;
    md2ppt)
      _arguments '-in[input markdown file]' '-out[output pptx file]' '-title[presentation title]'
      ;;
    info|validate)
      _arguments '-file[input pptx file]'
      ;;
    merge)
      _arguments '-out[output pptx file]'
      ;;
  esac
}

_pptcli "$@"`
}
