create table if not exists tv_show (
  id integer primary key,
  "name" varchar not null,
  tt_imdb varchar not null,
  popularity int not null default 0,

  created_at timestamp not null default (now()),
  updated_at timestamp not null default (now())
);

create index if not exists idx_tv_show_title on "tv_show" ("title");
create index if not exists idx_tv_show_updated_at on "tv_show" ("updated_at");

create table if not exists episodes (
  id integer primary key,
  tv_show_id integer not null,

  season integer not null,
  episode integer not null,
  released date,
  "name" varchar not null,
  plot text not null default '',
  avg_rating numeric(1,1) not null default 0,
  vote_count int not null default 0,

  foreign key (tv_show_id) references tv_show (id)
);





