package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/attributeerror/steam-api-service/configuration"
	"github.com/attributeerror/steam-api-service/handlers/response_models"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	ErrRequireAtLeastOneQuery = errors.New("missing at least one required parameter: vanity_url, or profile_id")
)

var GetSteamUserHex = func(sfGroup *singleflight.Group) func(c *gin.Context) {
	return func(c *gin.Context) {
		configuration := configuration.GetConfiguration()
		vanityUrl := c.DefaultQuery("vanity_url", "")
		profileId := c.DefaultQuery("profile_id", "")

		var groupId string
		if vanityUrl != "" {
			groupId = vanityUrl
		} else {
			groupId = profileId
		}

		response, err, _ := sfGroup.Do(groupId, func() (interface{}, error) {
			if vanityUrl == "" && profileId == "" {
				return nil, ErrRequireAtLeastOneQuery
			}

			if vanityUrl != "" {
				vanityUrlResp, err := http.Get(fmt.Sprintf("%s/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s", configuration.SteamApiBaseUrl, configuration.SteamApiKey, vanityUrl))
				if err != nil {
					return nil, err
				}
				defer vanityUrlResp.Body.Close()

				var resolveVanityUrlBody response_models.ResolveVanityURLResponse
				if err := json.NewDecoder(vanityUrlResp.Body).Decode(&resolveVanityUrlBody); err != nil {
					return nil, err
				}

				if resolveVanityUrlBody.Response.Success == 1 {
					profileId = resolveVanityUrlBody.Response.SteamId
				} else {
					return nil, fmt.Errorf("error from Steam API: %s", resolveVanityUrlBody.Response.Message)
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
