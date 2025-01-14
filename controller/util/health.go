package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 检查服务健康
// @Schemes 服务器健康检查
// @Description 检查服务器健康状态
// @Tags 工具
// @Accept json
// @Produce json
// @Success 200
// @Router /util/ping [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "ok",
	})
}
