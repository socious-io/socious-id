SELECT SUM(total_points) AS total_points, COUNT(*)::int, social_cause_category
FROM impact_points
WHERE user_id=$1
GROUP BY social_cause_category
