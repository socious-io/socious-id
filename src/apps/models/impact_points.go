package models

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	database "github.com/socious-io/pkg_database"
)

type ImpactPoint struct {
	ID uuid.UUID `db:"id" json:"id"`

	TotalPoints int     `db:"total_points" json:"total_points"`
	Value       float64 `db:"value" json:"value"`

	UserID   uuid.UUID      `db:"user_id" json:"user_id"`
	UserJson types.JSONText `db:"user" json:"-"`
	User     *User          `db:"-" json:"user"`

	SocialCause         string `db:"social_cause" json:"social_cause"`
	SocialCauseCategory string `db:"social_cause_category" json:"social_cause_category"`

	Type ImpactPointType `db:"type" json:"type"`

	AccessID *uuid.UUID       `db:"access_id" json:"access_id"`
	Meta     *json.RawMessage `db:"meta" json:"meta"`

	UniqueTag string `db:"unique_tag" json:"unique_tag"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	ClaimedAt *time.Time `db:"claimed_at" json:"claimed_at"`
}

type Badges struct {
	TotalPoints         int    `db:"total_points" json:"total_points"`
	Count               int    `db:"count" json:"count"`
	SocialCauseCategory string `db:"social_cause_category" json:"social_cause_category"`
	IsClaimed           bool   `db:"is_claimed" json:"is_claimed"`
}

type ImpactPointStats struct {
	TotalPoints int     `db:"total_points" json:"total_points"`
	TotalValues float64 `db:"total_values" json:"total_values"`

	TotalPerType []struct {
		Type        ImpactPointType `db:"type" json:"type"`
		TotalPoints int             `db:"total_points" json:"total_points"`
		TotalValues float64         `db:"total_values" json:"total_values"`
	} `db:"-" json:"total_per_type"`
	TotalPerTypeJson types.JSONText `db:"total_per_type" json:"-"`
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
		ip.UniqueTag,
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
		typeFilter   []string = []string{}
	)

	if len(p.Filters) > 0 {
		for _, filter := range p.Filters {
			if filter.Key == "type" {
				typeFilter = strings.Split(filter.Value, ",")
			}
		}
	}

	if err := database.QuerySelect("impact_points/get_all", &fetchList, userID, pq.Array(typeFilter), p.Limit, p.Offet); err != nil {
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

func GetImpactPointStats(userID uuid.UUID) (*ImpactPointStats, error) {
	stats := new(ImpactPointStats)
	if err := database.Get(stats, "impact_points/get_stats", userID); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(stats.TotalPerTypeJson, &stats.TotalPerType); err != nil {
		return nil, err
	}
	return stats, nil
}

func GetImpactBadges(userID uuid.UUID) ([]Badges, error) {
	var badges = []Badges{}
	if err := database.QuerySelect("impact_points/get_badges", &badges, userID); err != nil {
		return nil, err
	}
	return badges, nil
}

func ClaimAllImpactPoints(ctx context.Context, userID uuid.UUID) error {
	if _, err := database.Query(ctx, "impact_points/claim_all", userID); err != nil {
		return err
	}
	return nil
}
