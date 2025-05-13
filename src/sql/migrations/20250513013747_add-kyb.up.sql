CREATE TYPE kyb_verification_status_type AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE kyb_verifications (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  organization_id UUID UNIQUE NOT NULL,
  status kyb_verification_status_type DEFAULT 'PENDING' NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
  CONSTRAINT kyb_verifications_user_org_unique UNIQUE (user_id, organization_id)
);

CREATE TABLE kyb_verification_documents (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  verification_id UUID NOT NULL,
  document UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_verification FOREIGN KEY (verification_id) REFERENCES kyb_verifications(id) ON DELETE CASCADE,
  CONSTRAINT fk_media FOREIGN KEY (document) REFERENCES media(id) ON DELETE CASCADE
);

ALTER TYPE status_type ADD VALUE 'PENDING';