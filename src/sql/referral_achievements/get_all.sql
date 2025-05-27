SELECT id, COUNT(*) OVER () as total_count 
FROM referral_achievements
WHERE referrer_id=$1
LIMIT $2 OFFSET $3