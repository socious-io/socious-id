UPDATE users
SET
    status = 'ACTIVE',
    email_verified_at = CASE WHEN $2 = 'EMAIL' THEN true ELSE email_verified_at END,
    phone_verified_at = CASE WHEN $2 = 'PHONE' THEN true ELSE phone_verified_at END
    identity_verified_at = CASE WHEN $2 = 'IDENTITY' THEN true ELSE identity_verified_at END
WHERE id = $1;