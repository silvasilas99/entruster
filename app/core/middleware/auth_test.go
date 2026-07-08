package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silvasilas99/entruster/config"
)

func TestJWTAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock server to simulate the Mock User API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+config.StaticToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := UserInfo{
			ID:    "usr_test",
			Name:  "Test User",
			Email: "test@hospital.com",
			Role:  "Tester",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}))
	defer mockServer.Close()

	// Override the mock API URL for tests
	originalURL := config.UserMockApiUrl
	config.UserMockApiUrl = mockServer.URL
	defer func() { config.UserMockApiUrl = originalURL }()

	router := gin.New()
	router.Use(JWTAuth())
	router.GET("/test", func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, user)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token123")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Valid Static Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+config.StaticToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var user UserInfo
		if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}

		if user.ID != "usr_test" || user.Name != "Test User" {
			t.Errorf("Expected user ID 'usr_test' and name 'Test User', got %s, %s", user.ID, user.Name)
		}
	})
}
