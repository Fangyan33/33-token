create table if not exists account (
  id uuid primary key,
  status text not null,
  display_name text not null,
  default_contact_email text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists user_identity (
  id uuid primary key,
  account_id uuid not null references account(id),
  login_email text not null,
  auth_provider text not null,
  status text not null,
  last_login_at timestamptz,
  created_at timestamptz not null default now()
);

create table if not exists api_key (
  id uuid primary key,
  account_id uuid not null references account(id),
  key_prefix text not null,
  key_hash text not null,
  status text not null,
  created_at timestamptz not null default now(),
  disabled_at timestamptz,
  last_used_at timestamptz
);

create index if not exists idx_user_identity_account_id on user_identity(account_id);
create index if not exists idx_api_key_account_id on api_key(account_id);
