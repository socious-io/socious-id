INSERT INTO referral_achievements
    (referrer_id, referee_id, achievement_type, reward_amount, meta)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;