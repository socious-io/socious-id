SELECT ip.*,
    row_to_json(u.*) AS user
FROM impact_points ip
JOIN users u ON u.id=ip.user_id
WHERE ip.id IN (?)