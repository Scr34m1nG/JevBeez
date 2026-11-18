package scraper

import (
	"strings"
	"testing"
)

func TestExtractMetaTags(t *testing.T) {
	htmlData := `<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
	<meta property="og:title" content="Judul Postingan Heboh" />
	<meta property="og:description" content="Ini adalah caption postingan medsos yang sangat penting untuk diverifikasi." />
	<meta name="twitter:description" content="Caption twitter alternatif" />
	<meta name="author" content="Budi" />
</head>
<body><p>Hello World</p></body>
</html>`

	tags, err := ExtractMetaTags(strings.NewReader(htmlData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tags["og:title"] != "Judul Postingan Heboh" {
		t.Errorf("expected og:title 'Judul Postingan Heboh', got '%s'", tags["og:title"])
	}

	if tags["og:description"] != "Ini adalah caption postingan medsos yang sangat penting untuk diverifikasi." {
		t.Errorf("expected og:description match, got '%s'", tags["og:description"])
	}

	if tags["twitter:description"] != "Caption twitter alternatif" {
		t.Errorf("expected twitter:description match, got '%s'", tags["twitter:description"])
	}
}

func TestExtractParagraphText(t *testing.T) {
	oembedHTML := `<blockquote class="twitter-tweet">
<p lang="in" dir="ltr">BREAKING: Terjadi gempa susulan berkekuatan 6.5 SR di pesisir selatan &amp; warga diminta waspada! <a href="https://t.co/xyz">pic.twitter.com/xyz</a></p>
&mdash; Akun Berita (@berita) <a href="https://twitter.com/berita/status/123">November 19, 2026</a>
</blockquote>`

	extracted := ExtractParagraphText(oembedHTML)
	expectedContains := "BREAKING: Terjadi gempa susulan berkekuatan 6.5 SR di pesisir selatan & warga diminta waspada!"

	if !strings.Contains(extracted, expectedContains) {
		t.Errorf("expected extracted text to contain '%s', got '%s'", expectedContains, extracted)
	}

	if strings.Contains(extracted, "&amp;") {
		t.Errorf("expected HTML entities to be decoded, found &amp;")
	}

	if strings.Contains(extracted, "<a href=") {
		t.Errorf("expected HTML tags to be stripped, found <a href=")
	}
}
