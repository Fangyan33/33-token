CREATE TABLE upstream_credential_ref (
    id UUID PRIMARY KEY,
    credential_key TEXT NOT NULL
);

CREATE TABLE model_route (
    id UUID PRIMARY KEY,
    upstream_credential_ref_id UUID NOT NULL REFERENCES upstream_credential_ref(id)
);

CREATE TABLE admin_audit_log (
    id UUID PRIMARY KEY,
    operator_id UUID NOT NULL,
    action_type TEXT NOT NULL
);
