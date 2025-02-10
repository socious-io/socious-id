package models

import (
	"context"
	"time"

	"github.com/lib/pq"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Username        string    `db:"username" json:"username"`
	Password        *string   `db:"password" json:"-"`
	PasswordExpired bool      `db:"password_expired" json:"password_expired"`

	// @TODO: make it enum
	Status string `db:"status" json:"status"`

	Email     string  `db:"email" json:"email"`
	EmailText *string `db:"email_text" json:"email_text"`
	Phone     *string `db:"phone" json:"phone"`

	FirstName         *string `db:"first_name" json:"first_name"`
	LastName          *string `db:"last_name" json:"last_name"`
	Mission           *string `db:"mission" json:"mission"`
	Bio               *string `db:"bio" json:"bio"`
	DescriptionSearch *string `db:"description_search" json:"description_search"`

	City              *string `db:"city" json:"city"`
	Country           *string `db:"country" json:"country"`
	Address           *string `db:"address" json:"address"`
	GeonameId         *int64  `db:"geoname_id" json:"geoname_id"`
	MobileCountryCode *string `db:"mobile_country_code" json:"mobile_country_code"`

	Avatar     *uuid.UUID `db:"avatar" json:"avatar"`
	CoverImage *uuid.UUID `db:"cover_image" json:"cover_image"`

	SocialCauses pq.StringArray `db:"social_causes" json:"social_causes"`

	IdentityVerifiedAt *time.Time `db:"identity_verified_at" json:"identity_verified_at"`
	EmailVerifiedAt    *time.Time `db:"email_verified_at" json:"email_verified_at"`
	PhoneVerifiedAt    *time.Time `db:"phone_verified_at" json:"phone_verified_at"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (User) FetchQuery() string {
	return "users/fetch"
}

func (u *User) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/register",
		u.FirstName, u.LastName, u.Username, u.Email, u.Password,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Verify(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/verify",
		u.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) ExpirePassword(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/expire_password",
		u.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) UpdatePassword(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/update_password",
		u.ID, u.Password,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) UpdateProfile(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/update_profile",
		u.ID, u.FirstName, u.LastName, u.Bio, u.Phone, u.Username,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func GetUser(id uuid.UUID) (*User, error) {
	u := new(User)
	if err := database.Fetch(u, id.String()); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_email", email); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByUsername(username string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_username", username); err != nil {
		return nil, err
	}
	return u, nil
}
