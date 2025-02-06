INSERT INTO users (first_name, last_name, username, email, password, status) 
VALUES ($1, $2, $3, $4, $5, 'ACTIVE')
RETURNING *;