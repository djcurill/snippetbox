-- create a sessions table so we can begin to use stateful
-- http sessions

create table sessions(
  token text primary key,
  data bytea not null,
  expiry timestamptz not null
);

create index sessions_expiry_idx on sessions (expiry);
