INSERT INTO cards(identity_id, customer, card)
VALUES ($1, $2, $3)
ON CONFLICT(user_id, chain) UPDATE
SET
    address=EXCLUDED.address
    updated_at=NOW()
RETURNING *;