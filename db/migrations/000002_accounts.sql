CREATE TABLE account (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL
);

CREATE TABLE user_identity (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL
);

CREATE TABLE api_key (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    key_hash TEXT NOT NULL
);
