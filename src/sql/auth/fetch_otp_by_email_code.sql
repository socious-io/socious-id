SELECT 
  o.*,
  row_to_json(a) AS auth_session,
  row_to_json(u) AS user
FROM otps o
LEFT JOIN auth_sessions a ON a.id=o.auth_session_id
JOIN users u ON u.id=o.user_id
WHERE u.email=$1 AND o.code=$2