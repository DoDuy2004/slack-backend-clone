package service

import (
	"encoding/json"
	"time"

	"github.com/DoDuy2004/slack-clone-backend/internal/repository"
	"github.com/DoDuy2004/slack-clone-backend/internal/websocket"
	"github.com/google/uuid"
)

type ReadReceiptService interface {
	MarkChannelAsRead(userID, channelID uuid.UUID) error
	MarkDMAsRead(userID, dmID uuid.UUID) error
}

type readReceiptService struct {
	channelRepo repository.ChannelRepository
	dmRepo      repository.DMRepository
	hub         *websocket.Hub
}

func NewReadReceiptService(channelRepo repository.ChannelRepository, dmRepo repository.DMRepository, hub *websocket.Hub) ReadReceiptService {
	return &readReceiptService{
		channelRepo: channelRepo,
		dmRepo:      dmRepo,
		hub:         hub,
	}
}

func (s *readReceiptService) MarkChannelAsRead(userID, channelID uuid.UUID) error {
	if err := s.channelRepo.UpdateLastRead(channelID, userID); err != nil {
		return err
	}

	// Broadcast read receipt
	s.broadcastReadReceipt(userID, &channelID, nil)
	return nil
}

func (s *readReceiptService) MarkDMAsRead(userID, dmID uuid.UUID) error {
	if err := s.dmRepo.UpdateLastRead(dmID, userID); err != nil {
		return err
	}

	// Broadcast read receipt
	s.broadcastReadReceipt(userID, nil, &dmID)
	return nil
}

func (s *readReceiptService) broadcastReadReceipt(userID uuid.UUID, channelID, dmID *uuid.UUID) {
	payload := map[string]interface{}{
		"user_id":   userID,
		"timestamp": time.Now(),
	}
	if channelID != nil {
		payload["channel_id"] = *channelID
	}
	if dmID != nil {
		payload["dm_id"] = *dmID
	}

	data, _ := json.Marshal(payload)
	s.hub.Broadcast(&websocket.WSMessage{
		Type:      "message.read_receipt",
		Payload:   data,
		ChannelID: channelID,
		DMID:      dmID,
	})
}
