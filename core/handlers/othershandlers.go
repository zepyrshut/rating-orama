package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (rp Repository) Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":    "Template engine is working! We are in the index page!",
		"SomeData": []string{"Some", "Data", "Here"},
	})
}

func (rp Repository) Another(c *fiber.Ctx) error {
	return c.Render("another", fiber.Map{
		"Title":    "Template engine is working! We are in the another page!",
		"SomeData": []string{"Some", "Data", "Here"},
	})
}

func (rp Repository) Ping(c *fiber.Ctx) error {
	rp.App.Info("Ping!")
	return c.SendString("Pong!")
}

func (rp Repository) Panic(c *fiber.Ctx) error {
	panic("Panic!")
}
