package netsec

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"
)

const (
	restrictedDialTimeout   = 30 * time.Second
	restrictedDialKeepAlive = 30 * time.Second
)

// NewRestrictedHTTPClient builds an HTTP client that blocks private/internal
// IPs at dial time unless allowPrivateHosts is true.
//
// The SSRF guard is always applied regardless of whether http.DefaultTransport
// has been replaced — if the type assertion fails a fresh transport is used so
// the guard never silently falls back to an unrestricted client.
func NewRestrictedHTTPClient(timeout time.Duration, allowPrivateHosts bool) *http.Client {
	baseDialer := &net.Dialer{
		Timeout:   restrictedDialTimeout,
		KeepAlive: restrictedDialKeepAlive,
	}

	var transport *http.Transport
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = t.Clone()
	} else {
		transport = &http.Transport{}
	}
	transport.DialContext = restrictedDialContext(baseDialer, allowPrivateHosts)

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

// ValidateURLForHTTPFetch validates URL scheme and optionally blocks private hosts.
func ValidateURLForHTTPFetch(rawURL string, allowPrivateHosts bool) (*url.URL, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme %q: only http and https are allowed", parsed.Scheme)
	}
	if !allowPrivateHosts {
		if err := validateHost(context.Background(), parsed.Hostname()); err != nil {
			return nil, err
		}
	}
	return parsed, nil
}

func restrictedDialContext(
	base *net.Dialer,
	allowPrivateHosts bool,
) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("parse dial address: %w", err)
		}

		if allowPrivateHosts {
			return base.DialContext(ctx, network, addr)
		}

		if parsedIP, parseErr := netip.ParseAddr(host); parseErr == nil {
			ip := parsedIP.Unmap()
			if blocked, reason := isBlockedAddr(ip); blocked {
				return nil, fmt.Errorf("connection to %s blocked: %s", ip.String(), reason)
			}
			return base.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		}

		resolvedIPs, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", host, err)
		}
		if len(resolvedIPs) == 0 {
			return nil, fmt.Errorf("no addresses for %q", host)
		}

		chosen, err := validateAllPublic(resolvedIPs)
		if err != nil {
			return nil, err
		}
		return base.DialContext(ctx, network, net.JoinHostPort(chosen.String(), port))
	}
}

func validateHost(ctx context.Context, host string) error {
	if host == "" {
		return errors.New("URL host cannot be empty")
	}

	if parsedIP, err := netip.ParseAddr(host); err == nil {
		ip := parsedIP.Unmap()
		if blocked, reason := isBlockedAddr(ip); blocked {
			return fmt.Errorf("request to %q blocked: %s", host, reason)
		}
		return nil
	}

	resolvedIPs, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if err != nil {
		return fmt.Errorf("resolve host %q: %w", host, err)
	}
	if len(resolvedIPs) == 0 {
		return fmt.Errorf("resolve host %q: no addresses found", host)
	}

	for _, ip := range resolvedIPs {
		unmapped := ip.Unmap()
		if blocked, reason := isBlockedAddr(unmapped); blocked {
			return fmt.Errorf("request to %q blocked: %s", host, reason)
		}
	}
	return nil
}

func validateAllPublic(ips []netip.Addr) (netip.Addr, error) {
	for _, ip := range ips {
		unmapped := ip.Unmap()
		if blocked, reason := isBlockedAddr(unmapped); blocked {
			return netip.Addr{}, fmt.Errorf("connection to %s blocked: %s", unmapped.String(), reason)
		}
	}
	return ips[0].Unmap(), nil
}

func isBlockedAddr(ip netip.Addr) (bool, string) {
	if !ip.IsValid() {
		return true, "invalid IP address"
	}
	if ip.IsLoopback() {
		return true, "loopback range"
	}
	if ip.IsPrivate() {
		return true, "private range"
	}
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true, "link-local range"
	}
	if ip.IsMulticast() {
		return true, "multicast range"
	}
	if ip.IsUnspecified() {
		return true, "unspecified range"
	}
	for _, prefix := range blockedPrefixes() {
		if prefix.Contains(ip) {
			return true, "shared address range"
		}
	}
	return false, ""
}

func blockedPrefixes() []netip.Prefix {
	return []netip.Prefix{
		netip.MustParsePrefix("100.64.0.0/10"), // CGNAT
	}
}
