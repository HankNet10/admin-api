package admin

import (
	"context"
	"myadmin/model/config"
	"myadmin/util/redis"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var conigListModel config.ConfigList

func ConfigListList(c *gin.Context) {
	request := config.ConfigListParam{}
	c.BindQuery(&request)
	request.Sort = "-id"
	list, total := conigListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func ConfigListCreate(c *gin.Context) {
	request := config.ConfigList{}
	c.ShouldBindJSON(&request)
	err := request.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func ConfigListUpdate(c *gin.Context) {
	request := config.ConfigList{}
	c.ShouldBindJSON(&request)
	request.Save()

	cKey := "api:config:" + request.Name
	redis.Redis.Del(context.Background(), cKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func ConfigListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := conigListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在",
		})
		return
	}
	model.Delete()

	cKey := "api:config:" + model.Name
	redis.Redis.Del(context.Background(), cKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func ConfigSts(c *gin.Context) {
	// 如果存在oss加速首先返回加速信息
	var region string
	if os.Getenv("ALI_OSS_ACCELERATE") != "" {
		region = os.Getenv("ALI_OSS_ACCELERATE")
	} else {
		region = os.Getenv("ALI_OSS_REGION")
	}

	c.JSON(http.StatusOK, gin.H{
		"region":          region,
		"bucket":          os.Getenv("ALI_OSS_ORIGIN"),
		"accessKeyId":     os.Getenv("ALI_ACCESS_KEY_ID"),
		"accessKeySecret": os.Getenv("ALI_ACCESS_KEY_SECRET"),
	})
}
