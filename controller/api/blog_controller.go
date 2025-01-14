package api

import (
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/model/user"
	"myadmin/util/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用户up主发布的社区 还有 审核中 已拒绝 状态
func UserUpBlogList(c *gin.Context) {
	var request blog.BlogListParam
	c.ShouldBindQuery(&request)
	userID := c.MustGet("UserID").(string)
	request.UserID = userID
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Order == "" {
		request.Order = "-id"
	}
	request.Self = 1
	list, total := blog.BlogListModel.BlogSelectList(request)

	data := make([]apiResultBlogList, len(list))
	for i, v := range list {
		data[i] = newApiUpResultBlogList(*v)
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
func BlogSearch(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	var request blog.BlogListParam
	request.Page = page
	if request.Page < 1 {
		request.Page = 1
	}
	request.Limit = limit
	if limit > 1 || limit < 20 {
		request.Limit = 20
	}
	if request.Order == "" {
		request.Order = "-id"
	}
	wd := c.Query("wd")
	request.Wd = wd
	cKey := "api:blog:search-" + appID + ":" + strconv.Itoa(page) + ":" + strconv.Itoa(limit) + ":" + wd
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	list, total := blog.BlogListModel.BlogSelectList(request)
	data := make([]apiResultBlogList, len(list))
	for i, v := range list {
		data[i] = newApiResultBlogList(*v)
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
	redis.Set(cKey, jsonData, 60*time.Minute)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 图文列表
func BlogList(c *gin.Context) {
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
	top := c.Query("top")
	mustQ := strconv.Itoa(request.Page) + ":" + strconv.Itoa(request.Limit) + ":" + request.Order +
		":" + request.Type
	cKey := "api:blog:list:" + appID + ":" + selectQ + ":" + mustQ + ":" + top + ":" + request.Time
	if request.Star > 0 {
		cKey += ":s"
	}
	if request.Matchid > 0 {
		cKey += ":m" + strconv.Itoa(request.Matchid)
	}
	if request.Topicid > 0 {
		cKey += ":t" + strconv.Itoa(request.Topicid)
	}
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}

	list, total := blog.BlogListModel.BlogSelectList(request)

	data := make([]apiResultBlogList, len(list))
	for i, v := range list {
		data[i] = newApiResultBlogList(*v)
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

// @Summary 图文详情
// @Description 返回图文详情数据
// @Tags 图文
// @Accept json
// @Param id  query int true "视频ID"
// @Router /api/blog/info [get]
func BlogInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	uid, _ := strconv.Atoi(id)
	data := blog.BlogListModel.SelectInfo(uint(uid))

	var rdata apiResultBlogList
	if data != nil {
		rdata = newApiResultBlogList(*data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    rdata,
		"message": "",
	})
}

// 长文
func BlogArticleInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	uid, _ := strconv.Atoi(id)
	data := blog.BlogListModel.SelectArticle(uint(uid))

	if data != nil {
		rdata := struct {
			Essay string `json:"essay"`
		}{data.Essay}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    rdata,
			"message": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "",
	})
}

func BlogIsLike(c *gin.Context) {
	blogId := c.Query("blog_id")
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	userID := c.MustGet("UserID").(string)
	BlogStarModel := blog.BlogStar{}
	result := model.DataBase.Where("user_id = ? and blog_id = ?", userID, blogId).Limit(1).Find(&BlogStarModel)
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
	if BlogStarModel.ID != 0 {
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

// @Summary 创建收藏记录
// @Description 创建视频播放记录，单用户与单视频唯一。
// @Tags 图文
// @Security ApiKeyAuth
// @Accept json
// @Param data body UserBlogStarAddRequest true "参数列表"
// @Router /api/user/blog/star/add [put]
func BlogStarAdd(c *gin.Context) {
	var request UserBlogStarAddRequest
	c.ShouldBindJSON(&request)
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	userID := c.MustGet("UserID").(string)
	uID, _ := strconv.Atoi(userID)

	cacheKey := "middle-cgr-" + appID + ":api:blog:info?id=" + strconv.FormatUint(uint64(request.BlogID), 10)
	redis.Pull(cacheKey)

	var blodId uint
	if request.BlogID > 0 {
		blodId = request.BlogID
	} else if request.VlogId > 0 {
		blodId = request.VlogId
	}
	reuslt := blog.BlogAddUserStar(uint(uID), blodId)
	if reuslt == 0 {
		blog.BlogListAddFavorites(blodId)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}

// @Summary 获取点赞记录
// @Description 获取某个用户的点赞记录
// @Security ApiKeyAuth
// @Accept json
// @Param data query blog.BlogStarParam true "参数列表"
// @Router /api/user/blog/star/list [get]
func BlogStarList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, _ := strconv.Atoi(userID)

	var request blog.BlogStarParam
	c.ShouldBindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}

	list, total := blog.BlogListUserStar(uint(uID), request)
	data := make([]apiResultBlogList, len(list))
	for i, v := range list {
		data[i] = newApiResultBlogList(v.Blog)
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
// @Tags 图文
// @Security ApiKeyAuth
// @Accept json
// @Param id query string true "删除播放历史的id"
// @Param blog_id query string true "删除收藏的 blog id"
// @Router /api/user/blog/star/delete [delete]
func BlogStarDelete(c *gin.Context) {
	uID, _ := strconv.Atoi(c.MustGet("UserID").(string))
	blogId := ""
	if pbid, ok := c.GetQuery("blog_id"); ok {
		bid, _ := strconv.Atoi(pbid)
		blogId = pbid
		resultCode := blog.BlogDeleteUserStarByBlogId(uint(uID), uint(bid))
		if resultCode > 0 { //没有数据不需要删除列表数量
			blog.BlogListDeleteFavorites(uint(bid))
		}
	} else if vbid, ok := c.GetQuery("vlog_id"); ok {
		bid, _ := strconv.Atoi(vbid)
		blogId = vbid
		resultCode := blog.BlogDeleteUserStarByBlogId(uint(uID), uint(bid))
		if resultCode > 0 { //没有数据不需要删除列表数量
			blog.BlogListDeleteFavorites(uint(bid))
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
	cacheKey := "middle-cgr-" + appID + ":api:blog:info?id=" + blogId
	redis.Pull(cacheKey)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "删除完成",
	})

}

// 用户创建博客
func BlogCreateByUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.MustGet("UserID").(string))
	var request UserBlogCreateUserRequest
	c.ShouldBindJSON(&request)
	//查询是否有发布权限
	blackIp := user.BlackIp{}.SelectBlackIpByIP(c.ClientIP())
	if blackIp.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "状态异常",
		})
		return
	}

	m := user.UserUploader{}
	model.DataBase.Where("user_id = ?", userID).First(&m)
	if m.ID == 0 || m.State != 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "状态异常",
		})
		return
	}
	ckey := "createsqcount:" + c.MustGet("UserID").(string)
	cCount, _ := redis.Get(ckey) //控制发布次数10次 24小时
	vipuser := vipUserListModel.VIPUserSelectByUserId(uint(userID))
	requestCount := 0
	if cCount != "" {
		count, _ := strconv.Atoi(cCount)
		requestCount = count
		if vipuser.ID <= 0 {
			if count >= 10 {
				c.JSON(http.StatusOK, gin.H{
					"code":    400,
					"data":    nil,
					"message": "发布次数过多",
				})
				return
			}
		}
	}
	if request.Title == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "请输入标题",
		})
		return
	}
	//创建博客
	blogModel := blog.BlogList{}
	blogModel.Detail = request.Title
	blogModel.Essay = request.Essay
	blogModel.UserID = uint(userID)
	if request.MatchId > 0 {
		blogModel.MatchId = request.MatchId
	}
	if request.Topicid > 0 {
		blogModel.TopicId = request.Topicid
	}
	blogModel.Type = request.Type
	//社区帖子不自动发布
	// if vipuser.ID > 0 {
	// 	blogModel.Status = 1
	// } else {
	blogModel.Status = 2
	// }
	blogModel.Save()
	if blogModel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "创建失败",
		})
		return
	}
	for _, item := range request.Images {
		Image := blog.BlogImage{
			UserID: uint8(userID),
			BlogID: blogModel.ID,
			Path:   item,
		}
		Image.Save()
	}
	if request.VideoPath != "" {
		Video := blog.BlogVideo{
			UserID:  uint8(userID),
			BlogID:  blogModel.ID,
			OssName: request.VideoPath,
		}
		Video.Save()
	}
	redis.SetNoChangeTTl(ckey, requestCount+1, 24*time.Hour)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "已投稿,等待审核通过",
	})
}

func BlogUserDelete(c *gin.Context) {
	uID, _ := strconv.Atoi(c.MustGet("UserID").(string))
	blogId, ok := c.GetQuery("blog_id")
	blog := blog.BlogList{}
	if ok {
		dresult := model.DataBase.Where("user_id = ? AND id = ?", uID, blogId).First(&blog)
		if dresult.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "参数错误",
			})
			return
		}
		if blog.MatchId != 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "参加比赛的帖子,请私信联系客服删除",
			})
			return
		}
		model.DataBase.Model(&blog).Update("status", 4) //改为已删除
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "参数错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "已删除",
	})

}
