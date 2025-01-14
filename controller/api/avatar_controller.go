package api

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 更换头像
func UserUpdateAvatar(c *gin.Context) {
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
	// 获取用户输入参数
	var request UserAvatarRquestParam
	c.ShouldBindJSON(&request)
	println("requesrequestrequestrequestt")
	println(request.Id)
	if request.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "修改错误",
		})
		return
	}
	var UserAvatardata user.UserAvatar
	result := model.DataBase.First(&UserAvatardata, request.Id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "修改错误",
		})
		return
	}
	model.DataBase.Model(&userData).Update("avatar", UserAvatardata.Path)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "修改成功",
	})
}

// 获取用户头像列表
func UserGetAvatar(c *gin.Context) {
	request := user.UserAvatarListParam{}
	var UserAvatarList user.UserAvatar
	c.BindQuery(&request)
	request.Limit = 100
	request.Page = 1
	list, total := UserAvatarList.UserList(request)
	data := make([]apiResultAvatarList, len(list))
	for i, v := range list {
		data[i] = apiResultAvatarList{
			ID:   v.ID,
			Path: os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Path,
			Sort: v.Sort,
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
