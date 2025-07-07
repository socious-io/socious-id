package models

import (
	"context"
	"encoding/json"
	"socious-id/src/apps/utils"
	"socious-id/src/config"
	"strings"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type ReferralAchievement struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	ReferrerID      *uuid.UUID     `db:"referrer_id" json:"referrer_id"`
	RefereeID       uuid.UUID      `db:"referee_id" json:"referee_id"`
	AchievementType string         `db:"achievement_type" json:"achievement_type"`
	RewardAmount    float32        `db:"reward_amount" json:"reward_amount"`
	RewardClaimedAt *time.Time     `db:"reward_claimed_at" json:"reward_claimed_at"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	Meta            map[string]any `db:"-" json:"meta"`

	MetaJson *types.JSONText `db:"meta" json:"-"`
}

func (ReferralAchievement) TableName() string {
	return "referral_achievements"
}

func (ReferralAchievement) FetchQuery() string {
	return "referrals/fetch_achievement"
}

func (ra *ReferralAchievement) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(ra)
}

type ReferralStats struct {
	TotalCount                 int     `db:"total_count" json:"total_count"`
	TotalRewardAmount          float32 `db:"total_reward_amount" json:"total_reward_amount"`
	TotalUnclaimedRewardAmount float32 `db:"total_unclaimed_reward_amount" json:"total_unclaimed_reward_amount"`

	TotalPerAchievementType []struct {
		AchievementType string `db:"achievement_type" json:"achievement_type"`
		TotalCount      int    `db:"total_count" json:"total_count"`
	} `db:"-" json:"total_per_achievement_type"`
	TotalPerAchievementTypeJson types.JSONText `db:"total_per_achievement_type" json:"-"`
}

type Referral struct {
	Referee      *Identity `db:"-" json:"referee"`
	Achievements []struct {
		Type            string     `db:"type" json:"type"`
		RewardClaimedAt *time.Time `db:"reward_claimed_at" json:"reward_claimed_at"`
	} `db:"-" json:"achievements"`

	RefereeJson      types.JSONText `db:"referee" json:"-"`
	AchievementsJson types.JSONText `db:"achievements" json:"-"`
}

func (Referral) TableName() string {
	return "-"
}

func (Referral) FetchQuery() string {
	return "referrals/fetch"
}

func (r *Referral) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(r)
}

func (ra *ReferralAchievement) Create(ctx context.Context) error {
	var err error

	referralAchievements := []*ReferralAchievement{}

	if strings.HasPrefix(ra.AchievementType, "REF_") {
		referrerIdentity, err := GetReferrerIdentity(ra.RefereeID)
		if err == nil {
			ra.ReferrerID = &referrerIdentity.ID
			referralAchievements = append(referralAchievements, ra)
		}
	}

	referralAchievements = append(referralAchievements, &ReferralAchievement{
		RefereeID:       ra.RefereeID,
		AchievementType: strings.TrimPrefix(ra.AchievementType, "REF_"),
		Meta:            ra.Meta,
	})

	for _, ra := range referralAchievements {
		ra.MetaJson, err = utils.MapToJSONText(ra.Meta)
		if err != nil {
			return err
		}
		ra.RewardAmount = getRewardAmountByType(ra.AchievementType)

		rows, err := database.Query(
			ctx,
			"referrals/create_achievement",
			ra.ReferrerID, ra.RefereeID, ra.AchievementType, ra.RewardAmount, ra.MetaJson,
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
	}

	return database.Fetch(ra, referralAchievements[0].ID)
}

func GetReferralAchievements(identityID uuid.UUID, p database.Paginate) ([]ReferralAchievement, int, error) {
	var (
		referralAchievements = []ReferralAchievement{}
		fetchList            []database.FetchList
		ids                  []interface{}
	)

	if err := database.QuerySelect("referrals/get_all_achievements", &fetchList, identityID, p.Limit, p.Offet); err != nil {
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

func GetReferralStats(ctx context.Context, identityID uuid.UUID) (*ReferralStats, error) {
	stats := new(ReferralStats)
	if err := database.Get(stats, "referrals/get_stats", identityID); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(stats.TotalPerAchievementTypeJson, &stats.TotalPerAchievementType); err != nil {
		return nil, err
	}

	return stats, nil
}

func GetReferrals(identityID uuid.UUID, p database.Paginate) ([]Referral, int, error) {
	var (
		referrals = []Referral{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("referrals/get_all", &fetchList, identityID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return referrals, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}
	if err := database.Fetch(&referrals, ids...); err != nil {
		return nil, 0, err
	}
	return referrals, fetchList[0].TotalCount, nil
}

func getRewardAmountByType(t string) float32 {
	rewards := config.Config.ReferralAchievements.Rewards

	for _, reward := range rewards {
		if reward.Type == t {
			return reward.Amount
		}
	}

	return 0
}
