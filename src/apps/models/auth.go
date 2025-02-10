package models

import (
	"context"
	"fmt"
	"socious-id/src/apps/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	database "github.com/socious-io/pkg_database"
)

type Access struct {
	ID           uuid.UUID      `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Description  string         `db:"description" json:"description"`
	Scopes       pq.StringArray `db:"scopes" json:"scopes"`
	Logo         *uuid.UUID     `db:"logo" json:"logo"`
	ClientID     string         `db:"client_id" json:"client_id"`
	ClientSecret string         `db:"client_secret" json:"-"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}

type OTP struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Code          string     `db:"code" json:"code"`
	Type          OTPType    `db:"type" json:"type"`
	UserID        uuid.UUID  `db:"user_id" json:"user_id"`
	AuthSessionID uuid.UUID  `db:"auth_session_id" json:"auth_session_id"`
	ExpireAt      time.Time  `db:"expire_at" json:"expire_at"`
	VerifiedAt    *time.Time `db:"verified_at" json:"verified_at"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type AuthSession struct {
	ID          uuid.UUID `db:"id" json:"id"`
	AccessID    uuid.UUID `db:"access_id" json:"access_id"`
	RedirectURL string    `db:"redirect_url" json:"redirect_url"`

	Access     *Access        `db:"-" json:"access"`
	AccessJson types.JSONText `db:"access" json:"-"`

	ExpireAt   time.Time  `db:"expire_at" json:"expire_at"`
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

func (a *AuthSession) Create(ctx context.Context) error {
	// @IMPORTANT @TODO: read from DB
	s, err := GetAuthSession(uuid.New())
	if err != nil {
		return err
	}
	utils.Copy(s, a)
	return nil
}

func (o *OTP) Create(ctx context.Context) error {
	codeLength := 6
	if o.Type == SSOOTP {
		codeLength = 12
	}
	o.Code = fmt.Sprintf("%d", utils.GenerateRandomDigits(codeLength))
	rows, err := database.Query(
		ctx,
		"auth/create_otp",
		o.Type, o.RefID, o.Code,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(o); err != nil {
			return err
		}
	}
	return nil
}

func GetAuthSession(id uuid.UUID) (*AuthSession, error) {
	// @IMPORTANT @TODO: read from DB
	access, _ := GetAccessByClientID("")
	return &AuthSession{
		ID:          id,
		RedirectURL: "https://app.socious.io",
		Access:      access,
	}, nil
}

func GetAccessByClientID(clientID string) (*Access, error) {
	// @IMPORTANT @TODO: read from DB
	return &Access{
		ID:           uuid.New(),
		Name:         "Socious",
		Description:  "Socious",
		ClientID:     "test",
		ClientSecret: "test",
	}, nil
}
