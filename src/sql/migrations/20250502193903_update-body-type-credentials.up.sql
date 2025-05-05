ALTER TABLE verification_credentials
ALTER COLUMN body TYPE jsonb USING body::jsonb;

-- Impact Point user integration
ALTER TABLE users
ADD COLUMN impact_points int DEFAULT 0 NOT NULL;

CREATE OR REPLACE FUNCTION update_user_impact_points() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  UPDATE users SET impact_points = impact_points + NEW.total_points WHERE id=NEW.user_id;
  RETURN NEW;
END;
$$;
CREATE OR REPLACE TRIGGER update_user_impact_points AFTER INSERT ON impact_points FOR EACH ROW EXECUTE FUNCTION update_user_impact_points();

UPDATE users u
SET impact_points = sub.total
FROM (
  SELECT user_id, SUM(total_points) AS total
  FROM impact_points
  GROUP BY user_id
) AS sub
WHERE sub.user_id = u.id;