package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/DoDuy2004/slack-clone-backend/internal/models"
	"github.com/DoDuy2004/slack-clone-backend/internal/repository"
	"github.com/google/uuid"
)

type InviteService interface {
	GenerateInvite(userID, workspaceID uuid.UUID, expiresAt *time.Time, maxUses *int) (*models.WorkspaceInvite, error)
	JoinWorkspace(userID uuid.UUID, code string) (*models.Workspace, error)
}

type inviteService struct {
	inviteRepo    repository.InviteRepository
	workspaceRepo repository.WorkspaceRepository
}

func NewInviteService(inviteRepo repository.InviteRepository, workspaceRepo repository.WorkspaceRepository) InviteService {
	return &inviteService{
		inviteRepo:    inviteRepo,
		workspaceRepo: workspaceRepo,
	}
}

func (s *inviteService) GenerateInvite(userID, workspaceID uuid.UUID, expiresAt *time.Time, maxUses *int) (*models.WorkspaceInvite, error) {
	// 1. Verify user is admin/owner
	member, err := s.workspaceRepo.GetMember(workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil || (member.Role != "owner" && member.Role != "admin") {
		return nil, ErrUnauthorized
	}

	// 2. Generate random code
	code, err := generateRandomCode(12)
	if err != nil {
		return nil, err
	}

	invite := &models.WorkspaceInvite{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		InviterID:   userID,
		Code:        code,
		ExpiresAt:   expiresAt,
		MaxUses:     maxUses,
		Uses:        0,
	}

	if err := s.inviteRepo.Create(invite); err != nil {
		return nil, err
	}

	return invite, nil
}

func (s *inviteService) JoinWorkspace(userID uuid.UUID, code string) (*models.Workspace, error) {
	// 1. Find invite
	invite, err := s.inviteRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, errors.New("invalid or expired invite code")
	}

	// 2. Validate invite
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invite code has expired")
	}
	if invite.MaxUses != nil && invite.Uses >= *invite.MaxUses {
		return nil, errors.New("invite code has reached maximum uses")
	}

	// 3. Check if already a member
	existing, err := s.workspaceRepo.GetMember(invite.WorkspaceID, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return s.workspaceRepo.FindByID(invite.WorkspaceID) // Already a member, just return workspace
	}

	// 4. Add member
	if err := s.workspaceRepo.AddMember(invite.WorkspaceID, userID, "member"); err != nil {
		return nil, err
	}

	// 5. Increment usage
	if err := s.inviteRepo.IncrementUses(invite.ID); err != nil {
		// Log error but don't fail join
	}

	return s.workspaceRepo.FindByID(invite.WorkspaceID)
}

func generateRandomCode(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
