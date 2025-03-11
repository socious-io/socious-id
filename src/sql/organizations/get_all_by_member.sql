SELECT org.*,
    row_to_json(m1.*) AS logo,
    row_to_json(m2.*) AS cover
FROM organizations org
LEFT JOIN org_members om ON om.org_id=org.id
LEFT JOIN media m1 ON m1.id=org.logo_id
LEFT JOIN media m2 ON m2.id=org.cover_id
WHERE om.user_id=$1