UPDATE organizations SET
    shortname=$2,
    name=$3,
    bio=$4,
    description=$5,
    email=$6,
    phone=$7,
    city=$8,
    country=$9,
    address=$10,
    website=$11,
    mission=$12,
    culture=$13
WHERE id=$1
RETURNING *