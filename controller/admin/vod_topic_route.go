package admin

import (
	"myadmin/model/vod"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var vodTopicModel vod.VodTopic

func VodTopicList(c *gin.Context) {
	request := vod.VodTopicParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := vodTopicModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VodTopicCreate(c *gin.Context) {
	request := vod.VodTopic{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VodTopicUpdate(c *gin.Context) {
	var request vod.VodTopic
	c.ShouldBindJSON(&request)
	model := vodTopicModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "话题不存在",
		})
		return
	}
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func VodTopicDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vodTopicModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "话题不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
