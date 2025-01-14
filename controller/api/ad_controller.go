package api

import (
	"log"
	"myadmin/model"
	"myadmin/model/ad"
	"myadmin/util/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 广告列表
// @Description 返回广告列表数据
// @Tags 广告
// @Accept json
// @Router /api/ads [get]
func AdDetailList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	cKey := "api:ads:" + appID
	data := make(map[string][]apiResultAdList)
	if err := redis.Deserialize(cKey, &data); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    data,
			"message": "",
		})
		return
	}
	AppId, _ := strconv.Atoi(appID)
	var appid = []int{0, AppId}
	var list []*ad.AdDetail
	if result := model.DataBase.Model(ad.AdDetail{}).Where("app_id in ?", appid).Where("status = 1").Preload("Postion").Find(&list); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	for _, v := range list {
		adlist := newapiResultAdList(v)
		if item, ok := data[v.Postion.Value]; ok {
			data[v.Postion.Value] = append(item, adlist)
		} else {
			data[v.Postion.Value] = []apiResultAdList{adlist}
		}
	}

	if err := redis.Serialize(cKey, data, 2*time.Hour); err != nil {
		log.Panic("Redis set err", cKey, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

// @Summary 提交统计
// @Description 提交广告的点击数量
// @Tags 广告
// @Accept json
// @Param id  query int true "广告ID"
// @Router /api/adclick [get]
func AdClick(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	mw := ad.AdViews{
		DetailID: uint(id),
		Ip:       c.ClientIP(),
		Device:   c.Request.UserAgent(),
	}
	mw.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "",
	})
}
