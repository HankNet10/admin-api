package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserNameList(c *gin.Context) {
	request := user.UserNameListParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := user.UserNameListModel.List(request) //   .List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func UserNameUpdate(c *gin.Context) {
	request := user.UserName{}
	c.ShouldBindJSON(&request)
	id := request.UserId
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	userData := user.UserListModel.SelectByID(uint(rid))
	userData.Name = request.Name
	userData.Save()
	request.Status = 0
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func UserNameDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.UserNameListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "数据不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 批量通过此页
func UserNameAllowPage(c *gin.Context) {
	request := user.RequestUserAllowComment{}
	c.ShouldBindJSON(&request)
	var namelists []user.UserName
	model.DataBase.Model(user.UserName{}).Where("id in ?", request.Ids).Where("status = ?", 1).Find(&namelists)
	for _, item := range namelists {
		id, _ := strconv.Atoi(item.UserId)
		userData := user.UserListModel.SelectByID(uint(id))
		model.DataBase.Model(&userData).Update("name", item.Name)
		model.DataBase.Model(&item).Update("status", 0)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改完成",
	})
}
