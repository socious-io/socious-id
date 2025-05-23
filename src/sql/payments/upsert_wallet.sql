INSERT INTO wallets(identity_id, chain, chain_id, address)
VALUES ($1, $2, $3, $4)
ON CONFLICT(identity_id, chain) DO UPDATE
SET
    chain_id=EXCLUDED.chain_id,
    address=EXCLUDED.address,
    updated_at=NOW()
RETURNING *;