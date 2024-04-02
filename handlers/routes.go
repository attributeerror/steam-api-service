package handlers

import (
	"github.com/attributeerror/steam-api-service/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

func InitialiseRoutes(e *gin.Engine, steamService *services.SteamService, sfGroup *singleflight.Group) {
	e.GET("get_user_hex", GetSteamUserHex(steamService, sfGroup))
	e.GET("liveness", LivenessProbe())
}
