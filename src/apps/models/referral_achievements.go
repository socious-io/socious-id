package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type ReferralAchievement struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	IdentityID uuid.UUID      `db:"identity_id" json:"identity_id"`
	Type       string         `db:"type" json:"type"`
	Meta       types.JSONText `db:"meta" json:"meta"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
}

func (ReferralAchievement) TableName() string {
	return "referral_achievements"
}

func (ReferralAchievement) FetchQuery() string {
	return "referral_achievements/fetch"
}

func (ra *ReferralAchievement) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(ra)
}

func (ra *ReferralAchievement) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"referral_achievements/create",
		ra.IdentityID, ra.Type, ra.Meta,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := ra.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func GetReferralAchievements(identityID uuid.UUID, p database.Paginate) ([]ImpactPoint, int, error) {
	var (
		impactPoints = []ImpactPoint{}
		fetchList    []database.FetchList
		ids          []interface{}
	)

	if err := database.QuerySelect("referral_achievements/get_all", &fetchList, identityID, p.Limit, p.Offet); err != nil {
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
