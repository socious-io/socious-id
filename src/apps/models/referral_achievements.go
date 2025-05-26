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
	ID              uuid.UUID      `db:"id" json:"id"`
	ReferrerID      uuid.UUID      `db:"referrer_id" json:"referrer_id"`
	RefereeID       string         `db:"referee_id" json:"referee_id"`
	AchievementType string         `db:"achievement_type" json:"achievement_type"`
	Meta            types.JSONText `db:"meta" json:"meta"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
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
		ra.ReferrerID, ra.RefereeID, ra.AchievementType, ra.Meta,
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

func GetReferralAchievements(ReferrerID uuid.UUID, p database.Paginate) ([]ReferralAchievement, int, error) {
	var (
		referralAchievements = []ReferralAchievement{}
		fetchList            []database.FetchList
		ids                  []interface{}
	)

	if err := database.QuerySelect("referral_achievements/get_all", &fetchList, ReferrerID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return referralAchievements, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}
	if err := database.Fetch(&referralAchievements, ids...); err != nil {
		return nil, 0, err
	}
	return referralAchievements, fetchList[0].TotalCount, nil
}
