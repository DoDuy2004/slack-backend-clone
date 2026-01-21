package repository

import (
	"database/sql"
	"fmt"

	"github.com/DoDuy2004/slack-clone-backend/internal/database"
	"github.com/DoDuy2004/slack-clone-backend/internal/models"
	"github.com/google/uuid"
)

type WorkspaceRepository interface {
	Create(workspace *models.Workspace, ownerID uuid.UUID) error
	FindByID(id uuid.UUID) (*models.Workspace, error)
	FindBySlug(slug string) (*models.Workspace, error)
	ListByUserID(userID uuid.UUID) ([]*models.Workspace, error)
	Update(workspace *models.Workspace) error
	Delete(id uuid.UUID) error

	// Member operations
	AddMember(workspaceID, userID uuid.UUID, role string) error
	RemoveMember(workspaceID, userID uuid.UUID) error
	GetMember(workspaceID, userID uuid.UUID) (*models.WorkspaceMember, error)
	ListMembers(workspaceID uuid.UUID) ([]*models.WorkspaceMember, error)
}

type postgresWorkspaceRepository struct {
	db *database.DB
}

func NewWorkspaceRepository(db *database.DB) WorkspaceRepository {
	return &postgresWorkspaceRepository{db: db}
}

func (r *postgresWorkspaceRepository) Create(workspace *models.Workspace, ownerID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert Workspace
	query := `
		INSERT INTO workspaces (id, name, slug, icon_url, owner_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`
	err = tx.QueryRow(
		query,
		workspace.ID,
		workspace.Name,
		workspace.Slug,
		workspace.IconURL,
		ownerID,
	).Scan(&workspace.CreatedAt, &workspace.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// 2. Add Owner as first member
	memberQuery := `
		INSERT INTO workspace_members (id, workspace_id, user_id, role)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(memberQuery, uuid.New(), workspace.ID, ownerID, "owner")
	if err != nil {
		return fmt.Errorf("failed to add owner as member: %w", err)
	}

	return tx.Commit()
}

func (r *postgresWorkspaceRepository) FindByID(id uuid.UUID) (*models.Workspace, error) {
	ws := &models.Workspace{}
	query := `SELECT id, name, slug, icon_url, owner_id, created_at, updated_at FROM workspaces WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&ws.ID, &ws.Name, &ws.Slug, &ws.IconURL, &ws.OwnerID, &ws.CreatedAt, &ws.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (r *postgresWorkspaceRepository) FindBySlug(slug string) (*models.Workspace, error) {
	ws := &models.Workspace{}
	query := `SELECT id, name, slug, icon_url, owner_id, created_at, updated_at FROM workspaces WHERE slug = $1`
	err := r.db.QueryRow(query, slug).Scan(
		&ws.ID, &ws.Name, &ws.Slug, &ws.IconURL, &ws.OwnerID, &ws.CreatedAt, &ws.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (r *postgresWorkspaceRepository) ListByUserID(userID uuid.UUID) ([]*models.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.slug, w.icon_url, w.owner_id, w.created_at, w.updated_at
		FROM workspaces w
		JOIN workspace_members wm ON w.id = wm.workspace_id
		WHERE wm.user_id = $1
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*models.Workspace
	for rows.Next() {
		ws := &models.Workspace{}
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.Slug, &ws.IconURL, &ws.OwnerID, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

func (r *postgresWorkspaceRepository) Update(workspace *models.Workspace) error {
	query := `
		UPDATE workspaces
		SET name = $1, icon_url = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	_, err := r.db.Exec(query, workspace.Name, workspace.IconURL, workspace.ID)
	return err
}

func (r *postgresWorkspaceRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM workspaces WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postgresWorkspaceRepository) AddMember(workspaceID, userID uuid.UUID, role string) error {
	query := `
		INSERT INTO workspace_members (id, workspace_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (workspace_id, user_id) DO UPDATE SET role = EXCLUDED.role
	`
	_, err := r.db.Exec(query, uuid.New(), workspaceID, userID, role)
	return err
}

func (r *postgresWorkspaceRepository) RemoveMember(workspaceID, userID uuid.UUID) error {
	query := `DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, workspaceID, userID)
	return err
}

func (r *postgresWorkspaceRepository) GetMember(workspaceID, userID uuid.UUID) (*models.WorkspaceMember, error) {
	member := &models.WorkspaceMember{}
	query := `SELECT id, workspace_id, user_id, role, joined_at FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, workspaceID, userID).Scan(
		&member.ID, &member.WorkspaceID, &member.UserID, &member.Role, &member.JoinedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return member, err
}

func (r *postgresWorkspaceRepository) ListMembers(workspaceID uuid.UUID) ([]*models.WorkspaceMember, error) {
	query := `SELECT id, workspace_id, user_id, role, joined_at FROM workspace_members WHERE workspace_id = $1`
	rows, err := r.db.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.WorkspaceMember
	for rows.Next() {
		m := &models.WorkspaceMember{}
		if err := rows.Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}
