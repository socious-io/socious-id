UPDATE organizations
SET
    verified = (CASE WHEN $2 = 'NORMAL' THEN TRUE ELSE verified END),
    verified_impact = (CASE WHEN $2 = 'IMPACT' THEN TRUE ELSE verified_impact END)
WHERE id = $1;