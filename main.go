package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/heronhoga/talkey-be/config"
	"github.com/heronhoga/talkey-be/handler"
	"github.com/heronhoga/talkey-be/repository"
	"github.com/heronhoga/talkey-be/routes"
	"github.com/heronhoga/talkey-be/service"
	"github.com/heronhoga/talkey-be/util"
	"github.com/heronhoga/talkey-be/util/auth"
)

func main() {
	util.LoadEnv()
	app := fiber.New()

	// Connect to database
	db := config.ConnectDB()
	defer db.Close(context.Background())

	frontEndApp := os.Getenv("FRONTEND_APP")

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontEndApp,
		AllowHeaders:     "Origin, Content-Type, Accept, App-Key",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowCredentials: true,
	}))

	//dependency
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	routes.RegisterUserRoutes(app, userHandler)

	//generate PASETO keys and initialize them
	// auth.GenerateKey()
	auth.Init()

	app.Listen(":8000")
}
