package security

import (
	"testing"
)

func TestValidateAndNormalizeURL(t *testing.T) {
	tests := []struct {
		name      string
		rawURL    string
		wantURL   string
		expectErr bool
	}{
		{
			name:      "clean standard url",
			rawURL:    "https://news.com/article/1",
			wantURL:   "https://news.com/article/1",
			expectErr: false,
		},
		{
			name:      "strips tracking utm and fbclid params",
			rawURL:    "https://news.com/article/1?utm_source=twitter&utm_medium=social&fbclid=xyz123&page=2",
			wantURL:   "https://news.com/article/1?page=2",
			expectErr: false,
		},
		{
			name:      "rejects invalid scheme",
			rawURL:    "ftp://news.com/file.txt",
			expectErr: true,
		},
		{
			name:      "rejects invalid host without dot",
			rawURL:    "http://localhost/path",
			expectErr: true,
		},
		{
			name:      "rejects empty string",
			rawURL:    "",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ValidateAndNormalizeURL(tc.rawURL)
			if (err != nil) != tc.expectErr {
				t.Fatalf("expected error: %v, got error: %v", tc.expectErr, err)
			}
			if !tc.expectErr && got != tc.wantURL {
				t.Errorf("expected url: %s, got: %s", tc.wantURL, got)
			}
		})
	}
}
