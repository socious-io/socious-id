SELECT org.id, COUNT(*) OVER () as total_count 
FROM organizations org
LEFT JOIN org_members om ON om.org_id=org.id
WHERE om.user_id=$1
LIMIT $2 OFFSET $3