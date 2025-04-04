SELECT 
  id, COUNT(*) OVER () as total_count
FROM verification_credentials cv
WHERE cv.user_id = $1 AND cv.verification_id=$2 LIMIT $3 OFFSET $4