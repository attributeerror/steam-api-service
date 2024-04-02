package handlers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

func InitialiseRoutes(e *gin.Engine, sfGroup *singleflight.Group) {
	e.GET("get_user_hex", GetSteamUserHex(sfGroup))
}
