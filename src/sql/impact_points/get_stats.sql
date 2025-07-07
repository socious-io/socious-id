WITH ip AS (
  SELECT *
  FROM impact_points
  WHERE user_id = $1
)
SELECT
  COALESCE((SELECT SUM(ip.total_points) FROM ip), 0) AS total_points,
  COALESCE((SELECT SUM(ip.value) FROM ip), 0) AS total_values,
  COALESCE((
    SELECT json_agg(row_to_json(ipt))
    FROM (
      SELECT
        type,
        SUM(total_points) AS total_points,
        SUM(value) AS total_values
      FROM ip
      GROUP BY type
    ) AS ipt
  ), '[]') AS total_per_type;