package main

import (
	"log"
	"mock-app/common"
	_import "mock-app/import"
	"mock-app/route"
)

func main() {
	router := route.GetRouter()
	// 添加mock路由
	router.POST("/api/v1/mock/route", route.AddMockRouteHandler)
	// 删除mock路由
	router.DELETE("/api/v1/mock/route", route.DelMockRouteHandler)

	// 使用插件导入Mock数据
	// :plugin 插件名
	router.POST("/api/v1/import/plugins/:plugin", _import.ImportHandler)
	// 查询导入功能所有插件
	router.GET("/api/v1/import/plugins", _import.GetImportPlugins)
	// 删除某个导入插件
	router.DELETE("/api/v1/import/plugins/:plugin", _import.DelImportPlugin)
	// 启用或停用某个插件
	router.PUT("/api/v1/import/plugins/:plugin", _import.UpdateImportPlugin)

	router.Use(common.ErrHandler())

	err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}

}
