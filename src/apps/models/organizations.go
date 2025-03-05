package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Shortname string    `db:"shortname" json:"shortname"`

	Status StatusType `db:"status" json:"status"`

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

	Logo     *Media         `db:"-" json:"logo"`
	LogoJson types.JSONText `db:"logo" json:"-"`

	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

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

func (o *Organization) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"organizations/create",
		o.Shortname, o.Name, o.Bio, o.Description, o.Email, o.Phone,
		o.City, o.Country, o.Address, o.Website,
		o.Mission, o.Culture, o.Logo, o.Cover,
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

func (o *Organization) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"organizations/update",
		o.ID, o.Shortname, o.Name, o.Bio, o.Description, o.Email, o.Phone,
		o.City, o.Country, o.Address, o.Website,
		o.Mission, o.Culture, o.Logo, o.Cover,
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

func (o *Organization) AddMember(ctx context.Context, userId uuid.UUID) error {
	rows, err := database.Query(
		ctx,
		"organizations/add_member",
		o.ID, userId,
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

func (o *Organization) RemoveMember(ctx context.Context, userId uuid.UUID) error {
	rows, err := database.Query(
		ctx,
		"organizations/remove_member",
		o.ID, userId,
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

func (o *Organization) Remove(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"organizations/remove_member",
		o.ID,
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

func GetAllOrganizations(userId uuid.UUID, p database.Paginate) ([]Organization, int, error) {
	var (
		organizations = []Organization{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if err := database.QuerySelect("organizations/get_all", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return organizations, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&organizations, ids...); err != nil {
		return nil, 0, err
	}
	return organizations, fetchList[0].TotalCount, nil
}

func GetOrganization(id uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func GetOrganizationByMember(id uuid.UUID, userId uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Get(o, "organizations/get_by_member", id, userId); err != nil {
		return nil, err
	}
	return o, nil
}

func getManyOrganizations(ids []uuid.UUID, identity uuid.UUID) ([]Organization, error) {
	result := []Organization{}
	return result, fmt.Errorf("Not Implemented")
}

func GetOrganizationByShortname(shortname string, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, identity.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func shortnameExistsOrganization(shortname string) (bool, error) {
	return false, fmt.Errorf("Not Implemented")
}

func searchOrganizations(query string) ([]Organization, error) { // Do we need to implement this?
	result := []Organization{}
	return result, fmt.Errorf("Not Implemented")
}
