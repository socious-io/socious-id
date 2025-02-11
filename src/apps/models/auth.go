package models

import (
	"context"
	"fmt"
	"socious-id/src/apps/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type Access struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	ClientID     string    `db:"client_id" json:"client_id"`
	ClientSecret string    `db:"client_secret" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type AuthSession struct {
	ID          uuid.UUID `db:"id" json:"id"`
	RedirectURL string    `db:"redirect_url" json:"redirect_url"`

	AccessID   uuid.UUID      `db:"access_id" json:"access_id"`
	Access     *Access        `db:"-" json:"access"`
	AccessJson types.JSONText `db:"access" json:"-"`

	ExpireAt   time.Time  `db:"expire_at" json:"expire_at"`
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

type OTP struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Code string    `db:"code" json:"code"`
	Type OTPType   `db:"type" json:"type"`

	UserID   uuid.UUID      `db:"user_id" json:"user_id"`
	User     User           `db:"-" json:"user"`
	UserJson types.JSONText `db:"user" json:"-"`

	AuthSessionID   *uuid.UUID      `db:"auth_session_id" json:"auth_session_id"`
	AuthSession     *AuthSession    `db:"-" json:"auth_session"`
	AuthSessionJson *types.JSONText `db:"auth_session" json:"-"`

	ExpireAt   time.Time  `db:"expire_at" json:"expire_at"`
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

func (Access) TableName() string {
	return "accesses"
}

func (Access) FetchQuery() string {
	return "auth/fetch_access"
}

func (a *Access) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"auth/create_access",
		a.Name, a.Description, a.ClientID, a.ClientSecret,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(a); err != nil {
			return err
		}
	}
	return nil
}

func (AuthSession) TableName() string {
	return "auth_sessions"
}

func (AuthSession) FetchQuery() string {
	return "auth/fetch_auth_session"
}

func (a *AuthSession) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"auth/create_auth_session",
		a.RedirectURL, a.AccessID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(a); err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthSession) Verify(ctx context.Context) error {
	rows, err := database.Query(ctx, "auth/auth_session_verify", a.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(a); err != nil {
			return err
		}
	}
	return nil
}

func (OTP) TableName() string {
	return "otps"
}

func (OTP) FetchQuery() string {
	return "auth/fetch_otp"
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
		o.Type, o.UserID, o.AuthSessionID, o.Code,
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

func (o *OTP) Verify(ctx context.Context) error {
	if o.AuthSession != nil {
		if err := o.AuthSession.Verify(ctx); err != nil {
			return err
		}
	}
	rows, err := database.Query(ctx, "auth/otp_verify", o.ID)
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

func GetAccessByClientID(clientID string) (*Access, error) {
	a := new(Access)
	if err := database.Get(a, "auth/fetch_access_by_client_id"); err != nil {
		return nil, err
	}
	return a, nil
}

func GetAuthSession(id uuid.UUID) (*AuthSession, error) {
	a := new(AuthSession)
	if err := database.Fetch(a, id); err != nil {
		return nil, err
	}
	return a, nil
}

func GetOTPByCode(code string) (*OTP, error) {
	o := new(OTP)
	if err := database.Get(o, "auth/fetch_otp_by_code", code); err != nil {
		return nil, err
	}
	return o, nil
}
