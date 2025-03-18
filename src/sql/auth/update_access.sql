UPDATE accesses SET
    destination_synced_at=$2,
    source_synced_at=$3,
    updated_at=NOW()
WHERE id=$1
RETURNING *