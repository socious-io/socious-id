package models

import (
	"time"

	"github.com/google/uuid"
)

type VerificationCredential struct {
	ID uuid.UUID `db:"id" json:"id"`

	// @TODO: make it enum
	Status string `db:"status" json:"status"`

	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	ConnectionID  *string   `db:"connection_id" json:"connection_id"`
	ConnectionUrl *string   `db:"connection_url" json:"connection_url"`
	PresentID     *string   `db:"present_id" json:"present_id"`
	Body          *string   `db:"body" json:"body"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (VerificationCredential) TableName() string {
	return "verification_credentials"
}

func (VerificationCredential) FetchQuery() string {
	return "verification_credentials/fetch"
}
