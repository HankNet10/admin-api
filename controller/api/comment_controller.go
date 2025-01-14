package api

import (
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/model/seximg"
	"myadmin/model/user"
	"myadmin/model/vod"
	"myadmin/util/redis"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/feiin/sensitivewords"
	"github.com/gin-gonic/gin"
)

// @Summary 获取视频评论
// @Param data query user.CommentParam true "参数列表"
// @Router /api/blog/comment/list [get]
func CommentList(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request user.CommentParam
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
	var cid int
	var ctype int
	if request.SeximgID > 0 {
		cid = request.SeximgID
		ctype = 4
	} else if request.VlogID > 0 {
		cid = request.VlogID
		ctype = 3
	} else if request.BlogID > 0 {
		cid = request.BlogID
		ctype = 3
	} else if request.VodID > 0 {
		cid = request.VodID
		ctype = 1
	}
	request.AppID, _ = strconv.Atoi(appID)
	cacheKey := "api:commentlist:" + strconv.Itoa(cid) + ":" + strconv.Itoa(ctype) + ":" + strconv.Itoa(request.ParentId) + ":" + strconv.Itoa(request.Page) + ":" + strconv.Itoa(request.Limit) + ":" + request.Order
	if jsonData, err := redis.Get(cacheKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	// 获取敏感词 需要审核的话不用敏感词了
	// var listdirty []user.DirtyWord
	// query := model.DataBase.Model(user.DirtyWord{})
	// query.Find(&listdirty)
	// sensitive := sensitivewords.New()
	// for _, word := range listdirty {
	// 	sensitive.AddWord(word.Name)
	// }
	list, total := user.CommentModel.SelectList(request)
	data := make([]apiResultCommentList, len(list))
	for i, v := range list {
		var reUserName = ""
		if v.ReUserID != 0 && v.ReUser.Name != "" {
			reUserName = v.ReUser.Name
		}
		data[i] = apiResultCommentList{
			ID:         v.ID,
			CId:        v.CID,
			Comment:    v.Comment,
			User:       newApiResultUserShow(v.User),
			CreatedAt:  v.CreatedAt.Unix(),
			Count:      v.Count,
			ParentId:   v.ParentId,
			ReUserName: reUserName,
			ReUserID:   v.ReUserID,
			RID:        v.RID,
		}
	}
	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  data,
		},
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cacheKey, jsonData, 10*time.Minute)
	c.String(http.StatusOK, string(jsonData))
}

func CommentUserList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	var request user.CommentParam
	c.ShouldBindQuery(&request)
	request.UserId = userID
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Order == "" {
		request.Order = "-id"
	}
	list, total := user.CommentModel.SelectList(request)
	data := make([]apiResultCommentList, len(list))
	for i, v := range list {
		var reUserName = ""
		if v.ReUserID != 0 && v.ReUser.Name != "" {
			reUserName = v.ReUser.Name
		}
		data[i] = apiResultCommentList{
			ID:         v.ID,
			CId:        v.CID,
			CType:      v.Type,
			Comment:    v.Comment,
			User:       newApiResultUserShow(v.User),
			CreatedAt:  v.CreatedAt.Unix(),
			Count:      v.Count,
			ParentId:   v.ParentId,
			ReUserName: reUserName,
			ReUserID:   v.ReUserID,
			RID:        v.RID,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  data,
			"total": total,
		},
		"message": "评论成功",
	})
}

// @Summary 发布评论
// @Accept json
// @Security ApiKeyAuth
// @Router /api/blog/comment/add [post]
func CommentAdd(c *gin.Context) {
	userID, err := strconv.Atoi(c.MustGet("UserID").(string))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var request UserCommentAddRequest
	c.ShouldBindJSON(&request)
	if len(request.Comment) > 1500 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "评论字数太多了哦",
		})
		return
	}
	ckey := "commentip:" + c.ClientIP()
	cCount, _ := redis.Get(ckey) //控制24小时1个ip 15个评论
	requestCount := 0
	userinfo := user.SelectUserByID(uint(userID))
	if userinfo.DenyComment != 2 { //非优质用户
		if cCount != "" {
			count, _ := strconv.Atoi(cCount)
			requestCount = count
			if count > 15 {
				c.JSON(http.StatusOK, gin.H{
					"code":    400,
					"data":    nil,
					"message": "今日评论已达上限哦",
				})
				return
			}
		}
	}
	ckeyUser := "commentuser:" + strconv.Itoa(userID)
	cCountUser, _ := redis.Get(ckeyUser) //控制24小时1个用户 15个评论
	requestCountUser := 0
	if userinfo.DenyComment != 2 { //非优质用户
		if cCountUser != "" {
			cCountUser, _ := strconv.Atoi(cCountUser)
			requestCountUser = cCountUser
			if cCountUser > 15 {
				c.JSON(http.StatusOK, gin.H{
					"code":    400,
					"data":    nil,
					"message": "今日评论已达上限哦",
				})
				return
			}
		}
	}

	if userinfo.DenyComment == 1 { //用户已经被禁言
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "评论功能异常",
		})
		return
	}
	reg := regexp.MustCompile("^[a-zA-Z0-9, \\p{Han}]+$") //只有中文字和英文字母
	if !reg.MatchString(request.Comment) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "评论内容不能含有字符和特殊字母哦",
		})
		return
	}
	// if (time.Now().Unix() - userinfo.CreatedAt.Unix()) < 86400*2 { //注册第三天的上面才能评论
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":    400,
	// 		"data":    nil,
	// 		"message": "评论审核后即可显示哦",
	// 	})
	// 	return
	// }
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var cId uint = 0
	var cType uint = 1
	if request.SeximgID != 0 {
		cId = request.SeximgID
		cType = 4
	} else if request.BlogID != 0 {
		cId = request.BlogID
		cType = 3
	} else if request.VlogID != 0 {
		cId = request.VlogID
		cType = 3
	} else if request.VodID != 0 {
		cId = request.VodID
		cType = 1
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "评论异常",
		})
		return
	}
	appIdInt, _ := strconv.Atoi(appID)
	realModel := user.UserComment{
		UserID:  uint(userID),
		CID:     cId,
		Type:    cType,
		Comment: request.Comment,
		AppID:   uint(appIdInt),
	}
	if request.ReUserId != 0 && request.CommentId != 0 {
		realModel.ReUserID = request.ReUserId
		realModel.RID = request.CommentId
		model := user.CommentModel.SelectByID(uint(request.CommentId))
		if model.ParentId == 0 { //直接回复评论
			model.Count += 1
			model.Save()
			realModel.ParentId = request.CommentId //Pid与Rid相同
		} else { //回复评论下面的回复
			realModel.ParentId = model.ParentId //Pid与Rid不通
			modelParent := user.CommentModel.SelectByID(uint(model.ParentId))
			modelParent.Count += 1
			modelParent.Save()
		}
	}
	var status = 0 //默认都打开 根据黑名单ip判断是否打开
	var listBlackIp []user.BlackIp
	query := model.DataBase.Model(user.BlackIp{})
	query.Find(&listBlackIp)
	for _, Item := range listBlackIp {
		if Item.Ip == c.ClientIP() || strings.Contains(c.ClientIP(), Item.Ip) {
			status = 0
			break
		}
	}
	realModel.Status = uint(status) //默认都打开
	realModel.Ip = c.ClientIP()
	realModel.Save()
	switch cType { //添加评论数
	case 1:
		vod.VodListModel.AddComments(request.VodID)
	case 3:
		blog.BlogListModel.AddComments(cId)
	case 4:
		seximg.SeximgModel.AddComments(cId)
	}

	cacheKey := "api:commentlist:" + strconv.Itoa(int(cId)) + ":" + strconv.Itoa(int(cType)) + ":"
	redis.PullPrefix(cacheKey) //刷新列表缓存
	// redis.AddZsetValue("commenttext", request.Comment, 1.0)          //添加评论内容检测缓存 用来看是否重复刷屏内容
	// redis.AddZsetValue("commentuser", strconv.Itoa(userID), 1.0)     //添加用户ID检测缓存 用来看单个用户发了多少
	redis.SetNoChangeTTl(ckey, requestCount+1, 24*time.Hour)         //添加发布IP限制
	redis.SetNoChangeTTl(ckeyUser, requestCountUser+1, 24*time.Hour) //添加用户发布限制

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id": realModel.ID,
		},
		"message": "评论成功,审核后显示",
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": 200,
	// 	"data": gin.H{
	// 		"id": realModel.ID,
	// 	},
	// 	"message": "评论成功",
	// })
}
func CommentDlete(c *gin.Context) {
	uID, _ := strconv.Atoi(c.MustGet("UserID").(string))
	id, ok := c.GetQuery("id")
	comment := user.UserComment{}
	if ok {
		idd, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "参数错误",
			})
			return
		}
		dresult := model.DataBase.First(&comment, idd)
		if dresult.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "参数错误",
			})
			return
		}
		if comment.UserID != uint(uID) {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "参数错误",
			})
			return
		}
		model.DataBase.Model(&comment).Update("status", 0) //改为未审核
		cacheKey := "api:commentlist:" + strconv.Itoa(int(comment.CID)) + ":" + strconv.Itoa(int(comment.Type)) + ":"
		redis.PullPrefix(cacheKey) //删除缓存
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

// 通过接口给宝塔cron调用实现定时任务 删除违规评论
// 5分钟任务  有敏感词的删除 同时刷新缓存
func JobDleteComment(c *gin.Context) {
	var comments []user.UserComment
	model.DataBase.Where("created_at > ?", time.Now().Add(-5*time.Minute).Format("2006-01-02 15:04:05")).Where("status = ?", 1).Order("id desc").Limit(500).Find(&comments) //查到短视频的评论 都改成社区的评论ID
	var listdirty []user.DirtyWord
	query := model.DataBase.Model(user.DirtyWord{})
	query.Find(&listdirty)
	sensitive := sensitivewords.New()
	for _, word := range listdirty { //添加敏感词
		sensitive.AddWord(word.Name)
	}
	count := 0
	cIdArray := []string{} //用来刷新评论缓存

	for _, item := range comments {
		if sensitive.Check(item.Comment) { //包含敏感词的评论
			count++
			item.Status = 0
			item.Save()
			switch item.Type { //要扣一下视频社区评论数
			case 1:
				vod.VodListModel.DelComments(item.CID)
			case 3:
				blog.BlogListModel.DelComments(item.CID)
			case 4:
				seximg.SeximgModel.DelComments(item.CID)
			}
			//获得缓存key
			cIdArray = append(cIdArray, "api:commentlist:"+strconv.FormatUint(uint64(item.CID), 10)+":"+strconv.FormatUint(uint64(item.Type), 10)+":")
		}
	}
	for _, item := range cIdArray {
		cacheKey := item
		redis.PullPrefix(cacheKey) //删除修改的列表缓存
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "已完成,查到敏感词评论" + strconv.Itoa(count),
	})
}

// 通过接口给宝塔cron调用实现定时任务 删除违规评论
// 30分钟任务  查询评论超过重复5条的删除
func JobDleteComment1(c *gin.Context) {
	lists := redis.CommentRepeatList()
	count := 0
	cIdArray := []string{} //用来刷新评论缓存
	println(lists)
	for _, textitem := range lists {
		var comments []user.UserComment
		query := model.DataBase.Where("created_at > ?", time.Now().Add(-30*time.Minute).Format("2006-01-02 15:04:05"))
		query.Where("comment = ?", textitem.Member).Where("status = ?", 1).Order("id desc").Limit(500).Find(&comments)
		for _, item := range comments {
			count++
			item.Status = 0
			item.Save()
			switch item.Type { //要扣一下视频社区评论数
			case 1:
				vod.VodListModel.DelComments(item.CID)
			case 3:
				blog.BlogListModel.DelComments(item.CID)
			case 4:
				seximg.SeximgModel.DelComments(item.CID)
			}
			//获得缓存key
			cIdArray = append(cIdArray, "api:commentlist:"+strconv.FormatUint(uint64(item.CID), 10)+":"+strconv.FormatUint(uint64(item.Type), 10)+":")
		}
	}
	for _, item := range cIdArray {
		cacheKey := item
		redis.PullPrefix(cacheKey) //删除修改的列表缓存
	}
	redis.Pull("commenttext")
	c.JSON(http.StatusOK, gin.H{
		"message": "已完成,删除重复评论" + strconv.Itoa(count),
	})
}

// 通过接口给宝塔cron调用实现定时任务 删除违规评论
// 45分钟任务  查询评论超过10条的人
func JobDleteComment2(c *gin.Context) {
	lists := redis.CommentUserRepeatList()
	count := 0
	cUserCount := 0
	cIdArray := []string{} //用来刷新评论缓存
	for _, textitem := range lists {
		cUserCount++
		var comments []user.UserComment
		query := model.DataBase.Where("created_at > ?", time.Now().Add(-60*time.Minute).Format("2006-01-02 15:04:05"))
		query.Where("user_id = ?", textitem.Member).Where("status = ?", 1).Limit(500).Find(&comments)
		for _, item := range comments {
			count++
			item.Status = 0
			item.Save()
			switch item.Type { //要扣一下视频社区评论数
			case 1:
				vod.VodListModel.DelComments(item.CID)
			case 3:
				blog.BlogListModel.DelComments(item.CID)
			case 4:
				seximg.SeximgModel.DelComments(item.CID)
			}
			//获得缓存key
			cIdArray = append(cIdArray, "api:commentlist:"+strconv.FormatUint(uint64(item.CID), 10)+":"+strconv.FormatUint(uint64(item.Type), 10)+":")
		}
		userinfo := user.UserListModel.SelectByItem(textitem.Member)
		if userinfo != nil {
			userinfo.RejuectComment()
		}

	}
	for _, item := range cIdArray {
		cacheKey := item
		redis.PullPrefix(cacheKey) //删除修改的列表缓存
	}
	redis.Pull("commentuser")
	c.JSON(http.StatusOK, gin.H{
		"message": "已完成,删除评论" + strconv.Itoa(count) + "禁言" + strconv.Itoa(cUserCount),
	})
}
