package configuration

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var configOnce sync.Once
var configuration *Configuration

func GetConfiguration() *Configuration {
	configOnce.Do(func() {
		err := godotenv.Load()
		if err != nil && !os.IsNotExist(err) {
			panic(fmt.Errorf("error whilst loading .env file: %w", err))
		}

		var golangEnv string
		var steamApiBaseUrl string
		var steamApiKey string

		if envVar := getEnvironmentVariable("GOLANG_ENVIRONMENT", false); envVar != nil {
			golangEnv = strings.ToLower(envVar.Value)
		} else {
			golangEnv = "local"
		}

		if envVar := getEnvironmentVariable("STEAM_API_BASE_URL", true); envVar != nil {
			steamApiBaseUrl = envVar.Value
		}
		if envVar := getEnvironmentVariable("STEAM_API_KEY", true); envVar != nil {
			steamApiKey = envVar.Value
		}

		configuration = &Configuration{
			GoEnvironment:   golangEnv,
			SteamApiBaseUrl: steamApiBaseUrl,
			SteamApiKey:     steamApiKey,
		}
	})

	return configuration
}

type (
	Configuration struct {
		GoEnvironment   string
		SteamApiBaseUrl string
		SteamApiKey     string
	}

	EnvironmentVariable struct {
		Key   string
		Value string
	}
)

func getEnvironmentVariable(key string, required bool) *EnvironmentVariable {
	if value, exists := os.LookupEnv(key); exists {
		if value == "" && required {
			panic(fmt.Errorf("%v environment variable must be set", key))
		} else {
			return &EnvironmentVariable{
				Key:   key,
				Value: strings.TrimSpace(value),
			}
		}
	} else if required {
		panic(fmt.Errorf("%v environment variable must be set", key))
	}

	return nil
}
