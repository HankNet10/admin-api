package api

import (
	"myadmin/model"
	"myadmin/model/blog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//获取话题列表
func TopicListGet(c *gin.Context) {
	request := blog.BlogTopicParam{}
	c.BindQuery(&request)
	var count int64
	var list []blog.BlogTopic
	query := model.DataBase.Model(blog.BlogTopic{})
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
			Name      string `json:"name"`
			Image     string `json:"image"`
			ID        uint   `json:"id"`
			Count     uint   `json:"count"`
			Introduce string `json:"introduce"`
		}{item.Name, host + "/" + item.Image, item.ID, item.Count, item.Introduce}
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

func TopicDetailsGet(c *gin.Context) {
	request := blog.BlogTopicParam{}
	c.BindQuery(&request)
	if request.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "请填入ID",
		})
		return
	}
	TpModel := blog.BlogTopic{}
	model.DataBase.Where("id = ?", request.ID).First(&TpModel)
	if TpModel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "话题异常",
		})
		return
	}
	host := os.Getenv("ALI_OSS_DOMAIN")
	data := struct {
		ID        uint   `json:"id"`
		BgImage   string `json:"bgimage"`
		Image     string `json:"image"`
		Name      string `json:"name"`
		Count     uint   `json:"count"`
		Introduce string `json:"introduce"`
	}{TpModel.ID, host + "/" + TpModel.BgImage,
		host + "/" + TpModel.Image, TpModel.Name, TpModel.Count, TpModel.Introduce}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})

}
