create table if not exists upstream_credential_ref (
  id uuid primary key,
  provider text not null,
  credential_key text not null,
  display_name text not null,
  status text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists model_route (
  id uuid primary key,
  protocol text not null,
  public_model_name text not null,
  upstream_provider text not null,
  upstream_model_id text not null,
  upstream_credential_ref_id uuid not null references upstream_credential_ref(id),
  status text not null,
  request_adapter_type text,
  response_adapter_type text,
  priority integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists admin_audit_log (
  id uuid primary key,
  operator_id uuid not null,
  operator_type text not null,
  action_type text not null,
  target_type text not null,
  target_id text not null,
  before_snapshot jsonb,
  after_snapshot jsonb,
  reason text,
  created_at timestamptz not null default now()
);

create index if not exists idx_model_route_credential_ref on model_route(upstream_credential_ref_id);
create index if not exists idx_admin_audit_log_target on admin_audit_log(target_type, target_id);
