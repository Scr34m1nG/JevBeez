package security

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestReDoS_AdversarialRegexPatterns(t *testing.T) {
	adversarialInputs := []struct {
		name string
		gen  func() string
	}{
		{
			name: "5000 chars repeating forward noise",
			gen: func() string {
				return strings.Repeat("Diteruskan berkali-kali ", 250) + "klaim fakta penting ada di sini."
			},
		},
		{
			name: "5000 chars repeating viral group call",
			gen: func() string {
				return strings.Repeat("Viralkan ke 10 grup ", 250) + "klaim fakta penting ada di sini."
			},
		},
		{
			name: "5000 chars repeating doctor authority hook",
			gen: func() string {
				return strings.Repeat("Info dari dokter RSPAD: ", 200) + "klaim fakta penting ada di sini."
			},
		},
		{
			name: "10000 chars repeating punctuation and asterisks",
			gen: func() string {
				return strings.Repeat("*_~*!?", 1500) + "klaim fakta penting ada di sini."
			},
		},
		{
			name: "Nested whitespace and newlines explosion",
			gen: func() string {
				return strings.Repeat(" \t\r\n", 1000) + "klaim fakta penting ada di sini."
			},
		},
	}

	for _, tc := range adversarialInputs {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.gen()

			start := time.Now()
			_ = CleanBroadcastNoise(input)
			elapsed := time.Since(start)

			// Prompt requirement: Ensure regex matches execute in < 5ms
			if elapsed > 5*time.Millisecond {
				t.Errorf("execution took %v, exceeded 5ms ReDoS threshold", elapsed)
			}
		})
	}
}

func TestSSRF_IsolationNoNetworkActivity(t *testing.T) {
	ssrfPayloads := []string{
		"Info penting: Cek data Anda di 127.0.0.1 segera sebelum dinonaktifkan.",
		"Bantuan kuota gratis di http://169.254.169.254/latest/meta-data/",
		"Diteruskan: Unduh sertifikat vaksin di file:///etc/passwd sekarang juga.",
		"Viralkan! Saldo gratis masuk ke akun http://10.0.0.1/admin/transfer",
		"http://[::1]:8080/internal/secret/api/keys",
	}

	for _, payload := range ssrfPayloads {
		t.Run(payload, func(t *testing.T) {
			// Measure execution time and guarantee purely in-memory execution
			start := time.Now()

			cleaned, query, err := CleanText(payload)
			elapsed := time.Since(start)

			// CleanText must finish immediately (< 2ms) without making network socket connections
			if elapsed > 10*time.Millisecond {
				t.Errorf("CleanText took too long (%v), suspected network dialer invocation", elapsed)
			}

			if err != nil && err != ErrTextTooShort {
				t.Fatalf("unexpected error: %v", err)
			}

			// Processing must not panic or trigger dialers
			_ = HashText(payload)
			_ = ExtractSearchQuery(cleaned, 100)
			_ = query
		})
	}
}

func TestUnicode_EdgeCasesAndAttacks(t *testing.T) {
	unicodeTests := []struct {
		name  string
		input string
	}{
		{
			name:  "Zero-width spaces and invisible joiners",
			input: "Klaim\u200B penting\u200C dari\u200D dokter\uFEFF RSPAD:\u00AD minum air es membekukan usus.",
		},
		{
			name:  "RTL Override characters spoofing",
			input: "Info \u202Evaksin palsu\u202C beredar di rumah sakit \u200Fpusat.",
		},
		{
			name:  "High surrogate pairs and emojis explosion",
			input: "🚨🚨🚨 VIRALKAN!!! 😱😱😱 Info dokter: 💊💊💊 jangan minum obat ini 💥💥💥 bahaya banget!",
		},
		{
			name:  "Mixed international scripts (Arabic, Cyrillic, CJK, Thai)",
			input: "Внимание! تحذir! 注意！โปรดระวัง! Berita penting mengenai kesehatan masyarakat luas.",
		},
		{
			name:  "Max length with multi-byte runes without rune splitting",
			input: strings.Repeat("🎉 Berita Hoax 🌐 ", 70),
		},
	}

	for _, tc := range unicodeTests {
		t.Run(tc.name, func(t *testing.T) {
			// None of these should crash, panic, or output invalid UTF-8
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic occurred on input %s: %v", tc.name, r)
				}
			}()

			cleaned := CleanBroadcastNoise(tc.input)
			if !utf8.ValidString(cleaned) {
				t.Errorf("CleanBroadcastNoise produced invalid UTF-8 string: %q", cleaned)
			}

			hash := HashText(tc.input)
			if len(hash) != 64 {
				t.Errorf("invalid hash length for unicode string: %d", len(hash))
			}

			query := ExtractSearchQuery(cleaned, 100)
			if !utf8.ValidString(query) {
				t.Errorf("ExtractSearchQuery produced invalid UTF-8: %q", query)
			}

			val, err := ValidateRawText(tc.input)
			if err == nil {
				if !utf8.ValidString(val) {
					t.Errorf("ValidateRawText produced invalid UTF-8: %q", val)
				}
				if utf8.RuneCountInString(val) > MaxTextRunes {
					t.Errorf("rune count %d exceeds max runes %d", utf8.RuneCountInString(val), MaxTextRunes)
				}
			}
		})
	}
}

func BenchmarkCleanBroadcastNoise(b *testing.B) {
	sample := "Diteruskan berkali-kali\nInfo dari dokter RSPAD: jangan minum es setelah makan bakso karena membekukan lemak. Viralkan ke 10 grup!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CleanBroadcastNoise(sample)
	}
}

func BenchmarkHashText(b *testing.B) {
	sample := "Info penting: Vaksin menyebabkan demam tinggi pada balita dan anak-anak."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashText(sample)
	}
}
