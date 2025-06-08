SELECT id, COUNT(*) OVER () as total_count 
FROM (
	(SELECT id FROM users WHERE referred_by = $1)
		UNION
	(SELECT id FROM organizations WHERE referred_by = $1)
)
LIMIT $2 OFFSET $3