ALTER TABLE accesses 
    ADD COLUMN sync_url TEXT,
    ADD COLUMN destination_synced_at TIMESTAMP,
    ADD COLUMN source_synced_at TIMESTAMP;