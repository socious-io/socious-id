INSERT INTO users (first_name, last_name, username, email, password, cover_id, avatar_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;