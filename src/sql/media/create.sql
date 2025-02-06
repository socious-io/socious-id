INSERT INTO media(identity_id, url, filename)
VALUES($1, $2, $3)
RETURNING *