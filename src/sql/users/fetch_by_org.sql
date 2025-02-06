SELECT u.*
FROM users u
JOIN organizations o ON u.id=o.created_by
WHERE o.id=$1