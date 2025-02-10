package models

import (
	"time"

	"github.com/lib/pq"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Shortname string    `db:"shortname" json:"shortname"`

	// @TODO: make it enum
	Status string `db:"status" json:"status"`

	Email   *string `db:"email" json:"email"`
	Phone   *string `db:"phone" json:"phone"`
	Website *string `db:"website" json:"website"`

	Name        *string `db:"name" json:"name"`
	Bio         *string `db:"bio" json:"bio"`
	Description *string `db:"description" json:"description"`
	Mission     *string `db:"mission" json:"mission"`
	Culture     *string `db:"culture" json:"culture"`
	Size        *string `db:"size" json:"size"`

	Country           *string `db:"country" json:"country"`
	City              *string `db:"city" json:"city"`
	Address           *string `db:"address" json:"address"`
	GeonameId         *int    `db:"geoname_id" json:"geoname_id"`
	MobileCountryCode *string `db:"mobile_country_code" json:"mobile_country_code"`

	SocialCauses pq.StringArray `db:"social_causes" json:"social_causes"`

	ImpactPoints float64 `db:"impact_points" json:"impact_points"`

	Image      *uuid.UUID `db:"image" json:"image"`
	CoverImage *uuid.UUID `db:"cover_image" json:"cover_image"`

	VerifiedImpact bool `db:"verified_impact" json:"verified_impact"`
	Verified       bool `db:"verified" json:"verified"`

	CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (Organization) FetchQuery() string {
	return "organizations/fetch"
}

func (*Organization) Create() error {
	return nil
}

func (*Organization) Update() error {
	return nil
}

func (*Organization) Remove() error {
	return nil
}

func getAllOrganizations() ([]Organization, error) {
	result := []Organization{}
	return result, nil
}

func GetOrganization(id uuid.UUID, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func getManyOrganizations(ids []uuid.UUID, identity uuid.UUID) ([]Organization, error) {
	result := []Organization{}
	return result, nil
}

func GetOrganizationByShortname(shortname string, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, identity.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func shortnameExistsOrganization(shortname string) (bool, error) {
	return false, nil
}

func searchOrganizations(query string) ([]Organization, error) { // Do we need to implement this?
	result := []Organization{}
	return result, nil
}
