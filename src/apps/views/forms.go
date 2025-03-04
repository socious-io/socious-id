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

	Avatar     *uuid.UUID `json:"avatar" form:"avatar"`
	CoverImage *uuid.UUID `json:"cover_image" form:"cover_image"`
}

type OrganizationForm struct {
	Shortname string `db:"shortname" json:"shortname"`

	Email   *string `db:"email" json:"email"`
	Phone   *string `db:"phone" json:"phone"`
	Website *string `db:"website" json:"website"`

	Name        *string `db:"name" json:"name"`
	Bio         *string `db:"bio" json:"bio"`
	Description *string `db:"description" json:"description"`
	Mission     *string `db:"mission" json:"mission"`
	Culture     *string `db:"culture" json:"culture"`
	// Size        *string `db:"size" json:"size"`

	Country *string `db:"country" json:"country"`
	City    *string `db:"city" json:"city"`
	Address *string `db:"address" json:"address"`
	// GeonameId         *int    `db:"geoname_id" json:"geoname_id"`
	// MobileCountryCode *string `db:"mobile_country_code" json:"mobile_country_code"`

	// SocialCauses pq.StringArray `db:"social_causes" json:"social_causes"`

	// ImpactPoints float64 `db:"impact_points" json:"impact_points"`

	Image      *uuid.UUID `db:"image" json:"image"`             //logo JSONB
	CoverImage *uuid.UUID `db:"cover_image" json:"cover_image"` //cover JSONB
}
