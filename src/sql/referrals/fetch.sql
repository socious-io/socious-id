SELECT
	COALESCE(jsonb_agg(
		json_build_object(
			'type', ra.achievement_type, 
			'reward_claimed_at', ra.reward_claimed_at
		)
  	) FILTER (WHERE ra.id IS NOT NULL), '[]') AS achievements,
	row_to_json(i.*) as referee
FROM identities i
LEFT JOIN users u ON u.id=i.id
LEFT JOIN referral_achievements ra ON (ra.referrer_id = u.referred_by AND ra.referee_id=i.id)
WHERE i.id IN (?)
GROUP BY i.id;