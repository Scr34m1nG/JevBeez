package model

type Evidence struct {
	Source  string `json:"source"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

type Article struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	URL          string `json:"url"`
	WordCount    int    `json:"word_count"`
	SourceDomain string `json:"source_domain"`
	IsSocial     bool   `json:"is_social"`
	Platform     string `json:"platform"`
}

type JevJudgments struct {
	ContradictedByEvidence float64 `json:"contradicted_by_evidence"`
	SupportedByEvidence    float64 `json:"supported_by_evidence"`
	SensationalFraming     float64 `json:"sensational_framing"`
}

type VerificationResult struct {
	Verdict           string       `json:"verdict"`
	Confidence        float64      `json:"confidence"`
	Judgments         JevJudgments `json:"judgments"`
	ArticleTitle      string       `json:"article_title"`
	ArticleURL        string       `json:"article_url"`
	SourceDomain      string       `json:"source_domain"`
	Summary           string       `json:"summary"`
	RedFlags          []string     `json:"red_flags"`
	Evidences         []Evidence   `json:"evidences"`
	IsCached          bool         `json:"is_cached"`
	IsSocial          bool         `json:"is_social"`
	Platform          string       `json:"platform"`
}

type AnalyzeRequest struct {
	URL string `json:"url"`
}
