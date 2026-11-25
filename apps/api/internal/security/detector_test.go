package security

import (
	"testing"
)

func TestDetectExploit(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantBlock  bool
		wantReason string
	}{
		// Legitimate URLs
		{
			name:      "legitimate news url",
			url:       "https://news.detik.com/berita/d-12345/judul-berita-terbaru",
			wantBlock: false,
		},
		{
			name:      "legitimate threads share",
			url:       "https://www.threads.com/share/D7VmU4Wtk",
			wantBlock: false,
		},
		{
			name:      "legitimate twitter status",
			url:       "https://x.com/jack/status/20",
			wantBlock: false,
		},
		{
			name:      "legitimate php news url",
			url:       "https://www.kompas.com/read.php?id=99281",
			wantBlock: false,
		},
		{
			name:      "legitimate news url with ellipsis slug",
			url:       "https://news.detik.com/read/2026/01/01/isu-gempa.../",
			wantBlock: false,
		},
		{
			name:      "legitimate tech article with whoami in slug",
			url:       "https://linux.org/article/what-is-whoami",
			wantBlock: false,
		},

		// SSRF & Loopback / Private Network
		{
			name:       "ssrf localhost",
			url:        "http://localhost:8080/admin",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf 127.0.0.1",
			url:        "http://127.0.0.1:6379",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf cloud metadata aws",
			url:        "http://169.254.169.254/latest/meta-data/",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf private ip 192.168",
			url:        "http://192.168.1.1/router",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf private ip 10.x",
			url:        "http://10.0.0.5:9000",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf wildcard dns nip.io",
			url:        "http://127.0.0.1.nip.io/api",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf decimal ip integer",
			url:        "http://2130706433/",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf hex ip",
			url:        "http://0x7f000001/",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf file protocol",
			url:        "file:///etc/passwd",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},
		{
			name:       "ssrf gopher protocol",
			url:        "gopher://127.0.0.1:6379/_ping",
			wantBlock:  true,
			wantReason: ReasonSSRFAttempt,
		},

		// Path Traversal
		{
			name:       "path traversal plain",
			url:        "https://example.com/../../etc/passwd",
			wantBlock:  true,
			wantReason: ReasonPathTraversal,
		},
		{
			name:       "path traversal encoded",
			url:        "https://example.com/%2e%2e/%2e%2e/etc/passwd",
			wantBlock:  true,
			wantReason: ReasonPathTraversal,
		},
		{
			name:       "sensitive file env",
			url:        "https://target.com/.env",
			wantBlock:  true,
			wantReason: ReasonPathTraversal,
		},
		{
			name:       "sensitive file git config",
			url:        "https://target.com/.git/config",
			wantBlock:  true,
			wantReason: ReasonPathTraversal,
		},

		// Injection attacks
		{
			name:       "xss injection script tag",
			url:        "https://example.com/?q=<script>alert(1)</script>",
			wantBlock:  true,
			wantReason: ReasonInjectionAttempt,
		},
		{
			name:       "xss javascript pseudo protocol",
			url:        "javascript:alert(document.cookie)",
			wantBlock:  true,
			wantReason: ReasonInjectionAttempt,
		},
		{
			name:       "sqli union select",
			url:        "https://example.com/articles?id=1 union select 1,version()--",
			wantBlock:  true,
			wantReason: ReasonInjectionAttempt,
		},
		{
			name:       "sqli boolean based",
			url:        "https://example.com/items?cat=' or '1'='1",
			wantBlock:  true,
			wantReason: ReasonInjectionAttempt,
		},
		{
			name:       "command injection pipe whoami",
			url:        "https://example.com/ping?ip=1.1.1.1|whoami",
			wantBlock:  true,
			wantReason: ReasonInjectionAttempt,
		},
		{
			name:       "command injection semicolon id",
			url:        "https://example.com/ping?ip=1.1.1.1;id",
			wantBlock:  true,
			wantReason: ReasonInjectionAttempt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocked, reason, _ := DetectExploit(tt.url)
			if blocked != tt.wantBlock {
				t.Errorf("DetectExploit(%q) blocked = %v, want %v", tt.url, blocked, tt.wantBlock)
			}
			if tt.wantBlock && reason != tt.wantReason {
				t.Errorf("DetectExploit(%q) reason = %v, want %v", tt.url, reason, tt.wantReason)
			}
		})
	}
}
