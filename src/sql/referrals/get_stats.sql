WITH ra AS (
  SELECT *
  FROM referral_achievements
  WHERE referrer_id = $1
),
ri AS (
    (SELECT id FROM users WHERE referred_by = $1)
    UNION
    (SELECT id FROM organizations WHERE referred_by = $1)
)

SELECT
  (SELECT COUNT(*) FROM ri) AS total_count,
  COALESCE(
      (
        SELECT json_agg(row_to_json(raa))
        FROM (
          SELECT achievement_type, COUNT(*) AS total_count
          FROM ra
          GROUP BY achievement_type
        ) AS raa
      )
  , '[]') AS total_per_achievement_type,
  COALESCE((SELECT SUM(reward_amount) FROM ra), 0) AS total_reward_amount,
  COALESCE((SELECT SUM(reward_amount) FROM ra WHERE reward_claimed_at IS NULL), 0) AS total_unclaimed_reward_amount