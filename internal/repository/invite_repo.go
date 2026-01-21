package repository

import (
	"database/sql"

	"github.com/DoDuy2004/slack-clone-backend/internal/database"
	"github.com/DoDuy2004/slack-clone-backend/internal/models"
	"github.com/google/uuid"
)

type InviteRepository interface {
	Create(invite *models.WorkspaceInvite) error
	FindByCode(code string) (*models.WorkspaceInvite, error)
	IncrementUses(id uuid.UUID) error
}

type postgresInviteRepository struct {
	db *database.DB
}

func NewInviteRepository(db *database.DB) InviteRepository {
	return &postgresInviteRepository{db: db}
}

func (r *postgresInviteRepository) Create(invite *models.WorkspaceInvite) error {
	query := `
		INSERT INTO workspace_invites (id, workspace_id, inviter_id, code, expires_at, max_uses, uses)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`
	return r.db.QueryRow(
		query,
		invite.ID,
		invite.WorkspaceID,
		invite.InviterID,
		invite.Code,
		invite.ExpiresAt,
		invite.MaxUses,
		invite.Uses,
	).Scan(&invite.CreatedAt)
}

func (r *postgresInviteRepository) FindByCode(code string) (*models.WorkspaceInvite, error) {
	invite := &models.WorkspaceInvite{}
	query := `
		SELECT id, workspace_id, inviter_id, code, expires_at, max_uses, uses, created_at
		FROM workspace_invites
		WHERE code = $1
	`
	err := r.db.QueryRow(query, code).Scan(
		&invite.ID,
		&invite.WorkspaceID,
		&invite.InviterID,
		&invite.Code,
		&invite.ExpiresAt,
		&invite.MaxUses,
		&invite.Uses,
		&invite.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return invite, nil
}

func (r *postgresInviteRepository) IncrementUses(id uuid.UUID) error {
	query := `UPDATE workspace_invites SET uses = uses + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
