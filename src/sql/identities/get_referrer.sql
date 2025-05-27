SELECT *
FROM identities
WHERE id = COALESCE(
  (SELECT referred_by FROM users WHERE id=$1 LIMIT 1),
  (SELECT referred_by FROM organizations WHERE id=$1 LIMIT 1)
);