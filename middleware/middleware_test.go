package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	// Permite 2 requisições com burst 2
	limiter := NewIPRateLimiter(rate.Limit(2), 2)
	handler := RateLimitMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 2 primeiras chamadas devem passar (200)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Esperava status 200, recebeu %d na requisicao %d", rr.Code, i+1)
		}
	}

	// 3ª chamada imediata deve ser bloqueada (429)
	req := httptest.NewRequest("GET", "/health", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Esperava status 429 Too Many Requests, recebeu %d", rr.Code)
	}
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	secret := []byte("secreto_teste")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "123",
		"role": "admin",
	})
	tokenStr, _ := token.SignedString(secret)

	handler := JWTMiddleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-Id")
		userRole := r.Header.Get("X-User-Role")

		if userID != "123" || userRole != "admin" {
			t.Errorf("Headers de contexto nao foram propagados corretamente. UserID: %s, Role: %s", userID, userRole)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/orders/list", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava status 200 para token valido, recebeu %d", rr.Code)
	}
}