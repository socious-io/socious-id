SELECT u.*, row_to_json(m1.*) AS avatar, row_to_json(m2.*)
FROM users u
LEFT JOIN media m1 WHERE m1.id=u.avatar
LEFT JOIN media m1 WHERE m1.id=u.avatar
WHERE id IN(?)