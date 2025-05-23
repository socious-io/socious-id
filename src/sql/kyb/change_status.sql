UPDATE kyb_verifications
SET status=$2
WHERE id=$1
RETURNING *