package app

import (
	"gopher-toolbox/app"
)

type ExtendedApp struct {
	app.App
}

func NewExtendedApp(appName, version, envDirectory string) *ExtendedApp {
	app := app.New(appName, version, envDirectory)
	return &ExtendedApp{
		App: *app,
	}
}
