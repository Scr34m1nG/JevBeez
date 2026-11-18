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
