INSERT INTO oauth_connects (identity_id, provider, matrix_unique_id, access_token, refresh_token, meta, expired_at, is_confirmed, redirect_url, status)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (matrix_unique_id, provider) 
DO UPDATE SET
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    meta = EXCLUDED.meta,
    expired_at = EXCLUDED.expired_at,
    is_confirmed = EXCLUDED.is_confirmed,
    redirect_url = EXCLUDED.redirect_url,
    status = EXCLUDED.status,
    updated_at = NOW();