INSERT INTO otps (type, user_id, auth_session_id, code) 
VALUES ($1, $2, $3, $4)
RETURNING *;