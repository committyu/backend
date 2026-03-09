package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/internal/domain"
)

type GitHubClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewGitHubClient(id, secret string) *GitHubClient {
	return &GitHubClient{
		clientID:     id,
		clientSecret: secret,
		// タイムアウトを設定したHTTPクライアントを持つのが安全です
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// 1. アクセストークン取得用のレスポンス構造体
type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// 2. ユーザー情報取得用のレスポンス構造体
type gitHubUserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (c *GitHubClient) GetUser(ctx context.Context, code string) (*domain.User, error) {
	// --- STEP 1: Access Token の取得 ---
	token, err := c.getAccessToken(ctx, code)
	if err != nil {
		return nil, err
	}

	// --- STEP 2: ユーザー情報の取得 ---
	githubData, err := c.fetchGitHubUser(ctx, token)
	if err != nil {
		return nil, err
	}

	// --- STEP 3: Domainモデルへの変換 ---
	// 内部IDは新規発行（Usecase側で既存チェックする前提なら、ここでNewUserIDしてもOK）
	user := domain.NewUser(
		domain.NewUserID(),
		githubData.Login,
		githubData.Email,
		githubData.AvatarURL,
		githubData.ID,
		time.Now(),
	)

	return user, nil
}

// ヘルパー関数: アクセストークンを取得
func (c *GitHubClient) getAccessToken(ctx context.Context, code string) (string, error) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		c.clientID, c.clientSecret, code,
	)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	// GitHubにJSON形式で返してほしいと伝える
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
		return "", fmt.Errorf("github oauth error: failed to get access token")
	}

	return tResp.AccessToken, nil
}

// ヘルパー関数: ユーザー情報を取得
func (c *GitHubClient) fetchGitHubUser(ctx context.Context, token string) (*gitHubUserResponse, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

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