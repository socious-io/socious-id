INSERT INTO credentials (user_id, type)
VALUES ($1, $2)
RETURNING *