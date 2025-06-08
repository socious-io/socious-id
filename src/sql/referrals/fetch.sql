WITH ra AS (
  SELECT *
  FROM referral_achievements
  WHERE referee_id IN (?)
)
SELECT
	COALESCE(jsonb_agg(
		json_build_object(
		  'type', ra.achievement_type, 
		  'reward_claimed_at', ra.reward_claimed_at
		)
	), '[]') AS achievements,
	row_to_json(i.*) as referee
FROM ra
JOIN identities i ON i.id=ra.referee_id
GROUP BY i.id, ra.referee_id;