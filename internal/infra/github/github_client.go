package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend/internal/domain"
)

type GitHubClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
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
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tResp githubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tResp); err != nil {
		return "", err
	}

	if tResp.AccessToken == "" {
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
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: status %d", resp.StatusCode)
	}

	var uResp gitHubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&uResp); err != nil {
		return nil, err
	}

	return &uResp, nil
}