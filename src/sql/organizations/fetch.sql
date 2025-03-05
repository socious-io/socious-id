SELECT org.*
-- row_to_json(m_image.*) AS image,
-- row_to_json(m_cover.*) AS cover_image,
FROM organizations org
-- LEFT JOIN media m_image ON m_image.id=org.image
-- LEFT JOIN media m_cover ON m_cover.id=org.cover_image
WHERE org.id IN (?)