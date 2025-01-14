package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var uploadModel user.UserUploader

func UserUploaderList(c *gin.Context) {
	request := user.UserUploaderParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := uploadModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func UserUploaderUpdate(c *gin.Context) {
	var request user.UserUploader
	c.ShouldBindJSON(&request)
	Upmodel := uploadModel.SelectByID(request.ID)
	if Upmodel == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在",
		})
		return
	}
	userModel := user.UserList{}
	model.DataBase.Where("id = ?", Upmodel.UserID).First(&userModel)
	if userModel.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "用户异常",
		})
		return
	}
	if request.State == 1 { //通过

		model.DataBase.Model(&userModel).Update("avatar", Upmodel.ImgPath)
		model.DataBase.Model(&userModel).Update("introduction", Upmodel.Text)
		model.DataBase.Model(&userModel).Update("type", 3) //认证UP
	} else {
		model.DataBase.Model(&userModel).Update("type", 1) //注册用户
	}
	model.DataBase.Model(&Upmodel).Update("State", request.State)
	if request.State == 2 {
		model.DataBase.Model(&Upmodel).Update("ReReason", request.ReReason)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

// 批量通过审核
func AllowUp(c *gin.Context) {
	request := user.RequestUserAllowComment{}
	var uploaders []user.UserUploader
	c.ShouldBindJSON(&request)
	model.DataBase.Model(user.UserUploader{}).Where("id in ?", request.Ids).Where("State = ?", 0).Find(&uploaders)
	for _, item := range uploaders {
		userModel := user.UserList{}
		model.DataBase.Where("id = ?", item.UserID).First(&userModel)
		model.DataBase.Model(&userModel).Update("avatar", item.ImgPath)
		model.DataBase.Model(&userModel).Update("introduction", item.Text)
		model.DataBase.Model(&userModel).Update("type", 3) //认证UP
		model.DataBase.Model(&item).Update("State", 1)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func UserUploaderDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := uploadModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
