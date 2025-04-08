SELECT 
  cv.*,
  row_to_json(u.*) AS user
FROM verification_credentials cv 
LEFT JOIN users u ON u.id = cv.user_id
WHERE cv.user_id = $1