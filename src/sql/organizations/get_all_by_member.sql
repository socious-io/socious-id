SELECT org.*
FROM organizations org
LEFT JOIN org_members om ON om.org_id=org.id
WHERE om.user_id=$1