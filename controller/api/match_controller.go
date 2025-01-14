package api

import (
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/util/sugar"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//获取话题列表
func MatchListGet(c *gin.Context) {
	request := blog.BlogMatchParam{}
	c.BindQuery(&request)
	var count int64
	var list []blog.BlogMatch
	query := model.DataBase.Model(blog.BlogMatch{})
	query.Order("sort desc")
	query.Count(&count)
	result := query.Find(&list)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"total": 0,
				"list":  nil,
			},
			"message": "",
		})
		return
	}
	data := make([]interface{}, len(list))
	host := os.Getenv("ALI_OSS_DOMAIN")
	for index, item := range list {
		data[index] = struct {
			Name    string `json:"name"`
			Image   string `json:"image"`
			ID      uint   `json:"id"`
			BgImage string `json:"bgimage"`
			Time    string `json:"time"`
			Status  uint   `json:"status"`
		}{item.Name, host + "/" + item.Image, item.ID, host + "/" + item.BgImage, item.Time, uint(item.Status)}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": count,
			"list":  data,
		},
		"message": "",
	})

}

func MatchDetailsGet(c *gin.Context) {
	request := blog.BlogMatchParam{}
	c.BindQuery(&request)
	if request.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "请填入ID",
		})
		return
	}
	TpModel := blog.BlogMatch{}
	model.DataBase.Where("id = ?", request.ID).First(&TpModel)
	if TpModel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "比赛异常",
		})
		return
	}
	host := os.Getenv("ALI_OSS_DOMAIN")
	data := struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Image   string `json:"image"`
		BgImage string `json:"bgimage"`
		Rule    string `json:"rule"`
		Status  uint   `json:"status"`
		Time    string `json:"time"`
	}{TpModel.ID, TpModel.Name, host + "/" + TpModel.Image, host + "/" + TpModel.BgImage, TpModel.Rule, uint(TpModel.Status), TpModel.Time}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})

}

func MatchEndRankList(c *gin.Context) {
	request := blog.MatchRankParam{}
	c.BindQuery(&request)
	matchRank := blog.MatchRank{}.SelectByMatchID(uint(request.MatchId))
	if matchRank.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "比赛排行没有数据",
		})
		return
	}
	zans := strings.Split(matchRank.ZanCounts, ",")
	blogIds := strings.Split(matchRank.BlogIds, ",")
	list := []*blog.BlogList{}
	for index, id := range blogIds {
		blog := blog.BlogListModel.SelectInfo(sugar.StringToUint(id))
		blog.Favorites = sugar.StringToUint(zans[index])
		list = append(list, blog)
	}
	data := make([]apiResultBlogList, len(list))
	for i, v := range list {
		data[i] = newApiResultBlogList(*v)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  data,
			"total": 10,
		},
		"message": "",
	})
}
