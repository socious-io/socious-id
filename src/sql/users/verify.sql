UPDATE users
SET
    status = 'ACTIVE',
    email_verified_at = (CASE WHEN $2 = 'EMAIL' THEN NOW() ELSE email_verified_at END),
    phone_verified_at = (CASE WHEN $2 = 'PHONE' THEN NOW() ELSE phone_verified_at END),
    identity_verified_at = (CASE WHEN $2 = 'IDENTITY' THEN NOW() ELSE identity_verified_at END)
WHERE id = $1;