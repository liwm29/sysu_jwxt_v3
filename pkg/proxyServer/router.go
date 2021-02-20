package proxyServer

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "%s", "pang")
	})

}
