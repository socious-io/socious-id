package models

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Shortname   string    `db:"shortname" json:"shortname"`
	Name        *string   `db:"name" json:"name"`
	Bio         *string   `db:"bio" json:"bio"`
	Description *string   `db:"description" json:"description"`
	Email       *string   `db:"email" json:"email"`
	Phone       *string   `db:"phone" json:"phone"`

	City    *string `db:"city" json:"city"`
	Country *string `db:"country" json:"country"`
	Address *string `db:"address" json:"address"`
	Website *string `db:"website" json:"website"`

	Mission *string `db:"mission" json:"mission"`
	Culture *string `db:"culture" json:"culture"`

	LogoID   *uuid.UUID     `db:"logo_id" json:"logo_id"`
	Logo     *Media         `db:"-" json:"logo"`
	LogoJson types.JSONText `db:"logo" json:"-"`

	CoverID   *uuid.UUID     `db:"cover_id" json:"cover_id"`
	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	Status OrganizationStatusType `db:"status" json:"status"`

	VerifiedImpact bool `db:"verified_impact" json:"verified_impact"`
	Verified       bool `db:"verified" json:"verified"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type OrganizationMember struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OrganizationID uuid.UUID `db:"org_id" json:"org_id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
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
		o.Mission, o.Culture, o.CoverID, o.LogoID,
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
	return database.Fetch(o, o.ID)
}

func (o *Organization) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"organizations/update",
		o.ID, o.Shortname, o.Name, o.Bio, o.Description, o.Email, o.Phone,
		o.City, o.Country, o.Address, o.Website,
		o.Mission, o.Culture, o.CoverID, o.LogoID, o.Status, o.Verified, o.VerifiedImpact,
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
	return database.Fetch(o, o.ID)
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
		"organizations/remove",
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

func GetAllOrganizations(p database.Paginate) ([]Organization, int, error) {
	var (
		organizations = []Organization{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if err := database.QuerySelect("organizations/get_all", &fetchList, p.Limit, p.Offet); err != nil {
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

func GetOrganizationMembers(organizationId uuid.UUID) ([]OrganizationMember, error) {
	var members = []OrganizationMember{}
	if err := database.QuerySelect("organizations/get_members", &members, organizationId); err != nil {
		return nil, err
	}

	return members, nil
}

func GetOrganizationsByMember(userId uuid.UUID) ([]Organization, error) {
	var (
		organizations = []Organization{}
		ids           []interface{}
	)

	if err := database.QuerySelect("organizations/get_all_by_member", &ids, userId); err != nil {
		return nil, err
	}

	if len(ids) < 1 {
		return organizations, nil
	}

	if err := database.Fetch(&organizations, ids...); err != nil {
		return nil, err
	}

	return organizations, nil
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

func GetOrganizationByShortname(shortname string, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, identity.String()); err != nil {
		return nil, err
	}
	return o, nil
}
