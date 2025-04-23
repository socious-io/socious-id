INSERT INTO users (first_name, last_name, username, email, password, cover_id, avatar_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (email) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    username = EXCLUDED.username,
    password = EXCLUDED.password,
    cover_id = EXCLUDED.cover_id,
    avatar_id = EXCLUDED.avatar_id
WHERE users.status = 'INACTIVE' AND users.password IS NULL
RETURNING *;