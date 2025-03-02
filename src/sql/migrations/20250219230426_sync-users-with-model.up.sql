ALTER TYPE status_type RENAME TO status_type_old;
CREATE TYPE status_type AS ENUM (
    'ACTIVE',
    'INACTIVE',
    'SUSPENDED'
);

ALTER TABLE users
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE status_type USING status::text::status_type,
    ALTER COLUMN status SET DEFAULT 'INACTIVE',
    ALTER COLUMN username TYPE VARCHAR(200),
    ALTER COLUMN email TYPE VARCHAR(200),
    ALTER COLUMN first_name TYPE VARCHAR(70),
    ALTER COLUMN last_name TYPE VARCHAR(70),
    ADD COLUMN password_expired boolean DEFAULT false,
    ADD COLUMN email_text character varying(255),
    ADD COLUMN phone character varying(255),
    ADD COLUMN mission text,
    ADD COLUMN bio text,
    ADD COLUMN description_search text,
    ADD COLUMN city text,
    ADD COLUMN country character varying(3),
    ADD COLUMN address text,
    ADD COLUMN geoname_id integer,
    ADD COLUMN mobile_country_code character varying(16),
    ADD COLUMN avatar uuid,
    ADD COLUMN cover_image uuid,
    ADD COLUMN identity_verified_at date,
    ADD COLUMN email_verified_at timestamp without time zone,
    ADD COLUMN phone_verified_at timestamp without time zone,
    ADD COLUMN deleted_at timestamp without time zone;

ALTER TABLE otps
    RENAME COLUMN status TO type;