package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	configOnce                  sync.Once
	configuration               *Configuration
	defaultCachingExpiration    time.Duration = 3600 * time.Second
	defaultCachingPurgeDuration time.Duration = 600 * time.Second
)

func GetConfiguration() *Configuration {
	configOnce.Do(func() {
		err := godotenv.Load()
		if err != nil && !os.IsNotExist(err) {
			panic(fmt.Errorf("error whilst loading .env file: %w", err))
		}

		var port string
		var golangEnv string
		var steamApiBaseUrl string
		var steamApiKey string
		var cachingExpirationDuration time.Duration
		var cachingPurgeDuration time.Duration

		if envVar := getEnvironmentVariable("GOLANG_ENVIRONMENT", false); envVar != nil {
			golangEnv = strings.ToLower(envVar.Value)
		} else {
			golangEnv = "local"
		}

		if envVar := getEnvironmentVariable("CACHING_EXPIRATION_SECONDS", false); envVar != nil {
			parsedEnvVar, err := strconv.ParseInt(envVar.Value, 10, 64)
			if err != nil {
				cachingExpirationDuration = defaultCachingExpiration
			} else {
				cachingExpirationDuration = time.Duration(parsedEnvVar) * time.Second
			}
		} else {
			cachingExpirationDuration = defaultCachingExpiration
		}

		if envVar := getEnvironmentVariable("CACHING_PURGE_SECONDS", false); envVar != nil {
			parsedEnvVar, err := strconv.ParseInt(envVar.Value, 10, 64)
			if err != nil {
				cachingPurgeDuration = defaultCachingPurgeDuration
			} else {
				cachingPurgeDuration = time.Duration(parsedEnvVar) * time.Second
			}
		} else {
			cachingPurgeDuration = defaultCachingPurgeDuration
		}

		if envVar := getEnvironmentVariable("STEAM_API_BASE_URL", true); envVar != nil {
			steamApiBaseUrl = envVar.Value
		}
		if envVar := getEnvironmentVariable("STEAM_API_KEY", true); envVar != nil {
			steamApiKey = envVar.Value
		}
		if envVar := getEnvironmentVariable("PORT", false); envVar != nil {
			port = fmt.Sprintf(":%s", envVar.Value)
		} else {
			port = ":80"
		}

		configuration = &Configuration{
			Port:                      port,
			GoEnvironment:             golangEnv,
			SteamApiBaseUrl:           steamApiBaseUrl,
			SteamApiKey:               steamApiKey,
			CachingExpirationDuration: cachingExpirationDuration,
			CachingPurgeDuration:      cachingPurgeDuration,
		}
	})

	return configuration
}

type (
	Configuration struct {
		Port                      string
		GoEnvironment             string
		SteamApiBaseUrl           string
		SteamApiKey               string
		CachingExpirationDuration time.Duration
		CachingPurgeDuration      time.Duration
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
