package router

import (
	"myadmin/controller/callback"

	"github.com/gin-gonic/gin"
)

func initCallBack(r *gin.Engine, s string) {
	globeRoute := r.Group(s) // 全局通用工具类路由。
	{
		// 完成多部份上传。
		globeRoute.Any("/ossupload", callback.OssUpload)
		globeRoute.Any("/mtsnotify", callback.MtsNotify)
	}
}
