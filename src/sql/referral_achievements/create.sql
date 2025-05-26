INSERT INTO referral_achievements
    (referrer_id, referee_id, achievement_type, meta)
VALUES
    ($1, $2, $3, $4)
RETURNING *;