CREATE TYPE access_policy_type AS ENUM (
  'REQUIRE_ATLEAST_ONE_ORG',
  'PREVENT_USER_ACCOUNT_SELECTION'
);

ALTER TABLE accesses
ADD COLUMN policies access_policy_type[] NOT NULL DEFAULT '{}';