package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"backend/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func TestMain(m *testing.M) {
	logger.Init()
	os.Exit(m.Run())
}

func TestLoginRedirectsToGitHubWithState(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	h := NewAuthHandler(nil, nil, "client-id", "http://localhost:8080/api/auth/callback", "http://localhost:3000/auth/callback")

	if err := h.Login(e.NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	location, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusTemporaryRedirect || location.Host != "github.com" || location.Query().Get("state") == "" {
		t.Fatalf("unexpected OAuth redirect: status=%d location=%s", rec.Code, location)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Value != location.Query().Get("state") || !cookies[0].HttpOnly {
		t.Fatalf("unexpected state cookie: %#v", cookies)
	}
}

func TestCallbackRejectsInvalidState(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback?code=code&state=wrong", nil)
	req.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "expected"})
	rec := httptest.NewRecorder()
	h := NewAuthHandler(nil, nil, "client-id", "http://localhost:8080/api/auth/callback", "http://localhost:3000/auth/callback")

	if err := h.Callback(e.NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
