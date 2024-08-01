package main

import (
	"github.com/attributeerror/steam-api-service/configuration"
	"github.com/attributeerror/steam-api-service/handlers"
	"github.com/attributeerror/steam-api-service/repositories"
	"github.com/attributeerror/steam-api-service/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

func main() {
	config := configuration.GetConfiguration()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	sfGroup := &singleflight.Group{}

	steamRepository := &repositories.SteamRepository{}
	steamService := &services.SteamService{
		SteamRepository: steamRepository,
	}
	handlers.InitialiseRoutes(engine, steamService, sfGroup)

	engine.Run(config.Port)
}
