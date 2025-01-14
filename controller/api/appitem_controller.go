package api

import (
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/applicationad"
	"myadmin/util/redis"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 分组的列表
func AppItemDetailList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	cKey := "api:appads:1_" + appID
	data := make(map[string][]apiResultAppGroupList)
	if err := redis.Deserialize(cKey, &data); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    data,
			"message": "",
		})
		return
	}

	var list []*applicationad.ApplicationAd
	appid, _ := strconv.Atoi(appID)
	var appidList = []int{0, appid}
	result := model.DataBase.Model(applicationad.ApplicationAd{}).Where("app_id in ?", appidList).Where("status = 1").Where("count > 1").Preload("AppTypeItem").Find(&list)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	for _, v := range list {
		adlist := newapiResultAppGroupList(v)
		if item, ok := data[v.AppTypeItem.Value]; ok {
			data[v.AppTypeItem.Value] = append(item, adlist)
		} else {
			data[v.AppTypeItem.Value] = []apiResultAppGroupList{adlist}
		}
	}

	if err := redis.Serialize(cKey, data, 10*time.Minute); err != nil {
		log.Panic("Redis set err", cKey, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

// 不分组的列表在前端分组
func AppItemAllDetailList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	cKey := "api:appads:2_" + appID
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}

	var list []*applicationad.ApplicationAd
	appid, _ := strconv.Atoi(appID)
	var appidList = []int{0, appid}
	if result := model.DataBase.Model(applicationad.ApplicationAd{}).Where("app_id in ?", appidList).Where("status = 1").Where("count > 1").Preload("AppTypeItem").Find(&list); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	data := make([]apiResultApplicationList, len(list))

	for i, v := range list {
		data[i] = newapiResultApplicationList(v)
	}
	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"list": data,
		},
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 10*time.Minute)
	c.String(http.StatusOK, string(jsonData))

}

// 广告类型列表
func ApplicationTypeLists(c *gin.Context) {
	request := applicationad.ApplicationTypeParam{}
	request.Limit = 100
	request.Page = 1
	list, total := applicationad.ApplicationTypeModel.List(request)
	data := make([]apiResultAdTypeList, len(list))
	for i, v := range list {
		data[i] = newapiResultAdTypeList(v)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  data,
			"total": total,
		},
		"message": "",
	})
}

// @Summary 提交统计
// @Description 提交广告的点击数量
// @Tags 广告
// @Accept json
// @Param id  query int true "广告ID"
// @Router /api/adclick [get]
func APPItemClick(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	device := ""
	if strings.Contains(c.Request.UserAgent(), "okhttp") {
		device = "android"
	} else {
		device = "h5_ios"
	}
	mw := applicationad.ApplicationViews{
		DetailID:  uint(id),
		Ip:        c.ClientIP(),
		Device:    device,
		UserAgent: c.Request.UserAgent(),
	}
	mw.Save()
	modelA := applicationad.ApplicationAdModel.SelectByID(uint(id))
	count := modelA.Count - 1
	model.DataBase.Model(&modelA).Update("count", count) //认证UP
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "",
	})
}
func AppItemDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	modelA := applicationad.ApplicationAdModel.SelectByID(uint(id))
	if modelA != nil {
		adlist := newapiResultAppGroupList(modelA)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    adlist,
			"message": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    "",
			"message": "",
		})
	}

}

func AppItemBelong(c *gin.Context) {
	name := c.Query("name")
	println(name)
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数",
		})
		return
	}
	var list []*applicationad.ApplicationAd
	query := model.DataBase.Model(applicationad.ApplicationAd{})
	query.Where("belong = ?", name)
	query.Find(&list)
	data := make([]apiResultApplicationList, len(list))

	for i, v := range list {
		data[i] = newapiResultApplicationList(v)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": data,
		},
		"message": "",
	})
}

func AppItemView(c *gin.Context) {
	request := applicationad.ApplicationViewsParam{}
	c.BindQuery(&request)
	if request.DetailID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数",
		})
		return
	}
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := applicationad.ApplicationViewsModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
