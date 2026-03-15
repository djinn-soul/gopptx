package markdown

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func applyMarkdownLinkRuns(
	runs []elements.Run,
	destination string,
	resolveLink markdownRunLinkResolver,
) []elements.Run {
	if resolveLink == nil {
		return runs
	}
	link, ok := resolveLink(destination)
	if !ok {
		return runs
	}
	out := make([]elements.Run, 0, len(runs))
	for _, run := range runs {
		if strings.TrimSpace(run.Text) == "" {
			out = append(out, run)
			continue
		}
		out = append(out, run.WithHyperlink(link))
	}
	return out
}

func (p *markdownASTParser) resolveRunHyperlink(destination string) (action.Hyperlink, bool) {
	dest := strings.TrimSpace(destination)
	if dest == "" {
		return action.Hyperlink{}, false
	}
	if strings.HasPrefix(dest, "#") {
		// Policy: skip in-document anchors because PPT text-run hyperlinks do not
		// map markdown fragment anchors directly.
		return action.Hyperlink{}, false
	}

	parsed, err := url.Parse(dest)
	if err != nil {
		return action.Hyperlink{}, false
	}

	scheme := strings.ToLower(parsed.Scheme)
	switch scheme {
	case "http", "https":
		return action.NewHyperlink(action.HyperlinkURL(dest)), true
	case "mailto":
		email := strings.TrimPrefix(dest, "mailto:")
		email = strings.TrimSpace(email)
		if email == "" {
			return action.Hyperlink{}, false
		}
		return action.NewHyperlink(action.HyperlinkEmail(email)), true
	case "file":
		return action.NewHyperlink(action.HyperlinkFile(dest)), true
	case "":
	default:
		return action.Hyperlink{}, false
	}

	localPath := strings.TrimSpace(parsed.Path)
	if localPath == "" {
		return action.Hyperlink{}, false
	}

	baseDir := strings.TrimSpace(p.options.BaseDir)
	if baseDir == "" || filepath.IsAbs(localPath) {
		localPath = filepath.Clean(localPath)
		if strings.Contains(localPath, "..") {
			return action.Hyperlink{}, false
		}
		return action.NewHyperlink(action.HyperlinkFile(localPath)), true
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return action.Hyperlink{}, false
	}
	absLocalPath, err := filepath.Abs(filepath.Join(baseDir, localPath))
	if err != nil {
		return action.Hyperlink{}, false
	}
	relPath, err := filepath.Rel(absBaseDir, absLocalPath)
	if err != nil || relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return action.Hyperlink{}, false
	}
	localPath = absLocalPath

	return action.NewHyperlink(action.HyperlinkFile(localPath)), true
}
