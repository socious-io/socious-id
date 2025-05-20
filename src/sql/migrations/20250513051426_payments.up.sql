CREATE TABLE wallets (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
  chain text NOT NULL, --TODO: Define enum
  chain_id text,
  address text NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  
  UNIQUE(identity_id, chain)
);

CREATE TABLE cards (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
  customer text,
  card text,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);