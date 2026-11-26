package decision

import (
	"testing"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

func TestEvaluator(t *testing.T) {
	evaluator := NewEvaluator()
	article := &model.Article{
		Title:        "Beredar Kabar Megathrust Melanda Besok",
		URL:          "https://example.com/news/1",
		SourceDomain: "example.com",
	}

	evidences := []model.Evidence{
		{
			Source:  "kompas.com",
			Title:   "BMKG: Narasi Megathrust Melanda Besok Adalah Hoaks",
			Snippet: "BMKG membantah kabar yang viral...",
			URL:     "https://kompas.com/klarifikasi-1",
		},
	}

	t.Run("verdict hoax when contradicted is high", func(t *testing.T) {
		judgments := &model.JevJudgments{
			ContradictedByEvidence: 0.92,
			SupportedByEvidence:    0.05,
			SensationalFraming:     0.80,
		}
		result := evaluator.Evaluate(article, evidences, judgments)
		if result.Verdict != config.VerdictHoax {
			t.Errorf("expected %s, got %s", config.VerdictHoax, result.Verdict)
		}
		if result.Confidence < 90.0 {
			t.Errorf("expected high confidence, got %f", result.Confidence)
		}
	})

	t.Run("verdict valid when supported is high and sensational is low", func(t *testing.T) {
		judgments := &model.JevJudgments{
			ContradictedByEvidence: 0.05,
			SupportedByEvidence:    0.88,
			SensationalFraming:     0.15,
		}
		result := evaluator.Evaluate(article, evidences, judgments)
		if result.Verdict != config.VerdictValid {
			t.Errorf("expected %s, got %s", config.VerdictValid, result.Verdict)
		}
	})

	t.Run("verdict misleading when sensational is high and contradicted is low", func(t *testing.T) {
		judgments := &model.JevJudgments{
			ContradictedByEvidence: 0.20,
			SupportedByEvidence:    0.40,
			SensationalFraming:     0.85,
		}
		result := evaluator.Evaluate(article, evidences, judgments)
		if result.Verdict != config.VerdictMisleading {
			t.Errorf("expected %s, got %s", config.VerdictMisleading, result.Verdict)
		}
	})

	t.Run("verdict unverified when evidences are empty", func(t *testing.T) {
		judgments := &model.JevJudgments{
			ContradictedByEvidence: 0.50,
			SupportedByEvidence:    0.50,
			SensationalFraming:     0.50,
		}
		result := evaluator.Evaluate(article, []model.Evidence{}, judgments)
		if result.Verdict != config.VerdictUnverified {
			t.Errorf("expected %s, got %s", config.VerdictUnverified, result.Verdict)
		}
	})

	t.Run("evaluate threat zero-contact fast verdict", func(t *testing.T) {
		threatInput := &model.ThreatEvaluationInput{
			RawURL:      "https://bca.co.id.login.xyz/auth",
			ResolvedURL: "https://bca.co.id.login.xyz/auth",
			ThreatType:  config.VerdictPhishing,
			Confidence:  92.0,
			Reasons:     []string{"Domain tidak resmi"},
			TargetBrand: "bca",
			FormThreats: nil,
			IsDeepDive:  false,
		}

		result := evaluator.EvaluateThreat(threatInput)
		if !result.IsThreat {
			t.Errorf("expected IsThreat to be true")
		}
		if result.Verdict != config.VerdictPhishing {
			t.Errorf("expected %s, got %s", config.VerdictPhishing, result.Verdict)
		}
		if result.TargetBrand != "bca" {
			t.Errorf("expected target brand bca, got %s", result.TargetBrand)
		}
		if result.IsDeepDive {
			t.Errorf("expected IsDeepDive to be false")
		}
	})

	t.Run("evaluate threat deep-dive with form threats", func(t *testing.T) {
		threatInput := &model.ThreatEvaluationInput{
			RawURL:      "https://bca.co.id.login.xyz/auth",
			ResolvedURL: "https://bca.co.id.login.xyz/auth",
			ThreatType:  config.VerdictPhishing,
			Confidence:  95.0,
			Reasons:     []string{"Domain tidak resmi"},
			TargetBrand: "bca",
			FormThreats: []string{"Meminta PIN", "Exfiltrasi ke Telegram"},
			IsDeepDive:  true,
		}

		result := evaluator.EvaluateThreat(threatInput)
		if !result.IsThreat {
			t.Errorf("expected IsThreat to be true")
		}
		if len(result.FormThreats) != 2 {
			t.Errorf("expected 2 form threats, got %d", len(result.FormThreats))
		}
		if !result.IsDeepDive {
			t.Errorf("expected IsDeepDive to be true")
		}
	})
}

func TestEvaluateAds(t *testing.T) {
	e := NewEvaluator()
	res := e.EvaluateAds(
		"https://www.google.com/pagead/aclk?adurl=...",
		"https://www.ebay.com/itm/123",
		"ebay.com",
		[]string{"https://www.ebay.com/itm/123"},
	)

	if res.Verdict != config.VerdictAds {
		t.Fatalf("expected ADS, got %s", res.Verdict)
	}
	if res.IsThreat {
		t.Fatal("ADS must not be a threat")
	}
	if res.ArticleURL != "https://www.google.com/pagead/aclk?adurl=..." {
		t.Errorf("expected original url preserved, got %s", res.ArticleURL)
	}
	if res.ResolvedURL != "https://www.ebay.com/itm/123" {
		t.Errorf("expected resolved url, got %s", res.ResolvedURL)
	}
	if res.Confidence <= 0 {
		t.Error("expected non-zero confidence")
	}
	if len(res.Evidences) != 0 {
		t.Error("expected no evidences")
	}
}
