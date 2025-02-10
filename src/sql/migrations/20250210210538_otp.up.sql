CREATE TYPE otp_type AS ENUM (
    'FORGET_PASSWORD',
    'SSO',
    'VERIFICATION'
);

CREATE TABLE auth_sessions (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
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
  CONSTRAINT fk_auth_session FOREIGN KEY (auth_session_id) REFERENCES auth_sessions(id) ON DELETE CASCADE,
);

CREATE UNIQUE INDEX unique_sso_refid_code  ON otps (user_id, code);