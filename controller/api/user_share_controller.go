package api

import (
	"myadmin/model"
	"myadmin/model/user"
	"myadmin/model/vip"
	"myadmin/util/sugar"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 提交成功邀请消息
// @Description - 提交成功邀请消息 邀请人，被邀请人UUID，数据签名。防止接口被发现刷接口。
// @Tags 用户
// @Security ApiKeyAuth
// @Param data body UserSharePostRequest true "参数列表"
// @Accept json
// @Router /api/user/share [post]
func UserShareAdd(c *gin.Context) {
	// 获取用户输入参数
	// var request UserSharePostRequest
	// c.ShouldBindJSON(&request)
	// key := strconv.Itoa(int(request.ShareUserID)) + time.Now().Format("-20060102") + "-mogu-" + request.DeviceId
	// md5Size := md5.Sum([]byte(key))
	// md5hex := hex.EncodeToString(md5Size[:])
	// if md5hex[1:7] != request.AccessKey {
	// 	c.JSON(http.StatusOK, gin.H{"message": "fail"})
	// 	return
	// }
	// user.NewUserShare(request.ShareUserID, request.DeviceId)
	// 给用户增加邀请积分

	c.JSON(http.StatusOK, gin.H{"message": "分享已过期,请更新APP"})
}

//新的分享列表
func UserShareListNew(c *gin.Context) {
	request := user.UserInviteParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	var count int64
	var list []user.UserInviteList
	if request.InviteUserID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"list":  nil,
			"total": 0,
		})
		return
	}
	query := model.DataBase.Model(user.UserInviteList{})
	query.Order("id desc")
	query.Preload("BeUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Name")
	})

	if request.InviteUserID != 0 {
		query.Where("invite_user_id = ?", request.InviteUserID)
	}
	query.Count(&count)
	result := query.Offset((request.Page - 1) * request.Limit).Limit(request.Limit).Find(&list)
	data := make([]interface{}, len(list))
	for index, item := range list {
		data[index] = struct {
			Name string `json:"name"`
			Time uint   `json:"time"`
		}{item.BeUser.Name, uint(item.CreatedAt.Unix())}
	}
	nilarray := []string{}
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"data": gin.H{
				"total": count,
				"list":  nilarray,
			},
			"message": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": count,
			"list":  data,
		},
		"message": "",
	})
}

// @Summary  老版本分享 已过期 不处理 不可删除
// @Description 我的邀请列表或邀请的设备列表
// @Tags 用户
// @Accept json
// @Param data query user.UserShareParam true "参数列表"
// @Router /api/user/share [get]
func UserShareList(c *gin.Context) {
	// 获取用户输入参数
	// var request user.UserShareParam
	// c.ShouldBindQuery(&request)
	// userId, _ := strconv.Atoi(c.MustGet("UserID").(string))
	// request.ShareUserID = uint(userId)
	// request.Order = "-id"
	// list, count := user.ShareList(request)

	// data := make([]interface{}, len(list))
	// for index, item := range list {
	// 	data[index] = struct {
	// 		ShareId  uint   `json:"share_id"`
	// 		DeviceId string `json:"device_id"`
	// 		Time     uint   `json:"time"`
	// 	}{item.ShareUserID, item.DeviceID, uint(item.CreatedAt.Unix())}
	// }
	nilarray := []string{}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": 0,
			"list":  nilarray,
		},
		"message": "",
	})
}

// @Summary VIP列表数据
// @Description 返回VIP列表数据
// @Tags VIP
// @Accept json
// @Param data query vip.VipListParam true "参数列表"
// @Router /api/vip/list [get]
func VipList(c *gin.Context) {
	// 判断是否登录用户
	uid := sugar.StringToUint(c.MustGet("UserID").(string))

	var request vip.VipListParam
	c.ShouldBindQuery(&request)
	request.OnStatus = true
	request.Order = "-id"
	list, total := vip.VipListData(request)

	rlist := make([]apiResultVipList, len(list))
	// 获取播放域名
	// host := util.GetGinRequestHostUrl(c)
	host := os.Getenv("PLIST_DOMAIN")
	for index, item := range list {
		rlist[index] = newApiResultVipList(item)
		if uid > 0 {
			rlist[index].Unlock = vip.IsUnlockVip(uid, item.ID)
			rlist[index].PlayUrl = host + rlist[index].PlayUrl
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  rlist,
		},
		"message": "",
	})
}

// @Summary 解锁新的VIP
// @Description 解锁一个新的视频VIP
// @Tags VIP
// @Accept json
// @Param vipid query int true "要解锁的VIPid"
// @Router /api/vip/unlock [post]
func VipUnlock(c *gin.Context) {
	uid := sugar.StringToUint(c.MustGet("UserID").(string))
	vid := sugar.StringToUint(c.Query("vipid"))
	if vid < 1 || uid < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "错误数据  ",
		})
		return
	}
	if vip.IsUnlockVip(uid, vid) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "已解锁",
		})
		return
	}
	u := user.SelectUserByID(uid)
	if u.Integral < 10 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "余额不足",
		})
		return
	}
	vip.UnlockVip(uint(uid), uint(vid))
	u.Integral -= 10
	u.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "解锁完成",
	})
}
