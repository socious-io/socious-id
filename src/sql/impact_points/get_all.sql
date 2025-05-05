SELECT id, COUNT(*) OVER () as total_count 
FROM impact_points
WHERE
    user_id=$1 AND
    (cardinality($2::impact_points_type[])=0 OR type=ANY($2))
ORDER BY created_at DESC
LIMIT $3 OFFSET $4