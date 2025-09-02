package views

import (
	"encoding/json"
	"socious-id/src/apps/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type ConfirmForm struct {
	Confirmed  bool   `json:"confirmed" form:"confirmed" binding:"required"`
	IdentityId string `json:"identity_id" form:"identity_id" binding:"omitempty,uuid4"`
}

type ClientSecretForm struct {
	ClientSecret string `json:"client_secret" form:"client_secret" binding:"required"`
	ClientID     string `json:"client_id" form:"client_id" binding:"required"`
}

type AuthSessionForm struct {
	RedirectURL string    `json:"redirect_url" form:"redirect_url" binding:"required,url"`
	Policies    *[]string `json:"policies" form:"policies" binding:"omitempty,dive,required"`
}

type GetTokenForm struct {
	Code string `json:"code" form:"code" binding:"required"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required,jwt"`
}

type UserForm struct {
	Username string  `json:"username" form:"username" binding:"required"`
	Phone    *string `json:"phone" form:"phone" binding:"omitempty,e164"`

	FirstName *string `json:"first_name" form:"first_name"`
	LastName  *string `json:"last_name" form:"last_name"`
	Mission   *string `json:"mission" form:"mission"`
	Bio       *string `json:"bio" form:"bio"`

	City              *string `json:"city" form:"city"`
	Country           *string `json:"country" form:"country"`
	Address           *string `json:"address" form:"address"`
	MobileCountryCode *string `json:"mobile_country_code" form:"mobile_country_code"`

	AvatarID *uuid.UUID `json:"avatar_id" form:"avatar_id" binding:"omitempty,uuid4"`
	CoverID  *uuid.UUID `json:"cover_id" form:"cover_id" binding:"omitempty,uuid4"`
}

type UserCreateForm struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email" binding:"required,email"`
	Username  string `json:"username" form:"username"`
}

type UserUpdateStatusForm struct {
	Status models.UserStatusType `json:"status" form:"status" binding:"required,oneof=ACTIVE INACTIVE SUSPENDED"`
}

type OrganizationUpdateStatusForm struct {
	Status models.OrganizationStatusType `json:"status" form:"status" binding:"required,oneof=ACTIVE NOT_ACTIVE SUSPENDED PENDING"`
}

type OrganizationVerificationForm struct {
	Status models.OrganizationStatusType `json:"status" form:"status" binding:"required,oneof=ACTIVE NOT_ACTIVE SUSPENDED PENDING"`
}

type OrganizationForm struct {
	Shortname   string  `json:"shortname" form:"shortname" binding:"required"`
	Name        *string `json:"name" form:"name"`
	Bio         *string `json:"bio" form:"bio"`
	Description *string `json:"description" form:"description"`
	Email       *string `json:"email" form:"email" binding:"omitempty,email"`
	Phone       *string `json:"phone" form:"phone" binding:"omitempty,e164"`

	City    *string `json:"city" form:"city"`
	Country *string `json:"country" form:"country"`
	Address *string `json:"address" form:"address"`
	Website *string `json:"website" form:"website" binding:"omitempty,url"`

	Mission *string `json:"mission" form:"mission"`
	Culture *string `json:"culture" form:"culture"`

	LogoID  *uuid.UUID `json:"logo_id" form:"logo_id" binding:"omitempty,uuid4"`
	CoverID *uuid.UUID `json:"cover_id" form:"cover_id" binding:"omitempty,uuid4"`
}

type ImpactPointForm struct {
	UserID              uuid.UUID              `json:"user_id" form:"user_id" binding:"required,uuid4"`
	TotalPoints         int                    `json:"total_points" form:"total_points"` // pass through with value or 0 OR being missing
	SocialCause         string                 `json:"social_cause" form:"social_cause"`
	SocialCauseCategory string                 `json:"social_cause_category" form:"social_cause_category"`
	Value               float64                `json:"value" form:"value"` // pass through with value or 0 OR being missing
	Type                models.ImpactPointType `json:"type" form:"type" binding:"required,oneof=WORKSUBMIT DONATION VOLUNTEER VOTING OTHER"`
	AccessID            *uuid.UUID             `json:"access_id" form:"access_id" binding:"omitempty,uuid4"`
	Meta                *json.RawMessage       `json:"meta" form:"meta"`
	UniqueTag           string                 `json:"unique_tag" form:"unique_tag" binding:"required"`
}

type KYBVerificationForm struct {
	Documents []string `json:"documents" form:"documents" binding:"required,dive,uuid4"`
}

type AddWalletForm struct {
	Chain   string  `json:"chain" form:"chain" binding:"required"`
	ChainID *string `json:"chain_id" form:"chain_id"`
	Address string  `json:"address" form:"address" binding:"required"`
}

type AddCardForm struct {
	Token *string `json:"token" form:"token"`
}

type ReferralAchievementForm struct {
	RefereeID       uuid.UUID      `json:"referee_id" form:"referee_id" binding:"required,uuid4"`
	AchievementType string         `json:"achievement_type" form:"achievement_type" binding:"required"`
	Meta            types.JSONText `json:"meta" form:"meta"`
}

type CredentialForm struct {
	Type models.CredentialType `json:"type" form:"type" binding:"required,oneof=KYC BADGES"`
}
