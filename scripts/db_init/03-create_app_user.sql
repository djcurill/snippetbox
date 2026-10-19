-- step 03: create web user with less privileges than admin

\getenv app_pw SNIPPET_DB_PASSWORD

create user web with password :'app_pw';

grant usage on schema public to web;

grant select, insert, update, delete on all tables in schema public to web;

-- future grants on all new tables
alter default privileges in schema public
  grant select, insert, update, delete on tables
  to web;
