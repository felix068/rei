-- Init tables, the queries will be executed in db initialization

create table if not exists feeds(
id uuid not null unique primary key,
name text not null,
description text,
link text unique,
created_at timestamp not null,
updated_at timestamp not null
)
