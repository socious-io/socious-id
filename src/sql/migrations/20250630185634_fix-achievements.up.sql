ALTER TABLE referral_achievements
ALTER COLUMN referrer_id DROP NOT NULL;

DROP INDEX IF EXISTS idx_unique_achievements;
CREATE UNIQUE INDEX idx_unique_achievements
ON referral_achievements (
  COALESCE(referrer_id, '00000000-0000-0000-0000-000000000000'::uuid),
  referee_id,
  achievement_type
);