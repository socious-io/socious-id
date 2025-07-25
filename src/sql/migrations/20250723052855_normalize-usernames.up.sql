WITH cleaned_users AS (
  SELECT 
    id,
    username,
    REGEXP_REPLACE(LOWER(username), '[^a-z0-9._-]', '', 'g') AS cleaned
  FROM users
),
initial_fixes AS (
  SELECT
    id,
    LEFT(
      CASE
        WHEN LENGTH(cleaned) < 6 THEN
          cleaned || SUBSTR(FLOOR(RANDOM() * 1000000)::TEXT, 1, 6 - LENGTH(cleaned))
        ELSE cleaned
      END,
      24
    ) AS fixed_username
  FROM cleaned_users
),
duplicates AS (
  SELECT fixed_username
  FROM initial_fixes
  GROUP BY fixed_username
  HAVING COUNT(*) > 1
),
resolved_fixes AS (
  SELECT 
    f.id,
    CASE 
      WHEN d.fixed_username IS NOT NULL THEN 
        LEFT(f.fixed_username || SUBSTR(FLOOR(RANDOM() * 100)::TEXT, 1, 2), 24)
      ELSE f.fixed_username
    END AS final_username
  FROM initial_fixes f
  LEFT JOIN duplicates d ON f.fixed_username = d.fixed_username
)
UPDATE users u
SET username = r.final_username
FROM resolved_fixes r
WHERE u.id = r.id;