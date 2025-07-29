DELETE FROM kyb_verification_documents kd
WHERE
    kd.verification_id = $1;