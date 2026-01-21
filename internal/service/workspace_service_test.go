package service

import (
	"testing"

	"github.com/DoDuy2004/slack-clone-backend/internal/models"
	"github.com/DoDuy2004/slack-clone-backend/internal/models/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkspaceRepository is a mock implementation of WorkspaceRepository
type MockWorkspaceRepository struct {
	mock.Mock
}

func (m *MockWorkspaceRepository) Create(workspace *models.Workspace, ownerID uuid.UUID) error {
	args := m.Called(workspace, ownerID)
	return args.Error(0)
}

func (m *MockWorkspaceRepository) FindByID(id uuid.UUID) (*models.Workspace, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Workspace), args.Error(1)
}

func (m *MockWorkspaceRepository) FindBySlug(slug string) (*models.Workspace, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Workspace), args.Error(1)
}

func (m *MockWorkspaceRepository) ListByUserID(userID uuid.UUID) ([]*models.Workspace, error) {
	args := m.Called(userID)
	return args.Get(0).([]*models.Workspace), args.Error(1)
}

func (m *MockWorkspaceRepository) Update(workspace *models.Workspace) error {
	args := m.Called(workspace)
	return args.Error(0)
}

func (m *MockWorkspaceRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWorkspaceRepository) AddMember(workspaceID, userID uuid.UUID, role string) error {
	args := m.Called(workspaceID, userID, role)
	return args.Error(0)
}

func (m *MockWorkspaceRepository) RemoveMember(workspaceID, userID uuid.UUID) error {
	args := m.Called(workspaceID, userID)
	return args.Error(0)
}

func (m *MockWorkspaceRepository) GetMember(workspaceID, userID uuid.UUID) (*models.WorkspaceMember, error) {
	args := m.Called(workspaceID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkspaceMember), args.Error(1)
}

func (m *MockWorkspaceRepository) ListMembers(workspaceID uuid.UUID) ([]*models.WorkspaceMember, error) {
	args := m.Called(workspaceID)
	return args.Get(0).([]*models.WorkspaceMember), args.Error(1)
}

func TestCreateWorkspace(t *testing.T) {
	mockRepo := new(MockWorkspaceRepository)
	svc := NewWorkspaceService(mockRepo)

	userID := uuid.New()
	req := &dto.CreateWorkspaceRequest{
		Name: "Test Workspace",
		Slug: "test-workspace",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindBySlug", req.Slug).Return(nil, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.Workspace"), userID).Return(nil).Once()

		ws, err := svc.CreateWorkspace(userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, req.Name, ws.Name)
		assert.Equal(t, req.Slug, ws.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SlugAlreadyExists", func(t *testing.T) {
		mockRepo.On("FindBySlug", req.Slug).Return(&models.Workspace{}, nil).Once()

		ws, err := svc.CreateWorkspace(userID, req)

		assert.Error(t, err)
		assert.Equal(t, ErrWorkspaceExists, err)
		assert.Nil(t, ws)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateWorkspace(t *testing.T) {
	mockRepo := new(MockWorkspaceRepository)
	svc := NewWorkspaceService(mockRepo)

	userID := uuid.New()
	wsID := uuid.New()
	req := &dto.UpdateWorkspaceRequest{
		Name: stringPtr("New Name"),
	}

	t.Run("Success", func(t *testing.T) {
		existingWS := &models.Workspace{
			ID:      wsID,
			OwnerID: userID,
			Name:    "Old Name",
		}
		mockRepo.On("FindByID", wsID).Return(existingWS, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Workspace")).Return(nil).Once()

		ws, err := svc.UpdateWorkspace(userID, wsID, req)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", ws.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		existingWS := &models.Workspace{
			ID:      wsID,
			OwnerID: uuid.New(), // Different owner
		}
		mockRepo.On("FindByID", wsID).Return(existingWS, nil).Once()

		ws, err := svc.UpdateWorkspace(userID, wsID, req)

		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorized, err)
		assert.Nil(t, ws)
		mockRepo.AssertExpectations(t)
	})
}

func stringPtr(s string) *string {
	return &s
}
