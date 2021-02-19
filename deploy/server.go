package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liwm29/sysu_jwxt_v3/backend/proxyServer"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func index(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	ctx.File("../frontend/dist/index.html")
}

func main() {
	r := gin.Default()
	r.GET("/", index)
	proxyServer.InitRouter(r)
	r.Static("static", "../frontend/dist/static")
	r.StaticFile("favicon.ico", "../frontend/dist/favicon.ico")
	r.Run(":12345")
}
