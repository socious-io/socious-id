update verification_credentials SET
  user_id=$1
WHERE id=$1
RETURNING *