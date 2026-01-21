package handler

import (
	"net/http"
	"time"

	"github.com/DoDuy2004/slack-clone-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InviteHandler struct {
	inviteService service.InviteService
}

func NewInviteHandler(inviteService service.InviteService) *InviteHandler {
	return &InviteHandler{inviteService: inviteService}
}

type CreateInviteRequest struct {
	ExpiresAt *time.Time `json:"expires_at"`
	MaxUses   *int       `json:"max_uses"`
}

func (h *InviteHandler) Create(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(uuid.UUID)

	workspaceIDStr := c.Param("id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Optional fields, so we can ignore error if body is empty or malformed but we want JSON
	}

	invite, err := h.inviteService.GenerateInvite(userID, workspaceID, req.ExpiresAt, req.MaxUses)
	if err != nil {
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

func (h *InviteHandler) Join(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(uuid.UUID)

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invite code is required"})
		return
	}

	workspace, err := h.inviteService.JoinWorkspace(userID, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workspace)
}
