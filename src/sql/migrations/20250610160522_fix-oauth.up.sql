ALTER TABLE oauth_connects 
ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();