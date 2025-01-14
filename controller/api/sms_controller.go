package api

import (
	"math/rand"
	"myadmin/model/user"
	"myadmin/util/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ucloud/ucloud-sdk-go/services/usms"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

// @Summary 用户发送短信验证码
// @Description - 用户提交建议
// @Tags 用户
// @Accept json
// @Param data body UserSmsRequest true "参数"
// @Router /api/user/sms [post]
func SmsSubmit(c *gin.Context) {
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	var request UserSmsRequest
	c.ShouldBindJSON(&request)
	// 判断验证码是否正确。
	if code, err := redis.Pull("captcha-" + request.CaptchaId); err != nil || code != request.Captcha {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "验证码错误",
		})
		return
	}
	// 每天最多发5个
	if user.UserSmsModel.GetDayTotalByPhone(request.Phone) > 5 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "验证码已达上限",
		})
		return
	}
	var code = 0
	if lastSms := user.UserSmsModel.GetLastByPhone(request.Phone); lastSms != nil {
		if (time.Now().Unix() - lastSms.CreatedAt.Unix()) < 120 { // 120秒内只能发一次
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "发送频繁",
			})
			return
		}
		if lastSms.Used == 0 && (time.Now().Unix()-lastSms.CreatedAt.Unix()) < 900 {
			code, _ = strconv.Atoi(lastSms.Code) // 15分钟内验证码相同
		}
	}
	if code == 0 { // 没有旧密码则生成新的
		rand.Seed(time.Now().UnixMicro())
		code = rand.Intn(9000) + 1000
	}

	cfg := ucloud.NewConfig()
	cfg.BaseUrl = "https://api.ucloud.cn"
	cred := auth.NewCredential()
	cred.PublicKey = "9LzcmLtWEsCDHqUsxYSyc1JwZXcWYnfcU18ylxXAMA"
	cred.PrivateKey = "2amZNFFM4hcBBDo7Dr6CgHKfgn7EO4P6b8OqR93rvKzNIj9W4gbTt3FFWLnRoeCZp2"
	usmsClient := usms.NewClient(&cfg, &cred)
	req := usmsClient.NewSendUSMSMessageRequest()
	req.ProjectId = ucloud.String("org-ozwy0c")
	if appID == "1" {
		req.SigContent = ucloud.String("蘑菇")
	} else if appID == "2" {
		req.SigContent = ucloud.String("蓝莓计时器")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "错误的app id",
		})
		return
	}
	req.TemplateId = ucloud.String("UTA221029EJ6OIA")
	// 如果是蓝莓

	// 保存数据库
	sms := user.UserSms{}
	sms.Used = 0
	sms.Phone = request.Phone
	sms.Code = strconv.Itoa(code)

	req.PhoneNumbers = []string{sms.Phone}
	req.TemplateParams = []string{sms.Code}
	_, err := usmsClient.SendUSMSMessage(req)
	if err != nil {
		sms.Message = err.Error()
	}
	sms.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "发送失败,联系管理员",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "发送成功",
		})
	}
}
