package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"net/url"
	"strings"

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
	// 1. URL パラメータではなく、フォームデータとして定義
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("code", code)

	// 2. Body にデータを入れて POST
	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// 3. ヘッダーを適切に設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// --- ここから下はデバッグ用に一時的に追加 ---
	// レスポンスが 200 以外、または中身がエラーの場合に備えて
	var rawBody json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rawBody); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}
	
	// ログに出力して中身を確認
	fmt.Printf("DEBUG: GitHub Response: %s\n", string(rawBody))

	var tResp githubTokenResponse
	if err := json.Unmarshal(rawBody, &tResp); err != nil {
		return "", err
	}
	// --- デバッグ用ここまで ---

	if tResp.AccessToken == "" {
		// GitHubはエラー時も 200 OK で "error": "bad_verification_code" などを返してくることがあるため
		return "", fmt.Errorf("github oauth error: %s", string(rawBody))
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