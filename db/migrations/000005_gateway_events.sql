CREATE TABLE usage_event (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES account(id)
);

CREATE TABLE billing_event (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES account(id),
    usage_event_id UUID NOT NULL REFERENCES usage_event(id)
);
