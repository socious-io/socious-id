INSERT INTO accesses (name, description, client_id, client_secret)
VALUES ($1, $2, $3, $4)
RETURNING *