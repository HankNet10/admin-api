package util

import (
	"myadmin/util/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CachePost struct {
	Name string
	Data string
	Sec  uint
}

// @Summary 存储数据
// @Schemes 存储数据
// @Description redis临时存储数据
// @Tags 工具
// @Accept json
// @Produce json
// @Param data body CachePost true "参数列表"
// @Router /util/cache [post]
func CacheSet(c *gin.Context) {
	var request CachePost
	c.ShouldBindJSON(&request)
	request.Sec = 3600
	request.Name = c.ClientIP()
	redis.Set("apicache:"+request.Name, request.Data, time.Second*time.Duration(request.Sec))
	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"message": "完成",
		"ip":      request.Name,
	})
}

// @Summary 获取数据
// @Schemes 获取数据
// @Description 获取redis的key
// @Tags 工具
// @Accept json
// @Produce json
// @Param name query string true "查询IP地址"
// @Success 200 {object} ModelCaptcha
// @Router /util/cache [get]
func CacheGet(c *gin.Context) {
	name := c.ClientIP()
	data, _ := redis.Get("apicache:" + name)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"ip":   name,
	})
}
