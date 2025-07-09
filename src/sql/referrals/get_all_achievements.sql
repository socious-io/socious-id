SELECT id, COUNT(*) OVER () as total_count 
FROM referral_achievements
WHERE (referrer_id = $1 OR (referee_id = $1 AND referrer_id IS NULL))
LIMIT $2 OFFSET $3