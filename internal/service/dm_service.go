package service

import (
	"errors"

	"github.com/DoDuy2004/slack-clone-backend/internal/models"
	"github.com/DoDuy2004/slack-clone-backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrDMNotFound = errors.New("dm session not found")
)

type DMService interface {
	CreateDM(userID uuid.UUID, workspaceID uuid.UUID, participantIDs []uuid.UUID) (*models.DirectMessage, error)
	ListUserDMs(userID, workspaceID uuid.UUID) ([]*models.DirectMessage, error)
}

type dmService struct {
	dmRepo        repository.DMRepository
	workspaceRepo repository.WorkspaceRepository
	userRepo      repository.UserRepository
}

func NewDMService(
	dmRepo repository.DMRepository,
	workspaceRepo repository.WorkspaceRepository,
	userRepo repository.UserRepository,
) DMService {
	return &dmService{
		dmRepo:        dmRepo,
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

func (s *dmService) CreateDM(userID uuid.UUID, workspaceID uuid.UUID, participantIDs []uuid.UUID) (*models.DirectMessage, error) {
	// Add current user to participants list if not present
	found := false
	for _, pID := range participantIDs {
		if pID == userID {
			found = true
			break
		}
	}
	if !found {
		participantIDs = append(participantIDs, userID)
	}

	if len(participantIDs) < 2 {
		return nil, errors.New("a DM must have at least 2 participants")
	}

	// 1. Verify all users are members of the workspace
	for _, pID := range participantIDs {
		member, err := s.workspaceRepo.GetMember(workspaceID, pID)
		if err != nil || member == nil {
			return nil, errors.New("all participants must be members of the workspace")
		}
	}

	// 2. Check if DM already exists
	dm, err := s.dmRepo.FindByParticipants(workspaceID, participantIDs)
	if err != nil {
		return nil, err
	}

	if dm != nil {
		// DM already exists, but we want to return it with participants populated
		// ListUserDMs handles this, but for a single DM we might need another repo method or just fetch here
		return s.populateDM(dm)
	}

	// 3. Create new DM
	newDM := &models.DirectMessage{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
	}

	if err := s.dmRepo.Create(newDM, participantIDs); err != nil {
		return nil, err
	}

	return s.populateDM(newDM)
}

func (s *dmService) populateDM(dm *models.DirectMessage) (*models.DirectMessage, error) {
	// For a single DM creation/get, individual queries are fine.
	return dm, nil
}

func (s *dmService) ListUserDMs(userID, workspaceID uuid.UUID) ([]*models.DirectMessage, error) {
	// Verify workspace membership
	member, err := s.workspaceRepo.GetMember(workspaceID, userID)
	if err != nil || member == nil {
		return nil, ErrUnauthorized
	}

	return s.dmRepo.ListByUserID(workspaceID, userID)
}
