SELECT id, COUNT(*) OVER () as total_count 
FROM referral_achievements 
LIMIT $1 OFFSET $2