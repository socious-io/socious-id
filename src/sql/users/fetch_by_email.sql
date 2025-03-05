SELECT u.*, row_to_json(m1.*) AS avatar, row_to_json(m2.*) AS cover
FROM users u
LEFT JOIN media m1 ON m1.id=u.avatar_id
LEFT JOIN media m2 ON m2.id=u.cover_id
WHERE email=$1