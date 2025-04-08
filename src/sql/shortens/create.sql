INSERT INTO urls_shortens (long_url) VALUES ($1)
ON CONFLICT (long_url) DO NOTHING
RETURNING *