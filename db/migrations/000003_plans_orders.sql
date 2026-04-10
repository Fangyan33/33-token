create table if not exists plan (
  id uuid primary key,
  code text not null unique,
  name text not null,
  status text not null,
  billing_period_type text not null,
  quota_total bigint not null,
  rate_limit_policy jsonb not null default '{}'::jsonb,
  display_priority integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists plan_price_snapshot (
  id uuid primary key,
  plan_id uuid not null references plan(id),
  plan_code text not null,
  plan_name text not null,
  price_amount bigint not null,
  currency text not null,
  billing_period_type text not null,
  quota_total bigint not null,
  rate_limit_policy_snapshot jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists "order" (
  id uuid primary key,
  account_id uuid not null references account(id),
  plan_price_snapshot_id uuid not null references plan_price_snapshot(id),
  order_type text not null,
  status text not null,
  payment_provider text not null,
  payment_provider_order_id text,
  amount bigint not null,
  currency text not null,
  created_at timestamptz not null default now(),
  paid_at timestamptz,
  completed_at timestamptz
);

create table if not exists payment_event (
  id uuid primary key,
  order_id uuid not null references "order"(id),
  payment_provider text not null,
  provider_event_id text not null,
  event_type text not null,
  event_status text not null,
  raw_reference text,
  event_occurred_at timestamptz not null,
  received_at timestamptz not null default now()
);

create index if not exists idx_plan_price_snapshot_plan_id on plan_price_snapshot(plan_id);
create index if not exists idx_order_account_id on "order"(account_id);
create index if not exists idx_payment_event_order_id on payment_event(order_id);
