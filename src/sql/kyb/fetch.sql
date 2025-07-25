SELECT
	k.*,
	COALESCE(
		jsonb_agg (
			DISTINCT jsonb_build_object ('url', m.url, 'filename', m.filename)
		) FILTER ( WHERE m.id IS NOT NULL),
		'[]'
	) AS documents
FROM kyb_verifications k
LEFT JOIN kyb_verification_documents kd ON kd.verification_id = k.id
LEFT JOIN media m ON m.id = kd.document
WHERE k.id IN (?)
GROUP BY k.id;