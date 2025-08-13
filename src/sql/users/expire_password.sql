UPDATE users
SET
    password_expired=TRUE
WHERE id=$1