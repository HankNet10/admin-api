package admin

import (
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var userNoticeModel user.UserNotice

func UserNoticeList(c *gin.Context) {
	request := user.UserNoticeListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := userNoticeModel.AdminList(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
func UserNoticeUpdate(c *gin.Context) {
	var request user.UserNotice
	c.ShouldBindJSON(&request)
	model := userNoticeModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在",
		})
		return
	}
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}
func UserAdminNoticeCreate(c *gin.Context) {
	request := user.UserNoticeSendParam{}
	c.ShouldBindJSON(&request)
	if request.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请输入",
		})
		return
	}
	userNoticeModel.SendAdminNotice(request.Title, request.Text, uint(request.AppID), request.Link)
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
	})
}

func UserNoticeDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := userNoticeModel.SelectByID(uint(rid))
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
