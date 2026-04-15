package urlfetch

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	defaultDialTimeout   = 30 * time.Second
	defaultDialKeepAlive = 30 * time.Second
)

func blockedCIDRs() []string {
	return []string{
		"127.0.0.0/8",    // IPv4 loopback
		"::1/128",        // IPv6 loopback
		"0.0.0.0/8",      // IPv4 unspecified / "this" network
		"::/128",         // IPv6 unspecified
		"10.0.0.0/8",     // RFC 1918 private
		"172.16.0.0/12",  // RFC 1918 private
		"192.168.0.0/16", // RFC 1918 private
		"169.254.0.0/16", // IPv4 link-local (APIPA / AWS metadata)
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique-local (ULA)
		"100.64.0.0/10",  // RFC 6598 shared address space
		"::ffff:0:0/96",  // IPv4-mapped IPv6 addresses
	}
}

// checkIPBlocked returns an error if ip falls within a protected range.
func checkIPBlocked(ip net.IP) error {
	for _, cidr := range blockedCIDRs() {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
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
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Transport{}
	}
	t := defaultTransport.Clone()
	if allowPrivate {
		// Tests: use the default dialer unchanged.
		return t
	}

	base := &net.Dialer{
		Timeout:   defaultDialTimeout,
		KeepAlive: defaultDialKeepAlive,
	}

	t.DialContext = ssrfSafeDialContext(base)
	return t
}

func ssrfSafeDialContext(base *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("parse dial address: %w", err)
		}

		if ip := net.ParseIP(host); ip != nil {
			if err := checkIPBlocked(ip); err != nil {
				return nil, err
			}
			return base.DialContext(ctx, network, addr)
		}

		resolvedIPs, err := net.DefaultResolver.LookupHost(ctx, host)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", host, err)
		}
		if len(resolvedIPs) == 0 {
			return nil, fmt.Errorf("no addresses for %q", host)
		}
		if err := validateResolvedIPs(resolvedIPs); err != nil {
			return nil, err
		}
		return base.DialContext(ctx, network, net.JoinHostPort(resolvedIPs[0], port))
	}
}

func validateResolvedIPs(resolvedIPs []string) error {
	for _, addr := range resolvedIPs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		if err := checkIPBlocked(ip); err != nil {
			return err
		}
	}
	return nil
}

// denyPrivateHost is a lightweight pre-flight check (fast-fail before request
// setup). It does NOT replace ssrfSafeTransport — it still has the TOCTOU
// window — but it produces an early, readable error for the common case.
func denyPrivateHost(hostWithPort string) error {
	hostname, _, err := net.SplitHostPort(hostWithPort)
	if err != nil {
		hostname = hostWithPort
	}
	addrs, err := net.DefaultResolver.LookupHost(context.Background(), hostname)
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
