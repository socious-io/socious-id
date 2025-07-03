ALTER TABLE referral_achievements
ADD COLUMN reward_amount FLOAT DEFAULT 0 NOT NULL,
ADD COLUMN reward_claimed_at TIMESTAMP DEFAULT NULL;