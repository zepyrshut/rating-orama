create table if not exists tv_show (
  id serial unique primary key,
  "name" varchar not null,
  tt_imdb varchar not null,
  popularity int not null default 0,

  created_at timestamp not null default (now()),
  updated_at timestamp not null default (now())
);

create index if not exists idx_tv_show_name on "tv_show" ("name");
create index if not exists idx_tv_show_tt_imdb on "tv_show" ("tt_imdb");
create index if not exists idx_tv_show_updated_at on "tv_show" ("updated_at");

create table if not exists episodes (
  id serial unique primary key,
  tv_show_id integer not null,

  season integer not null,
  episode integer not null,
  released date,
  "name" varchar not null,
  plot text not null default '',
  avg_rating real not null default 0,
  vote_count int not null default 0,

  foreign key (tv_show_id) references tv_show (id)
);





