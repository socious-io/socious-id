UPDATE auth_sessions
SET policies = $2
WHERE id=$1