package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type ImpactPoint struct {
	ID uuid.UUID `db:"id" json:"id"`

	TotalPoints int `db:"total_points" json:"total_points"`

	UserID   uuid.UUID      `db:"user_id" json:"user_id"`
	UserJson types.JSONText `db:"user" json:"-"`
	User     *User          `db:"-" json:"user"`

	SocialCause         string `db:"social_cause" json:"social_cause"`
	SocialCauseCategory string `db:"social_cause_category" json:"social_cause_category"`

	Type ImpactPointType `db:"type" json:"type"`

	AccessID *uuid.UUID       `db:"access_id" json:"access_id"`
	Meta     *json.RawMessage `db:"meta" json:"meta"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Badges struct {
	TotalPoints         int    `db:"total_points" json:"total_points"`
	Count               int    `db:"count" json:"count"`
	SocialCauseCategory string `db:"social_cause_category" json:"social_cause_category"`
}

func (ImpactPoint) TableName() string {
	return "impact_points"
}

func (ImpactPoint) FetchQuery() string {
	return "impact_points/fetch"
}

func (ip *ImpactPoint) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"impact_points/create",
		ip.UserID,
		ip.TotalPoints,
		ip.SocialCause,
		ip.SocialCauseCategory,
		ip.Type,
		ip.AccessID,
		ip.Meta,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(ip); err != nil {
			return err
		}
	}
	return database.Fetch(ip, ip.ID)
}

func GetImpactPoints(userID uuid.UUID, p database.Paginate) ([]ImpactPoint, int, error) {
	var (
		impactPoints = []ImpactPoint{}
		fetchList    []database.FetchList
		ids          []interface{}
	)

	if err := database.QuerySelect("impact_points/get_all", &fetchList, userID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return impactPoints, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}
	if err := database.Fetch(&impactPoints, ids...); err != nil {
		return nil, 0, err
	}
	return impactPoints, fetchList[0].TotalCount, nil
}

func GetImpactBadges(userID uuid.UUID) ([]Badges, error) {
	var badges = []Badges{}
	if err := database.QuerySelect("impact_points/get_badges", &badges, userID); err != nil {
		return nil, err
	}
	return badges, nil
}
