package admin

import (
	"math/rand"
	"myadmin/model/user"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var vipListModel user.VipCodeList
var vipUserListModel user.VipUserList

func VipCodeList(c *gin.Context) {
	request := user.VipCodeListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := vipListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VipCodeListCreate(c *gin.Context) {
	request := user.VipCodeList{}
	c.ShouldBindJSON(&request)
	if len(request.Code) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "创建失败",
		})
		return
	}
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VipCodeListCreateMore(c *gin.Context) {
	request := user.VipCodeListCreateParam{}
	model := &user.VipCodeList{}
	c.ShouldBindJSON(&request)
	if request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "创建失败",
		})
		return
	}
	for i := 0; i < request.Amount; {
		code := generateRandomString(6)
		model = vipListModel.SelectByCode(code)
		if model.ID == 0 {
			model.Code = code
			model.Status = 0
			model.Save()
			i++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}
func generateRandomString(length int) string {
	// 定义字符集
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func VipUserList(c *gin.Context) {
	request := user.VipUserListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Order = "-id"
	list, total := vipUserListModel.VipUserList(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VipUserListCreate(c *gin.Context) {
	request := user.VipUserList{}
	c.ShouldBindJSON(&request)
	vipuser := vipUserListModel.VIPUserSelectByUserId(request.UserID)
	if vipuser.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户已是VIP",
		})
		return
	}
	vipuser1 := vipUserListModel.VIPUserSelectByVipId(request.VipID)
	if vipuser1.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "该VIP码已使用",
		})
		return
	}
	model := vipListModel.SelectByID(request.VipID)
	if model.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "VIP码错误",
		})
		return
	}
	request.VipCode = model.Code
	err := request.VIPUserSave()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	model.Status = 1
	model.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VipUserListUpdate(c *gin.Context) {
	request := user.VipUserList{}
	c.ShouldBindJSON(&request)
	request.VIPUserSave()
	//修改用户表的VIP状态
	user.UpdateIsVip(request.UserID, request.Status)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func VipUserListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, _ := strconv.Atoi(id)
	model := vipUserListModel.VIPUserSelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除失败",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
