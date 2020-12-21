package main

import (
	"net/http"

	"goblog/app/http/middlewares"
	"goblog/bootstrap"
)

func main() {
	// 数据库连接及数据库表初始化
	bootstrap.SetupDB()

	// 初始化路由对象
	router := bootstrap.SetupRoute()

	// 服务
	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
