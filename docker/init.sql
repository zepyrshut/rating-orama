CREATE TABLE IF NOT EXISTS "tv_show" (
  "show_id" integer PRIMARY KEY,
  "title" varchar NOT NULL,
  "runtime" integer NOT NULL,
  "popularity" integer NOT NULL DEFAULT 0,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "episodes" (
  "episode_id" integer PRIMARY KEY,
  "tv_show_id" integer NOT NULL,
  "season_number" integer NOT NULL,
  "title" varchar NOT NULL,
  "number" int NOT NULL,
  "aired" date NOT NULL,
  "avg_rating" numeric NOT NULL,
  "votes" int NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tv_show_show_id ON "tv_show" ("show_id");

CREATE INDEX IF NOT EXISTS idx_tv_show_title ON "tv_show" ("title");

CREATE INDEX IF NOT EXISTS idx_tv_show_updated_at ON "tv_show" ("updated_at");

CREATE INDEX IF NOT EXISTS idx_episodes_avg_rating ON "episodes" ("avg_rating");

ALTER TABLE "episodes" ADD FOREIGN KEY ("tv_show_id") REFERENCES "tv_show" ("show_id");
