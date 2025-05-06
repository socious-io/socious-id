SELECT
  type,
  SUM(total_points) AS total_points,
  SUM(value) AS total_values
FROM impact_points
WHERE user_id = $1
GROUP BY type