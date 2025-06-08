WITH ip AS (
  SELECT *
  FROM impact_points
  WHERE user_id = $1
)
SELECT
  (SELECT SUM(ip.total_points) FROM ip) AS total_points,
  (SELECT SUM(ip.value) FROM ip) AS total_values,
  (
    SELECT json_agg(row_to_json(ipt))
    FROM (
      SELECT
        type,
        SUM(total_points) AS total_points,
        SUM(value) AS total_values
      FROM ip
      GROUP BY type
    ) AS ipt
  ) AS total_per_type;
