INSERT INTO kyb_verifications(user_id, organization_id)
VALUES($1, $2)
ON CONFLICT (user_id, organization_id) DO UPDATE SET
    status = 'PENDING'
RETURNING *