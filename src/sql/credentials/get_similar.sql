SELECT * FROM credentials
WHERE
(
    (
        body->>'document_number' = $2 AND
        body->>'country' = $3
    ) OR
    (
        body->>'first_name' = $4 AND 
        body->>'last_name' = $5 AND 
        body->>'date_of_birth' = $6
    )
) AND id!=$1 AND status='VERIFIED';