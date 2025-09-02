WITH RECURSIVE chain AS (
  SELECT id, referred_by, 0 AS depth
  FROM users
  WHERE id = $1

  UNION ALL

  SELECT u.id, u.referred_by, c.depth + 1
  FROM users u
  JOIN chain c ON c.referred_by = u.id
)
SELECT id
FROM chain
WHERE id != $1
ORDER BY depth DESC;