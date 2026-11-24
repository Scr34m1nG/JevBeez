package security

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"jevbeez/api/internal/config"
)

var (
	ErrShortenerLoop    = errors.New("shortener redirect loop detected")
	ErrMaxHopsExceeded  = errors.New("maximum shortener redirect hops exceeded")
	ErrInvalidRedirect  = errors.New("invalid redirect location")
)

func IsShortener(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return false
	}

	host := strings.ToLower(u.Hostname())
	for _, domain := range config.ShortenerDomains {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

func ResolveShortener(ctx context.Context, rawURL string, maxHops int) (string, []string, error) {
	if !IsShortener(rawURL) {
		return rawURL, nil, nil
	}

	if maxHops <= 0 {
		maxHops = config.MaxShortenerHops
	}

	client := &http.Client{
		Transport: NewSafeTransport(),
		Timeout:   config.SandboxTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return resolveShortenerWithClient(ctx, rawURL, maxHops, client)
}

func resolveShortenerWithClient(ctx context.Context, rawURL string, maxHops int, client *http.Client) (string, []string, error) {
	currentURL := rawURL
	var hops []string
	visited := make(map[string]bool)

	safeClient := *client
	safeClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	for i := 0; i < maxHops; i++ {
		if visited[currentURL] {
			return currentURL, hops, ErrShortenerLoop
		}
		visited[currentURL] = true

		req, err := http.NewRequestWithContext(ctx, http.MethodHead, currentURL, nil)
		if err != nil {
			return currentURL, hops, err
		}
		req.Header.Set("User-Agent", "JevBeez-SafetyBot/1.0")

		resp, err := safeClient.Do(req)
		if err != nil {
			// Some servers reject HEAD requests; fallback to safe GET with immediate body close
			reqGet, getErr := http.NewRequestWithContext(ctx, http.MethodGet, currentURL, nil)
			if getErr != nil {
				return currentURL, hops, err
			}
			reqGet.Header.Set("User-Agent", "JevBeez-SafetyBot/1.0")
			resp, err = safeClient.Do(reqGet)
			if err != nil {
				return currentURL, hops, err
			}
		}

		if resp.Body != nil {
			_, _ = io.CopyN(io.Discard, resp.Body, 512)
			_ = resp.Body.Close()
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			loc := resp.Header.Get("Location")
			if loc == "" {
				break
			}

			parsedLoc, parseErr := url.Parse(loc)
			if parseErr != nil {
				return currentURL, hops, fmt.Errorf("%w: %s", ErrInvalidRedirect, loc)
			}

			baseParsed, _ := url.Parse(currentURL)
			nextURL := baseParsed.ResolveReference(parsedLoc).String()

			hops = append(hops, nextURL)
			currentURL = nextURL
			continue
		}

		break
	}

	if len(hops) >= maxHops {
		return currentURL, hops, ErrMaxHopsExceeded
	}

	return currentURL, hops, nil
}
