package api

import (
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/config"
	"myadmin/util/redis"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 配置列表
func ConfigList(c *gin.Context) {
	name := c.Query("name")

	var data *config.ConfigList
	cKey := "api:config:" + name
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	if result := model.DataBase.Where("name = ?", name).First(&data); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	image := ""
	if data.Image != "" {
		image = os.Getenv("ALI_OSS_DOMAIN") + "/" + data.Image
	}
	result := gin.H{
		"Name":  data.Name,
		"Value": data.Value,
		"Image": image,
	}
	resultData := gin.H{
		"code":    200,
		"data":    result,
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 2*time.Hour)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 配置多个列表
func ConfigGroupList(c *gin.Context) {
	name := c.Query("name")

	var datas []config.ConfigList
	cKey := "api:config:" + name
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	if result := model.DataBase.Where("name = ?", name).Order("sort desc").Find(&datas); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	lists := make([]interface{}, len(datas))
	host := os.Getenv("ALI_OSS_DOMAIN")
	for index, item := range datas {
		lists[index] = struct {
			Name  string `json:"Name"`
			Value string `json:"Value"`
			Image string `json:"Image"`
		}{item.Name, item.Value, host + "/" + item.Image}
	}

	result := gin.H{
		"code": 200,
		"data": gin.H{
			"list": lists,
		},
		"message": "",
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 2*time.Hour)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary cdn加速域名配置
// @Description - 注意不要配置敏感信息
// @Tags 配置
// @Accept json
// @Router /api/cdnhost [get]
func ConfigCdnHost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    os.Getenv("ALI_OSS_DOMAIN"),
		"message": "",
	})
}
