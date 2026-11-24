package security

import (
	"testing"
	"jevbeez/api/internal/config"
)

func TestCheckURLHeuristics_MalwareAPK(t *testing.T) {
	urls := []string{
		"https://undangan-nikah-digital.xyz/surat_undangan.apk",
		"http://cek-resi-jnt.top/paket_tracking.exe",
		"https://surat-tilang.buzz/dokumen_tilang.apk?v=1",
	}

	for _, u := range urls {
		res := CheckURLHeuristics(u)
		if !res.IsThreat {
			t.Errorf("expected %s to be flagged as threat", u)
		}
		if res.ThreatType != config.VerdictMalware {
			t.Errorf("expected threat type %s, got %s", config.VerdictMalware, res.ThreatType)
		}
		if len(res.Reasons) == 0 {
			t.Errorf("expected reasons not to be empty")
		}
	}
}

func TestCheckURLHeuristics_BrandSpoofingAndSubdomainTrap(t *testing.T) {
	tests := []struct {
		url          string
		expectedType string
		brand        string
	}{
		{"https://bca.co.id.login-verifikasi.xyz/akun", config.VerdictPhishing, "bca"},
		{"http://promo-gebyar-bri.top/hadiah", config.VerdictPhishing, "bri"},
		{"https://dana-kaget-saldo.cfd/klaim", config.VerdictPhishing, "dana"},
		{"https://mandiri-online-secure.xyz/auth", config.VerdictPhishing, "mandiri"},
	}

	for _, tt := range tests {
		res := CheckURLHeuristics(tt.url)
		if !res.IsThreat {
			t.Errorf("expected %s to be flagged as threat", tt.url)
		}
		if res.ThreatType != tt.expectedType {
			t.Errorf("expected %s, got %s for url %s", tt.expectedType, res.ThreatType, tt.url)
		}
		if res.TargetBrand != tt.brand {
			t.Errorf("expected target brand %s, got %s", tt.brand, res.TargetBrand)
		}
	}
}

func TestCheckURLHeuristics_AtSymbolDeception(t *testing.T) {
	url := "https://bca.co.id@evil-attacker.com/login"
	res := CheckURLHeuristics(url)
	if !res.IsThreat {
		t.Errorf("expected @ trick to be flagged as threat")
	}
	if res.ThreatType != config.VerdictPhishing {
		t.Errorf("expected PHISHING, got %s", res.ThreatType)
	}
}

func TestCheckURLHeuristics_RawIP(t *testing.T) {
	url := "http://185.220.101.5/bca/login.php"
	res := CheckURLHeuristics(url)
	if !res.IsThreat {
		t.Errorf("expected raw IP host with brand path to be flagged as threat")
	}
	if res.ThreatType != config.VerdictPhishing {
		t.Errorf("expected PHISHING, got %s", res.ThreatType)
	}
}

func TestCheckURLHeuristics_Punycode(t *testing.T) {
	url := "https://xn--bca-9la.com/login"
	res := CheckURLHeuristics(url)
	if !res.IsThreat {
		t.Errorf("expected punycode domain to be flagged as threat")
	}
}

func TestCheckURLHeuristics_LegitimateURL(t *testing.T) {
	legitURLs := []string{
		"https://www.bca.co.id/id/tentang-bca",
		"https://klikbca.com",
		"https://bankmandiri.co.id/promo",
		"https://news.detik.com/berita/123/presiden-resmikan-jalan",
		"https://kompas.com/read/2026/01/01/berita",
	}

	for _, u := range legitURLs {
		res := CheckURLHeuristics(u)
		if res.IsThreat {
			t.Errorf("expected legitimate url %s NOT to be flagged as threat, but got reasons: %v", u, res.Reasons)
		}
	}
}
