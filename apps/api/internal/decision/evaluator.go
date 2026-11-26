package decision

import (
	"fmt"
	"math"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(article *model.Article, evidences []model.Evidence, judgments *model.JevJudgments) *model.VerificationResult {
	result := &model.VerificationResult{
		ArticleTitle: article.Title,
		ArticleURL:   article.URL,
		SourceDomain: article.SourceDomain,
		Judgments:    *judgments,
		Evidences:    evidences,
		RedFlags:     []string{},
		IsCached:     false,
		IsSocial:     article.IsSocial,
		Platform:     article.Platform,
	}

	if len(evidences) == 0 {
		result.Verdict = config.VerdictUnverified
		result.Confidence = 50.0
		result.Summary = "Tidak ditemukan bukti rujukan atau klarifikasi yang cukup di media kredibel mengenai artikel ini."
		return result
	}

	c := judgments.ContradictedByEvidence
	s := judgments.SupportedByEvidence
	sf := judgments.SensationalFraming

	switch {
	case c >= config.ThresholdContradictedHigh:
		result.Verdict = config.VerdictHoax
		result.Confidence = round(c * 100)
		result.Summary = "Klaim dalam artikel ini bertentangan dengan fakta atau telah dibantah oleh sumber rujukan kredibel."
		result.RedFlags = append(result.RedFlags, "Klaim inti bertentangan langsung dengan bukti yang ditemukan.")
		if sf >= config.ThresholdSensationalLow {
			result.RedFlags = append(result.RedFlags, "Gaya penulisan terindikasi provokatif atau hiperbolik.")
		}

	case s >= config.ThresholdSupportedHigh && sf <= config.ThresholdSensationalLow:
		result.Verdict = config.VerdictValid
		result.Confidence = round(s * 100)
		result.Summary = "Informasi dalam artikel ini terverifikasi dan didukung oleh laporan media arus utama."

	case sf >= config.ThresholdSensationalHigh && c <= config.ThresholdContradictedLow:
		result.Verdict = config.VerdictMisleading
		result.Confidence = round(sf * 100)
		result.Summary = "Artikel menggunakan judul atau framing yang melebih-lebihkan fakta asli (klikbait misleading)."
		result.RedFlags = append(result.RedFlags, "Judul atau narasi memelintir fakta dari sumber aslinya.")

	default:
		result.Verdict = config.VerdictUnverified
		result.Confidence = round(math.Max(c, math.Max(s, sf)) * 100)
		result.Summary = fmt.Sprintf("Bukti yang tersedia belum cukup konklusif untuk menyatakan artikel ini 100%% valid atau hoax (Probabilitas bantahan: %.0f%%, dukungan: %.0f%%).", c*100, s*100)
	}

	return result
}

func round(val float64) float64 {
	return math.Round(val*10) / 10
}

func (e *Evaluator) EvaluateAds(rawURL, destURL, destDomain string, chain []string) *model.VerificationResult {
	return &model.VerificationResult{
		Verdict:      config.VerdictAds,
		Confidence:   95.0,
		ArticleTitle: "Iklan Berbayar: " + destDomain,
		ArticleURL:   rawURL,
		ResolvedURL:  destURL,
		RedirectHops: chain,
		SourceDomain: destDomain,
		Summary:      "Tautan ini adalah iklan berbayar yang mengarah ke marketplace " + destDomain + ". Karena bukan klaim berita, JevBeez tidak memverifikasi faktanya. Pastikan transaksi hanya dilakukan di dalam platform resmi.",
		RedFlags:     []string{"Tautan iklan berbayar, bukan artikel berita."},
		Evidences:    []model.Evidence{},
		IsThreat:     false,
	}
}

func (e *Evaluator) EvaluateThreat(in *model.ThreatEvaluationInput) *model.VerificationResult {
	title := in.PageTitle
	if title == "" {
		switch in.ThreatType {
		case config.VerdictMalware:
			title = "Peringatan Bahaya: Unduhan Berbahaya (Malware/APK)"
		case config.VerdictPhishing:
			title = "Peringatan Bahaya: Indikasi Tautan Phishing"
		default:
			title = "Peringatan Waspada: Indikasi Situs Penipuan (Scam)"
		}
	}

	summary := "Tautan ini terdeteksi memiliki karakteristik kejahatan siber yang dapat membahayakan keamanan data Anda."
	switch in.ThreatType {
	case config.VerdictMalware:
		summary = "Tautan ini mengarah langsung ke pengunduhan file eksekusi berbahaya (seperti aplikasi .APK atau .EXE). Modus ini umum digunakan untuk membobol rekening atau membaca SMS/OTP korban."
	case config.VerdictPhishing:
		summary = "Tautan ini terdeteksi meniru entitas atau instansi resmi untuk mengelabui korban agar menyerahkan akun, kata sandi, atau data perbankan."
	case config.VerdictScam:
		summary = "Tautan ini terindikasi kuat sebagai situs penipuan dengan rekayasa sosial, iming-iming hadiah palsu, atau pemerasan digital."
	}

	redFlags := append([]string{}, in.Reasons...)
	redFlags = append(redFlags, in.FormThreats...)

	return &model.VerificationResult{
		Verdict:       in.ThreatType,
		Confidence:    in.Confidence,
		ArticleTitle:  title,
		ArticleURL:    in.RawURL,
		ResolvedURL:   in.ResolvedURL,
		RedirectHops:  in.RedirectHops,
		Summary:       summary,
		RedFlags:      redFlags,
		ThreatType:    in.ThreatType,
		ThreatReasons: in.Reasons,
		TargetBrand:   in.TargetBrand,
		FormThreats:   in.FormThreats,
		IsThreat:      true,
		IsDeepDive:    in.IsDeepDive,
		Evidences:     []model.Evidence{},
	}
}
