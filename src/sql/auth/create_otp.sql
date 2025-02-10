INSERT INTO otps (type, ref_id, code) 
VALUES ($1, $2, $3)
RETURNING *;