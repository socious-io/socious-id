SELECT id, COUNT(*) OVER () as total_count 
FROM cards
WHERE identity_id=$1
LIMIT $2 OFFSET $3;