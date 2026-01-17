package handler

import (
	"net/http"
	"time"

	"github.com/DoDuy2004/slack-clone/backend/internal/models/dto"
	"github.com/DoDuy2004/slack-clone/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.authService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Set cookies
	h.setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear cookies
	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// TODO: Implement refresh logic with token rotation
	c.JSON(http.StatusOK, gin.H{"message": "Refresh endpoint - TODO"})
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	// Set access token cookie
	c.SetCookie(
		"access_token",
		accessToken,
		int(15*time.Minute.Seconds()), // 15 mins
		"/",
		"",
		false, // Secure: set to true in production
		true,  // HttpOnly
	)

	// Set refresh token cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(7*24*time.Hour.Seconds()), // 7 days
		"/",
		"",
		false, // Secure: set to true in production
		true,  // HttpOnly
	)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}
