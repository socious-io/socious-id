INSERT INTO impact_points (user_id, total_points, social_cause, social_cause_category, type, access_id, meta)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

