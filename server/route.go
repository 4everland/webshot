package server

import (
	"github.com/4everland/screenshot/api"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.Engine) *gin.Engine {
	router.GET("ping", func(ctx *gin.Context) { ctx.String(200, "pong") })

	router.GET("screenshot", api.Screenshot)

	return router
}
