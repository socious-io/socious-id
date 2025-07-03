SELECT 
  c.*,
  row_to_json(u.*) AS user
FROM credentials c 
LEFT JOIN users u ON u.id = c.user_id
WHERE c.user_id = $1 AND c.type = $2
ORDER BY c.created_at DESC
LIMIT 1