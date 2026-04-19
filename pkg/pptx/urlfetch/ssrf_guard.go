package urlfetch

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"time"
)

const (
	defaultDialTimeout   = 30 * time.Second
	defaultDialKeepAlive = 30 * time.Second
)

// checkAddrBlocked reports whether ip (already unwrapped via Unmap) falls
// within a protected address range. Uses netip predicates rather than a
// manual CIDR list so that IPv4-mapped IPv6 addresses are handled correctly
// after the caller calls ip.Unmap().
func checkAddrBlocked(ip netip.Addr) error {
	switch {
	case !ip.IsValid(),
		ip.IsLoopback(),
		ip.IsPrivate(),
		ip.IsLinkLocalUnicast(),
		ip.IsLinkLocalMulticast(),
		ip.IsMulticast(),
		ip.IsUnspecified():
		return fmt.Errorf("connection to %s blocked: protected address range", ip)
	}
	for _, pfx := range blockedPrefixes() {
		if pfx.Contains(ip) {
			return fmt.Errorf("connection to %s blocked: protected address range", ip)
		}
	}
	return nil
}

func blockedPrefixes() []netip.Prefix {
	return []netip.Prefix{
		netip.MustParsePrefix("100.64.0.0/10"), // CGNAT / RFC 6598 shared address space
	}
}

// ssrfSafeTransport returns a transport whose DialContext checks the resolved
// IP at connection time — after Go's own DNS resolution — eliminating the
// TOCTOU window present in pre-request checks.
//
// The guard is always applied: if http.DefaultTransport has been replaced with
// a non-*http.Transport type a fresh transport is used so the guard never
// silently falls back to an unrestricted client.
//
// Set allowPrivate=true only in tests that use httptest.NewServer.
func ssrfSafeTransport(allowPrivate bool) *http.Transport {
	var t *http.Transport
	if dt, ok := http.DefaultTransport.(*http.Transport); ok {
		t = dt.Clone()
	} else {
		t = &http.Transport{}
	}
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

		// Literal IP: unmap IPv4-mapped IPv6 before checking.
		if parsedIP, parseErr := netip.ParseAddr(host); parseErr == nil {
			ip := parsedIP.Unmap()
			if err := checkAddrBlocked(ip); err != nil {
				return nil, err
			}
			return base.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		}

		// Hostname: resolve, validate every returned address, then dial the first.
		resolvedIPs, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", host, err)
		}
		if len(resolvedIPs) == 0 {
			return nil, fmt.Errorf("no addresses for %q", host)
		}
		chosen, err := validateAllPublicAndPick(resolvedIPs)
		if err != nil {
			return nil, err
		}
		return base.DialContext(ctx, network, net.JoinHostPort(chosen.String(), port))
	}
}

// validateAllPublicAndPick rejects the address list if any entry resolves to a
// protected range and returns the first (unmapped) address otherwise.
func validateAllPublicAndPick(ips []netip.Addr) (netip.Addr, error) {
	for _, ip := range ips {
		unmapped := ip.Unmap()
		if err := checkAddrBlocked(unmapped); err != nil {
			return netip.Addr{}, err
		}
	}
	return ips[0].Unmap(), nil
}

// denyPrivateHost is a lightweight pre-flight check (fast-fail before request
// setup). It does NOT replace ssrfSafeTransport — it still has the TOCTOU
// window — but it produces an early, readable error for the common case.
func denyPrivateHost(hostWithPort string) error {
	hostname, _, err := net.SplitHostPort(hostWithPort)
	if err != nil {
		hostname = hostWithPort
	}
	addrs, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", hostname)
	if err != nil {
		return fmt.Errorf("resolve host %q: %w", hostname, err)
	}
	for _, addr := range addrs {
		ip := addr.Unmap()
		if err := checkAddrBlocked(ip); err != nil {
			return fmt.Errorf("request to %q: %w", hostname, err)
		}
	}
	return nil
}
