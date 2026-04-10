create table if not exists account_subscription_state (
  id uuid primary key,
  account_id uuid not null unique references account(id),
  current_order_id uuid references "order"(id),
  current_plan_price_snapshot_id uuid references plan_price_snapshot(id),
  status text not null,
  billing_period_start timestamptz,
  billing_period_end timestamptz,
  activated_at timestamptz,
  expired_at timestamptz,
  paused_at timestamptz,
  updated_at timestamptz not null default now()
);

create table if not exists account_cycle_summary (
  id uuid primary key,
  account_id uuid not null references account(id),
  billing_period_start timestamptz not null,
  billing_period_end timestamptz not null,
  quota_total bigint not null,
  quota_used bigint not null default 0,
  quota_remaining bigint not null,
  input_tokens_total bigint not null default 0,
  output_tokens_total bigint not null default 0,
  total_tokens_total bigint not null default 0,
  status text not null,
  updated_at timestamptz not null default now()
);

create unique index if not exists idx_account_cycle_summary_account_period
  on account_cycle_summary(account_id, billing_period_start, billing_period_end);
