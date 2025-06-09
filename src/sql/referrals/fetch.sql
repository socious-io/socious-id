WITH ra AS (
  SELECT *
  FROM referral_achievements
  WHERE referee_id IN (?)
)
SELECT
  COALESCE(jsonb_agg(
    json_build_object(
      'type', ra_sub.achievement_type, 
      'reward_claimed_at', ra_sub.reward_claimed_at
    )
  ), '[]') AS achievements,
  row_to_json(i.*) as referee
FROM ra AS ra_sub
JOIN identities i ON i.id = ra_sub.referee_id
GROUP BY i.id;