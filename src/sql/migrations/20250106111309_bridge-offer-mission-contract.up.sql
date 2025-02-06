-- Add triggers
CREATE OR REPLACE FUNCTION upsert_contract_on_offer() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    contract RECORD;
    project RECORD;
    contract_name TEXT;
    contract_status TEXT;
    v_commitment_period TEXT;
    v_commitment_period_count INTEGER;
BEGIN
    SELECT * INTO contract
    FROM contracts
    WHERE offer_id = NEW.id;

    SELECT * INTO project
    FROM projects
    WHERE id = NEW.project_id AND payment_type IS NOT NULL;

    contract_status := CASE
        WHEN NEW.status='PENDING' THEN 'CREATED'
        WHEN NEW.status='APPROVED' THEN 'CLIENT_APPROVED'
        WHEN NEW.status='HIRED' THEN 'SIGNED'
        WHEN NEW.status='WITHDRAWN' THEN 'CLIENT_CANCELED'
        WHEN NEW.status='CANCELED' THEN 'PROVIDER_CANCELED'
        ELSE NULL
    END;

    v_commitment_period := CASE
        WHEN NEW.total_hours IS NOT NULL THEN 'HOURLY'
        WHEN NEW.weekly_limit IS NOT NULL THEN 'WEEKLY'
        ELSE NULL
    END; --Fix: Fix calculation

    v_commitment_period_count := CASE
        WHEN NEW.total_hours IS NOT NULL THEN NEW.total_hours
        WHEN NEW.weekly_limit IS NOT NULL THEN NEW.weekly_limit
        ELSE NULL
    END; --Fix: Fix calculation

    IF project.id IS NULL OR contract_status IS NULL THEN
        RETURN NEW; -- Exit the function
    END IF;

    contract_name := CASE
        WHEN char_length(NEW.offer_message) > 128 THEN substring(NEW.offer_message FROM 1 FOR 125) || '...'
        ELSE NEW.offer_message
    END;
	
    IF contract.id IS NOT NULL THEN
        UPDATE contracts
        SET
            offer_id=NEW.id,
            name=contract_name,
            description=NEW.offer_message,
            status=COALESCE(contract_status::contract_status, status),
            type=project.payment_type::payment_type::text::contract_type,
            currency=NEW.currency,
            crypto_currency=NEW.crypto_currency_address,
            total_amount=COALESCE(NEW.assignment_total, 0),
            payment_type=NEW.payment_mode,
            project_id=NEW.project_id,
            client_id=NEW.recipient_id,
            provider_id=NEW.offerer_id,
            applicant_id=NEW.applicant_id,
            currency_rate=COALESCE(NEW.offer_rate, 1),
            commitment=COALESCE(v_commitment_period_count, 1),--FIX
            commitment_period=COALESCE(v_commitment_period::text::contract_commitment_period, 'HOURLY'),
            commitment_period_count=COALESCE(v_commitment_period_count, 1)
        WHERE id = contract.id;
    ELSE
        INSERT INTO
        contracts
        (
            offer_id,
            name,
            description,
            status,
            type,
            currency,
            crypto_currency,
            total_amount,
            payment_type,
            project_id,
            client_id,
            provider_id,
            applicant_id,
            currency_rate,
            commitment,
            commitment_period,
            commitment_period_count
        )
        VALUES (
            NEW.id,
            contract_name,
            NEW.offer_message,
            contract_status::contract_status,
            project.payment_type::payment_type::text::contract_type,
            NEW.currency,
            NEW.crypto_currency_address,
            COALESCE(NEW.assignment_total, 0),
            NEW.payment_mode,
            NEW.project_id,
            NEW.recipient_id,
            NEW.offerer_id,
            NEW.applicant_id,
            COALESCE(NEW.offer_rate, 1),
            COALESCE(v_commitment_period_count, 1), --FIX
            COALESCE(v_commitment_period::text::contract_commitment_period, 'HOURLY'),
            COALESCE(v_commitment_period_count, 1)
        );
    END IF;

    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION upsert_contract_on_mission() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    contract RECORD;
    offer RECORD;
    project RECORD;
    contract_name TEXT;
    contract_status TEXT;
    v_commitment_period TEXT;
    v_commitment_period_count INTEGER;
BEGIN
    SELECT * INTO contract
    FROM contracts
    WHERE offer_id = NEW.offer_id;

    SELECT * INTO project
    FROM projects
    WHERE id = NEW.project_id AND payment_type IS NOT NULL;

    SELECT * INTO offer
    FROM offers
    WHERE id = NEW.offer_id;

    contract_status := CASE
        WHEN NEW.status='ACTIVE' THEN 'SIGNED'
        WHEN NEW.status='COMPLETE' THEN 'APPLIED'
        WHEN NEW.status='CONFIRMED' THEN 'COMPLETED'
        WHEN NEW.status='CANCELED' THEN 'CLIENT_CANCELED'
        WHEN NEW.status='KICKED_OUT' THEN 'PROVIDER_CANCELED'
        ELSE NULL
    END;

    
    IF project.id IS NULL OR contract_status IS NULL OR offer.id IS NULL THEN
        RETURN NEW; -- Exit the function
    END IF;

    v_commitment_period := CASE
        WHEN offer.total_hours IS NOT NULL THEN 'HOURLY'
        WHEN offer.weekly_limit IS NOT NULL THEN 'WEEKLY'
        ELSE NULL
    END;

    v_commitment_period_count := CASE
        WHEN offer.total_hours IS NOT NULL THEN offer.total_hours
        WHEN offer.weekly_limit IS NOT NULL THEN offer.weekly_limit
        ELSE NULL
    END;

    contract_name := CASE
        WHEN char_length(offer.offer_message) > 128 THEN substring(offer.offer_message FROM 1 FOR 125) || '...'
        ELSE offer.offer_message
    END;

    IF contract.id IS NOT NULL THEN
        UPDATE contracts
        SET
            mission_id=NEW.id,
            status=contract_status::contract_status
        WHERE id = contract.id;
    ELSE
        INSERT INTO
        contracts
        (
            mission_id,
            offer_id,
            name,
            description,
            status,
            type,
            currency,
            crypto_currency,
            currency_rate,
            total_amount,
            payment_type,
            project_id,
            client_id,
            provider_id,
            applicant_id,
            commitment,
            commitment_period,
            commitment_period_count
        )
        VALUES (
            NEW.id,
            NEW.offer_id,
            contract_name,
            offer.offer_message,
            contract_status::contract_status,
            project.payment_type::payment_type::text::contract_type,
            offer.currency,
            offer.crypto_currency_address,
            COALESCE(offer.offer_rate, 1),
            COALESCE(offer.assignment_total, 0),
            offer.payment_mode,
            NEW.project_id,
            NEW.assignee_id,
            NEW.assigner_id,
            NEW.applicant_id,
            COALESCE(v_commitment_period_count, 1), --FIX
            COALESCE(v_commitment_period::text::contract_commitment_period, 'HOURLY'),
            COALESCE(v_commitment_period_count, 1)
        );
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER upsert_contract_on_offer AFTER INSERT OR UPDATE ON offers FOR EACH ROW EXECUTE FUNCTION upsert_contract_on_offer();
CREATE OR REPLACE TRIGGER upsert_contract_on_mission AFTER INSERT OR UPDATE ON missions FOR EACH ROW EXECUTE FUNCTION upsert_contract_on_mission();

-- Migrating the Offers and Missions
UPDATE offers SET id=id;
UPDATE missions SET id=id;