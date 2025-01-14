package api

import (
	"myadmin/model"
	"myadmin/model/blog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AwardDetailsGet(c *gin.Context) {
	request := blog.AwardMatchParam{}
	c.BindQuery(&request)
	if request.MatchID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "请填入ID",
		})
		return
	}
	var count int64
	var list []blog.AwardMatch
	query := model.DataBase.Model(blog.AwardMatch{})
	query.Order("sort desc")
	query.Where("match_id = ?", request.MatchID)
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
			Title     string `json:"title"`
			Image     string `json:"image"`
			Introduce string `json:"introduce"`
		}{item.Name, host + "/" + item.Image, item.ReMark}
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
