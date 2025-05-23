SELECT *
FROM identities
WHERE id = COALESCE(
  (SELECT id FROM users WHERE username=$1 LIMIT 1),
  (SELECT id FROM organizations WHERE shortname=$1 LIMIT 1)
);