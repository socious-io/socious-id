INSERT INTO kyb_verification_documents(verification_id, document)
VALUES($1, $2)
RETURNING *