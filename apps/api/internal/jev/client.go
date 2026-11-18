package jev

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

type Client struct {
	apiKey string
	client *http.Client
}

type systemOneRequest struct {
	Model     string                     `json:"model"`
	State     JevState                   `json:"state"`
	Questions map[string]QuestionPayload `json:"questions"`
}

type noulAnswer struct {
	Type string  `json:"type"`
	Noul float64 `json:"noul"`
}

type systemOneResponse struct {
	Model   string                `json:"model"`
	Answers map[string]noulAnswer `json:"answers"`
	Error   string                `json:"error"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{Timeout: config.HTTPClientTimeout},
	}
}

func (c *Client) Evaluate(ctx context.Context, article *model.Article, evidences []model.Evidence) (*model.JevJudgments, error) {
	if c.apiKey == "" {
		return nil, errors.New("typesafe api key not configured")
	}

	payload := systemOneRequest{
		Model:     config.JevDefaultModel,
		State:     BuildState(article, evidences),
		Questions: BuildQuestions(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.JevEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("typesafe jev api returned status %d", resp.StatusCode)
	}

	var res systemOneResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	judgments := &model.JevJudgments{
		ContradictedByEvidence: res.Answers["contradicted_by_evidence"].Noul,
		SupportedByEvidence:    res.Answers["supported_by_evidence"].Noul,
		SensationalFraming:     res.Answers["sensational_framing"].Noul,
	}

	return judgments, nil
}
