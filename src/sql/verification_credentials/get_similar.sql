SELECT * FROM verification_credentials
WHERE 
(
    body->>'document_number' = :document_number AND
    body->>'country' = :country
) OR
(
    body->>'first_name' = :first_name AND 
    body->>'last_name' = :last_name AND 
    body->>'date_of_birth' = :date_of_birth
);