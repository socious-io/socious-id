package views

import (
	"encoding/json"
	"socious-id/src/apps/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
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
	RedirectURL string `json:"redirect_url" form:"redirect_url" validate:"required"`
}

type GetTokenForm struct {
	Code string `json:"code" form:"code" validate:"required"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required"`
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
	Status models.StatusType `json:"status" form:"status" validate:"required"`
}

type OrganizationUpdateStatusForm struct {
	Status models.OrganizationStatusType `json:"status" form:"status" validate:"required"`
}

type OrganizationVerificationForm struct {
	Status models.OrganizationStatusType `json:"status" form:"status" validate:"required"`
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

type ImpactPointForm struct {
	UserID              uuid.UUID              `json:"user_id" form:"user_id" validate:"required"`
	TotalPoints         int                    `json:"total_points" form:"total_points"`
	SocialCause         string                 `json:"social_cause" form:"social_cause"`
	SocialCauseCategory string                 `json:"social_cause_category" form:"social_cause_category"`
	Value               float64                `json:"value" form:"value"`
	Type                models.ImpactPointType `json:"type" form:"type" validate:"required,oneof=WORKSUBMIT DONATION VOLUNTEER OTHER"`
	AccessID            *uuid.UUID             `json:"access_id" form:"access_id"`
	Meta                *json.RawMessage       `json:"meta" form:"meta"`
}

type KYBVerificationForm struct {
	Documents []string `json:"documents"`
}

type ReferralAchievementForm struct {
	IdentityID string         `json:"identity_id"`
	Type       string         `json:"type"`
	Meta       types.JSONText `json:"meta"`
}
