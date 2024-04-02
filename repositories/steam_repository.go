package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/attributeerror/steam-api-service/configuration"
	"github.com/attributeerror/steam-api-service/models"
)

type (
	SteamRepository struct{}
)

func (*SteamRepository) ResolveVanityUrl(query string) (*models.ResolveVanityURLResponse, error) {
	configuration := configuration.GetConfiguration()

	vanityUrlResp, err := http.Get(fmt.Sprintf("%s/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s", configuration.SteamApiBaseUrl, configuration.SteamApiKey, query))
	if err != nil {
		return nil, err
	}
	defer vanityUrlResp.Body.Close()

	var resolveVanityUrlBody models.ResolveVanityURLResponse
	if err := json.NewDecoder(vanityUrlResp.Body).Decode(&resolveVanityUrlBody); err != nil {
		return nil, err
	}

	if resolveVanityUrlBody.Response.Success == 1 {
		return &resolveVanityUrlBody, nil
	} else {
		return nil, fmt.Errorf("error from Steam API: %s", resolveVanityUrlBody.Response.Message)
	}
}
