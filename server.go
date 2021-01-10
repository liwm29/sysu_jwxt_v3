package main

import (
	"github.com/gin-gonic/gin"
	// 注意这里,涉及go mod的本地包导入,如果是子文件夹的包,不需要在子包也go mod init
	// 引入子包的方法: 比如go mod init module_name,则引入module_name/package_dir_name,即子包的文件夹名,但引入后,用的是go文件头部package_name
	"server/backend/proxyServer"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func index(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	ctx.File("frontend/dist/index.html")
}

func main() {
	r := gin.Default()
	r.GET("/", index)
	proxyServer.InitRouter(r)
	r.Static("static", "frontend/dist/static")
	r.StaticFile("favicon.ico", "frontend/dist/favicon.ico")
	r.Run(":12345")
}
