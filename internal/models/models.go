package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	Username      string     `json:"username" db:"username"`
	PasswordHash  string     `json:"-" db:"password_hash"` // Never send password hash to client
	FullName      *string    `json:"full_name,omitempty" db:"full_name"`
	AvatarURL     *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Status        string     `json:"status" db:"status"` // online, offline, away
	StatusMessage *string    `json:"status_message,omitempty" db:"status_message"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastSeenAt    *time.Time `json:"last_seen_at,omitempty" db:"last_seen_at"`
}

type Workspace struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	IconURL   *string   `json:"icon_url,omitempty" db:"icon_url"`
	OwnerID   uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WorkspaceMember struct {
	ID          uuid.UUID `json:"id" db:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id" db:"workspace_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Role        string    `json:"role" db:"role"` // owner, admin, member
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
}

type Channel struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id" db:"workspace_id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	IsPrivate   bool       `json:"is_private" db:"is_private"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Virtual fields
	UnreadCount int `json:"unread_count" db:"-"`
}

type ChannelMember struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	ChannelID  uuid.UUID  `json:"channel_id" db:"channel_id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	JoinedAt   time.Time  `json:"joined_at" db:"joined_at"`
	LastReadAt *time.Time `json:"last_read_at,omitempty" db:"last_read_at"`
}

type DirectMessage struct {
	ID          uuid.UUID `json:"id" db:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id" db:"workspace_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// Virtual fields
	Participants []*User `json:"participants,omitempty" db:"-"`
	UnreadCount  int     `json:"unread_count" db:"-"`
}

type DMParticipant struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	DMID       uuid.UUID  `json:"dm_id" db:"dm_id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	LastReadAt *time.Time `json:"last_read_at,omitempty" db:"last_read_at"`
}

type Message struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Content         string     `json:"content" db:"content"`
	SenderID        *uuid.UUID `json:"sender_id,omitempty" db:"sender_id"`
	ChannelID       *uuid.UUID `json:"channel_id,omitempty" db:"channel_id"`
	DMID            *uuid.UUID `json:"dm_id,omitempty" db:"dm_id"`
	ParentMessageID *uuid.UUID `json:"parent_message_id,omitempty" db:"parent_message_id"`
	EditedAt        *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`

	// Virtual fields (not in DB, populated by queries)
	Sender      *User        `json:"sender,omitempty" db:"-"`
	Reactions   []Reaction   `json:"reactions,omitempty" db:"-"`
	Attachments []Attachment `json:"attachments,omitempty" db:"-"`
	ReplyCount  int          `json:"reply_count,omitempty" db:"-"`
}

type Reaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Emoji     string    `json:"emoji" db:"emoji"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Virtual field
	User *User `json:"user,omitempty" db:"-"`
}

type Attachment struct {
	ID         uuid.UUID `json:"id" db:"id"`
	MessageID  uuid.UUID `json:"message_id" db:"message_id"`
	FileName   string    `json:"file_name" db:"file_name"`
	FileURL    string    `json:"file_url" db:"file_url"`
	FileType   *string   `json:"file_type,omitempty" db:"file_type"`
	FileSize   *int64    `json:"file_size,omitempty" db:"file_size"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

type WorkspaceInvite struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id" db:"workspace_id"`
	InviterID   uuid.UUID  `json:"inviter_id" db:"inviter_id"`
	Code        string     `json:"code" db:"code"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	MaxUses     *int       `json:"max_uses,omitempty" db:"max_uses"`
	Uses        int        `json:"uses" db:"uses"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}
