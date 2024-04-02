package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/attributeerror/steam-api-service/handlers/response_models"
	"github.com/attributeerror/steam-api-service/services"
	"github.com/attributeerror/steam-api-service/utilities"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	ErrRequireAtLeastOneQuery = errors.New("missing at least one required parameter: vanity_url, profile_id, or query")
	ErrInvalidSteamUrl        = errors.New("'query' parameter provided, but it wasn't a valid Steam URL")
)

var GetSteamUserHex = func(steamService *services.SteamService, sfGroup *singleflight.Group) func(c *gin.Context) {
	return func(c *gin.Context) {
		vanityUrl := c.DefaultQuery("vanity_url", "")
		profileId := c.DefaultQuery("profile_id", "")
		query := c.DefaultQuery("query", "")

		var groupId string = ""
		if vanityUrl != "" {
			groupId = vanityUrl
		} else if profileId != "" {
			groupId = profileId
		} else if query != "" {
			split := utilities.Filter(strings.Split(query, "/"), func(elem string) bool {
				return strings.TrimSpace(elem) != ""
			})
			if len(split) > 2 {
				if !strings.HasPrefix(query, "https://steamcommunity.com") {
					query = "invalid"
				} else {
					typeKw := split[len(split)-2]
					if typeKw == "profiles" {
						profileId = split[len(split)-1]
						query = profileId
					} else if typeKw == "id" {
						vanityUrl = split[len(split)-1]
						query = vanityUrl
					} else {
						query = "invalid"
					}
				}
			} else {
				if is_numeric := regexp.MustCompile(`\d`).MatchString(query); is_numeric {
					profileId = query
				} else {
					vanityUrl = query
				}
			}
		}

		response, err, _ := sfGroup.Do(groupId, func() (interface{}, error) {
			if query == "invalid" {
				return nil, ErrInvalidSteamUrl
			}
			if vanityUrl == "" && profileId == "" {
				return nil, ErrRequireAtLeastOneQuery
			}

			if vanityUrl != "" {
				var err error = nil
				profileId, err = steamService.ResolveVanityUrlToProfileId(vanityUrl)

				if err != nil {
					return nil, err
				}
			}

			profileIdAsInt, err := strconv.ParseInt(profileId, 0, 64)
			if err != nil {
				return nil, err
			}

			steamHex := fmt.Sprintf("%X", profileIdAsInt)
			return &response_models.GetUserHexResponse{
				ProfileId: profileId,
				VanityUrl: vanityUrl,
				SteamHex:  steamHex,
			}, nil
		})

		if err != nil {
			var statusCode int
			if errors.Is(err, ErrRequireAtLeastOneQuery) {
				statusCode = http.StatusBadRequest
			} else {
				statusCode = http.StatusInternalServerError
			}

			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
			return
		}

		if responseModel, ok := response.(*response_models.GetUserHexResponse); ok {
			c.JSON(http.StatusOK, responseModel)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unknown error occurred whilst parsing the response. Please try again later.",
		})
	}
}
