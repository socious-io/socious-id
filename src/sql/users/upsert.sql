INSERT INTO users (
    first_name,
    last_name,
    username,
    email,
    password,
    password_expired
)
VALUES ($1, $2, $3, $4, NULL, TRUE)
ON CONFLICT (email) DO UPDATE
SET email = users.email
RETURNING *;