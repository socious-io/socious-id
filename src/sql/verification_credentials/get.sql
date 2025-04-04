SELECT 
  id, COUNT(*) OVER () as total_count
FROM verification_credentials cv
WHERE cv.user_id = $1 LIMIT $2 OFFSET $3