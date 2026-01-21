package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDMRequest struct {
	ParticipantIDs []uuid.UUID `json:"participant_ids" binding:"required"`
}

type DMResponse struct {
	ID           uuid.UUID        `json:"id"`
	WorkspaceID  uuid.UUID        `json:"workspace_id"`
	Participants []UserSummary    `json:"participants"`
	LastMessage  *MessageResponse `json:"last_message,omitempty"`
	UnreadCount  int              `json:"unread_count"`
	CreatedAt    time.Time        `json:"created_at"`
}
