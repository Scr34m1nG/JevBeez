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
}
