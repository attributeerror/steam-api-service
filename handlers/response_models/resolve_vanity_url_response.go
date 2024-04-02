package response_models

type (
	ResolveVanityURLResponse struct {
		Response struct {
			SteamId string `json:"steamid,omitempty"`
			Success int    `json:"success,omitempty"`
			Message string `json:"message,omitempty"`
		} `json:"response,omitempty"`
	}
)
