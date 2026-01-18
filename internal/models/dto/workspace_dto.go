package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateWorkspaceRequest struct {
	Name    string  `json:"name" binding:"required,min=3,max=50"`
	Slug    string  `json:"slug" binding:"required,min=3,max=50,lowercase"`
	IconURL *string `json:"icon_url,omitempty"`
}

type UpdateWorkspaceRequest struct {
	Name    *string `json:"name,omitempty" binding:"omitempty,min=3,max=50"`
	IconURL *string `json:"icon_url,omitempty"`
}

type WorkspaceResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IconURL   *string   `json:"icon_url,omitempty"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkspaceMemberResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  *string   `json:"full_name,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}
