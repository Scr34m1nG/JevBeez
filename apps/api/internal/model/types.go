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
	CredentialHarvesting   float64 `json:"credential_harvesting,omitempty"`
	SocialEngineering      float64 `json:"social_engineering,omitempty"`
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
	IsThreat          bool         `json:"is_threat"`
	ThreatType        string       `json:"threat_type,omitempty"`
	ThreatReasons     []string     `json:"threat_reasons,omitempty"`
	TargetBrand       string       `json:"target_brand,omitempty"`
	ResolvedURL       string       `json:"resolved_url,omitempty"`
	RedirectHops      []string     `json:"redirect_hops,omitempty"`
	FormThreats       []string     `json:"form_threats,omitempty"`
	IsDeepDive        bool         `json:"is_deep_dive,omitempty"`
}

type AnalyzeRequest struct {
	URL      string `json:"url"`
	DeepDive bool   `json:"deep_dive,omitempty"`
}

type ThreatEvaluationInput struct {
	RawURL       string
	ResolvedURL  string
	RedirectHops []string
	ThreatType   string
	Confidence   float64
	Reasons      []string
	TargetBrand  string
	FormThreats  []string
	PageTitle    string
	IsDeepDive   bool
}

type BlockedIPInfo struct {
	IP        string `json:"ip"`
	Reason    string `json:"reason"`
	Payload   string `json:"payload,omitempty"`
	BlockedAt string `json:"blocked_at"`
	ExpiresAt string `json:"expires_at"`
}
