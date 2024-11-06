-- name: CreateTVShow :one
insert into "tv_show" (name, tt_imdb)
values ($1, $2)
returning *;

-- name: CreateEpisodes :one
insert into "episodes" (tv_show_id, season, episode, released, name, plot, avg_rating, vote_count)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: CheckTVShowExists :one
select * from "tv_show"
where tt_imdb = $1;

-- name: GetEpisodes :many
select * from "episodes"
where tv_show_id = $1;

-- name: IncreasePopularity :exec
update "tv_show" set popularity = popularity + 1
where tt_imdb = $1;

-- name: TvShowAverageRating :one
select avg(avg_rating) from "episodes"
where tv_show_id = $1;

-- name: TvShowMedianRating :one
select percentile_cont(0.5) within group (order by avg_rating) from "episodes"
where tv_show_id = $1;

-- name: SeasonAverageRating :one
select avg(avg_rating) from "episodes"
where tv_show_id = $1 and season = $2;

-- name: SeasonMedianRating :one
select percentile_cont(0.5) within group (order by avg_rating) from "episodes"
where tv_show_id = $1 and season = $2;