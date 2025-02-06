-- ALTER TABLE contracts
-- ADD CONSTRAINT unique_client_provider_status_project UNIQUE (client_id, provider_id, status, project_id);

CREATE UNIQUE INDEX unique_client_provider_status_project
ON contracts (client_id, provider_id, project_id)
WHERE status = 'CREATED';