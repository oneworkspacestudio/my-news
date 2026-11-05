package news

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

func publicIP(ip net.IP) bool {
	a, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	a = a.Unmap()
	if !a.IsGlobalUnicast() || a.IsPrivate() || a.IsLoopback() || a.IsLinkLocalUnicast() {
		return false
	}
	for _, s := range []string{"0.0.0.0/8", "100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "240.0.0.0/4", "2001:db8::/32", "2001::/32", "2002::/16", "64:ff9b::/96"} {
		if netip.MustParsePrefix(s).Contains(a) {
			return false
		}
	}
	return true
}
func validURL(raw string) (*url.URL, error) {
	u, e := url.Parse(raw)
	if e != nil || u.Hostname() == "" || u.User != nil || len(raw) > 2048 || (u.Scheme != "https" && u.Scheme != "http") || u.Port() != "" && u.Port() != "443" && u.Port() != "80" {
		return nil, errors.New("Use a public HTTP or HTTPS feed URL")
	}
	if ip := net.ParseIP(u.Hostname()); ip != nil && !publicIP(ip) {
		return nil, errors.New("Private network addresses are not allowed")
	}
	return u, nil
}
func newClient() *http.Client {
	tr := &http.Transport{Proxy: nil, MaxIdleConns: 12, MaxIdleConnsPerHost: 2, IdleConnTimeout: 30 * time.Second, TLSHandshakeTimeout: 8 * time.Second, ResponseHeaderTimeout: 12 * time.Second, DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, e := net.SplitHostPort(address)
		if e != nil {
			return nil, e
		}
		ips, e := net.DefaultResolver.LookupIP(ctx, "ip", host)
		if e != nil {
			return nil, e
		}
		if len(ips) == 0 {
			return nil, errors.New("No public address")
		}
		for _, ip := range ips {
			if !publicIP(ip) {
				return nil, errors.New("Private network addresses are blocked")
			}
		}
		d := net.Dialer{Timeout: 8 * time.Second}
		var last error
		for _, ip := range ips {
			c, e := d.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
			if e == nil {
				return c, nil
			}
			last = e
		}
		return nil, last
	}}
	return &http.Client{Transport: tr, Timeout: 20 * time.Second, CheckRedirect: func(r *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("Too many redirects")
		}
		_, e := validURL(r.URL.String())
		return e
	}}
}
func fetch(ctx context.Context, c *http.Client, raw string, max int64) ([]byte, string, error) {
	u, e := validURL(raw)
	if e != nil {
		return nil, "", e
	}
	r, e := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if e != nil {
		return nil, "", e
	}
	r.Header.Set("User-Agent", "MyNews/0.1 (personal RSS reader)")
	r.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml, text/xml, text/html;q=0.8, */*;q=0.1")
	res, e := c.Do(r)
	if e != nil {
		return nil, "", errors.New("Publisher could not be reached")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, "", errors.New("Publisher returned " + res.Status)
	}
	b, e := io.ReadAll(io.LimitReader(res.Body, max+1))
	if e != nil {
		return nil, "", e
	}
	if int64(len(b)) > max {
		return nil, "", errors.New("Publisher response is too large")
	}
	return b, strings.ToLower(res.Header.Get("Content-Type")), nil
}
