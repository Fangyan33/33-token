CREATE TABLE plan (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL UNIQUE
);

CREATE TABLE plan_price_snapshot (
    id UUID PRIMARY KEY,
    plan_id UUID NOT NULL REFERENCES plan(id)
);

CREATE TABLE "order" (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES account(id),
    plan_price_snapshot_id UUID NOT NULL REFERENCES plan_price_snapshot(id)
);

CREATE TABLE payment_event (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES "order"(id)
);
