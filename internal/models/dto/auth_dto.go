package dto

import (
	"time"

	"github.com/DoDuy2004/slack-clone/backend/internal/models"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User models.User `json:"user"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  *string   `json:"full_name"`
	AvatarURL *string   `json:"avatar_url"`
	Status    string    `json:"status"`
}
