UPDATE users
SET
    first_name=$2,
    last_name=$3,
    bio=$4,
    phone=$5,
    username=$6,
    cover_id=$7,
    avatar_id=$8
WHERE id=$1
RETURNING *