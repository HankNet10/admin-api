package api

import (
	"encoding/json"
	"log"
	"myadmin/model/blog"
	"myadmin/util/redis"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 短视频列表
// @Description 返回短视频列表
// @Tags 短视频
// @Accept json
// @Param data query blog.BlogListParam true "参数列表"
// @Router /api/vlog/list [get]
func VlogList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request blog.BlogListParam
	c.ShouldBindQuery(&request)
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Order == "" {
		request.Order = "-id"
	}
	request.AppId, _ = strconv.Atoi(appID)
	// 处理博客缓存key
	selectQ := c.Query("user_id")
	mustQ := strconv.Itoa(request.Page) + "_" + strconv.Itoa(request.Limit) + "_" + request.Order +
		"_" + request.Type
	cKey := "api:vlog:list-" + appID + ":" + selectQ + "_" + mustQ
	if request.Star > 0 {
		cKey += "-s"
	}
	if request.Matchid > 0 {
		cKey += "-m" + strconv.Itoa(request.Matchid)
	}
	if request.Topicid > 0 {
		cKey += "-t" + strconv.Itoa(request.Topicid)
	}
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	request.Type = "1"
	list, total := blog.BlogListModel.BlogSelectList(request)

	data := make([]apiResultVlogList, len(list))
	var host = os.Getenv("PLIST_DOMAIN")
	for i, v := range list {
		data[i] = newApiResultVlogList(*v)
		data[i].PlayUrl = host + data[i].PlayUrl
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
	redis.Set(cKey, jsonData, 10*time.Minute)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 获取收藏记录
// @Description 获取某个用户的收藏列表
// @Tags 短视频 已经和blog 融合
func VlogStarList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	var request blog.BlogStarParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 50 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Type = 1
	list, total := blog.BlogListUserStar(uint(uID), request)
	data := make([]apiResultVlogStar, len(list))
	host := os.Getenv("PLIST_DOMAIN")
	for i, v := range list {
		vlog := newApiResultVlogList(v.Blog)
		vlog.PlayUrl = host + vlog.PlayUrl
		// vlog.Cover = host + "/" + vlog.Cover
		data[i] = apiResultVlogStar{
			Id:        v.ID,
			UpdatedAt: v.UpdatedAt.Unix(),
			Vlog:      vlog,
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
