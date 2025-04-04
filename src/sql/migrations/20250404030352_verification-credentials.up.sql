CREATE TYPE verification_status_type AS ENUM (
    "CREATED",
	"REQUESTED",
	"VERIFIED",
	"FAILED"
);

CREATE TABLE verification_credentials (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    status verification_status_type NOT NULL DEFAULT "CREATED",
    user_id UUID NOT NULL,
    connection_id VARCHAR(255),
    connection_url TEXT,
    present_id VARCHAR(255),
    body TEXT,
    validation_error TEXT,
    connection_at TIMESTAMP,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
