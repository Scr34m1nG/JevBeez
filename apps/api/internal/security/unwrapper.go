package security

import (
	"net/url"
	"strings"

	"jevbeez/api/internal/config"
)

type UnwrapResult struct {
	FinalURL string
	Chain    []string
	Wrappers []string
	Depth    int
}

func UnwrapRedirectChain(rawURL string) *UnwrapResult {
	res := &UnwrapResult{FinalURL: rawURL}

	current := rawURL
	visited := map[string]bool{rawURL: true}

	for depth := 0; depth < config.MaxUnwrapDepth; depth++ {
		next, ok := extractWrappedURL(current)
		if !ok || visited[next] {
			break
		}

		visited[next] = true
		res.Chain = append(res.Chain, next)
		res.FinalURL = next
		res.Depth++
		current = next
	}

	for _, candidate := range append([]string{rawURL}, res.Chain...) {
		host := wrapperHost(candidate)
		if host == "" {
			continue
		}
		if !containsString(res.Wrappers, host) {
			res.Wrappers = append(res.Wrappers, host)
		}
	}

	return res
}

func extractWrappedURL(rawURL string) (string, bool) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return "", false
	}

	query := parsed.Query()
	for _, param := range config.WrapperParams {
		values, ok := query[param]
		if !ok {
			continue
		}
		for _, value := range values {
			if candidate, ok := parseEmbeddedURL(value); ok {
				return candidate, true
			}
		}
	}

	return "", false
}

func parseEmbeddedURL(value string) (string, bool) {
	candidate := strings.TrimSpace(value)
	for i := 0; i < 2; i++ {
		decoded, err := url.QueryUnescape(candidate)
		if err != nil || decoded == candidate {
			break
		}
		candidate = strings.TrimSpace(decoded)
	}

	parsed, err := url.Parse(candidate)
	if err != nil {
		return "", false
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", false
	}
	if !strings.Contains(parsed.Hostname(), ".") {
		return "", false
	}

	return parsed.String(), true
}

func wrapperHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return ""
	}

	host := strings.ToLower(parsed.Hostname())
	target := host + strings.ToLower(parsed.Path)
	for _, pattern := range config.WrapperHostPatterns {
		if strings.Contains(target, pattern) {
			return host
		}
	}
	return ""
}

func IsEcommerceURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return false
	}

	host := strings.ToLower(parsed.Hostname())
	for _, domain := range config.EcommerceDomains {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

func containsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
