package urlfetch

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// blockedRanges lists IP address ranges that are never allowed as fetch targets.
// Covers loopback, private, link-local, unspecified, shared, and IPv4-mapped ranges.
var blockedRanges []*net.IPNet

func init() {
	cidrs := []string{
		"127.0.0.0/8",    // IPv4 loopback
		"::1/128",         // IPv6 loopback
		"0.0.0.0/8",      // IPv4 unspecified / "this" network
		"::/128",          // IPv6 unspecified
		"10.0.0.0/8",     // RFC 1918 private
		"172.16.0.0/12",  // RFC 1918 private
		"192.168.0.0/16", // RFC 1918 private
		"169.254.0.0/16", // IPv4 link-local (APIPA / AWS metadata)
		"fe80::/10",       // IPv6 link-local
		"fc00::/7",        // IPv6 unique-local (ULA)
		"100.64.0.0/10",  // RFC 6598 shared address space
		"::ffff:0:0/96",  // IPv4-mapped IPv6 addresses
	}
	for _, cidr := range cidrs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			blockedRanges = append(blockedRanges, block)
		}
	}
}

// checkIPBlocked returns an error if ip falls within a protected range.
func checkIPBlocked(ip net.IP) error {
	for _, block := range blockedRanges {
		if block.Contains(ip) {
			return fmt.Errorf("connection to %s blocked: protected address range", ip)
		}
	}
	return nil
}

// ssrfSafeTransport clones http.DefaultTransport and replaces its DialContext
// with one that checks the resolved IP at connection time — after Go's own DNS
// resolution — eliminating the TOCTOU window present in pre-request checks.
//
// Set allowPrivate=true only in tests that use httptest.NewServer.
func ssrfSafeTransport(allowPrivate bool) *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	if allowPrivate {
		// Tests: use the default dialer unchanged.
		return t
	}

	base := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("parse dial address: %w", err)
		}

		// Fast path: addr is already an IP literal (HTTP transport pre-resolved it).
		if ip := net.ParseIP(host); ip != nil {
			if err := checkIPBlocked(ip); err != nil {
				return nil, err
			}
			return base.DialContext(ctx, network, addr)
		}

		// Slow path: still a hostname — resolve here so we control which
		// address is dialed and can verify every candidate.
		ips, err := net.DefaultResolver.LookupHost(ctx, host)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", host, err)
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no addresses for %q", host)
		}
		for _, a := range ips {
			ip := net.ParseIP(a)
			if ip == nil {
				continue
			}
			if err := checkIPBlocked(ip); err != nil {
				return nil, err
			}
		}
		// All candidates are safe; dial the first one.
		return base.DialContext(ctx, network, net.JoinHostPort(ips[0], port))
	}
	return t
}

// denyPrivateHost is a lightweight pre-flight check (fast-fail before request
// setup). It does NOT replace ssrfSafeTransport — it still has the TOCTOU
// window — but it produces an early, readable error for the common case.
func denyPrivateHost(hostWithPort string) error {
	hostname, _, err := net.SplitHostPort(hostWithPort)
	if err != nil {
		hostname = hostWithPort
	}
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return fmt.Errorf("resolve host %q: %w", hostname, err)
	}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		if err := checkIPBlocked(ip); err != nil {
			return fmt.Errorf("request to %q: %w", hostname, err)
		}
	}
	return nil
}
