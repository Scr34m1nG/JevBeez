package security

import (
	"context"
	"errors"
	"net"
	"net/http"
	"jevbeez/api/internal/config"
)

var (
	ErrPrivateIPBlocked = errors.New("access to private or local network is prohibited")
	ErrInvalidIPAddress = errors.New("cannot resolve target host to an IP address")
)

var privateIPBlocks []*net.IPNet

func init() {
	cidrs := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fe80::/10",
		"fc00::/7",
	}

	for _, cidr := range cidrs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return true
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func NewSafeTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout: config.DialerTimeout,
	}

	return &http.Transport{
		ResponseHeaderTimeout: config.ResponseHeaderTimeout,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
			if err != nil {
				return nil, err
			}
			if len(ips) == 0 {
				return nil, ErrInvalidIPAddress
			}

			for _, ip := range ips {
				if IsPrivateIP(ip) {
					return nil, ErrPrivateIPBlocked
				}
			}

			chosenIP := ips[0].String()
			targetAddr := net.JoinHostPort(chosenIP, port)
			return dialer.DialContext(ctx, network, targetAddr)
		},
	}
}

func NewSafeHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewSafeTransport(),
		Timeout:   config.HTTPClientTimeout,
	}
}
