package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// トークンストア（メモリ内）
var tokenStore = struct {
	sync.RWMutex
	tokens map[string]time.Time
}{tokens: make(map[string]time.Time)}

const tokenExpiry = 24 * time.Hour

// AdminAuth はパスワードを検証してトークンを発行する
func AdminAuth(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	// 環境変数からハッシュ値を取得
	expectedHash := os.Getenv("ADMIN_PASSWORD_HASH")
	if expectedHash == "" {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Admin password not configured",
		})
		return
	}

	// 入力パスワードをSHA-256でハッシュ化
	hash := sha256.Sum256([]byte(req.Password))
	inputHash := hex.EncodeToString(hash[:])

	if inputHash != expectedHash {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   "Invalid password",
		})
		return
	}

	// トークンを生成
	token := generateToken()
	
	// トークンを保存
	tokenStore.Lock()
	tokenStore.tokens[token] = time.Now().Add(tokenExpiry)
	tokenStore.Unlock()

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Token:   token,
	})
}

// AdminVerifyToken はトークンの有効性を検証する
func AdminVerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	// "Bearer " プレフィックスを除去
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	tokenStore.RLock()
	expiry, exists := tokenStore.tokens[token]
	tokenStore.RUnlock()

	if !exists || time.Now().After(expiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// AdminAuthMiddleware は管理APIの認証ミドルウェア
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		// "Bearer " プレフィックスを除去
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		tokenStore.RLock()
		expiry, exists := tokenStore.tokens[token]
		tokenStore.RUnlock()

		if !exists || time.Now().After(expiry) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CleanupExpiredTokens は期限切れトークンを削除（定期実行用）
func CleanupExpiredTokens() {
	tokenStore.Lock()
	defer tokenStore.Unlock()
	
	now := time.Now()
	for token, expiry := range tokenStore.tokens {
		if now.After(expiry) {
			delete(tokenStore.tokens, token)
		}
	}
}
