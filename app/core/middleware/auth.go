package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/silvasilas99/entruster/config"
)

// UserInfo represents the user details returned by the mock API and saved in the context.
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// JWTAuth matches the Authorization header token, validates it against JWTSecret or StaticToken,
// calls the mock user API, and registers the user in the Gin context.
func JWTAuth() gin.HandlerFunc {
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			ctx.Abort()
			return
		}

		tokenStr := parts[1]
		var claims jwt.MapClaims
		isValid := false

		// 1. Try to parse and validate token as JWT using the secret key
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})

		if err == nil && token.Valid {
			if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
				claims = mapClaims
				isValid = true
			}
		}

		// 2. Fallback: If it's exactly the configured static token, consider it valid
		if !isValid && tokenStr == config.StaticToken {
			isValid = true
			// Try to extract claims without verifying signature (since we already matched the static token)
			parser := jwt.NewParser()
			tokenUnverified, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
			if err == nil {
				if mapClaims, ok := tokenUnverified.Claims.(jwt.MapClaims); ok {
					claims = mapClaims
				}
			}
		}

		if !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authorization token"})
			ctx.Abort()
			return
		}

		// 3. Make HTTP request to the mock/fictitious user API
		var user *UserInfo
		user, err = fetchUserInfoFromMockAPI(httpClient, config.UserMockApiUrl, tokenStr)
		if err != nil {
			log.Printf("[JWTAuth] Mock user API error: %v. Using context-based fallback user.", err)
			
			// Extract claims from token as fallback values
			sub, _ := claims["sub"].(string)
			name, _ := claims["name"].(string)
			role, _ := claims["role"].(string)

			if sub == "" {
				sub = "usr_unknown"
			}
			if name == "" {
				name = "Fallback User"
			}
			if role == "" {
				role = "Guest"
			}

			user = &UserInfo{
				ID:    sub,
				Name:  name,
				Email: strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@hospital.com",
				Role:  role,
			}
		}

		// 4. Save user information in Gin context
		ctx.Set("currentUser", user)

		log.Printf("[JWTAuth] User %s (%s) authenticated successfully for path: %s", user.Name, user.ID, ctx.Request.URL.Path)

		ctx.Next()
	}
}

// fetchUserInfoFromMockAPI performs the actual HTTP call to retrieve mock user data.
func fetchUserInfoFromMockAPI(client *http.Client, url, token string) (*UserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mock API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
