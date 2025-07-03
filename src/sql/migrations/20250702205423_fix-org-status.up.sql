UPDATE organizations
SET status=(CASE WHEN verified=TRUE OR verified_impact=TRUE THEN 'ACTIVE'::organization_status_type ELSE 'NOT_ACTIVE'::organization_status_type END)
WHERE status!='SUSPENDED';