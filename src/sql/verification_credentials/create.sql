INSERT INTO verification_credentials (user_id)
VALUES ($1)
RETURNING *