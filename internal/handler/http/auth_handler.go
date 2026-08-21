package http

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"backend/internal/infra/auth"
	"backend/internal/pkg/logger"
	usecaseAuth "backend/internal/usecase/auth"

	"github.com/labstack/echo/v4"
)

const (
	oauthStateCookieName  = "github_oauth_state"
	accessTokenCookieName = "access_token"
)

type AuthHandler struct {
	loginUc             *usecaseAuth.LoginUsecase
	tokenUc             *usecaseAuth.GenerateTokenUsecase
	githubClientID      string
	githubRedirectURL   string
	frontendCallbackURL string
}

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

func NewAuthHandler(loginUc *usecaseAuth.LoginUsecase, tokenUc *usecaseAuth.GenerateTokenUsecase, githubClientID, githubRedirectURL, frontendCallbackURL string) *AuthHandler {
	return &AuthHandler{loginUc: loginUc, tokenUc: tokenUc, githubClientID: githubClientID, githubRedirectURL: githubRedirectURL, frontendCallbackURL: frontendCallbackURL}
}

// Login godoc
// @Summary GitHub OAuthログインを開始
// @Description GitHubの認可画面へリダイレクトします。
// @Tags auth
// @Success 307
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [get]
func (h *AuthHandler) Login(c echo.Context) error {
	if h.githubClientID == "" || h.githubRedirectURL == "" {
		logger.Error("github oauth configuration missing")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "github oauth is not configured"})
	}

	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		logger.Error("failed to generate oauth state", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)
	c.SetCookie(&http.Cookie{
		Name: oauthStateCookieName, Value: state, Path: "/api/auth", HttpOnly: true,
		Secure: requestIsHTTPS(c), SameSite: http.SameSiteLaxMode,
		MaxAge: 600, Expires: time.Now().Add(10 * time.Minute),
	})

	authorizeURL := &url.URL{Scheme: "https", Host: "github.com", Path: "/login/oauth/authorize"}
	query := authorizeURL.Query()
	query.Set("client_id", h.githubClientID)
	query.Set("redirect_uri", h.githubRedirectURL)
	query.Set("scope", "user:email")
	query.Set("state", state)
	authorizeURL.RawQuery = query.Encode()
	return c.Redirect(http.StatusTemporaryRedirect, authorizeURL.String())
}

// Callback godoc
// @Summary GitHub OAuthコールバック
// @Description GitHub認証後にJWTをHttpOnly Cookieへ保存し、フロントへリダイレクトします。
// @Tags auth
// @Param code query string true "GitHub OAuth code"
// @Param state query string true "OAuth state"
// @Success 307
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/callback [get]
func (h *AuthHandler) Callback(c echo.Context) error {
	if oauthError := c.QueryParam("error"); oauthError != "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "github authorization failed: " + oauthError})
	}

	stateCookie, err := c.Cookie(oauthStateCookieName)
	if err != nil || stateCookie.Value == "" || c.QueryParam("state") != stateCookie.Value {
		logger.Error("invalid github oauth state")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid oauth state"})
	}
	c.SetCookie(&http.Cookie{Name: oauthStateCookieName, Value: "", Path: "/api/auth", MaxAge: -1, HttpOnly: true})

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "code is required"})
	}
	user, err := h.loginUc.Execute(c.Request().Context(), code)
	if err != nil {
		logger.Error("login usecase failed", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
	token, err := h.tokenUc.Execute(user.ID())
	if err != nil {
		logger.Error("token generation failed", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
	}
	frontendURL, err := url.Parse(h.frontendCallbackURL)
	if err != nil || frontendURL.Scheme == "" || frontendURL.Host == "" {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "frontend callback is not configured"})
	}

	c.SetCookie(&http.Cookie{
		Name: accessTokenCookieName, Value: token, Path: "/", HttpOnly: true,
		Secure: requestIsHTTPS(c), SameSite: http.SameSiteLaxMode,
		MaxAge: int(auth.TokenExpireDuration.Seconds()), Expires: time.Now().Add(auth.TokenExpireDuration),
	})
	logger.Info("login success", "user_id", user.ID())
	return c.Redirect(http.StatusTemporaryRedirect, frontendURL.String())
}

func requestIsHTTPS(c echo.Context) bool {
	return c.Request().TLS != nil || c.Request().Header.Get(echo.HeaderXForwardedProto) == "https"
}
