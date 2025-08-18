INSERT INTO auth_sessions (redirect_url, access_id, policies)
VALUES ($1, $2, $3)
RETURNING *