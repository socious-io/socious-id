SELECT
    SUM(total_points) AS total_points,
    COUNT(*)::int,
    social_cause_category,
    BOOL_AND(claimed_at IS NOT NULL) AS is_claimed
FROM impact_points
WHERE user_id=$1
GROUP BY social_cause_category
