CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TYPE status_type AS ENUM (
    'ACTIVE',
    'NOT_ACTIVE',
    'SUSPENDED'
);

CREATE TABLE users (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  status status_type DEFAULT 'NOT_ACTIVE',
  username VARCHAR(32) UNIQUE NOT NULL,
  email VARCHAR(32) UNIQUE NOT NULL,
  first_name VARCHAR(32),
  last_name VARCHAR(32),
  password TEXT,
  updated_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW()
);