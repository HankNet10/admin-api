package util

import (
	"myadmin/util/ip"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 查询IP信息
// @Schemes 查询IP信息
// @Description 查询IP信息
// @Param addr query string true "查询IP地址"
// @Tags 工具
// @Accept json
// @Produce json
// @Success 200
// @Router /util/ip [get]
func IP(c *gin.Context) {
	addr := c.Query("addr")
	c.JSON(http.StatusOK, ip.Find(addr))
}

var starttime time.Time

func init() {
	starttime = time.Now()
}

func Runtime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"starttime": starttime.Format("2006-01-02 15:04:05"),
	})
}
