package handler

import (
	"net/http"

	"github.com/DoDuy2004/slack-clone-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReadReceiptHandler struct {
	readService service.ReadReceiptService
}

func NewReadReceiptHandler(readService service.ReadReceiptService) *ReadReceiptHandler {
	return &ReadReceiptHandler{readService: readService}
}

func (h *ReadReceiptHandler) MarkChannelAsRead(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(uuid.UUID)

	channelIDStr := c.Param("id")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	if err := h.readService.MarkChannelAsRead(userID, channelID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ReadReceiptHandler) MarkDMAsRead(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(uuid.UUID)

	dmIDStr := c.Param("id")
	dmID, err := uuid.Parse(dmIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DM ID"})
		return
	}

	if err := h.readService.MarkDMAsRead(userID, dmID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
