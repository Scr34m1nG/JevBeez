package security

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidScheme = errors.New("url scheme must be http or https")
	ErrEmptyHost     = errors.New("url host cannot be empty")
	ErrInvalidHost   = errors.New("url host must be a valid domain")
	ErrSelfURL       = errors.New("cannot analyze JevBeez itself, please provide another website URL")
)

var trackingParams = map[string]struct{}{
	"fbclid":                 {},
	"gclid":                  {},
	"dclid":                  {},
	"msclkid":                {},
	"ref":                    {},
	"source":                 {},
	"_ga":                    {},
	"_gl":                    {},
	"mc_cid":                 {},
	"mc_eid":                 {},
	"yclid":                  {},
	"_hsenc":                 {},
	"_hsmi":                  {},
	"vero_id":                {},
	"vero_conv":              {},
	"nr_email_referer":       {},
	"igshid":                 {},
	"si":                     {},
	"mibextid":               {},
	"xmt":                    {},
	"igsh":                   {},
	"s":                      {},
	"t":                      {},
	"context":                {},
}

func ValidateAndNormalizeURL(rawURL string) (string, error) {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidScheme
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return "", ErrEmptyHost
	}

	if !strings.Contains(host, ".") {
		return "", ErrInvalidHost
	}

	if host == "jevbeez.hityo.cloud" || strings.HasSuffix(host, ".jevbeez.hityo.cloud") ||
		host == "jevbeez.com" || strings.HasSuffix(host, ".jevbeez.com") {
		return "", ErrSelfURL
	}

	query := parsed.Query()
	for param := range query {
		lower := strings.ToLower(param)
		if strings.HasPrefix(lower, "utm_") {
			query.Del(param)
			continue
		}
		if _, exists := trackingParams[lower]; exists {
			query.Del(param)
		}
	}

	normalized := url.URL{
		Scheme:   scheme,
		Host:     parsed.Host,
		Path:     parsed.Path,
		RawQuery: query.Encode(),
	}

	res := normalized.String()
	res = strings.TrimSuffix(res, "?")
	return res, nil
}
