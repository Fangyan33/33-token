CREATE TABLE account_subscription_state (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES account(id),
    status TEXT NOT NULL
);

CREATE TABLE account_cycle_summary (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES account(id),
    quota_remaining BIGINT NOT NULL
);
