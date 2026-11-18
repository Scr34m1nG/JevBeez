package scraper

import (
	"html"
	"io"
	"strings"

	nethtml "golang.org/x/net/html"
)

func ExtractMetaTags(r io.Reader) (map[string]string, error) {
	tags := make(map[string]string)
	tokenizer := nethtml.NewTokenizer(r)

	for {
		tt := tokenizer.Next()
		switch tt {
		case nethtml.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return tags, nil
			}
			return tags, err
		case nethtml.StartTagToken, nethtml.SelfClosingTagToken:
			tok := tokenizer.Token()
			if tok.Data == "meta" {
				var key, val string
				for _, attr := range tok.Attr {
					attrKey := strings.ToLower(attr.Key)
					if attrKey == "property" || attrKey == "name" {
						key = strings.ToLower(attr.Val)
					} else if attrKey == "content" {
						val = attr.Val
					}
				}
				if key != "" && val != "" {
					tags[key] = strings.TrimSpace(html.UnescapeString(val))
				}
			} else if tok.Data == "body" {
				// Once in body, stop scanning meta tags to conserve CPU
				return tags, nil
			}
		}
	}
}

func ExtractParagraphText(rawHTML string) string {
	doc, err := nethtml.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return ""
	}

	var paragraphs []string
	var walk func(*nethtml.Node)
	walk = func(n *nethtml.Node) {
		if n.Type == nethtml.ElementNode && n.Data == "p" {
			var sb strings.Builder
			var collectText func(*nethtml.Node)
			collectText = func(child *nethtml.Node) {
				if child.Type == nethtml.TextNode {
					sb.WriteString(child.Data)
				}
				for c := child.FirstChild; c != nil; c = c.NextSibling {
					collectText(c)
				}
			}
			collectText(n)
			text := strings.TrimSpace(sb.String())
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)
	joined := strings.Join(paragraphs, "\n\n")
	return html.UnescapeString(strings.TrimSpace(joined))
}
