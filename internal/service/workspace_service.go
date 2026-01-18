package service

import (
	"errors"

	"github.com/DoDuy2004/slack-clone/backend/internal/models"
	"github.com/DoDuy2004/slack-clone/backend/internal/models/dto"
	"github.com/DoDuy2004/slack-clone/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrWorkspaceExists   = errors.New("workspace with this slug already exists")
	ErrUnauthorized      = errors.New("unauthorized")
)

type WorkspaceService interface {
	CreateWorkspace(userID uuid.UUID, req *dto.CreateWorkspaceRequest) (*models.Workspace, error)
	GetWorkspace(id uuid.UUID) (*models.Workspace, error)
	GetWorkspaceBySlug(slug string) (*models.Workspace, error)
	ListUserWorkspaces(userID uuid.UUID) ([]*models.Workspace, error)
	UpdateWorkspace(userID uuid.UUID, wsID uuid.UUID, req *dto.UpdateWorkspaceRequest) (*models.Workspace, error)
	DeleteWorkspace(userID uuid.UUID, wsID uuid.UUID) error
}

type workspaceService struct {
	workspaceRepo repository.WorkspaceRepository
}

func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository) WorkspaceService {
	return &workspaceService{
		workspaceRepo: workspaceRepo,
	}
}

func (s *workspaceService) CreateWorkspace(userID uuid.UUID, req *dto.CreateWorkspaceRequest) (*models.Workspace, error) {
	// Check if slug exists
	existing, err := s.workspaceRepo.FindBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrWorkspaceExists
	}

	workspace := &models.Workspace{
		ID:      uuid.New(),
		Name:    req.Name,
		Slug:    req.Slug,
		IconURL: req.IconURL,
		OwnerID: userID,
	}

	if err := s.workspaceRepo.Create(workspace, userID); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (s *workspaceService) GetWorkspace(id uuid.UUID) (*models.Workspace, error) {
	ws, err := s.workspaceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, ErrWorkspaceNotFound
	}
	return ws, nil
}

func (s *workspaceService) GetWorkspaceBySlug(slug string) (*models.Workspace, error) {
	ws, err := s.workspaceRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, ErrWorkspaceNotFound
	}
	return ws, nil
}

func (s *workspaceService) ListUserWorkspaces(userID uuid.UUID) ([]*models.Workspace, error) {
	return s.workspaceRepo.ListByUserID(userID)
}

func (s *workspaceService) UpdateWorkspace(userID uuid.UUID, wsID uuid.UUID, req *dto.UpdateWorkspaceRequest) (*models.Workspace, error) {
	ws, err := s.workspaceRepo.FindByID(wsID)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, ErrWorkspaceNotFound
	}

	// Only owner can update
	if ws.OwnerID != userID {
		return nil, ErrUnauthorized
	}

	if req.Name != nil {
		ws.Name = *req.Name
	}
	if req.IconURL != nil {
		ws.IconURL = req.IconURL
	}

	if err := s.workspaceRepo.Update(ws); err != nil {
		return nil, err
	}

	return ws, nil
}

func (s *workspaceService) DeleteWorkspace(userID uuid.UUID, wsID uuid.UUID) error {
	ws, err := s.workspaceRepo.FindByID(wsID)
	if err != nil {
		return err
	}
	if ws == nil {
		return ErrWorkspaceNotFound
	}

	// Only owner can delete
	if ws.OwnerID != userID {
		return ErrUnauthorized
	}

	return s.workspaceRepo.Delete(wsID)
}
