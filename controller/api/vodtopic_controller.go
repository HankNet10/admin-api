package api

import (
	"myadmin/model"
	"myadmin/model/vod"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//获取专题列表
func VodTopicListGet(c *gin.Context) {
	request := vod.VodTopicParam{}
	c.BindQuery(&request)
	var count int64
	var list []vod.VodTopic
	query := model.DataBase.Model(vod.VodTopic{})
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
			Introduce string `json:"introduce"`
		}{item.Name, host + "/" + item.BgImage, item.ID, item.Introduce}
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

func VodTopicDetailsGet(c *gin.Context) {
	request := vod.VodTopicParam{}
	c.BindQuery(&request)
	if request.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "请填入ID",
		})
		return
	}
	TpModel := vod.VodTopic{}
	model.DataBase.Where("id = ?", request.ID).First(&TpModel)
	if TpModel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "专题异常",
		})
		return
	}
	host := os.Getenv("ALI_OSS_DOMAIN")
	data := struct {
		ID        uint   `json:"id"`
		BgImage   string `json:"image"`
		Name      string `json:"name"`
		Introduce string `json:"introduce"`
	}{TpModel.ID, host + "/" + TpModel.BgImage,
		TpModel.Name, TpModel.Introduce}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})

}
