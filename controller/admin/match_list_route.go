package admin

import (
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/util/sugar"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var blogMatchModel blog.BlogMatch

func BlogMatchList(c *gin.Context) {
	request := blog.BlogMatchParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := blogMatchModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func BlogMatchCreate(c *gin.Context) {
	request := blog.BlogMatch{}
	c.ShouldBindJSON(&request)
	request.Status = 1
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func BlogMatchUpdate(c *gin.Context) {
	var request blog.BlogMatch
	c.ShouldBindJSON(&request)
	model := blogMatchModel.SelectByID(request.ID)
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

func BlogMatchDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := blogMatchModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "比赛不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 结束 开启比赛状态
func BlogMatchStatusUpdate(c *gin.Context) {
	request := blog.BlogMatch{}
	c.ShouldBindJSON(&request)
	if request.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "比赛异常",
		})
		return
	}
	if request.Status == 0 {
		//结束比赛
		var list []*blog.BlogList
		var count int64
		model.DataBase.Model(blog.BlogList{}).Where("match_id = ?", request.ID).Order("favorites desc").Count(&count).Limit(int(request.RankCount)).Find(&list)
		if count < int64(request.RankCount) {
			c.JSON(http.StatusOK, gin.H{
				"message": "数据不够比赛结束",
			})
			return
		}
		zans := ""
		blogIds := ""
		for index, item := range list {
			zans += strconv.Itoa(int(item.Favorites))
			blogIds += strconv.Itoa(int(item.ID))
			if index < int(request.RankCount)-1 {
				zans += ","
				blogIds += ","
			}
		}
		matchRank := blog.MatchRank{}.SelectByMatchID(request.ID)
		matchRank.BlogIds = blogIds
		matchRank.ZanCounts = zans
		matchRank.MatchId = request.ID
		matchRank.Save()
		model := blogMatchModel.SelectByID(request.ID)
		model.Status = 0
		model.Save()
		c.JSON(http.StatusOK, gin.H{
			"message": "结束成功",
		})
		return
	}
	if request.Status == 1 {
		//开启比赛
		model := blogMatchModel.SelectByID(request.ID)
		model.Status = 1
		model.Save()
		c.JSON(http.StatusOK, gin.H{
			"message": "开启成功",
		})
	}

}

func MatchRankList(c *gin.Context) {
	request := blog.MatchRankParam{}
	c.BindQuery(&request)
	matchRank := blog.MatchRank{}.SelectByMatchID(uint(request.MatchId))
	if matchRank.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "该比赛还未结束",
		})
		return
	}
	zans := strings.Split(matchRank.ZanCounts, ",")

	blogIds := strings.Split(matchRank.BlogIds, ",")
	list := []blog.BlogList{}
	for index, id := range blogIds {
		blog := blog.BlogList{}.SelectByID(sugar.StringToUint(id))
		blog.Favorites = sugar.StringToUint(zans[index])
		list = append(list, *blog)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    list,
	})
}
