package api

import (
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/sexnovel"
	"myadmin/util/redis"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var sexnovelTypeModel sexnovel.SexnovelType

// @Summary 小说分类列表
// @Tags 小说
// @Router /api/sexnoveltype/list [get]
func SexnovelTypeList(c *gin.Context) {
	list := sexnovelTypeModel.AllType()
	dataMap := make(map[uint][]apiResultTypeListChild)

	for _, v := range list {
		child := apiResultTypeListChild{
			ID:   v.ID,
			Name: v.Name,
		}
		if dataMap[v.Parent] != nil {
			dataMap[v.Parent] = append(dataMap[v.Parent], child)
		} else {
			dataMap[v.Parent] = []apiResultTypeListChild{child}
		}
	}

	data := make([]apiResultTypeList, len(dataMap[0]))
	for i, v := range dataMap[0] {
		data[i] = apiResultTypeList{
			ID:    v.ID,
			Name:  v.Name,
			Child: dataMap[v.ID],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": data,
		},
		"message": "",
	})
}

// @Summary 小说推荐
// @Description 返回推荐的标签列表
// @Tags 小说
// @Accept json
// @Router /api/sexnovel/label [get]
func SexnovelLabel(c *gin.Context) {
	list := sexnovel.SexnovelLabelModel.AllLabel()
	data := make([]apiResultSexnovelLabel, len(list))
	for i, vl := range list {
		data[i] = apiResultSexnovelLabel{
			ID:   vl.ID,
			Name: vl.Name,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": data,
		},
		"message": "",
	})
}

// @Summary 小说列表
// @Tags 小说
// @Param data query sexnovel.SexnovelParam true "参数列表"
// @Router /api/sexnovel/list [get]
func SexnovelList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request sexnovel.SexnovelParam
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
	// 处理小说缓存key
	selectQ := c.Query("user_id")
	top := c.Query("top")
	labelsId := c.Query("labelsId")
	mustQ := strconv.Itoa(request.Page) + ":" + strconv.Itoa(request.Limit) + ":" + request.Order +
		":" + strconv.Itoa(request.Typeid) + ":" + labelsId
	cKey := "api:sexnovel:list:" + appID + ":" + selectQ + ":" + mustQ + ":" + top

	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}

	list, total := sexnovel.SexnovelModel.SexnovelSelectList(request)

	data := make([]apiResultSexnovel, len(list))
	for i, v := range list {
		data[i] = newApiResultSexnovel(*v)
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

// @Summary 小说章节列表
// @Tags 小说
// @Param data query sexnovel.SexnovelChapterParam true "参数列表"
// @Router /api//sexnovel/chapte/list [get]
func SexnovelChapterList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request sexnovel.SexnovelChapterParam
	c.ShouldBindQuery(&request)
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Page < 1 {
		request.Page = 1
	}
	// 处理小说缓存key
	selectQ := c.Query("user_id")
	pid := c.Query("pid")
	mustQ := strconv.Itoa(request.Page) + ":" + strconv.Itoa(request.Limit)
	cKey := "api:sexnovel:chapterlist:" + appID + ":" + selectQ + ":" + pid + ":" + mustQ

	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}

	list, total := sexnovel.SexnovelChapterModel.List(request)

	data := make([]apiResultSexnovelChapterShow, len(list))
	for i, v := range list {
		data[i] = newApiResultSexnovelChapterShow(*v)
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

// @Summary 小说详情
// @Description 返回小说详情数据
// @Tags 小说
// @Accept json
// @Param id query string true "小说详情的id"
// @Router /api/sexnovel/info [get]
func SexnovelInfo(c *gin.Context) {
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
	cKey := "api:sexnovel:info:" + appID + ":" + id

	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	sid, _ := strconv.Atoi(id)
	data := sexnovel.SexnovelModel.SelectInfo(uint(sid))

	var rdata apiResultSexnovel
	if data != nil {
		sexnovel.SexnovelAddLooks(uint(sid))
		rdata = newApiResultSexnovel(*data)
	}
	resultData := gin.H{
		"code":    200,
		"data":    rdata,
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 2*time.Hour)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 小说内容
// @Description 返回小说内容数据
// @Tags 小说
// @Accept json
// @Param id query string true "小说内容的章节id"
// @Router /api/sexnovel/content/info [get]
func SexnovelContent(c *gin.Context) {
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
	islong, _ := strconv.ParseBool(c.Query("islong"))
	cKey := "api:sexnovel:content:" + appID + ":" + c.Query("islong") + ":" + id
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	data := sexnovel.SexnovelContentModel.SelectContentInfo(islong, uint(sid))
	resultData := gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 72*time.Hour)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 创建小说观看记录
// @Description 创建小说观看记录，单用户与单小说唯一。
// @Tags 小说
// @Security ApiKeyAuth
// @Accept json
// @Param data body UserSexnovelHistroyAddRequest true "参数列表"
// @Router /api/user/sexnovel/history/add [put]
func SexnovelHistoryAdd(c *gin.Context) {
	var request UserSexnovelHistroyAddRequest
	c.ShouldBindJSON(&request)

	// 处理播放记录

	uID, _ := strconv.Atoi(c.MustGet("UserID").(string))
	if uID < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "记录完成",
		})
		return
	}
	sexnovel.SexnovelAddUserHistroy(uint(uID), request.SexnovelID)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}

// @Summary 获取小说观看记录
// @Description 创建小说观看记录，单用户与单小说唯一。
// @Tags 小说
// @Security ApiKeyAuth
// @Param data query sexnovel.SexnovelHistoryParam true "参数列表"
// @Accept json
// @Router /api/user/sexnovel/history/list [get]
func SexnovelHistoryList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var request sexnovel.SexnovelHistoryParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := sexnovel.SexnovelListUserHistory(uint(uID), request)
	data := make([]apiResultSexnovelHistory, len(list))
	for i, v := range list {
		data[i] = apiResultSexnovelHistory{
			Id:        v.ID,
			UpdatedAt: v.UpdatedAt.Unix(),
			Sexnovel:  newapiResultSexnovelList(v.Sexnovel),
		}
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

// @Summary 删除小说观看记录
// @Description 删除小说观看记录 id 为观看记录id:'1,2,3' id=0 则为 删除所有
// @Tags 小说
// @Security ApiKeyAuth
// @Param id query string true "删除观看历史的id"
// @Accept json
// @Router /api/user/sexnovel/history/delete [delete]
func SexnovelHistoryDelete(c *gin.Context) {
	var hids []uint
	ids := c.Query("id")
	ids, _ = url.QueryUnescape(ids)
	strSlice := strings.Split(ids, ",")
	for _, s := range strSlice {
		u, _ := strconv.ParseUint(s, 10, 64)
		hids = append(hids, uint(u))
	}
	//hid, _ := strconv.Atoi(c.Query("id"))
	if len(hids) < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "删除失败",
		})
		return
	}
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	b := sexnovel.SexnovelDeleteUserHistroy(uint(uID), hids)
	if b {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "删除完成",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "删除失败",
		})
	}
}

// @Summary 是否点赞收藏
// @Description 返回isLike
// @Tags 小说
// @Param sexnovel_id query string true "小说的id"
// @Router /api/user/sexnovel/star/isLike [get]
func SexnovelIsLike(c *gin.Context) {
	sexnovelId := c.Query("sexnovel_id")
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
	SexnovelStarModel := sexnovel.SexnovelStar{}
	result := model.DataBase.Where("user_id = ? and sexnovel_id = ?", rid, sexnovelId).Limit(1).Find(&SexnovelStarModel)
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
	if SexnovelStarModel.ID != 0 {
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
// @Tags 小说
// @Security ApiKeyAuth
// @Accept json
// @Param data body UserSexnovelStarAddRequest true "参数列表"
// @Router /api/user/sexnovel/star/add [put]
func SexnovelStarAdd(c *gin.Context) {
	var request UserSexnovelStarAddRequest
	c.ShouldBindJSON(&request)
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var sexnovelId uint
	if request.SexnovelID > 0 {
		sexnovelId = request.SexnovelID
	}
	reuslt := sexnovel.SexnovelAddUserStar(uint(uID), sexnovelId)
	if reuslt == 0 {
		sexnovel.SexnovelAddFavorites(sexnovelId)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}

// @Summary 获取点赞收藏记录
// @Description 获取某个用户的点赞收藏记录
// @Tags 小说
// @Security ApiKeyAuth
// @Accept json
// @Param data query sexnovel.SexnovelStarParam true "参数列表"
// @Router /api/user/sexnovel/star/list [get]
func SexnovelStarList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var request sexnovel.SexnovelStarParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}

	list, total := sexnovel.SexnovelUserStar(uint(uID), request)
	data := make([]apiResultSexnovel, len(list))
	for i, v := range list {
		data[i] = newApiResultSexnovel(v.Sexnovel)
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
// @Description 创建小说点赞收藏记录，单用户与单小说唯一。
// @Tags 小说
// @Security ApiKeyAuth
// @Accept json
// @Param sexnovel_id query string true "删除点赞收藏的 sexnovel id"
// @Router /api/user/sexnovel/star/delete [delete]
func SexnovelStarDelete(c *gin.Context) {
	uID, err := strconv.Atoi(c.MustGet("UserID").(string))
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var hids []uint
	ids := c.Query("sexnovel_id")
	ids, _ = url.QueryUnescape(ids)
	strSlice := strings.Split(ids, ",")
	for _, s := range strSlice {
		u, _ := strconv.ParseUint(s, 10, 64)
		hids = append(hids, uint(u))
	}
	if len(hids) < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "删除失败",
		})
		return
	} else {
		resultCode := sexnovel.SexnovelDeleteUserStarBySexnovelId(uint(uID), hids)
		if resultCode { //没有数据不需要删除列表数量
			sexnovel.SexnovelDeleteFavorites(hids)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "删除完成",
	})

}
