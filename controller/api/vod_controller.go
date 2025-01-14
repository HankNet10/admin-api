package api

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"myadmin/model"
	"myadmin/model/vod"
	"myadmin/util/redis"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 视频分类
// @Description 返回视频分类数据
// @Tags 长视频
// @Accept json
// @Router /api/vod/type [get]
func TypeList(c *gin.Context) {
	list := vod.VodTypeModel.AllType()

	//
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

// @Summary 视频首页推荐
// @Description 返回推荐的标签列表
// @Tags 长视频
// @Accept json
// @Router /api/vod/index [get]
func VodIndex(c *gin.Context) {
	list := vod.VodLabelModel.AllIndex()
	data := make([]apiResultVodLabel, len(list))
	for i, vl := range list {
		data[i] = apiResultVodLabel{
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

// @Summary 视频列表
// @Description 返回视频列表数据
// @Tags 长视频
// @Accept json
// @Param data query vod.VodListParam true "参数列表"
// @Router /api/vod/list [get]
func VodList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request vod.VodListParam
	c.ShouldBindQuery(&request)
	request.AppId, _ = strconv.Atoi(appID)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Order == "" {
		request.Order = "-id"
	}
	types := c.Query("types")
	if types != "" {
		request.Types = strings.Split(types, ",")
	}
	labels := c.Query("labels")
	if labels != "" {
		request.Labels = strings.Split(labels, ",")
	}
	users := c.Query("users")
	if users != "" {
		request.Users = strings.Split(users, ",")
	}
	// 处理视频缓存key
	selectQ := c.Query("types") + "_" + c.Query("labels") + "_" + c.Query("users")
	mustQ := strconv.Itoa(request.Page) + "_" + strconv.Itoa(request.Limit) + "_" + request.Order
	cKey := "api:vod:list-" + appID + ":" + selectQ + "_" + mustQ
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}

	list, total := vod.VodListModel.SelectList(request)
	ndata := make([]apiResultVodList, len(list))
	for i, v := range list {
		ndata[i] = newapiResultVodList(*v)
		if appID == "1" {
			ndata[i].Views = ndata[i].Views*10 + uint(rand.Intn(9))
		} else if appID == "2" {
			ndata[i].Views = ndata[i].Views*5 + uint(rand.Intn(9))
		}
	}

	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"list":  ndata,
			"total": total,
		},
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 25*time.Minute)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 专题视频列表
func VodTopicList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request vod.VodListParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Order == "" {
		request.Order = "-id"
	}
	list, total := vod.VodListModel.SelectTopicList(request)
	ndata := make([]apiResultTopicVodList, len(list))
	for i, v := range list {
		ndata[i] = newapiResultTopicVodList(*v)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  ndata,
			"total": total,
		},
		"message": "",
	})
}

// @Summary 视频详情
// @Description 返回视频详情数据
// @Tags 长视频
// @Accept json
// @Param id  query int true "视频ID"
// @Router /api/vod/info [get]
func VodInfo(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	vodID := c.Query("id")
	if vodID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}

	// 处理递增缓存
	cKey := "api:vod:info-" + appID + ":" + vodID
	newValue, err := redis.Redis.Incr(context.TODO(), cKey+"-incr").Result()
	if err == nil {
		if newValue < 200 { // 播放两百次则穿透缓存
			if jsonData, err := redis.Redis.Get(context.TODO(), cKey).Result(); err == nil {
				c.String(http.StatusOK, string(jsonData))
				return
			}
		} else {
			redis.Redis.Del(context.TODO(), cKey+"-incr")
		}
	} else {
		newValue = 1
	}
	rid, err := strconv.Atoi(vodID)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	// 开始拼装返回数据
	vodInfo := vod.VodListModel.SelectInfo(uint(rid))
	if vodInfo == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "视频不存在",
		})
		return
	}
	// 添加播放次数
	model.DataBase.Model(&vodInfo).Update("views", vodInfo.Views+uint(newValue))
	// 组合用户数据
	users := make([]apiResultUserShow, len(vodInfo.UserList))
	for i, v := range vodInfo.UserList {
		users[i] = newApiResultUserShow(v.User)
	}
	// 组合标签数据
	labels := make([]apiResultVodLabel, len(vodInfo.LabelList))
	for i, v := range vodInfo.LabelList {
		labels[i] = apiResultVodLabel{
			ID:   v.Label.ID,
			Name: v.Label.Name,
		}
	}
	// 组合返回数据
	data := apiResultVodInfo{
		ID:        vodInfo.ID,
		Title:     vodInfo.Title,
		Number:    vodInfo.Number,
		PlayUrl:   os.Getenv("PLIST_DOMAIN") + "/play/" + strconv.Itoa(int(vodInfo.ID)) + "/vod.plist",
		Favorites: vodInfo.Favorites,
		Views:     vodInfo.Views,
		Comments:  vodInfo.Comments,
		Cover:     os.Getenv("ALI_OSS_DOMAIN") + "/" + vodInfo.Cover,
		Users:     users,
		Labels:    labels,
		CreatedAt: vodInfo.CreatedAt.Unix(),
	}
	if appID == "1" {
		data.Favorites = data.Favorites*5 + uint(rand.Intn(9))
		data.Views = data.Views*10 + uint(rand.Intn(9)) + uint(newValue)
	} else if appID == "2" {
		data.Favorites = data.Favorites*3 + uint(rand.Intn(9))
		data.Views = data.Views*5 + uint(rand.Intn(9)) + uint(newValue)
	}

	// 返回组合数据并缓存
	resultData := gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 2*time.Hour)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 视频搜索
// @Description 返回视频详情数据
// @Tags 长视频
// @Accept json
// @Param page  query int true "视频搜索文字"
// @Param limit  query int true "视频搜索文字"
// @Param wd  query string true "视频搜索文字"
// @Router /api/vod/clever [get]
func VodClever(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	wd := c.Query("wd")
	if wd == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "填写搜索内容",
		})
		return
	}
	if len(wd) == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "最少输入2个字符查询哦",
		})
		return
	}
	redis.AddZsetValue("search", wd, 1.0)
	cacheKey := "searchvod" + appID + ":api:" + strconv.Itoa(page) + ":" + strconv.Itoa(limit) + ":" + wd
	if jsonData, err := redis.Get(cacheKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	var no = regexp.MustCompile("^[a-zA-Z]{2,6}(-|[0-9])[0-9]")
	var request vod.VodListParam
	request.Page = page
	if request.Page < 1 {
		request.Page = 1
	}
	request.Limit = limit
	if limit > 1 || limit < 20 {
		request.Limit = 20
	}
	wdLen := no.FindStringIndex(wd)
	if len(wdLen) > 0 { //判断是否是番号。
		request.Number = wd
	} else {
		request.Title = wd
	}
	request.Order = "-id"
	list, total := vod.VodListModel.SelectList(request)

	data := make([]apiResultVodList, len(list))
	for i, v := range list {
		data[i] = newapiResultVodList(*v)
		if appID == "1" {
			data[i].Views = data[i].Views*10 + uint(rand.Intn(9))
		} else if appID == "2" {
			data[i].Views = data[i].Views*5 + uint(rand.Intn(9))
		}
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
	redis.Set(cacheKey, jsonData, 120*time.Minute)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 视频详情推荐
// @Description 返回视频详情数据
// @Tags 长视频
// @Accept json
// @Param id  query int true "视频ID"
// @Router /api/vod/recommend [get]
func VodRecommend(c *gin.Context) {
	vodID := c.Query("id")
	if vodID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	if _, err := strconv.Atoi(vodID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "请求超时",
		})
		return
	}

	list, total := vod.VodListModel.SelectList(vod.VodListParam{
		Page:  rand.Intn(2500) + 1,
		Limit: 10,
	})

	data := make([]apiResultVodList, len(list))
	for i, v := range list {
		data[i] = newapiResultVodList(*v)
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

// @Summary 24小时上传的视频数量
// @Description 24小时上传的视频数量
// @Tags 长视频
// @Accept json
// @Router /api/vod/24upload [get]
func Vod24Upload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": vod.VodListModel.UploadAfter24(),
		},
		"message": "",
	})
}

// @Summary 播放热榜
// @Description 根据type统计各种播放热榜100条
// @Tags 长视频
// @Accept json
// @Param type  query string true "type： day 当天热榜1小时更新,week 一周热榜,month 月榜，own 总榜"
// @Router /api/vod/hotspot [get]
func Hotspot(c *gin.Context) {
	datetype, exist := c.GetQuery("type")
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "",
		})
		return
	}
	var day int
	switch datetype {
	case "day":
		day = 1
	case "week":
		day = 7
	case "month":
		day = 30
	case "own":
		day = 0
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "",
		})
		return
	}

	var request = vod.VodListParam{
		Order: "-view",
		Limit: 100,
	}
	if day > 0 {
		now := time.Now()
		request.AfterCreatedAt = now.AddDate(0, 0, -day)
	}

	list, _ := vod.VodListModel.SelectList(request)
	data := make([]apiResultVodList, len(list))
	for i, v := range list {
		data[i] = newapiResultVodList(*v)
		if appID == "1" {
			data[i].Views = data[i].Views*10 + uint(rand.Intn(9))
		} else if appID == "2" {
			data[i].Views = data[i].Views*5 + uint(rand.Intn(9))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

// @Summary 创建播放记录
// @Description 创建视频播放记录，单用户与单视频唯一。
// @Tags 长视频
// @Security ApiKeyAuth
// @Accept json
// @Param data body UserVodHistroyAddRequest true "参数列表"
// @Router /api/user/vod/history/add [put]
func VodHistoryAdd(c *gin.Context) {
	var request UserVodHistroyAddRequest
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
	vod.VodAddUserHistroy(uint(uID), request.VodID, request.VodTime)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}

// @Summary 获取播放记录
// @Description 创建视频播放记录，单用户与单视频唯一。
// @Tags 长视频
// @Security ApiKeyAuth
// @Param data query vod.VodHistoryParam true "参数列表"
// @Accept json
// @Router /api/user/vod/history/list [get]
func VodHistoryList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var request vod.VodHistoryParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := vod.VodListUserHistory(uint(uID), request)
	data := make([]apiResultVodHistory, len(list))
	for i, v := range list {
		data[i] = apiResultVodHistory{
			Id:        v.ID,
			LastSeen:  v.VideoTime,
			UpdatedAt: v.UpdatedAt.Unix(),
			Vod:       newapiResultVodList(v.Vod),
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

// @Summary 删除播放记录
// @Description 删除播放记录 id 为播放记录id id=0 则为 删除所有
// @Tags 长视频
// @Security ApiKeyAuth
// @Param id query string true "删除播放历史的id"
// @Accept json
// @Router /api/user/vod/history/delete [delete]
func VodHistoryDelete(c *gin.Context) {
	hid, _ := strconv.Atoi(c.Query("id"))
	if hid < 0 {
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

	b := vod.VodDeleteUserHistroy(uint(uID), uint(hid))
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

// @Summary 创建收藏记录
// @Description 创建视频播放记录，单用户与单视频唯一。
// @Tags 长视频
// @Security ApiKeyAuth
// @Accept json
// @Param data body UserVodStarAddRequest true "参数列表"
// @Router /api/user/vod/star/add [put]
func VodStarAdd(c *gin.Context) {
	var request UserVodStarAddRequest
	c.ShouldBindJSON(&request)

	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	vod.VodAddUserStar(uint(uID), request.VodID)
	vod.VlogListAddFavorites(request.VodID)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}

// @Summary 获取收藏记录
// @Description 获取某个用户的收藏列表
// @Tags 长视频
// @Security ApiKeyAuth
// @Accept json
// @Param data query vod.VodHistoryParam true "参数列表"
// @Router /api/user/vod/star/list [get]
func VodStarList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var request vod.VodStarParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := vod.VodListUserStar(uint(uID), request)
	data := make([]apiResultVodStar, len(list))
	for i, v := range list {
		data[i] = apiResultVodStar{
			Id:        v.ID,
			UpdatedAt: v.UpdatedAt.Unix(),
			Vod:       newapiResultVodList(v.Vod),
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

// @Summary 删除收藏列表
// @Description 创建视频播放记录，单用户与单视频唯一。
// @Tags 长视频
// @Security ApiKeyAuth
// @Param data query vod.VodHistoryParam true "参数列表"
// @Accept json
// @Param id query string true "删除播放历史的id"
// @Param vod_id query string true "删除收藏的 vod id"
// @Router /api/user/vod/star/delete [delete]
func VodStarDelete(c *gin.Context) {
	uID, err := strconv.Atoi(c.MustGet("UserID").(string))
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var resStatus bool
	if pid, ok := c.GetQuery("id"); ok {
		sid, _ := strconv.Atoi(pid)
		resStatus = vod.VodDeleteUserStar(uint(uID), uint(sid))
	} else if pvid, ok := c.GetQuery("vod_id"); ok {
		vid, _ := strconv.Atoi(pvid)
		resStatus = vod.VodDeleteUserStarByVodId(uint(uID), uint(vid))
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "参数错误",
		})
		return
	}
	if resStatus {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "删除完成",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "删除失败",
		})
	}
}

// @Summary 全部标签
// @Description 返回视频标签数据
// @Tags 长视频
// @Accept json
// @Router /api/vodlabel/all [get]
func VodLabelAll(c *gin.Context) {
	cKey := "api:vodlabel:all"
	data := struct {
		Total int64               `json:"total"`
		List  []apiResultVodLabel `json:"list"`
	}{}

	if err := redis.Deserialize(cKey, &data); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    data,
			"message": "",
		})
		return
	}

	list, total := vod.VodLabelModel.List(vod.VodLabelParam{
		Page:  1,
		Limit: 999,
	})

	labels := make([]apiResultVodLabel, len(list))
	for i, v := range list {
		labels[i] = apiResultVodLabel{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	data.List = labels
	data.Total = total

	if err := redis.Serialize(cKey, data, 2*time.Hour); err != nil {
		log.Panic("Redis set err:", cKey, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

// @Summary 关注的视频列表
// @Description 返回视频列表数据
// @Tags 长视频
// @Accept json
// @Param data query vod.VodListParam true "参数列表"
// @Router /api/vod/list/follow [get]
func VodFollowList(c *gin.Context) {
	var request vod.VodListParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}

	if request.Order == "" {
		request.Order = "-id"
	}

	userID, err := strconv.Atoi(c.MustGet("UserID").(string))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	list, total := vod.VodListModel.FollowList(uint(userID), request)

	data := make([]apiResultVodList, len(list))
	for i, v := range list {
		data[i] = newapiResultVodList(*v)
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

func interfaceToString(i interface{}) string {
	// 使用类型断言将接口值转换为字符串
	if str, ok := i.(string); ok {
		return str
	}
	return ""
}
func HotSearchVodList(c *gin.Context) {
	lists := redis.HotSeachList(int64(1), int64(8))
	results := make([]interface{}, len(lists))
	for i, v := range lists {
		results[i] = struct {
			Text string `json:"text"`
		}{interfaceToString(v.Member)}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  results,
			"total": 0,
		},
		"message": "",
	})
}
