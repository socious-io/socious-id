UPDATE referral_achievements ra
SET reward_claimed_at = NOW()
WHERE (referrer_id=$1 OR (referee_id=$1 AND referrer_id IS NULL))

