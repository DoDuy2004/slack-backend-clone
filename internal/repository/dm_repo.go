package repository

import (
	"database/sql"

	"github.com/DoDuy2004/slack-clone-backend/internal/database"
	"github.com/DoDuy2004/slack-clone-backend/internal/models"
	"github.com/google/uuid"
)

type DMRepository interface {
	Create(dm *models.DirectMessage, participantIDs []uuid.UUID) error
	FindByParticipants(workspaceID uuid.UUID, userIDs []uuid.UUID) (*models.DirectMessage, error)
	ListByUserID(workspaceID, userID uuid.UUID) ([]*models.DirectMessage, error)
	GetByID(id uuid.UUID) (*models.DirectMessage, error)
	IsParticipant(dmID, userID uuid.UUID) (bool, error)
	UpdateLastRead(dmID, userID uuid.UUID) error
}

type postgresDMRepository struct {
	db *database.DB
}

func NewDMRepository(db *database.DB) DMRepository {
	return &postgresDMRepository{db: db}
}

func (r *postgresDMRepository) Create(dm *models.DirectMessage, participantIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO direct_messages (id, workspace_id) VALUES ($1, $2) RETURNING created_at`
	if err := tx.QueryRow(query, dm.ID, dm.WorkspaceID).Scan(&dm.CreatedAt); err != nil {
		return err
	}

	for _, userID := range participantIDs {
		participantQuery := `INSERT INTO dm_participants (dm_id, user_id) VALUES ($1, $2)`
		if _, err := tx.Exec(participantQuery, dm.ID, userID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresDMRepository) FindByParticipants(workspaceID uuid.UUID, userIDs []uuid.UUID) (*models.DirectMessage, error) {
	// Find a DM that has exactly these participants in this workspace
	query := `
		SELECT dm.id, dm.workspace_id, dm.created_at
		FROM direct_messages dm
		JOIN dm_participants dp ON dm.id = dp.dm_id
		WHERE dm.workspace_id = $1
		AND dp.user_id = ANY($2)
		GROUP BY dm.id, dm.workspace_id, dm.created_at
		HAVING COUNT(DISTINCT dp.user_id) = $3
		AND (SELECT COUNT(*) FROM dm_participants WHERE dm_id = dm.id) = $3
	`
	dm := &models.DirectMessage{}
	err := r.db.QueryRow(query, workspaceID, userIDs, len(userIDs)).Scan(&dm.ID, &dm.WorkspaceID, &dm.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func (r *postgresDMRepository) ListByUserID(workspaceID, userID uuid.UUID) ([]*models.DirectMessage, error) {
	query := `
		SELECT dm.id, dm.workspace_id, dm.created_at
		FROM direct_messages dm
		JOIN dm_participants dp ON dm.id = dp.dm_id
		WHERE dm.workspace_id = $1 AND dp.user_id = $2
		ORDER BY dm.created_at DESC
	`
	rows, err := r.db.Query(query, workspaceID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dms []*models.DirectMessage
	for rows.Next() {
		dm := &models.DirectMessage{}
		if err := rows.Scan(&dm.ID, &dm.WorkspaceID, &dm.CreatedAt); err != nil {
			return nil, err
		}
		dms = append(dms, dm)
	}

	if len(dms) > 0 {
		if err := r.attachParticipants(dms); err != nil {
			return nil, err
		}
		if err := r.attachUnreadCounts(dms, userID); err != nil {
			return nil, err
		}
	}

	return dms, nil
}

func (r *postgresDMRepository) attachParticipants(dms []*models.DirectMessage) error {
	dmIDs := make([]uuid.UUID, len(dms))
	dmMap := make(map[uuid.UUID]*models.DirectMessage)
	for i, dm := range dms {
		dmIDs[i] = dm.ID
		dmMap[dm.ID] = dm
		dm.Participants = []*models.User{}
	}

	query := `
		SELECT dp.dm_id, u.id, u.username, u.full_name, u.avatar_url, u.status, u.status_message
		FROM dm_participants dp
		JOIN users u ON dp.user_id = u.id
		WHERE dp.dm_id = ANY($1)
	`
	rows, err := r.db.Query(query, dmIDs)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dmID uuid.UUID
		u := &models.User{}
		if err := rows.Scan(&dmID, &u.ID, &u.Username, &u.FullName, &u.AvatarURL, &u.Status, &u.StatusMessage); err != nil {
			return err
		}
		if dm, ok := dmMap[dmID]; ok {
			dm.Participants = append(dm.Participants, u)
		}
	}
	return nil
}

func (r *postgresDMRepository) attachUnreadCounts(dms []*models.DirectMessage, userID uuid.UUID) error {
	dmIDs := make([]uuid.UUID, len(dms))
	dmMap := make(map[uuid.UUID]*models.DirectMessage)
	for i, dm := range dms {
		dmIDs[i] = dm.ID
		dmMap[dm.ID] = dm
	}

	query := `
		SELECT dp.dm_id, COUNT(m.id)
		FROM dm_participants dp
		LEFT JOIN messages m ON dp.dm_id = m.dm_id 
			AND m.created_at > COALESCE(dp.last_read_at, '1970-01-01')
			AND m.sender_id != $2
			AND m.deleted_at IS NULL
		WHERE dp.dm_id = ANY($1) AND dp.user_id = $2
		GROUP BY dp.dm_id
	`
	rows, err := r.db.Query(query, dmIDs, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dmID uuid.UUID
		var count int
		if err := rows.Scan(&dmID, &count); err != nil {
			return err
		}
		if dm, ok := dmMap[dmID]; ok {
			dm.UnreadCount = count
		}
	}
	return nil
}

func (r *postgresDMRepository) GetByID(id uuid.UUID) (*models.DirectMessage, error) {
	dm := &models.DirectMessage{}
	query := `SELECT id, workspace_id, created_at FROM direct_messages WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&dm.ID, &dm.WorkspaceID, &dm.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func (r *postgresDMRepository) IsParticipant(dmID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM dm_participants WHERE dm_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, dmID, userID).Scan(&exists)
	return exists, err
}

func (r *postgresDMRepository) UpdateLastRead(dmID, userID uuid.UUID) error {
	query := `
		UPDATE dm_participants
		SET last_read_at = CURRENT_TIMESTAMP
		WHERE dm_id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, dmID, userID)
	return err
}
