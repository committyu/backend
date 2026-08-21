package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend/internal/domain"
	infraAuth "backend/internal/infra/auth"

	"github.com/labstack/echo/v4"
)

func TestAuthMiddlewareReadsJWTFromCookie(t *testing.T) {
	const secret = "test-secret"
	previous := os.Getenv("JWT_SECRET")
	t.Setenv("JWT_SECRET", secret)
	t.Cleanup(func() { _ = os.Setenv("JWT_SECRET", previous) })

	userID := domain.NewUserID()
	token, err := infraAuth.NewJWTService().Generate(userID.String())
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	e.GET("/me", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Get("userID").(string))
	}, AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != userID.String() {
		t.Fatalf("status=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddlewareDoesNotUseBearerHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	e := echo.New()
	e.GET("/me", func(c echo.Context) error { return c.NoContent(http.StatusOK) }, AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer ignored")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
