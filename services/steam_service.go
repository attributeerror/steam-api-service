package services

import (
	"github.com/attributeerror/steam-api-service/models"
	"github.com/attributeerror/steam-api-service/repositories"
)

type (
	SteamService struct {
		SteamRepository *repositories.SteamRepository
	}
)

func (s *SteamService) ResolveVanityUrl(query string) (*models.ResolveVanityURLResponse, error) {
	return s.SteamRepository.ResolveVanityUrl(query)
}

func (s *SteamService) ResolveVanityUrlToProfileId(query string) (string, error) {
	vanityUrlResponse, err := s.ResolveVanityUrl(query)
	if err != nil {
		return "", err
	}

	return vanityUrlResponse.Response.SteamId, nil
}
