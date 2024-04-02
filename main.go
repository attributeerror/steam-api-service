package main

import (
	"fmt"
	"os"

	"github.com/attributeerror/steam-api-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/sync/singleflight"
)

func main() {
	if err := loadDotEnv(); err != nil {
		panic(fmt.Errorf("error whilst loading .env file: %w", err))
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	sfGroup := singleflight.Group{}
	handlers.InitialiseRoutes(engine, &sfGroup)

	port, _ := loadenvvar("PORT", false)
	if port == nil {
		engine.Run(":80")
	} else {
		engine.Run(fmt.Sprintf(":%s", *port))
	}
}

func loadDotEnv() error {
	err := godotenv.Load()

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	return nil
}

func loadenvvar(key string, required bool) (*string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return &value, nil
	} else if required {
		return nil, fmt.Errorf("required env var not set: %s", key)
	}

	return nil, nil
}
