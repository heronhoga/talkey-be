package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/heronhoga/talkey-be/config"
	"github.com/heronhoga/talkey-be/util"
)

func main() {
	util.LoadEnv()
	app := fiber.New()

	//db config
	config.ConnectDB()

	frontEndApp := os.Getenv("FRONTEND_APP")

	//cors config
	app.Use(cors.New(cors.Config{
        AllowOrigins: frontEndApp,
        AllowHeaders: "Origin, Content-Type, Accept, App-Key",
        AllowMethods: "GET,POST,PUT,DELETE",
        AllowCredentials: true,
        }))

	app.Listen(":8000")
}