import imdb

from models.tv_show import TVShow, Episode

ia = imdb.Cinemagoer()


def get_tv_show_episodes(tt_id):
    tv_show_episodes = ia.get_movie(tt_id)
    ia.update(tv_show_episodes, 'episodes')

    runtime = tv_show_episodes['runtimes'][0] if 'runtimes' in tv_show_episodes else 0
    tv_show = TVShow(tv_show_episodes.getID(), tv_show_episodes['original title'], runtime)

    episodes = []

    for season in tv_show_episodes['episodes']:
        for episode in tv_show_episodes['episodes'][season]:
            one_episode = Episode(
                tv_show_episodes['episodes'][season][episode].getID(),
                tv_show_episodes['episodes'][season][episode].get('title', "#{}.{}".format(season, episode)),
                tv_show_episodes['episodes'][season][episode].get('episode', episode),
                tv_show_episodes['episodes'][season][episode].get('original air date', 0),
                tv_show_episodes['episodes'][season][episode].get('rating', 0),
                tv_show_episodes['episodes'][season][episode].get('votes', 0),
                tv_show_episodes['episodes'][season][episode].get('season', season))

            episodes.append(one_episode.to_dict())

    tv_show.add_episodes(episodes)
    return tv_show



