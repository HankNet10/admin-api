package api

import (
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/seximg"
	"myadmin/util/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var seximgTypeModel seximg.SeximgType

// @Summary 色图分类列表
// @Tags 色图
// @Param data query seximg.SeximgTypeParam true "参数列表"
// @Router /api/seximgtype/list [get]
func SeximgTypeList(c *gin.Context) {
	request := seximg.SeximgTypeParam{}
	c.BindQuery(&request)
	list, total := seximgTypeModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  list,
			"total": total,
		},
		"message": "",
	})
}

// @Summary 色图列表
// @Tags 色图
// @Param data query seximg.SeximgParam true "参数列表"
// @Router /api/seximg/list [get]
func SeximgList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request seximg.SeximgParam
	c.ShouldBindQuery(&request)
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Page < 1 {
		request.Page = 1
	}
	request.Top = 1
	if request.Order == "" {
		request.Order = "-id"
	}
	request.AppId, _ = strconv.Atoi(appID)
	// 处理色图缓存key
	selectQ := c.Query("user_id")
	top := c.Query("top")
	mustQ := strconv.Itoa(request.Page) + ":" + strconv.Itoa(request.Limit) + ":" + request.Order +
		":" + strconv.Itoa(request.Typeid)
	cKey := "api:seximg:list:" + appID + ":" + selectQ + ":" + mustQ + ":" + top + ":" + request.Time
	if request.Typeid > 0 {
		cKey += ":t" + strconv.Itoa(request.Typeid)
	}
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}

	list, total := seximg.SeximgModel.SeximgSelectList(request)

	data := make([]apiResultSeximg, len(list))
	for i, v := range list {
		data[i] = newApiResultSeximg(*v)
	}
	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"list":  data,
			"total": total,
		},
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 2*time.Hour)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 色图详情
// @Description 返回色图详情数据
// @Tags 色图
// @Accept json
// @Param id query string true "色图详情的id"
// @Router /api/seximg/info [get]
func SeximgInfo(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	sid, _ := strconv.Atoi(id)
	data := seximg.SeximgModel.SelectInfo(uint(sid))

	var rdata apiResultSeximg
	if data != nil {
		seximg.SeximgAddLooks(uint(sid))
		rdata = newApiResultSeximg(*data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    rdata,
		"message": "",
	})
}

// @Summary 是否点赞收藏
// @Description 返回isLike
// @Tags 色图
// @Param seximg_id query string true "色图的id"
// @Router /api/user/seximg/star/isLike [get]
func SeximgIsLike(c *gin.Context) {
	seximgId := c.Query("seximg_id")
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	userID := c.MustGet("UserID").(string)
	rid, err := strconv.Atoi(userID)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	SeximgStarModel := seximg.SeximgStar{}
	result := model.DataBase.Where("user_id = ? and seximg_id = ?", rid, seximgId).Limit(1).Find(&SeximgStarModel)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"isLike": 0,
			},
			"message": "完成",
		})
	}
	isLike := 0
	if SeximgStarModel.ID != 0 {
		isLike = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"isLike": isLike,
		},
		"message": "完成",
	})
}

// @Summary 创建点赞收藏记录
// @Description 创建点赞收藏记录
// @Tags 色图
// @Security ApiKeyAuth
// @Accept json
// @Param data body UserSeximgStarAddRequest true "参数列表"
// @Router /api/user/seximg/star/add [put]
func SeximgStarAdd(c *gin.Context) {
	var request UserSeximgStarAddRequest
	c.ShouldBindJSON(&request)
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	cacheKey := "middle-cgr-" + appID + ":api:seximg:info?id=" + strconv.FormatUint(uint64(request.SeximgID), 10)
	redis.Pull(cacheKey)

	var seximgId uint
	if request.SeximgID > 0 {
		seximgId = request.SeximgID
	}
	reuslt := seximg.SeximgAddUserStar(uint(uID), seximgId)
	if reuslt == 0 {
		seximg.SeximgAddFavorites(seximgId)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}

// @Summary 获取点赞收藏记录
// @Description 获取某个用户的点赞收藏记录
// @Tags 色图
// @Security ApiKeyAuth
// @Accept json
// @Param data query seximg.SeximgStarParam true "参数列表"
// @Router /api/user/seximg/star/list [get]
func SeximgStarList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var request seximg.SeximgStarParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}

	list, total := seximg.SeximgUserStar(uint(uID), request)
	data := make([]apiResultSeximg, len(list))
	for i, v := range list {
		data[i] = newApiResultSeximg(v.Seximg)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  data,
		},
		"message": "",
	})
}

// @Summary 删除点赞收藏列表
// @Description 删除点赞收藏记录，单用户与单色图唯一。
// @Tags 色图
// @Security ApiKeyAuth
// @Accept json
// @Param seximg_id query string true "删除点赞收藏的 seximg id"
// @Router /api/user/seximg/star/delete [delete]
func SeximgStarDelete(c *gin.Context) {
	uID, err := strconv.Atoi(c.MustGet("UserID").(string))
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	seximgId := ""
	if pbid, ok := c.GetQuery("seximg_id"); ok {
		bid, _ := strconv.Atoi(pbid)
		seximgId = pbid
		resultCode := seximg.SeximgDeleteUserStarBySeximgId(uint(uID), uint(bid))
		if resultCode > 0 { //没有数据不需要删除列表数量
			seximg.SeximgDeleteFavorites(uint(bid))
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "参数错误,请升级版本",
		})
		return
	}
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	cacheKey := "middle-cgr-" + appID + ":api:seximg:info?id=" + seximgId
	redis.Pull(cacheKey)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "删除完成",
	})

}
