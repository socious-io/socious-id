UPDATE verification_credentials SET
  body=$2,
  status='VERIFIED',
  verified_at=NOW(),
  updated_at=NOW()
WHERE id=$1
RETURNING *