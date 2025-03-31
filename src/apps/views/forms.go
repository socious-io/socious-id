package views

import (
	"socious-id/src/apps/models"

	"github.com/google/uuid"
)

type ConfirmForm struct {
	Confirmed  bool   `json:"confirmed" form:"confirmed"`
	IdentityId string `json:"identity_id" form:"identity_id"`
}

type ClientSecretForm struct {
	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"`
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
}

type AuthSessionForm struct {
	ClientSecretForm
	RedirectURL string `json:"redirect_url" form:"redirect_url" validate:"required"`
}

type GetTokenForm struct {
	ClientSecretForm
	Code string `json:"code" form:"code" validate:"required"`
}

type RefreshTokenForm struct {
	ClientSecretForm
	RefreshToken string `json:"code" form:"code" validate:"required"`
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

type UserUpdateStatusForm struct {
	ClientSecretForm
	Status models.StatusType `json:"status" form:"status" validate:"required"`
}

type OrganizationForm struct {
	Shortname   string  `json:"shortname" form:"shortname"`
	Name        *string `json:"name" form:"name"`
	Bio         *string `json:"bio" form:"bio"`
	Description *string `json:"description" form:"description"`
	Email       *string `json:"email" form:"email"`
	Phone       *string `json:"phone" form:"phone"`

	City    *string `json:"city" form:"city"`
	Country *string `json:"country" form:"country"`
	Address *string `json:"address" form:"address"`
	Website *string `json:"website" form:"website"`

	Mission *string `json:"mission" form:"mission"`
	Culture *string `json:"culture" form:"culture"`

	LogoID  *uuid.UUID `json:"logo_id" form:"logo_id"`
	CoverID *uuid.UUID `json:"cover_id" form:"cover_id"`
}
