package security

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MaxTextRunes   = 1000
	MinTextRunes   = 25
	MinTextWords   = 4
	MaxSearchRunes = 100
	MinSearchRunes = 10
)

var (
	ErrEmptyText    = errors.New("teks kosong")
	ErrTextTooShort = errors.New("teks terlalu pendek (minimal 25 karakter dan 4 kata)")
	ErrTextTooLong  = errors.New("teks melebihi batas maksimal")
)

var (
	reFormatting  = regexp.MustCompile(`[*_~` + "`" + `]`)
	reForward     = regexp.MustCompile(`(?i)(pesan\s+ini\s+)?(diteruskan(\s+berkali-kali|\s+seperti\s+yang\s+diterima)?|forwarded(\s+as\s+received|\s+many\s+times)?)\s*:?`)
	reAuthority   = regexp.MustCompile(`(?i)(kabar|info|peringatan|pesan)\s+dari\s+(dokter\s+)?[a-z0-9\-]+\s*:?`)
	reViralGroup  = regexp.MustCompile(`(?i)\b(viralkan|sebarkan|share|bagikan)\s+ke\s+(semua\s+)?(\d+\s+)?(grup|group|kontak|teman|keluarga|orang|wa)\b[!.]*`)
	reViralSingle = regexp.MustCompile(`(?i)\b(viralkan|sebarkan)\b[!.]*`)
	reFooters     = regexp.MustCompile(`(?i)(semoga\s+bermanfaat|mohon\s+maaf\s+sekadar\s+mengingatkan,?\s*jangan\s+berhenti\s+di\s+anda)[!.]*`)
	reUrls        = regexp.MustCompile(`(?i)(https?://\S+|www\.\S+|wa\.me/\S+|bit\.ly/\S+|t\.me/\S+)`)
	rePhones      = regexp.MustCompile(`(?i)(\+?62|08)[0-9\s\-]{8,15}\b`)
)

func ValidateRawText(text string) (string, error) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", ErrEmptyText
	}

	runes := []rune(trimmed)
	if len(runes) < MinTextRunes || len(strings.Fields(trimmed)) < MinTextWords {
		return "", ErrTextTooShort
	}

	if len(runes) > MaxTextRunes {
		runes = runes[:MaxTextRunes]
	}

	return string(runes), nil
}

func CleanBroadcastNoise(rawText string) string {
	text := reFormatting.ReplaceAllString(rawText, "")
	text = reForward.ReplaceAllString(text, " ")
	text = reAuthority.ReplaceAllString(text, " ")
	text = reViralGroup.ReplaceAllString(text, " ")
	text = reViralSingle.ReplaceAllString(text, " ")
	text = reFooters.ReplaceAllString(text, " ")
	text = reUrls.ReplaceAllString(text, " ")
	text = rePhones.ReplaceAllString(text, " ")
	text = stripEmojisAndSymbols(text)

	words := strings.Fields(text)
	cleaned := strings.Join(words, " ")
	cleaned = strings.Trim(cleaned, " :,-\t\r\n")

	return cleaned
}

func ExtractSearchQuery(text string, maxLen int) string {
	if maxLen <= 0 {
		maxLen = MaxSearchRunes
	}

	firstSentence := text
	for _, delim := range []string{".", "?", "!", "\n"} {
		if idx := strings.Index(firstSentence, delim); idx > 15 {
			firstSentence = firstSentence[:idx]
			break
		}
	}
	firstSentence = strings.TrimSpace(firstSentence)

	runes := []rune(firstSentence)
	if len(runes) < MinSearchRunes {
		runes = []rune(text)
	}

	if len(runes) >= maxLen {
		sub := string(runes[:maxLen-1])
		if lastSpace := strings.LastIndex(sub, " "); lastSpace > 20 {
			sub = sub[:lastSpace]
		}
		firstSentence = strings.TrimSpace(sub)
	} else {
		firstSentence = strings.TrimSpace(string(runes))
	}

	firstSentence = strings.Trim(firstSentence, ".,-;:?! ")
	return firstSentence
}

func HashText(text string) string {
	norm := strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(text)), " "))
	h := sha256.Sum256([]byte(norm))
	return hex.EncodeToString(h[:])
}

func CleanText(rawText string) (cleaned string, query string, err error) {
	validated, err := ValidateRawText(rawText)
	if err != nil {
		return "", "", err
	}

	cleaned = CleanBroadcastNoise(validated)
	if utf8.RuneCountInString(cleaned) < MinTextRunes || len(strings.Fields(cleaned)) < MinTextWords {
		return "", "", ErrTextTooShort
	}

	query = ExtractSearchQuery(cleaned, MaxSearchRunes)
	return cleaned, query, nil
}

func CleanBroadcastText(rawText string) (cleanedText, searchQuery string) {
	c, q, _ := CleanText(rawText)
	return c, q
}

func stripEmojisAndSymbols(input string) string {
	var b strings.Builder
	b.Grow(len(input))

	for _, r := range input {
		if r > 0x1F000 && r <= 0x1FAFF {
			continue
		}
		if (r >= 0x2600 && r <= 0x27BF) || (r >= 0xFE00 && r <= 0xFE0F) {
			continue
		}
		if unicode.Is(unicode.So, r) || unicode.Is(unicode.Sk, r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
