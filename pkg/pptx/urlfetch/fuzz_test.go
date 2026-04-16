package urlfetch

import (
	"net/netip"
	"testing"
)

func FuzzWebParserParse(f *testing.F) {
	f.Add("<html><body><main><h1>Title</h1><p>Some text here.</p></main></body></html>", "https://example.com")
	f.Add("<article><h2>Sub</h2><ul><li>item one</li><li>item two</li></ul></article>", "")
	f.Add("<html><body><article><table><tr><th>A</th><th>B</th></tr><tr><td>1</td><td>2</td></tr></table></article></body></html>", "https://example.com/page")
	f.Add("<html><body><main><blockquote>quote text</blockquote></main></body></html>", "")
	f.Add("<html><body><main><img src=\"/img.png\" alt=\"desc\"/></main></body></html>", "https://example.com")
	f.Add("<html><body><main><a href=\"https://link.com\">link text</a></main></body></html>", "https://example.com")
	f.Add("", "")
	f.Add("<html></html>", "not-a-url")
	f.Fuzz(func(_ *testing.T, html, pageURL string) {
		p := NewWebParser()
		_, _ = p.Parse(html, pageURL)
	})
}

func FuzzCheckAddrBlocked(f *testing.F) {
	f.Add("127.0.0.1")
	f.Add("192.168.1.1")
	f.Add("10.0.0.0")
	f.Add("8.8.8.8")
	f.Add("::1")
	f.Add("fe80::1")
	f.Add("100.64.0.1")
	f.Add("0.0.0.0")
	f.Add("::ffff:192.168.0.1")
	f.Fuzz(func(_ *testing.T, addrStr string) {
		parsed, err := netip.ParseAddr(addrStr)
		if err != nil {
			return
		}
		// Fuzz the blocking predicate directly to avoid network calls.
		_ = checkAddrBlocked(parsed.Unmap())
	})
}
