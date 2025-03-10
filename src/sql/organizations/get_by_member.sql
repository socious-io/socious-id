SELECT org.*
FROM organizations org
LEFT JOIN org_members om ON om.org_id=org.id
WHERE org.id=$1 AND om.user_id=$2