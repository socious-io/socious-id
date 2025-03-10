SELECT id, COUNT(*) OVER () as total_count 
FROM organizations 
LIMIT $1 OFFSET $2