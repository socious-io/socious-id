UPDATE impact_points
SET claimed_at=NOW()
WHERE user_id=$1 AND claimed_at=NULL