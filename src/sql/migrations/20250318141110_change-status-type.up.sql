ALTER TYPE status_type RENAME TO status_type_old2;
CREATE TYPE status_type AS ENUM (
    'ACTIVE',
    'NOT_ACTIVE',
    'SUSPENDED'
);

UPDATE organizations 
SET status = 'NOT_ACTIVE' 
WHERE status = 'INACTIVE';

ALTER TABLE organizations
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE status_type USING status::text::status_type,
    ALTER COLUMN status SET DEFAULT 'NOT_ACTIVE';