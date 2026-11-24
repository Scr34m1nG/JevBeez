package jev

import (
	"strings"
	"jevbeez/api/internal/model"
)

type JevState struct {
	Article  ArticleContext    `json:"article"`
	Evidence []EvidenceContext `json:"evidence"`
}

type ArticleContext struct {
	Title   string `json:"title"`
	Source  string `json:"source_domain"`
	Excerpt string `json:"excerpt"`
}

type EvidenceContext struct {
	Source  string `json:"source"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

type QuestionPayload struct {
	Type         string `json:"type"`
	Instructions string `json:"instructions"`
}

func BuildState(article *model.Article, evidences []model.Evidence) JevState {
	excerpt := article.Content
	words := strings.Fields(excerpt)
	if len(words) > 300 {
		excerpt = strings.Join(words[:300], " ") + "..."
	}

	var evList []EvidenceContext
	for _, ev := range evidences {
		evList = append(evList, EvidenceContext{
			Source:  ev.Source,
			Title:   ev.Title,
			Snippet: ev.Snippet,
		})
	}

	return JevState{
		Article: ArticleContext{
			Title:   article.Title,
			Source:  article.SourceDomain,
			Excerpt: excerpt,
		},
		Evidence: evList,
	}
}

func BuildQuestions() map[string]QuestionPayload {
	return map[string]QuestionPayload{
		"contradicted_by_evidence": {
			Type:         "noul",
			Instructions: "Does the evidence directly contradict the main claim or event reported in the article?",
		},
		"supported_by_evidence": {
			Type:         "noul",
			Instructions: "Does the evidence confirm that the main event or fact in the article actually happened as described?",
		},
		"sensational_framing": {
			Type:         "noul",
			Instructions: "Does the article use sensational, alarmist, or clickbait framing that exaggerates or misrepresents the evidence?",
		},
	}
}

func BuildThreatQuestions() map[string]QuestionPayload {
	return map[string]QuestionPayload{
		"credential_harvesting": {
			Type:         "noul",
			Instructions: "Does this webpage attempt to harvest credentials, PIN, passwords, OTP, or sensitive financial information from the user?",
		},
		"social_engineering": {
			Type:         "noul",
			Instructions: "Does this webpage use deceptive urgency, fake prizes, account suspension threats, or impersonate an institution to trick the visitor?",
		},
	}
}
