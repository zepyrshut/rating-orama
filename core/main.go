package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/zepyrshut/rating-orama/app"
	"github.com/zepyrshut/rating-orama/db"
	"github.com/zepyrshut/rating-orama/handlers"
	"github.com/zepyrshut/rating-orama/router"
	"os"
)

var application *app.Application
var isProduction = os.Getenv("IS_PRODUCTION") == "true"

const version = "0.1.0"
const appName = "Rating Orama Core"
const author = "Pedro Pérez"

func main() {

	engine := html.New("./views", ".html")

	fmt.Println(isProduction)

	fiberApp := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	application = app.NewApp(isProduction)

	dbPool := db.NewDBPool(application.Datasource)
	defer dbPool.Close()

	repo := handlers.NewRepo(dbPool, application)
	handlers.NewHandlers(repo)
	router.Routes(fiberApp)

	application.Logger.Info("API is running", "version", version, "app", appName, "author", author)
	err := fiberApp.Listen("0.0.0.0:3000")
	if err != nil {
		application.Logger.Error(err.Error())
		return
	}

}
