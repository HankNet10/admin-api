package router

import (
	"myadmin/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitAdminRouter() *gin.Engine {
	// gin 允许全局跨域
	r := gin.Default()
	r.Use(CORSMiddleware()) // gin 允许全局跨域

	if gin.Mode() != "release" { //处理swag
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.BasePath = "/"
	}

	// api接口
	initApi(r, "/api")
	// 播放相关
	initPlay(r, "/play")
	// 回调
	initCallBack(r, "/callback")
	// 工具
	initUtil(r, "/util")
	// 后台接口
	initAdmin(r, "/admin")

	return r
}
