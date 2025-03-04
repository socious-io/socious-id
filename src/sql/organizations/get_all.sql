SELECT id, COUNT(*) OVER () as total_count 
FROM organizations 
LIMIT $2 OFFSET $3