package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationMemeber struct {
	ID             uuid.UUID `db:"id" json:"id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	OrganizationID uuid.UUID `db:"user_id" json:"organization_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
