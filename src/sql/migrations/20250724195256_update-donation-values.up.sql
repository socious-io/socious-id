UPDATE impact_points i
    SET value=(i.total_points)::double precision
WHERE type='DONATION';

ALTER TYPE impact_points_type ADD VALUE 'VOTING';