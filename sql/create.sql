-- Init tables, the queries will be executed in db initialization

create table if not exists feeds(
id uuid not null unique primary key,
name text not null,
description text,
link text unique,
rss_link text unique,
created_at timestamp not null,
updated_at timestamp not null
);

create table if not exists posts(
id uuid not null unique primary key,
name text not null,
description text,
link text unique,
feed_id uuid not null,
created_at timestamp not null,
updated_at timestamp not null
);
