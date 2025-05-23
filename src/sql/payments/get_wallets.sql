SELECT id, COUNT(*) OVER () as total_count 
FROM wallets
WHERE identity_id=$1;