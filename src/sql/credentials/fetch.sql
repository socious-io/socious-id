SELECT 
  v.*,
  row_to_json(u.*) AS user
FROM credentials v
LEFT JOIN users u ON u.id = v.user_id
WHERE v.id IN (?)