create table users (
	id serial primary key,
	name text not null,
	email_address text unique not null
);
