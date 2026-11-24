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
