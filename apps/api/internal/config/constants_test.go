package config

import (
	"testing"
)

func TestThreatConstants(t *testing.T) {
	if VerdictPhishing != "PHISHING" {
		t.Errorf("expected VerdictPhishing to be 'PHISHING', got %s", VerdictPhishing)
	}
	if VerdictMalware != "MALWARE" {
		t.Errorf("expected VerdictMalware to be 'MALWARE', got %s", VerdictMalware)
	}
	if VerdictScam != "SCAM" {
		t.Errorf("expected VerdictScam to be 'SCAM', got %s", VerdictScam)
	}

	if len(DangerousExtensions) == 0 {
		t.Errorf("expected DangerousExtensions not to be empty")
	}
	if len(HighRiskTLDs) == 0 {
		t.Errorf("expected HighRiskTLDs not to be empty")
	}
	if len(ShortenerDomains) == 0 {
		t.Errorf("expected ShortenerDomains not to be empty")
	}
	if len(BrandProfiles) == 0 {
		t.Errorf("expected BrandProfiles not to be empty")
	}

	if MaxShortenerHops != 5 {
		t.Errorf("expected MaxShortenerHops to be 5, got %d", MaxShortenerHops)
	}
	if MaxSandboxPayloadBytes != 100*1024 {
		t.Errorf("expected MaxSandboxPayloadBytes to be 102400, got %d", MaxSandboxPayloadBytes)
	}
}

func TestUnwrapperConfigValues(t *testing.T) {
	if VerdictAds != "ADS" {
		t.Fatalf("expected ADS verdict, got %s", VerdictAds)
	}
	if MaxUnwrapDepth != 3 {
		t.Fatalf("expected MaxUnwrapDepth 3, got %d", MaxUnwrapDepth)
	}
	if MinHexSubdomainLength != 16 {
		t.Fatalf("expected MinHexSubdomainLength 16, got %d", MinHexSubdomainLength)
	}

	params := map[string]bool{}
	for _, p := range WrapperParams {
		params[p] = true
	}
	for _, want := range []string{"adurl", "to", "u", "q", "url", "redirect"} {
		if !params[want] {
			t.Errorf("missing wrapper param %q", want)
		}
	}

	shops := map[string]bool{}
	for _, d := range EcommerceDomains {
		shops[d] = true
	}
	for _, want := range []string{"ebay.com", "amazon.com", "shopee.co.id", "tokopedia.com", "lazada.co.id"} {
		if !shops[want] {
			t.Errorf("missing ecommerce domain %q", want)
		}
	}

	if len(WrapperHostPatterns) == 0 {
		t.Fatal("WrapperHostPatterns must not be empty")
	}
}
