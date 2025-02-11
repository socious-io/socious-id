CREATE TYPE otp_type AS ENUM (
    'FORGET_PASSWORD',
    'SSO',
    'VERIFICATION'
);

CREATE TABLE accesses (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  description TEXT,
  client_id TEXT NOT NULL UNIQUE,
  client_secret TEXT NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auth_sessions (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  access_id UUID NOT NULL,
  redirect_url TEXT NOT NULL,
  expire_at TIMESTAMP DEFAULT (NOW() + '00:20:00'::interval) NOT NULL,
  verified_at TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_access FOREIGN KEY (access_id) REFERENCES accesses(id) ON DELETE CASCADE
);

CREATE TABLE otps (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  status otp_type DEFAULT 'SSO',
  user_id UUID NOT NULL,
  auth_session_id UUID NOT NULL,
  code VARCHAR(32) NOT NULL,
  expire_at TIMESTAMP DEFAULT (NOW() + '00:10:00'::interval) NOT NULL,
  verified_at TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_auth_session FOREIGN KEY (auth_session_id) REFERENCES auth_sessions(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_sso_refid_code  ON otps (user_id, code);
CREATE UNIQUE INDEX unique_code_per_auth_session ON otps(auth_session_id, code) WHERE auth_session_id IS NOT NULL;
