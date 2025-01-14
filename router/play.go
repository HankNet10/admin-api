package router

import (
	"myadmin/controller/play"

	"github.com/gin-gonic/gin"
)

func initPlay(r *gin.Engine, s string) {
	playRoute := r.Group(s) // 全局通用工具类路由。
	{
		// 长视频的播放地址
		playRoute.GET("/:id/vod.plist", play.VodHlsM3u8Enc)
		playRoute.GET("/:id/vod.plist.m3u8", play.VodHlsM3u8)
		playRoute.GET("/:id/vod.enc", play.VodHlsKey)
		// 短视频的播放地址
		playRoute.GET("/:id/vlog.plist", play.BlogHlsM3u8Enc)
		playRoute.GET("/:id/vlog.plist.m3u8", play.BlogHlsM3u8)
		playRoute.GET("/:id/vlog.enc", play.BlogHlsKey)
		// 社区视频的播放地址
		playRoute.GET("/:id/blog.plist", play.BlogHlsM3u8Enc)
		playRoute.GET("/:id/blog.plist.m3u8", play.BlogHlsM3u8)
		playRoute.GET("/:id/blog.enc", play.BlogHlsKey)
	}
}
