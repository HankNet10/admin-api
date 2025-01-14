package admin

import (
	"context"
	"myadmin/model"
	"myadmin/model/applicationad"
	"myadmin/util/redis"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ApplicationAdList(c *gin.Context) {
	request := applicationad.ApplicationAdParam{}
	c.BindQuery(&request)
	request.Sort = "-id"
	list, total := applicationad.ApplicationAdModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func ApplicationAdCreate(c *gin.Context) {
	request := applicationad.ApplicationAd{}
	c.ShouldBindJSON(&request)
	err := request.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if request.AppID == 0 {
		cKey := "api:appads:1_1"
		cKey2 := "api:appads:2_1"
		cKey3 := "api:appads:1_2"
		cKey4 := "api:appads:2_2"
		redis.Redis.Del(context.Background(), cKey)
		redis.Redis.Del(context.Background(), cKey2)
		redis.Redis.Del(context.Background(), cKey3)
		redis.Redis.Del(context.Background(), cKey4)
	} else if request.AppID > 0 {
		cKey := "api:appads:1_" + strconv.Itoa(int(request.AppID))
		cKey2 := "api:appads:2_" + strconv.Itoa(int(request.AppID))
		redis.Redis.Del(context.Background(), cKey)
		redis.Redis.Del(context.Background(), cKey2)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func ApplicationAdUpdate(c *gin.Context) {
	request := applicationad.ApplicationAd{}
	c.ShouldBindJSON(&request)
	if result := model.DataBase.Save(request); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "修改成功！",
		})
		if request.AppID == 0 {
			cKey := "api:appads:1_1"
			cKey2 := "api:appads:2_1"
			cKey3 := "api:appads:1_2"
			cKey4 := "api:appads:2_2"
			redis.Redis.Del(context.Background(), cKey)
			redis.Redis.Del(context.Background(), cKey2)
			redis.Redis.Del(context.Background(), cKey3)
			redis.Redis.Del(context.Background(), cKey4)
		} else if request.AppID > 0 {
			cKey := "api:appads:1_" + strconv.Itoa(int(request.AppID))
			cKey2 := "api:appads:2_" + strconv.Itoa(int(request.AppID))
			redis.Redis.Del(context.Background(), cKey)
			redis.Redis.Del(context.Background(), cKey2)
		}
	}
}

func ApplicationAdDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := applicationad.ApplicationAdModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})

	cKey := "api:appads"
	redis.Redis.Del(context.Background(), cKey)
}

func ApplicationAdViewsList(c *gin.Context) {
	request := applicationad.ApplicationViewsParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := applicationad.ApplicationViewsModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
