CREATE TYPE sync_method_type AS ENUM (
  'MQ',
  'HTTP'
);

ALTER TABLE accesses
ADD COLUMN sync_method sync_method_type NOT NULL DEFAULT 'MQ'::sync_method_type;