-- Add Offer and Mission to contract
ALTER TABLE contracts
    ADD COLUMN offer_id UUID,
    ADD COLUMN mission_id UUID,
    ADD CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE SET NULL;

ALTER TYPE contract_status ADD VALUE 'APPLIED';
ALTER TYPE contract_status ADD VALUE 'COMPLETED';