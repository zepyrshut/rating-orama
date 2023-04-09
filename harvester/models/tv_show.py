class TVShow:
    def __init__(self, tt_show_id, title, runtime):
        self.tt_show_id = tt_show_id
        self.title = title
        self.runtime = runtime
        self.episodes = []

    def add_episodes(self, episodes):
        self.episodes = episodes

    def to_dict(self):
        return {
            'tt_show_id': self.tt_show_id,
            'title': self.title,
            'runtime': self.runtime,
            'episodes': [episode.to_dict() for episode in self.episodes]
        }


class Episode:
    def __init__(self, tt_episode_id, title, number, aired, avg_rating, votes, season_id):
        self.tt_episode_id = tt_episode_id
        self.title = title
        self.number = number
        self.aired = aired
        self.avg_rating = avg_rating
        self.votes = votes
        self.season_id = season_id

    def to_dict(self):
        return {
            'tt_episode_id': self.tt_episode_id,
            'title': self.title,
            'number': self.number,
            'aired': self.aired,
            'avg_rating': self.avg_rating,
            'votes': self.votes,
            'season_id': self.season_id
        }
