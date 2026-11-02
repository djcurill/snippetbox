create table users(
  id integer generated always as identity primary key,
  name varchar(255) not null,
  email varchar(255) not null,
  hashed_password text not null,
  created timestamptz not null,
  constraint users_uc_email unique(email)
)
