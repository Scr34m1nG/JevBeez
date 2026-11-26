package security

import (
	"strings"
	"testing"
	"unicode/utf8"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/decision"
	"jevbeez/api/internal/model"
)

func TestCleanBroadcastNoise(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Forwarded and doctor authority hook",
			input:    "Diteruskan berkali-kali\nInfo dari dokter RSPAD: jangan minum es setelah makan bakso karena membekukan lemak. Viralkan ke 10 grup!",
			expected: "jangan minum es setelah makan bakso karena membekukan lemak.",
		},
		{
			name:     "Forwarded single and IDI doctor hook with bold asterisks",
			input:    "*Diteruskan*\n*Kabar dari dokter IDI*: Vaksin influenza menyebabkan kemandulan permanen. *Viralkan!*",
			expected: "Vaksin influenza menyebabkan kemandulan permanen.",
		},
		{
			name:     "Viral call to all groups and salutation footer",
			input:    "Pemerintah membagikan bantuan sosial BBM gratis sebesar 5 juta rupiah per KTP.\n\nSebarkan ke semua grup! Semoga bermanfaat.",
			expected: "Pemerintah membagikan bantuan sosial BBM gratis sebesar 5 juta rupiah per KTP.",
		},
		{
			name:     "Chain message with reminder and multiple noise markers",
			input:    "Diteruskan seperti yang diterima:\nPeringatan dari Kemenkes: Air kelapa hijau campur garam bisa menyembuhkan stroke seketika. Mohon maaf sekadar mengingatkan, jangan berhenti di anda!",
			expected: "Air kelapa hijau campur garam bisa menyembuhkan stroke seketika.",
		},
		{
			name:     "Clean input without noise remains intact",
			input:    "Presiden meresmikan jalan tol baru sepanjang 50 kilometer di Jawa Barat pada hari Senin.",
			expected: "Presiden meresmikan jalan tol baru sepanjang 50 kilometer di Jawa Barat pada hari Senin.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CleanBroadcastNoise(tc.input)
			if got != tc.expected {
				t.Errorf("\nInput:    %q\nExpected: %q\nGot:      %q", tc.input, tc.expected, got)
			}
		})
	}
}

func TestExtractSearchQuery(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLen    int
		checkLess bool
	}{
		{
			name:      "Short claim stays under 100",
			input:     "Minum air es setelah makan bakso membekukan lemak di usus",
			maxLen:    100,
			checkLess: true,
		},
		{
			name:      "Long article snippet is truncated cleanly at word boundary under 100 chars",
			input:     "Kementerian Komunikasi dan Digital mengumumkan bahwa seluruh nomor WhatsApp yang tidak terdaftar ulang di situs resmi komdigi-verifikasi.xyz akan dinonaktifkan dalam waktu 24 jam ke depan secara permanen.",
			maxLen:    100,
			checkLess: true,
		},
		{
			name:      "First sentence extraction from multi-sentence claim",
			input:     "Vaksin cacar monyet mengandung zat berbahaya bagi janin. Segera konsultasikan ke puskesmas terdekat bila Anda mengalaminya.",
			maxLen:    100,
			checkLess: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			query := ExtractSearchQuery(tc.input, tc.maxLen)
			runeCount := utf8.RuneCountInString(query)

			if runeCount >= 100 {
				t.Errorf("expected query < 100 chars, got %d (%q)", runeCount, query)
			}
			if strings.HasSuffix(query, " ") || strings.HasSuffix(query, ",") || strings.HasSuffix(query, ".") {
				t.Errorf("query has trailing punctuation or space: %q", query)
			}
			if len(query) == 0 {
				t.Errorf("extracted query should not be empty")
			}
		})
	}
}

func TestValidateRawText_LengthLimits(t *testing.T) {
	t.Run("Rejects empty string", func(t *testing.T) {
		_, err := ValidateRawText("   \n\t  ")
		if err != ErrEmptyText {
			t.Errorf("expected ErrEmptyText, got %v", err)
		}
	})

	t.Run("Rejects input shorter than 25 runes", func(t *testing.T) {
		_, err := ValidateRawText("Halo halo halo")
		if err != ErrTextTooShort {
			t.Errorf("expected ErrTextTooShort, got %v", err)
		}
	})

	t.Run("Rejects input with less than 4 words", func(t *testing.T) {
		_, err := ValidateRawText("Supercalifragilisticexpialidocious satu dua")
		if err != ErrTextTooShort {
			t.Errorf("expected ErrTextTooShort, got %v", err)
		}
	})

	t.Run("Accepts valid input with >= 25 runes and >= 4 words", func(t *testing.T) {
		val, err := ValidateRawText("Ini adalah berita penting mengenai kesehatan masyarakat.")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val == "" {
			t.Errorf("expected non-empty validated text")
		}
	})

	t.Run("Truncates input exceeding 1000 runes cleanly", func(t *testing.T) {
		longText := strings.Repeat("Kata fakta dan berita terpercaya ", 50)
		if utf8.RuneCountInString(longText) <= 1000 {
			t.Fatalf("test fixture must be > 1000 runes, got %d", utf8.RuneCountInString(longText))
		}

		validated, err := ValidateRawText(longText)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		runeCount := utf8.RuneCountInString(validated)
		if runeCount != 1000 {
			t.Errorf("expected exact 1000 runes cutoff, got %d", runeCount)
		}
	})
}

func TestHashText_Deduplication(t *testing.T) {
	text1 := "  Info penting: Vaksin menyebabkan demam tinggi pada balita.  "
	text2 := "INFO PENTING:  vaksin menyebabkan demam tinggi pada balita.\n"
	text3 := "Info penting: Vaksin menyebabkan demam tinggi pada orang dewasa."

	hash1 := HashText(text1)
	hash2 := HashText(text2)
	hash3 := HashText(text3)

	if hash1 != hash2 {
		t.Errorf("expected identical hashes for normalized text, got %q vs %q", hash1, hash2)
	}
	if hash1 == hash3 {
		t.Errorf("expected different hashes for different texts, got %q", hash1)
	}
	if len(hash1) != 64 {
		t.Errorf("expected SHA-256 hex length of 64, got %d", len(hash1))
	}
}

func TestZeroEvidence_ShortCircuitBypassesJev(t *testing.T) {
	evaluator := decision.NewEvaluator()
	article := &model.Article{
		Title:        "Klaim tanpa bukti yang dapat ditemukan di internet",
		Content:      "Narasi ini tidak memiliki grounding atau rujukan artikel media manapun.",
		URL:          "text:a1b2c3d4",
		WordCount:    10,
		SourceDomain: "direct_text",
	}

	evidences := []model.Evidence{} // 0 evidence returned by search provider

	// When 0 evidence is returned, Jev evaluation is bypassed and Evaluator directly handles it
	result := evaluator.Evaluate(article, evidences, &model.JevJudgments{})

	if result.Verdict != config.VerdictUnverified {
		t.Errorf("expected verdict UNVERIFIED when evidence is empty, got %s", result.Verdict)
	}
	if result.Confidence != 50.0 {
		t.Errorf("expected confidence 50.0 when unverified, got %f", result.Confidence)
	}
}
