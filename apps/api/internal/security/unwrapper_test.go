package security

import "testing"

func TestUnwrapRedirectChain(t *testing.T) {
	userLink := "https://www.google.com/pagead/aclk?sa=L&adurl=https%3A%2F%2Fvk.com%2Faway.php%3Fto%3Dhttps%3A%2F%2Fbaf9105c4584fac832929f959728944d.cavemin.com%2Fsttodm%2F%3F373737%7Cn5Gcuw&num=1"

	tests := []struct {
		name         string
		rawURL       string
		wantFinal    string
		wantDepth    int
		wantWrappers []string
	}{
		{
			name:         "google ads to vk away to obfuscated domain",
			rawURL:       userLink,
			wantFinal:    "https://baf9105c4584fac832929f959728944d.cavemin.com/sttodm/?373737|n5Gcuw",
			wantDepth:    2,
			wantWrappers: []string{"www.google.com", "vk.com"},
		},
		{
			name:      "plain url returns itself",
			rawURL:    "https://news.com/article/1",
			wantFinal: "https://news.com/article/1",
			wantDepth: 0,
		},
		{
			name:      "double encoded value is still unwrapped",
			rawURL:    "https://l.facebook.com/l.php?u=https%253A%252F%252Fshopee.co.id%252Fitem%252F1",
			wantFinal: "https://shopee.co.id/item/1",
			wantDepth: 1,
		},
		{
			name:      "non-url value is ignored",
			rawURL:    "https://google.com/url?q=javascript%3Aalert(1)",
			wantFinal: "https://google.com/url?q=javascript%3Aalert(1)",
			wantDepth: 0,
		},
		{
			name:      "nested self reference terminates at depth cap",
			rawURL:    "https://x.com/r?u=https%3A%2F%2Fx.com%2Fr%3Fu%3Dhttps%253A%252F%252Fx.com%252Fr",
			wantFinal: "https://x.com/r",
			wantDepth: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := UnwrapRedirectChain(tc.rawURL)
			if got.FinalURL != tc.wantFinal {
				t.Errorf("expected final %q, got %q", tc.wantFinal, got.FinalURL)
			}
			if got.Depth != tc.wantDepth {
				t.Errorf("expected depth %d, got %d", tc.wantDepth, got.Depth)
			}
			if len(got.Chain) != tc.wantDepth {
				t.Errorf("expected chain length %d, got %d", tc.wantDepth, len(got.Chain))
			}
			for _, w := range tc.wantWrappers {
				found := false
				for _, gw := range got.Wrappers {
					if gw == w {
						found = true
					}
				}
				if !found {
					t.Errorf("expected wrapper host %q in %v", w, got.Wrappers)
				}
			}
		})
	}
}

func TestIsEcommerceURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://www.ebay.com/itm/123", true},
		{"https://shopee.co.id/product/1/2", true},
		{"https://tokopedia.com/store", true},
		{"https://notebay.com/item", false},
		{"https://news.com/article", false},
	}
	for _, tc := range tests {
		if got := IsEcommerceURL(tc.url); got != tc.want {
			t.Errorf("IsEcommerceURL(%q) = %v, want %v", tc.url, got, tc.want)
		}
	}
}
