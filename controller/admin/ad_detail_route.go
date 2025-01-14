package admin

import (
	"context"
	"myadmin/model"
	"myadmin/model/ad"
	"myadmin/util/redis"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdDetailList(c *gin.Context) {
	request := ad.AdDetailParam{}
	c.BindQuery(&request)
	request.Sort = "-id"
	list, total := ad.AdDetailModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AdDetailCreate(c *gin.Context) {
	request := ad.AdDetail{}
	c.ShouldBindJSON(&request)
	err := request.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	cKey := "api:ads:" + strconv.Itoa(int(request.AppID))
	redis.Redis.Del(context.Background(), cKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func AdDetailUpdate(c *gin.Context) {
	request := ad.AdDetail{}
	c.ShouldBindJSON(&request)
	if result := model.DataBase.Save(request); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "修改成功！",
		})

		cKey := "api:ads:" + strconv.Itoa(int(request.AppID))
		redis.Redis.Del(context.Background(), cKey)
	}
}

func AdDetailDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := ad.AdDetailModel.SelectByID(uint(rid))
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

	cKey := "api:ads:" + strconv.Itoa(int(model.AppID))
	redis.Redis.Del(context.Background(), cKey)
}

func AdViewsList(c *gin.Context) {
	request := ad.AdViewsParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := ad.AdViewsModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AdDetailReplaceAction(c *gin.Context) {
	request := ad.AdDetail{}
	c.ShouldBindJSON(&request)
	var adlists []ad.AdDetail
	model.DataBase.Where("belong = ?", request.Belong).Find(&adlists)
	for _, item := range adlists {
		item.Action = request.Action
		item.Save()
	}
	redis.PullPrefix("api:ads:") //删除缓存
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}
