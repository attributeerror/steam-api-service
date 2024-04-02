package response_models

type (
	GetUserHexResponse struct {
		ProfileId string `json:"profile_id,omitempty"`
		VanityUrl string `json:"vanity_url,omitempty"`
		SteamHex  string `json:"steam_hex,omitempty"`
	}
)
