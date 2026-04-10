create table if not exists usage_event (
  id uuid primary key,
  request_id text not null,
  account_id uuid not null references account(id),
  api_key_id uuid references api_key(id),
  protocol text not null,
  model_name text not null,
  upstream_provider text not null,
  upstream_model_id text,
  request_started_at timestamptz not null,
  request_finished_at timestamptz not null,
  latency_ms bigint not null,
  result_status text not null,
  error_type text,
  input_tokens bigint not null default 0,
  output_tokens bigint not null default 0,
  total_tokens bigint not null default 0,
  billing_period_start timestamptz,
  billing_period_end timestamptz,
  created_at timestamptz not null default now()
);

create table if not exists billing_event (
  id uuid primary key,
  idempotency_key text not null unique,
  account_id uuid not null references account(id),
  usage_event_id uuid not null references usage_event(id),
  billing_period_start timestamptz,
  billing_period_end timestamptz,
  settlement_type text not null,
  quota_delta bigint not null,
  before_quota_remaining bigint not null,
  after_quota_remaining bigint not null,
  result_status text not null,
  failure_reason text,
  created_at timestamptz not null default now(),
  settled_at timestamptz
);

create index if not exists idx_usage_event_account_id on usage_event(account_id);
create index if not exists idx_billing_event_account_id on billing_event(account_id);
