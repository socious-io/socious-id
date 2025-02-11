INSERT INTO auth_sessions (redirect_url, access_id)
VALUES ($1, $2)
RETURNING *