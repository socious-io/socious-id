-- Generalizing Project Length
ALTER TABLE projects
    DROP COLUMN service_length,
    DROP COLUMN service_total_hours,
    DROP COLUMN service_price;
DROP TYPE service_length;
ALTER TYPE project_length ADD VALUE '1_3_DAYS';
ALTER TYPE project_length ADD VALUE '1_WEEK';
ALTER TYPE project_length ADD VALUE '2_WEEKS';
ALTER TYPE project_length ADD VALUE '1_MONTH';
