CREATE TYPE contract_status AS ENUM (
  'CREATED',
  'CLIENT_APPROVED',
  'SIGNED',
  'PROVIDER_CANCELED',
  'CLIENT_CANCELED'
);

CREATE TYPE contract_type AS ENUM (
  'VOLUNTEER',
  'PAID'
);

CREATE TYPE contract_commitment_period AS ENUM (
  'HOURLY',
  'DAILY',
  'WEEKLY',
  'MONTHLY'
);

CREATE TABLE contracts (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT,
  status contract_status DEFAULT 'CREATED',
  type contract_type NOT NULL,
  total_amount FLOAT DEFAULT 0,
  currency_rate FLOAT DEFAULT 1,
  commitment_period contract_commitment_period NOT NULL,
  commitment integer DEFAULT 1,
  commitment_period_count integer DEFAULT 1,
  payment_type payment_mode_type DEFAULT 'FIAT',
  currency payment_currency DEFAULT 'USD',
  crypto_currency TEXT,
  provider_id UUID NOT NULL,
  client_id UUID NOT NULL,
  project_id UUID,
  applicant_id UUID,
  payment_id UUID,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_identity_provider FOREIGN KEY (provider_id) REFERENCES identities(id) ON DELETE CASCADE,
  CONSTRAINT fk_identity_client FOREIGN KEY (client_id) REFERENCES identities(id) ON DELETE CASCADE,
  CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
  CONSTRAINT fk_applicant FOREIGN KEY (applicant_id) REFERENCES applicants(id) ON DELETE SET NULL
);
