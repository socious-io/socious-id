package views

import (
	"github.com/google/uuid"
)

type ConfirmForm struct {
	Confirmed bool `json:"confirmed" form:"confirmed"`
}

type AuthSessionForm struct {
	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"`
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
	RedirectURL  string `json:"redirect_url" form:"redirect_url" validate:"required"`
}

type GetTokenForm struct {
	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"`
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
	Code         string `json:"code" form:"code" validate:"required"`
}

type UserForm struct {
	Username string  `json:"username" form:"username"`
	Phone    *string `json:"phone" form:"phone"`

	FirstName *string `json:"first_name" form:"first_name"`
	LastName  *string `json:"last_name" form:"last_name"`
	Mission   *string `json:"mission" form:"mission"`
	Bio       *string `json:"bio" form:"bio"`

	City              *string `json:"city" form:"city"`
	Country           *string `json:"country" form:"country"`
	Address           *string `json:"address" form:"address"`
	MobileCountryCode *string `json:"mobile_country_code" form:"mobile_country_code"`

	AvatarID *uuid.UUID `json:"avatar_id" form:"avatar_id"`
	CoverID  *uuid.UUID `json:"cover_id" form:"cover_id"`
}

type OrganizationForm struct {
	Shortname   string  `db:"shortname" json:"shortname"`
	Name        *string `db:"name" json:"name"`
	Bio         *string `db:"bio" json:"bio"`
	Description *string `db:"description" json:"description"`
	Email       *string `db:"email" json:"email"`
	Phone       *string `db:"phone" json:"phone"`

	City    *string `db:"city" json:"city"`
	Country *string `db:"country" json:"country"`
	Address *string `db:"address" json:"address"`
	Website *string `db:"website" json:"website"`

	Mission *string `db:"mission" json:"mission"`
	Culture *string `db:"culture" json:"culture"`

	LogoID  *uuid.UUID `db:"logo_id" json:"logo_id"`
	CoverID *uuid.UUID `db:"cover_id" json:"cover_id"`
}
