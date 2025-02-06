-- Add Offer and Mission to contract
ALTER TABLE contracts
    ADD COLUMN requirement_description text;
    
CREATE TABLE contract_requirements_files (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  contract_id UUID NOT NULL,
  document UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_contract FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE,
  CONSTRAINT fk_media FOREIGN KEY (document) REFERENCES media(id) ON DELETE CASCADE
);