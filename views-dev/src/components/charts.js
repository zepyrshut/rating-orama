export function initCharts (tvShowParsed) {
  new EpisodeChart(tvShowParsed, 'episodes', 0)
  new SeasonsChart(tvShowParsed, 'seasons')
}

export function loadSpecificSeason (tvShowParsed, seasonId) {
  new EpisodeChart(tvShowParsed, 'episodes', seasonId - 1)
}

let chartInstance

class SeasonsChart {
  constructor (tvShowParsed, canvasId) {
    this.tvShowParsed = tvShowParsed
    this.canvasId = canvasId
    this.createChart()
  }

  createChart () {
    const seasons = this.tvShowParsed.seasons
    const labels = seasons.map((season) => `Season ${season.number}`)
    const averageRating = seasons.map((season) => season.avg_rating)
    const medianRating = seasons.map((season) => season.median_rating)
    const votes = seasons.map((season) => season.votes)
    const title = 'Seasons'
    const ctx = document.getElementById(this.canvasId)

    new Chart(ctx, {
      type: 'bar',
      data: {
        labels,
        datasets: [
          {
            label: 'Average rating',
            data: averageRating,
            backgroundColor: 'rgba(75, 192, 192, 0.5)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
            Range: 10,
            yAxisID: 'y-axis-ratings'
          },
          {
            label: 'Median rating',
            data: medianRating,
            backgroundColor: 'rgba(255, 206, 86, 0.5)',
            borderColor: 'rgba(255, 206, 86, 1)',
            borderWidth: 1,
            Range: 10,
            yAxisID: 'y-axis-ratings'

          },
          {
            label: 'Votes',
            data: votes,
            type: 'line',
            tension: 0.4,
            fill: false,
            borderColor: 'rgba(255, 99, 132, 1)',
            borderWidth: 2,
            yAxisID: 'y-axis-votes'
          }
        ]
      },
      options: {
        animation: {
          duration: 0
        },
        y: {
          stacked: true
        },
        scales: {
          'y-axis-ratings': {
            min: 0,
            max: 10,
            type: 'linear',
            display: true,
            position: 'left',
            beginAtZero: true
          },
          'y-axis-votes': {
            type: 'linear',
            display: true,
            position: 'right',
            beginAtZero: true
          }
        }
      },
      plugins: {
        title: {
          display: true,
          text: title
        }
      }
    })
  }
}

class EpisodeChart {
  constructor (tvShowParsed, canvasId, seasonId) {
    this.tvShowParsed = tvShowParsed
    this.canvasId = canvasId
    this.seasonId = seasonId
    this.createChart()
  }

  createChart () {
    if (chartInstance) {
      chartInstance.destroy()
    }

    const episodes = this.tvShowParsed.seasons[this.seasonId].episodes
    const labels = episodes.map((episode) => `Ep. ${episode.number}`)
    const ratings = episodes.map((episode) => episode.avg_rating)
    const votes = episodes.map((episode) => episode.votes)
    const title = `Episodes of season ${this.seasonId + 1}`
    const ctx = document.getElementById(this.canvasId)

    chartInstance = new Chart(ctx, {
      type: 'bar',
      data: {
        labels,
        datasets: [
          {
            label: 'Average rating',
            data: ratings,
            backgroundColor: 'rgba(75, 192, 192, 0.5)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
            Range: 10,
            yAxisID: 'y-axis-ratings'
          },
          {
            label: 'Votes',
            data: votes,
            type: 'line',
            tension: 0.4,
            fill: false,
            borderColor: 'rgba(255, 99, 132, 1)',
            borderWidth: 2,
            yAxisID: 'y-axis-votes'
          }
        ]
      },
      options: {
        animation: {
          duration: 0
        },
        scales: {
          'y-axis-ratings': {
            min: 0,
            max: 10,
            type: 'linear',
            display: true,
            position: 'left',
            beginAtZero: true
          },
          'y-axis-votes': {
            type: 'linear',
            display: true,
            position: 'right',
            beginAtZero: true
          }
        },
        plugins: {
          tooltip: {
            callbacks: {
              title: function (context) {
                const index = context[0].dataIndex
                const episode = episodes[index]
                return `${episode.title} (${episode.aired.split('T')[0]})`
              }
            }
          }
        }
      }
    })
  }
}
