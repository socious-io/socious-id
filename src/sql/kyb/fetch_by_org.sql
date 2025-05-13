SELECT k.*,
(SELECT
	jsonb_agg(
		json_build_object(
		  'url', m.url, 
		  'filename', m.filename
		)
	)
	FROM media m
	WHERE m.id = kd.document
) AS documents
FROM kyb_verifications k
LEFT JOIN kyb_verification_documents kd ON kd.verification_id=k.id
WHERE k.organization_id=$1