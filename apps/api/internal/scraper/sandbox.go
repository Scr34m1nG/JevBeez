package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/security"
)

func InspectFormThreats(ctx context.Context, targetURL string) ([]string, string, error) {
	client := security.NewSafeHTTPClient()
	client.Timeout = config.SandboxTimeout
	return InspectFormThreatsWithClient(ctx, targetURL, client)
}

func InspectFormThreatsWithClient(ctx context.Context, targetURL string, client *http.Client) ([]string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "JevBeez-ThreatSandbox/1.0 (Safe Inspector; +https://jevbeez.com)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Protect against decompression or oversized payloads
	limitedReader := io.LimitReader(resp.Body, config.MaxSandboxPayloadBytes)
	doc, err := html.Parse(limitedReader)
	if err != nil {
		return nil, "", err
	}

	targetParsed, _ := url.Parse(targetURL)
	targetHost := ""
	if targetParsed != nil {
		targetHost = strings.ToLower(targetParsed.Hostname())
	}

	var threats []string
	threatSet := make(map[string]bool)
	var title string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = strings.TrimSpace(n.FirstChild.Data)
			}

			if n.Data == "form" {
				action := ""
				for _, attr := range n.Attr {
					if attr.Key == "action" {
						action = strings.ToLower(strings.TrimSpace(attr.Val))
						break
					}
				}

				if action != "" {
					if strings.Contains(action, "api.telegram.org") {
						threatSet["Formulir mengirim data curian secara langsung ke Telegram Bot (modus exfiltrasi phising)."] = true
					} else if parsedAction, parseErr := url.Parse(action); parseErr == nil && parsedAction.Host != "" {
						actionHost := strings.ToLower(parsedAction.Hostname())
						if targetHost != "" && actionHost != targetHost && !strings.HasSuffix(actionHost, "."+targetHost) {
							threatSet[fmt.Sprintf("Formulir mengarahkan kiriman data ke domain eksternal berbeda (%s).", actionHost)] = true
						}
					}
				}

				// Inspect inputs inside this form
				inspectFormInputs(n, threatSet)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	for t := range threatSet {
		threats = append(threats, t)
	}

	return threats, title, nil
}

func inspectFormInputs(formNode *html.Node, threatSet map[string]bool) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			inputType := ""
			inputName := ""
			inputPlaceholder := ""

			for _, attr := range n.Attr {
				switch attr.Key {
				case "type":
					inputType = strings.ToLower(strings.TrimSpace(attr.Val))
				case "name":
					inputName = strings.ToLower(strings.TrimSpace(attr.Val))
				case "placeholder":
					inputPlaceholder = strings.ToLower(strings.TrimSpace(attr.Val))
				}
			}

			combined := inputName + " " + inputPlaceholder

			if inputType == "password" {
				threatSet["Terdapat kolom input kata sandi (password)."] = true
			}

			if strings.Contains(combined, "pin") || strings.Contains(combined, "mpin") {
				threatSet["Meminta input PIN rahasia transaksi / m-Banking."] = true
			}

			if strings.Contains(combined, "otp") || strings.Contains(combined, "token") {
				threatSet["Meminta kode OTP / One-Time Password yang seharusnya rahasia."] = true
			}

			if strings.Contains(combined, "cvv") || strings.Contains(combined, "cvc") || strings.Contains(combined, "kartu") || strings.Contains(combined, "card") {
				threatSet["Meminta data nomor kartu debit/kredit atau kode CVV."] = true
			}

			if strings.Contains(combined, "rekening") || strings.Contains(combined, "account") || strings.Contains(combined, "nik") || strings.Contains(combined, "ktp") {
				threatSet["Meminta data identitas pribadi atau nomor rekening bank."] = true
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(formNode)
}
