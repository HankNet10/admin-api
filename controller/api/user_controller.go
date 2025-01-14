package api

import (
	"context"
	"encoding/json"
	"log"
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/model/suggest"
	"myadmin/model/user"
	"myadmin/model/vod"
	"myadmin/util/bcrypt"
	"myadmin/util/jwt"
	"myadmin/util/redis"
	"myadmin/util/sugar"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
)

// @Summary 注册用户
// @Description - 注册用户信息
// @Tags 用户
// @Accept json
// @Router /api/user/register [post]
func UserRegister(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	request := UserRegisterRequest{}
	c.ShouldBindJSON(&request)
	// 判断验证码是否正确。
	if appID == "" {
		appID = "1"
	}

	if appID == "1" { //蓝莓暂时不维护 只更新蘑菇
		code, err := redis.Pull("captcha-" + request.CaptchaId)
		if err != nil || code != request.Captcha {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "验证码错误",
			})
			return
		}
	}

	//检测恶意注册
	ckey1 := "creat1user:" + appID + "-" + c.ClientIP()
	cCount1, _ := redis.Get(ckey1) //控制5分钟1个ip 1个注册
	if cCount1 != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "注册过于频繁,稍后再来注册",
		})
		return
	}
	//检测恶意注册
	ckey := "creatuser:" + appID + "-" + c.ClientIP()
	cCount, _ := redis.Get(ckey) //控制24小时1个ip 10个注册
	requestCount := 0
	if cCount != "" {
		count, _ := strconv.Atoi(cCount)
		requestCount = count
		if count >= 10 {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "注册过于频繁,稍后再来注册",
			})
			return
		}
	}

	// 保存历史密码
	historyPasswd := user.UserPasswd{
		Phone:    request.Phone,
		Password: request.Password,
	}
	model.DataBase.Save(&historyPasswd)
	// 判断 手机号 长度
	if len(request.Phone) > 11 || len(request.Phone) < 7 { //支持国际手机 7到11位
		c.JSON(http.StatusOK, gin.H{
			"code":    403,
			"data":    nil,
			"message": "手机格式错误",
		})
		return
	}
	// 判断 用户密码长度
	if len(request.Password) < 6 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "密码长度不对",
		})
		return
	}
	var isInvite = false
	inviteUser := user.UserList{}

	if len(request.InviteCode) != 6 && len(request.InviteCode) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "错误邀请码,无邀请码可不填写,不影响观看体验",
		})
		return
	} else if len(request.InviteCode) == 6 {
		model.DataBase.Where("invite_code = ? and app_id = ?", request.InviteCode, appID).First(&inviteUser)
		if inviteUser.ID != 0 {
			uskey := "usershare:" + appID + "-" + strconv.Itoa(int(inviteUser.ID))
			usValue, _ := redis.Get(uskey)
			if usValue != "" {
				if usValue != c.ClientIP() {
					isInvite = true
				}
			}
		} //没查到就算啦也放行。不添加邀请关系即可
	}

	// 查询老用户
	m := user.UserList{}
	model.DataBase.Where("phone = ? and app_id = ?", request.Phone, appID).Find(&m)
	if m.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户已存在",
		})
		return
	}
	// 创建用户
	createUser := user.UserList{
		Phone:    request.Phone,
		Type:     1, //注册用户
		AppID:    sugar.StringToUint(appID),
		Ip:       c.ClientIP(),
		Password: bcrypt.GeneratePassword(request.Password),
	}
	result := model.DataBase.Save(&createUser)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "创建用户失败，请提交反馈",
		})
		return
	}
	//生成自己的InviteCode
	createUser.InviteCode = sugar.GetInvCodeByUIDUnique(uint64(createUser.ID))
	model.DataBase.Model(&createUser).Update("invite_code", createUser.InviteCode)
	//添加邀请关系
	if isInvite {
		InviteModel := user.UserInviteList{
			InviteUserID:   inviteUser.ID,
			InviteCode:     inviteUser.InviteCode,
			BeInviteUserID: createUser.ID,
		}
		model.DataBase.Save(&InviteModel)
		//添加邀请关系

		inviteUser.Integral += 10 //添加邀请者积分
		model.DataBase.Model(&inviteUser).Update("integral", inviteUser.Integral)

		createUser.Integral += 10 //添加新用户积分
		model.DataBase.Model(&createUser).Update("integral", createUser.Integral)

	}
	redis.SetNoChangeTTl(ckey, requestCount+1, 24*time.Hour)
	redis.Set(ckey1, 1, 5*time.Minute)
	// 返回用户的Token
	exp := time.Now().AddDate(1, 0, 0).Unix()
	sub := strconv.Itoa(int(createUser.ID))
	if token, err := jwt.EncodeToken(sub, exp); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "生成token失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": jwt.AccessToken{
				AccessToken: token,
				TokenType:   jwt.TokenType,
				ExpiresIn:   exp,
			},
			"user":    createUser,
			"message": "",
		})
	}

}

// **用户注销
func UserUnRegister(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	query := model.DataBase.Model(blog.BlogList{})
	query.Where("user_id = ?", uID)
	var list []*blog.BlogList
	// 获取总数据
	var count int64
	query.Count(&count)
	query.Find(&list)
	//先查社区有没有帖子 有帖子不能注销，请先处理发帖
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "您还有未删除的社区帖,请先处理。",
		})
		return
	}
	userinfo := user.UserListModel.SelectByID(uint(uID))

	//修改phone为000000000000 12个0 /不删除用户 避免其他数据异常
	userinfo.Phone = "000000000000"
	userinfo.Password = bcrypt.GeneratePassword("000000000000")
	userinfo.Type = 1
	userinfo.Name = "默认用户"
	userinfo.Avatar = ""
	userinfo.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "您的账户已注销",
	})
	//Password: bcrypt.GeneratePassword(000000000000) //密码修改为12个0
	//修改用户名为注销用户，改用户类型为1 已注册用户

}

// @Summary 重设密码
// @Tags 用户
// @Accept json
// @Param data body UserRegisterRequest  true "参数列表"
// @Router /api/user/resetpasswd [post]
func UserResetPasswd(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	request := UserRegisterRequest{}
	c.ShouldBindJSON(&request)
	// 判断用户手机号格式
	if len(request.Phone) != 11 {
		c.JSON(http.StatusOK, gin.H{
			"code":    403,
			"data":    nil,
			"message": "手机格式错误",
		})
		return
	}
	// 保存历史密码
	historyPasswd := user.UserPasswd{
		Phone:    request.Phone,
		Password: request.Password,
	}
	model.DataBase.Save(&historyPasswd)
	// 验证密码长度
	if len(request.Password) < 6 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "密码长度不对",
		})
		return
	}
	// 验证短信验证码
	if lastSms := user.UserSmsModel.GetLastByPhone(request.Phone); lastSms != nil {
		if lastSms.Used != 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "短信验证失败",
			})
			return
		}
		if lastSms.Code != request.Verification {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "短信验证错误",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "短信验证获取失败",
		})
		return
	}
	// 判断用户是否存在
	m := user.UserList{}
	model.DataBase.Where("phone = ? and app_id = ?", request.Phone, appID).Find(&m)
	if m.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户不存在,请先注册。",
		})
		return
	}
	model.DataBase.Model(&m).Update("password", bcrypt.GeneratePassword(request.Password))

	// 返回token
	exp := time.Now().AddDate(1, 0, 0).Unix()
	sub := strconv.Itoa(int(m.ID))
	if token, err := jwt.EncodeToken(sub, exp); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "生成token失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": jwt.AccessToken{
				AccessToken: token,
				TokenType:   jwt.TokenType,
				ExpiresIn:   exp,
			},
			"message": "",
		})
	}
}

// @Summary 用户登录
// @Description - 用户登录
// @Tags 用户
// @Accept json
// @Param data body UserLoginRequest  true "参数列表"
// @Router /api/user/login [post]
func UserLogin(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	request := UserLoginRequest{}
	c.ShouldBindJSON(&request)
	// 判断验证码是否正确。
	code, err := redis.Pull("captcha-" + request.CaptchaId)
	if err != nil || code != request.Captcha {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "验证码错误",
		})
		return
	}
	// 判断用户手机号是否存在
	m := user.UserList{}
	model.DataBase.Where("phone = ? and app_id = ?", request.Phone, appID).Find(&m)
	if m.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户或密码错误",
		})
		return
	}
	// 验证密码是否正确
	if !bcrypt.ComparePassword(request.Password, m.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户密码错误",
		})
		return
	}
	//修改IP
	model.DataBase.Model(&m).Update("ip", c.ClientIP())
	// 	创建新的用户
	// 返回用户的Token
	exp := time.Now().AddDate(1, 0, 0).Unix()
	sub := strconv.Itoa(int(m.ID))
	if token, err := jwt.EncodeToken(sub, exp); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "生成token失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": jwt.AccessToken{
				AccessToken: token,
				TokenType:   jwt.TokenType,
				ExpiresIn:   exp,
			},
			"user":    m,
			"message": "",
		})
	}
}

// @Summary 用户详情
// @Description - 用户详细信息
// @Tags 用户
// @Security ApiKeyAuth
// @Accept json
// @Router /api/user/info [get]
func UserInfo(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	userModel := user.UserList{}
	rid, err := strconv.Atoi(userID)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	userData := userModel.SelectByID(uint(rid))
	if userData == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户状态异常",
		})
		return
	}
	avatar := ""
	if userData.Avatar != "" {
		avatar = os.Getenv("ALI_OSS_DOMAIN") + "/" + userData.Avatar
	}
	data := gin.H{
		"id":         userData.ID,
		"type":       userData.Type,
		"name":       userData.Name,
		"avatar":     avatar,
		"follower":   user.UserFollowModel.Follower(userData.ID),
		"integral":   userData.Integral,
		"inviteCode": userData.InviteCode,
		"phone":      userData.Phone,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

// @Summary 修改用户
// @Description - 修改用户信息
// @Tags 用户
// @Security ApiKeyAuth
// @Param data body UserEditRequest true "参数列表"
// @Accept json
// @Router /api/user/edit [put]
func UserEdit(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	rid, err := strconv.Atoi(userID)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	userData := user.UserListModel.SelectByID(uint(rid))
	if userData == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户状态异常",
		})
		return
	}
	// if userData.Type == 3 {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":    400,
	// 		"data":    nil,
	// 		"message": "认证UP主,不能再次修改昵称,请联系官方修改",
	// 	})
	// 	return
	// }
	// 获取用户输入参数
	var request UserEditRequest
	c.ShouldBindJSON(&request)
	reqstNameLen := len(request.Name)
	if reqstNameLen > 0 && reqstNameLen < 48 {
		// reg := regexp.MustCompile("^[a-zA-Z0-9\\p{Han}]+$") //只有中文字和英文字母
		// if !reg.MatchString(request.Name) {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"code":    400,
		// 		"data":    nil,
		// 		"message": "姓名中含有违规字词,请修改",
		// 	})
		// 	return
		// }

		// var list []user.DirtyWord
		// query := model.DataBase.Model(user.DirtyWord{})
		// query.Find(&list)
		// if containsSensitiveWords(request.Name, list) {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"code":    400,
		// 		"data":    nil,
		// 		"message": "姓名中含有违规字词,请修改",
		// 	})
		// 	return
		// }
		// updateUser.Name = request.Name
		// //改动到审核姓名逻辑
		nameItem := user.UserNameListModel.SelectByUserID(userID)
		if nameItem == nil {
			nameItem := user.UserName{}
			nameItem.UserId = userID
			nameItem.Name = request.Name
			nameItem.Status = 1
			nameItem.Save()
		} else {
			nameItem.Name = request.Name
			nameItem.Status = 1
			nameItem.Save()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "提交昵称成功,审核通过后显示",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "昵称长度不合规",
		})
	}
	// userData.Updates(updateUser)
	// c.JSON(http.StatusOK, gin.H{
	// 	"code":    200,
	// 	"data":    nil,
	// 	"message": "修改成功",
	// })
}
func containsSensitiveWords(text string, sensitiveWords []user.DirtyWord) bool {
	for _, word := range sensitiveWords {
		if strings.Contains(text, word.Name) {
			return true
		}
	}
	return false
}

// @Summary 其他用户详情
// @Description - 查看其他的用户详细信息 - 可以登录也可以不登陆
// @Tags 用户
// @Security ApiKeyAuth
// @Accept json
// @Param id query int true "参数列表"
// @Router /api/user/other [get]
func UserOther(c *gin.Context) {
	followid, _ := strconv.Atoi(c.Query("id"))
	if followid < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "查看的用户ID不存在",
		})
		return
	}
	userModel := user.UserList{}
	userData := userModel.SelectByID(uint(followid))
	if userData == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "查看的用户不存在",
		})
		return
	}

	// 判断是否关注当前用户。
	userID := c.MustGet("UserID").(string)
	floowed := false
	if uid, _ := strconv.Atoi(userID); uid > 0 {
		floowed = user.UserFollowModel.IsFollow(uint(uid), uint(followid))
	}

	data := gin.H{
		"id":           userData.ID,
		"type":         userData.Type,
		"name":         userData.Name,
		"avatar":       os.Getenv("ALI_OSS_DOMAIN") + "/" + userData.Avatar,
		"follower":     user.UserFollowModel.Follower(userData.ID),
		"gender":       userData.Gender,
		"birthday":     userData.Birthday,
		"weight":       userData.Weight,
		"height":       userData.Height,
		"introduction": userData.Introduction,
		"is_floowed":   floowed,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

// @Summary 女优列表
// @Description 返回用户信息
// @Tags 用户
// @Accept json
// @Param data query user.UserListParam true "参数列表"
// @Router /api/user/actor [get]
func UserActor(c *gin.Context) {

	request := user.UserListParam{}
	request.Limit = 999
	request.Page = 1
	request.Type = "11"
	list, total := user.UserListModel.SelectList(request)

	labels := make([]apiPinyinUserShow, len(list))
	arg := pinyin.NewArgs()
	for i, v := range list {
		var py string
		pydata := pinyin.Pinyin(v.Name, arg)
		if len(pydata) > 0 {
			if len(pydata[0]) > 0 {
				py = pydata[0][0]
			}
		}
		labels[i] = apiPinyinUserShow{
			Pinyin: py,
			User:   newApiResultUserShow(*v),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  labels,
		},
		"message": "",
	})
}

// @Summary 搜索用户列表
// @Description 返回用户信息
// @Tags 用户
// @Accept json
// @Param data query user.UserListParam true "参数列表"
// @Router /api/user/list [get]
func UserList(c *gin.Context) {
	request := user.UserListParam{}
	c.ShouldBindQuery(&request)
	if request.Limit > 20 || request.Limit < 0 {
		request.Limit = 20
	}
	if request.Page < 1 {
		request.Page = 1
	}
	list, total := user.UserListModel.SelectMainList(request)

	labels := make([]apiResultUserShow, len(list))
	for i, v := range list {
		labels[i] = newApiResultUserShow(*v)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  labels,
		},
		"message": "",
	})
}

// @Summary 我关注的，关注我的用户列表
// @Description 返回用户信息
// @Tags 用户
// @Accept json
// @Param data query user.UserFollowParam true "参数列表"
// @Router /api/user/follow [get]
func UserFollow(c *gin.Context) {
	// 我关注的，关注我的用户列表
	request := user.UserFollowParam{}
	c.ShouldBindQuery(&request)
	request.Limit = 20
	if request.Page < 1 {
		request.Page = 1
	}
	list, total := user.SelectFollower(request)

	labels := make([]apiResultUserShow, len(list))
	if request.UserID > 0 {
		for i, v := range list {
			labels[i] = newApiResultUserShow(v.FUser)
		}
	} else if request.FollowID > 0 {
		for i, v := range list {
			labels[i] = newApiResultUserShow(v.User)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  labels,
		},
		"message": "",
	})
}

// 抽奖
func UserGetGift(c *gin.Context) {
	// 我关注的，关注我的用户列表
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "抽奖错误,请更新APP,或者登录使用。",
		})
		return
	}
	userModel := user.UserList{}
	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "参数错误",
		})
		return
	}
	result := model.DataBase.First(&userModel, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "抽奖错误,请更新APP,或者登录使用。",
		})
		return
	}
	if userModel.Integral < 50 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "积分不足",
		})
		return
	}
	if userModel.Integral >= 50 {
		userModel.Integral -= 50
		userModel.Save()
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "成功",
		})

	}
}

// @Summary 关注其他用户
// @Description 关注其他用户
// @Tags 用户
// @Accept json
// @Security ApiKeyAuth
// @Param data body UserFollowAddRequest true "参数列表"
// @Router /api/user/follow/add [put]
func UserAddFollow(c *gin.Context) {
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	request := UserFollowAddRequest{}
	c.ShouldBindJSON(&request)
	if request.FollowID < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "输入有效的关注ID",
		})
		return
	}

	user.UserFollowModel.AddFollow(uint(uid), request.FollowID)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "完成",
	})
}

// @Summary 用户列表
// @Description 返回用户信息
// @Tags 用户
// @Accept json
// @Param id query int true "参数列表"
// @Security ApiKeyAuth
// @Router /api/user/follow/delete [delete]
func UserDeleteFollow(c *gin.Context) {
	fid, _ := strconv.Atoi(c.Query("id"))
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	if user.UserFollowModel.DeleteFollow(uint(uid), uint(fid)) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "取消成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "取消失败",
		})
	}
}

// @Summary 用户提交建议
// @Description - 用户提交建议
// @Tags 用户
// @Security ApiKeyAuth
// @Accept json
// @Param suggest body string true "suggest"
// @Router /api/suggest/add [post]
func SuggestAdd(c *gin.Context) {
	request := struct {
		Suggest string
		Device  string
		Version string
		UserId  string
	}{}
	c.ShouldBindJSON(&request)
	m := suggest.SuggestList{
		Ip:      c.ClientIP(),
		Comment: request.Suggest,
		UserId:  request.UserId,
		Version: request.Version,
		Device:  request.Device,
	}
	m.Save()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "",
	})
}

// 验证用户数据 是否为
func UserVerify(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	rid, err := strconv.Atoi(userID)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "参数错误",
		})
		return
	}
	var playcount int64
	query := model.DataBase.Model(vod.VodHistory{}).Where("user_id = ?", rid).Limit(20)
	query.Count(&playcount)
	var starcount int64
	query2 := model.DataBase.Model(vod.VodStar{}).Where("user_id = ?", rid).Limit(20)
	query2.Count(&starcount)
	//收藏记录大于5,播放历史记录大于15 判断优质用户 自动通过评论
	userModel := user.UserList{}
	userData := userModel.SelectByID(uint(rid))
	if playcount >= 15 && starcount >= 5 && userData.DenyComment != 1 {
		user.AllowComment(uint(rid))
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "",
	})
}

// @Description - 用户提交举报
func ReportAdd(c *gin.Context) {
	request := UserReportRequest{}
	c.ShouldBindJSON(&request)
	m := user.Report{
		Text:           request.Text,
		CType:          request.CType,
		RType:          request.RType,
		ReporterUserID: request.ReportId,
		ReportedUserID: request.ReportedId,
		CId:            request.CId,
	}
	m.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "感谢您的举报,24小时内会进行核实",
	})
}

// @Summary 用户签到
// @Description 返回操作结果
// @Tags 签到
// @Accept json
// @Router /api/user/signday [get]
func SignDay(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	ckey := "signday:" + appID + "-" + strconv.Itoa(uid)
	if jsonData, err := redis.Get(ckey); err == nil {
		c.String(http.StatusOK, string(jsonData))
		return
	}
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := now.AddDate(0, 0, offset) //time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	weekMondayInt, _ := strconv.Atoi(weekStartDate.Format("20060102"))
	list, _ := user.UserSignListModel.SelectByUserID(uid, weekMondayInt, 7)
	labels := []uint{0, 0, 0, 0, 0, 0, 0}
	for _, value := range list {
		newTime, err := time.Parse("20060102", strconv.Itoa(int(value.SignDay)))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "数据错误",
			})
			return
		}
		days := newTime.Weekday()
		if days <= 0 {
			days = 7
		}
		labels[days-1] = 1
	}

	msg := "成功"
	resultData := gin.H{
		"code": 200,
		"data": gin.H{
			"total": 7,
			"list":  labels,
		},
		"message": msg,
	}
	jsonData, err := json.Marshal(resultData)
	if err != nil {
		log.Panic("jsonData set err", err)
	}

	// 计算当天 0 点的时间
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

	// 计算到 0 点的秒数
	duration := midnight.Sub(now)
	seconds := int64(duration.Seconds())

	redis.Set(ckey, jsonData, time.Duration(seconds)*time.Second)
	c.String(http.StatusOK, string(jsonData))
}

func SignIn(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	ckey := "signday:" + appID + "-" + strconv.Itoa(uid)
	userList := user.UserListModel.SelectByID(uint(uid))
	if userList == nil || userList.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "获取用户信息错误",
		})
		return
	} else {
		now := time.Now()
		day, _ := strconv.Atoi(now.Format("20060102"))
		offset := int(time.Monday - now.Weekday())
		if offset > 0 {
			offset = -6
		}
		weekStartDate := now.AddDate(0, 0, offset) //time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		weekMondayInt, _ := strconv.Atoi(weekStartDate.Format("20060102"))
		list, _ := user.UserSignListModel.SelectByUserID(uid, weekMondayInt, 1)
		userSign := user.UserSignList{}
		userSign.UserId = uint(uid)
		userSign.SignDay = uint(day)
		userSign.SignType = 1
		if len(list) > 0 && list[0].SignDay == uint(day-1) {
			userSign.Days = list[0].Days + 1
		} else if len(list) > 0 && list[0].SignDay == uint(day) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "今日已签到",
			})
			return
		} else {
			userSign.Days = 1
		}
		userSign.Save()
		integral := uint(10)
		if userSign.Days == 7 {
			integral = 20
		}
		userList.Integral += integral //添加用户积分
		model.DataBase.Model(&userList).Update("integral", userList.Integral)
		msg := "签到成功,获得积分" + strconv.FormatUint(uint64(integral), 10)
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"integral": integral,
			"message":  msg,
		})
		redis.Redis.Del(context.Background(), ckey)
	}
}

// @Summary 添加用户分享记录
// @Description 添加用户分享记录
// @Tags 添加用户分享记录
// @Accept json
// @Router /api/user/newusershare [get]
func NewUserShare(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	ckey1 := "usershare:" + appID + "-" + strconv.Itoa(uid)
	redis.Set(ckey1, c.ClientIP(), 12*time.Hour)
	user.UserShareModel.NewUserShare(uint(uid), c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "完成",
	})
}

func TaskListList(c *gin.Context) {
	request := user.TaskListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Order == "" {
		request.Order = "-sort"
	}
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	request.Status = 1
	request.AppID, _ = strconv.Atoi(appID)
	list, total := user.TaskListModel.TaskList(request)
	data := make([]interface{}, len(list))
	host := os.Getenv("ALI_OSS_DOMAIN")
	for index, item := range list {
		data[index] = struct {
			ID        uint   `json:"id"`
			Title     string `json:"title"`
			Image     string `json:"image"`
			BgImage   string `json:"bgimage"`
			Introduce string `json:"introduce"`
		}{item.ID, item.Title, host + "/" + item.Image, host + "/" + item.BgImage, item.Introduce}
	}
	count := user.TaskUserModel.TaskUserToDayCount()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":          data,
			"total":         total,
			"taskfinishnum": count,
		},
		"message": "",
	})
}

func TaskListInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	taskID, _ := strconv.Atoi(id)
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	TlModel := user.TaskListModel.TaskListSelectByID(uint(taskID))
	TuModel := user.TaskUserModel.TaskUserSelectByTaskID(uint(taskID), uint(uid))
	var rdata apiResultTaskUser
	if TuModel != nil {
		rdata = newApiResultTaskUser(*TuModel)
	}
	host := os.Getenv("ALI_OSS_DOMAIN")
	data := struct {
		ID        uint              `json:"id"`
		BgImage   string            `json:"bgimage"`
		Image     string            `json:"image"`
		Title     string            `json:"title"`
		Introduce string            `json:"introduce"`
		Taskuser  apiResultTaskUser `json:"taskuser"`
	}{TlModel.ID, host + "/" + TlModel.BgImage,
		host + "/" + TlModel.Image, TlModel.Title, TlModel.Introduce, rdata}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "",
	})
}

func TaskUserInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "缺少参数id",
		})
		return
	}
	taskID, _ := strconv.Atoi(id)
	uid, _ := strconv.Atoi(c.MustGet("UserID").(string))
	data := user.TaskUserModel.TaskUserSelectByTaskID(uint(taskID), uint(uid))
	var rdata apiResultTaskUser
	if data != nil {
		rdata = newApiResultTaskUser(*data)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    rdata,
		"message": "",
	})
}

func newApiResultTaskUser(v user.TaskUser) apiResultTaskUser {
	images := make([]apiResultTaskUserImages, len(v.Images))
	videopath := ""
	for i, vv := range v.Images {
		if vv.Type == 0 {
			images[i] = apiResultTaskUserImages{
				ID:   vv.ID,
				Type: vv.Type,
				Path: os.Getenv("ALI_OSS_DOMAIN") + "/" + vv.Path,
			}
		} else if vv.Type == 1 {
			images = append(images[:len(v.Images)-1])
			videopath = os.Getenv("ALI_OSS_DOMAIN") + "/" + vv.Path
		}
	}
	return apiResultTaskUser{
		ID:        v.ID,
		Images:    images,
		Videopath: videopath,
		CreatedAt: v.CreatedAt.Unix(),
		Content:   v.Content,
		Status:    uint(v.Status),
		Refuse:    v.Refuse,
	}
}

func TaskUserAdd(c *gin.Context) {
	var request TaskUserCreateUserRequest
	c.ShouldBindJSON(&request)
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	if request.TaskID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "参数错误",
		})
		return
	}
	userID := c.MustGet("UserID").(string)
	appid, _ := strconv.Atoi(appID)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	taskUser := user.TaskUser{}
	taskUser.TaskID = request.TaskID
	taskUser.AppID = uint(appid)
	taskUser.UserID = uint(uID)
	taskUser.Content = request.Content
	taskUser.Status = 1
	taskUser.Refuse = ""
	taskUser.TaskUserSave()
	if taskUser.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "创建失败",
		})
		return
	}
	for _, item := range request.Images {
		Image := user.TaskImage{
			UserID:     uint(uID),
			Type:       0,
			TaskUserID: taskUser.ID,
			Path:       item,
		}
		Image.TaskImageSave()
	}
	if request.Videopath != "" {
		Image := user.TaskImage{
			UserID:     uint(uID),
			Type:       1,
			TaskUserID: taskUser.ID,
			Path:       request.Videopath,
		}
		Image.TaskImageSave()
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "添加完成",
	})
}
