ALTER TABLE users
ADD COLUMN referred_by uuid;

ALTER TABLE users
ADD CONSTRAINT fk_referral_user FOREIGN KEY (referred_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE organizations
ADD COLUMN referred_by uuid;

ALTER TABLE organizations
ADD CONSTRAINT fk_referrer_user FOREIGN KEY (referred_by) REFERENCES users(id) ON DELETE SET NULL;