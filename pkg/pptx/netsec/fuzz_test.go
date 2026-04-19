package netsec

import (
	"net/netip"
	"testing"
)

func FuzzValidateURLForHTTPFetch(f *testing.F) {
	f.Add("https://example.com/path", true)
	f.Add("https://example.com/path", false)
	f.Add("http://localhost/admin", false)
	f.Add("ftp://example.com", false)
	f.Add("https://192.168.1.1/", false)
	f.Add("http://[::1]/", false)
	f.Add("", false)
	f.Add("//example.com", false)
	f.Add("javascript:alert(1)", false)
	f.Add("https://100.64.0.1/", false) // CGNAT
	f.Fuzz(func(_ *testing.T, rawURL string, allowPrivate bool) {
		_, _ = ValidateURLForHTTPFetch(rawURL, allowPrivate)
	})
}

func FuzzIsBlockedAddr(f *testing.F) {
	f.Add("127.0.0.1")
	f.Add("192.168.0.1")
	f.Add("10.0.0.1")
	f.Add("172.16.0.1")
	f.Add("8.8.8.8")
	f.Add("::1")
	f.Add("fe80::1")
	f.Add("100.64.0.1")
	f.Add("0.0.0.0")
	f.Add("255.255.255.255")
	f.Add("::ffff:192.168.1.1")
	f.Fuzz(func(_ *testing.T, addrStr string) {
		parsed, err := netip.ParseAddr(addrStr)
		if err != nil {
			return
		}
		_, _ = isBlockedAddr(parsed.Unmap())
	})
}
