UPDATE users
SET
    password=$2,
    password_expired=FALSE
WHERE id=$1