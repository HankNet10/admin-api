package admin

import (
	"myadmin/model/blog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var awardMatchModel blog.AwardMatch

func AwardMatchList(c *gin.Context) {
	request := blog.AwardMatchParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := awardMatchModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AwardMatchCreate(c *gin.Context) {
	request := blog.AwardMatch{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func AwardMatchUpdate(c *gin.Context) {
	var request blog.AwardMatch
	c.ShouldBindJSON(&request)
	model := blogMatchModel.SelectByID(request.MatchId)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "比赛不存在",
		})
		return
	}
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func AwardMatchDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := blog.AwardMatch{}.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "奖品不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
