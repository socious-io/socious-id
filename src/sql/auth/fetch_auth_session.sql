SELECT
  sess.*,
  row_to_json(a.*) AS access
FROM auth_sessions sess
JOIN accesses a ON sess.access_id = a.id
WHERE id IN(?)