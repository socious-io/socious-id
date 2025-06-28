UPDATE credentials SET
  present_id=$2,
  status='REQUESTED',
  updated_at=NOW()
WHERE id=$1
RETURNING *