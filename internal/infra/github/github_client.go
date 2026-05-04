package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain"
	"backend/internal/pkg/logger"
)

const compareRequestInterval = 100 * time.Millisecond

type GitHubClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

type gitHubEventResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Before string `json:"before"`
		Head   string `json:"head"`
	} `json:"payload"`
}

type gitHubCompareResponse struct {
	TotalCommits int `json:"total_commits"`
}

var _ domain.GitHubService = (*GitHubClient)(nil)

type Config struct {
	ClientID     string
	ClientSecret string
}

func NewGitHubClient(cfg Config) *GitHubClient {
	return &GitHubClient{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type gitHubUserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// OAuth code → GitHub user
func (c *GitHubClient) GetUser(ctx context.Context, code string) (*domain.User, error) {

	token, err := c.getAccessToken(ctx, code)
	if err != nil {
		return nil, err
	}

	githubUser, err := c.fetchGitHubUser(ctx, token)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(
		domain.NewUserID(),
		githubUser.Login,
		githubUser.Email,
		githubUser.AvatarURL,
		githubUser.ID,
		time.Now(),
	)

	return user, nil
}

// OAuth code → access_token
func (c *GitHubClient) getAccessToken(ctx context.Context, code string) (string, error) {

	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("code", code)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		logger.Error("failed to create github token request", "error", err)
		return "", fmt.Errorf("create github token request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("github oauth request failed", "error", err)
		return "", fmt.Errorf("github oauth request failed: %w", err)
	}
	defer resp.Body.Close()

	var tResp githubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tResp); err != nil {
		logger.Error("failed to decode github token response", "error", err)
		return "", fmt.Errorf("decode github token response failed: %w", err)
	}

	if tResp.AccessToken == "" {
		logger.Error("github oauth returned empty token")
		return "", fmt.Errorf("github oauth error")
	}

	return tResp.AccessToken, nil
}

// access_token → GitHub user
func (c *GitHubClient) fetchGitHubUser(ctx context.Context, token string) (*gitHubUserResponse, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		logger.Error("failed to create github user request", "error", err)
		return nil, fmt.Errorf("create github user request failed: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("github user api request failed", "error", err)
		return nil, fmt.Errorf("github user api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("github api returned non-200 status",
			"status", resp.StatusCode,
		)
		return nil, fmt.Errorf("github api error: status %d", resp.StatusCode)
	}

	var uResp gitHubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&uResp); err != nil {
		logger.Error("failed to decode github user response", "error", err)
		return nil, fmt.Errorf("decode github user response failed: %w", err)
	}

	if uResp.Login == "" {
		logger.Error("github user response missing login")
		return nil, fmt.Errorf("github user response missing login")
	}

	return &uResp, nil
}

func (c *GitHubClient) GetPushEvents(ctx context.Context, username string, lastCommitCheckedAt time.Time) ([]domain.GitHubPushEvent, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.github.com/users/%s/events?per_page=100", username),
		nil,
	)
	if err != nil {
		logger.Error("failed to create github events request", "error", err)
		return nil, fmt.Errorf("create github events request failed: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("github events api request failed", "error", err)
		return nil, fmt.Errorf("github events api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := githubAPIError(resp, "events")
		logger.Error("github events api returned error", "error", err)
		return nil, err
	}

	var eventResp []gitHubEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResp); err != nil {
		logger.Error("failed to decode github events response", "error", err)
		return nil, fmt.Errorf("decode github events response failed: %w", err)
	}

	pushEvents := make([]domain.GitHubPushEvent, 0)
	compareRequestCount := 0

	for _, event := range eventResp {
		if event.Type != "PushEvent" {
			continue
		}

		if !event.CreatedAt.After(lastCommitCheckedAt) {
			continue
		}

		if event.Payload.Before == "" || event.Payload.Head == "" {
			logger.Error(
				"push event compare refs missing",
				"event_id", event.ID,
				"repo", event.Repo.Name,
			)
			continue
		}

		owner, repo, err := splitRepositoryName(event.Repo.Name)
		if err != nil {
			logger.Error(
				"invalid push event repo name",
				"event_id", event.ID,
				"repo", event.Repo.Name,
				"error", err,
			)
			continue
		}

		if compareRequestCount > 0 {
			if err := waitForCompareRequest(ctx); err != nil {
				return nil, err
			}
		}

		commitCount, err := c.getCommitCount(ctx, owner, repo, event.Payload.Before, event.Payload.Head)
		compareRequestCount++
		if err != nil {
			logger.Error(
				"get commit count failed",
				"event_id", event.ID,
				"repo", event.Repo.Name,
				"before", event.Payload.Before,
				"head", event.Payload.Head,
				"error", err,
			)
			continue
		}

		pushEvents = append(pushEvents, domain.GitHubPushEvent{
			ID:          event.ID,
			CreatedAt:   event.CreatedAt,
			CommitCount: commitCount,
		})
	}

	return pushEvents, nil
}

func (c *GitHubClient) getCommitCount(ctx context.Context, owner, repo, before, head string) (int, error) {
	compareURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/compare/%s...%s",
		url.PathEscape(owner),
		url.PathEscape(repo),
		url.PathEscape(before),
		url.PathEscape(head),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, compareURL, nil)
	if err != nil {
		logger.Error("failed to create github compare request", "error", err)
		return 0, fmt.Errorf("create github compare request failed: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("github compare api request failed", "error", err)
		return 0, fmt.Errorf("github compare api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := githubAPIError(resp, "compare")
		logger.Error("github compare api returned error", "error", err)
		return 0, err
	}

	var compareResp gitHubCompareResponse
	if err := json.NewDecoder(resp.Body).Decode(&compareResp); err != nil {
		logger.Error("failed to decode github compare response", "error", err)
		return 0, fmt.Errorf("decode github compare response failed: %w", err)
	}

	return compareResp.TotalCommits, nil
}

func splitRepositoryName(name string) (string, string, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected owner/repo, got %q", name)
	}

	return parts[0], parts[1], nil
}

func waitForCompareRequest(ctx context.Context) error {
	timer := time.NewTimer(compareRequestInterval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func githubAPIError(resp *http.Response, apiName string) error {
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		return fmt.Errorf("github %s api rate limited: status=%d retry_after=%ss", apiName, resp.StatusCode, retryAfter)
	}

	if resp.Header.Get("X-RateLimit-Remaining") == "0" {
		resetAt := "unknown"
		if resetUnix, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64); err == nil {
			resetAt = time.Unix(resetUnix, 0).Format(time.RFC3339)
		}

		return fmt.Errorf("github %s api rate limit exceeded: status=%d reset_at=%s", apiName, resp.StatusCode, resetAt)
	}

	return fmt.Errorf("github %s api error: status %d", apiName, resp.StatusCode)
}
