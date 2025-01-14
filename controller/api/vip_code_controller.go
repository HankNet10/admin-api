package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"myadmin/model/user"
	"myadmin/util/redis"

	"github.com/gin-gonic/gin"
)

var vipListModel user.VipCodeList
var vipUserListModel user.VipUserList

// @Summary 激活VIP
// @Description 激活VIP
// @Tags VIP码
// @Security ApiKeyAuth
// @Accept json
// @Param data body VipWatchRequest true "参数列表"
// @Router /api/user/vip/watch [post]
func VipWatch(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request VipWatchRequest
	c.ShouldBindJSON(&request)
	if request.Vid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
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
	vipuser := vipUserListModel.VIPUserSelectByUserId(uint(uID))
	cKey := "api:video:watch-" + appID + ":" + strconv.Itoa(int(vipuser.UserID))
	vids := ""
	if jsonData, err := redis.Get(cKey); err == nil {
		vids = string(jsonData)
	}
	vidlist := strings.Split(vids, ",")
	total := len(vidlist)
	message := ""
	if len(vids) > 0 {
		isadd := true
		for _, value := range vidlist {
			id, _ := strconv.Atoi(value)
			if request.Vid == id {
				isadd = false
			}
		}
		if isadd {
			total += 1
			if total > 10 {
				message = "今日已免费观看10个VIP视频,请明天继续观看"
			} else {
				vids = vids + "," + strconv.Itoa(request.Vid)
			}
		}
	} else {
		vids = strconv.Itoa(request.Vid)
		total = 1
	}
	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"list":  vids,
			"total": total,
		},
		"message": message,
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	now := time.Now()
	// 将时间设置为今天23:59:59
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	aa := time.Duration(endOfDay.Unix()-time.Now().Unix()) * time.Second
	redis.Set(cKey, vids, aa)
	c.String(http.StatusOK, string(jsonData))
}

// @Summary 激活VIP
// @Description 激活VIP
// @Tags VIP码
// @Security ApiKeyAuth
// @Accept json
// @Param data body VipActivationRequest true "参数列表"
// @Router /api/user/vipcode/list [post]
func VipActivation(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request VipActivationRequest
	c.ShouldBindJSON(&request)
	if request.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
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
	model := vipUserListModel.VIPUserSelectByUserId(uint(uID))
	if model.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "您已激活VIP,请勿重复激活",
		})
		return
	}
	vipModel := vipListModel.SelectByCode(request.Code)
	if vipModel.ID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "无效的VIP码",
		})
		return
	}
	vipuser1 := vipUserListModel.VIPUserSelectByCode(request.Code)
	if vipuser1.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "无效的VIP码",
		})
		return
	}
	vipuser := user.VipUserList{}
	vipuser.UserID = uint(uID)
	vipuser.VipID = vipModel.ID
	vipuser.VipCode = vipModel.Code
	vipuser.Status = 1
	vipuser.VIPUserSave()
	vipModel.Status = 1
	vipModel.Save()

	cKey := "api:userinfo:vipcode-" + appID + ":" + userID
	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"code": request.Code,
		},
		"message": "",
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}
	redis.Set(cKey, jsonData, 60*time.Minute)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "激活成功",
	})
}

// @Summary 获取已激活的VIP码
// @Description 获取已激活的VIP码
// @Tags VIP码
// @Accept json
// @Param data query user.VipUserParam true "参数列表"
// @Router /api/user/vipcode/list [get]
func VipUserInfo(c *gin.Context) {
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
	cKey := "api:userinfo:vipcode-" + appID + ":" + userID
	if jsonData, err := redis.Get(cKey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	model := vipUserListModel.VIPUserSelectByUserId(uint(uID))

	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"code": model.VipCode,
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
