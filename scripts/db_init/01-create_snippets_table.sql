create table snippets (
	id integer generated always as identity primary key,
	title varchar(100) not null,
	content text not null,
	created timestamptz not null,
	expires timestamptz not null
)
