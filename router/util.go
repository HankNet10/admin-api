package router

import (
	"myadmin/controller/util"

	"github.com/gin-gonic/gin"
)

func initUtil(r *gin.Engine, s string) {
	globeRoute := r.Group(s) // 全局通用工具类路由。
	{
		globeRoute.GET("/ping", util.Health)       // 健康检查
		globeRoute.GET("/captcha", util.Captcha)   // 验证码
		globeRoute.GET("/ip", util.IP)             // 获取IP信息
		globeRoute.GET("/starttime", util.Runtime) // 启动时间
	}
}
