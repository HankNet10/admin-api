package admin

import (
	"context"
	"myadmin/model"
	"myadmin/model/user"
	"myadmin/util/redis"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AppListList(c *gin.Context) {
	request := user.AppListParam{}
	c.BindQuery(&request)
	request.Sort = "-id"
	list, total := user.App.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AppListCreate(c *gin.Context) {
	request := user.AppList{}
	c.ShouldBindJSON(&request)
	err := request.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	cKey := "api:app_list:" + strconv.Itoa(int(request.ID))
	redis.Redis.Del(context.Background(), cKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func AppListUpdate(c *gin.Context) {
	request := user.AppList{}
	c.ShouldBindJSON(&request)
	if result := model.DataBase.Save(request); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "修改成功！",
		})

		cKey := "api:app_list:" + strconv.Itoa(int(request.ID))
		redis.Redis.Del(context.Background(), cKey)
	}
}

func AppListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.App.SelectByID(uint(rid))
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

	cKey := "api:app_list:" + strconv.Itoa(int(model.ID))
	redis.Redis.Del(context.Background(), cKey)
}
