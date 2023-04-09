package app

import (
	"golang.org/x/exp/slog"
	"os"
)

type Application struct {
	*slog.Logger
	Environment
}

type Environment struct {
	Datasource   string
	HarvesterApi string
}

func NewApp(isProduction bool) *Application {
	if isProduction {
		return &Application{
			newStructuredLogger(),
			Environment{
				Datasource:   os.Getenv("DATASOURCE"),
				HarvesterApi: os.Getenv("HARVESTER_API"),
			},
		}
	} else {
		return &Application{
			newStructuredLogger(),
			Environment{
				Datasource:   "postgres://postgres:postgres@localhost:5432/postgres",
				HarvesterApi: "http://localhost:5000/tv-show/%s",
			},
		}
	}
}
