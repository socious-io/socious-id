SELECT 
  o.*,
  row_to_json(a) AS auth_session
FROM otps
LEFT JOIN auth_sessions a ON a.id=o.auth_session_id
WHERE code=$1