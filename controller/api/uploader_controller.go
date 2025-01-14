package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io"
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 得到Up主状态
func UserGetUploader(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	rid, err := strconv.Atoi(userID)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	m := user.UserUploader{}
	model.DataBase.Where("user_id = ?", rid).First(&m)
	if m.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"status": m.State,
				"reason": m.ReReason,
			},
			"message": "成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"status": -1,
		},
		"message": "成功",
	})
}

// 申请UP主
func UserUploaderPost(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	uID, err := strconv.Atoi(userID)
	if err != nil || uID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	request := UserUploaderPostRequest{}
	c.ShouldBindJSON(&request)
	// 验证简介
	if len(request.Introduce) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "请填写个人简介",
		})
		return
	}
	// 得到用户
	userModel := user.UserList{}
	model.DataBase.Where("id = ?", uID).First(&userModel)
	if userModel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户异常",
		})
		return
	}
	// 验证短信验证码
	if lastSms := user.UserSmsModel.GetLastByPhone(userModel.Phone); lastSms != nil {
		if lastSms.Used != 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "短信验证失败",
			})
			return
		}
		if lastSms.Code != request.Verify {
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
	// 拿用户已经传过的数据覆盖 再次提交
	m := user.UserUploader{}
	model.DataBase.Where("user_id = ?", uID).First(&m)
	m.UserID = uint(uID)
	m.ImgPath = request.ImagePath
	m.Phone = userModel.Phone
	m.Text = request.Introduce
	m.State = 0
	m.ReReason = ""
	model.DataBase.Save(&m)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "提交成功",
	})
}

// 得到oss服务器签名直传过期时间
func getSTsEmpriTime() string {
	return time.Now().Add(5*time.Minute).UTC().Format("2006-01-02T15:04:05") + "Z"
	// return time.Unix(1698065068, 0).UTC().Format("2006-01-02T15:04:05") + "Z"
}
func getSTsEmpriTimeUnit() int64 {
	return time.Now().Add(5 * time.Minute).Unix()
	// return time.Unix(1698065068, 0).Unix()

}

type ConfigStruct struct {
	Conditions [][]string `json:"conditions"`
	Expiration string     `json:"expiration"`
}

func UserSTSTOken(c *gin.Context) {
	// userID := c.MustGet("UserID").(string)
	// m := user.UserUploader{}
	// model.DataBase.Where("user_id = ?", userID).First(&m)
	// if m.ID == 0 || m.State != 1 {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":    400,
	// 		"data":    nil,
	// 		"message": "认证UP主才能发送图片哦",
	// 	})
	// 	return
	// }

	var upload_dir string = "userupload/"

	expireStS := getSTsEmpriTime()
	var config ConfigStruct
	config.Expiration = expireStS
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)
	result, _ := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(os.Getenv("ALI_ACCESS_KEY_SECRET")))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	expireReturn := getSTsEmpriTimeUnit()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"accessid":  os.Getenv("ALI_ACCESS_KEY_ID"),
			"expire":    expireReturn,
			"dir":       upload_dir,
			"signature": string(signedStr),
			"policy":    string(debyte),
		},
		"message": "成功",
	})

}
